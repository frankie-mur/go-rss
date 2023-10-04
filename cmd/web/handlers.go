package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/frankie-mur/go-rss/internal/database"
	"github.com/google/uuid"
)

func (app *application) readinessHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "OK",
	}
	err := respondWithJSON(w, 200, data)
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) errorHandler(w http.ResponseWriter, r *http.Request) {
	err := respondWithError(w, 500, "Internal Server Error")
	if err != nil {
		log.Fatal(err)
	}
}

type createUserRequest struct {
	Name string `json:"name"`
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload createUserRequest
	/* Need to also check len(payload.Name) because default behavior
	does not error if field is not present */
	if err := decoder.Decode(&payload); err != nil || len(payload.Name) == 0 {
		respondWithError(w, http.StatusBadRequest, "Bad request")
		return
	}
	//Using name we can create a new user
	user, err := app.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      payload.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func (app *application) getUserByApiKeyHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	//Validate authorization header is in correct format
	apikey, err := validateAuthHeader(authHeader)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	//Get the user
	user, err := app.DB.GetUserByApiKey(context.Background(), apikey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	//Send back user to client
	respondWithJSON(w, http.StatusOK, user)
}
