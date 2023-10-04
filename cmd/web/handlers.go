package main

import (
	"context"
	"encoding/json"
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

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var payload createUserRequest
	/* Need to also check len(payload.Name) because default behavior
	does not error if field is not present */
	if err := decoder.Decode(&payload); err != nil || len(payload.Name) == 0 {
		respondWithError(w, http.StatusBadRequest, "Bad request")
		return
	}

	user, err := app.DB.CreateUser(context.TODO(), database.CreateUserParams{
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
