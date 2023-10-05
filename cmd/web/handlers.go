package main

import (
	"context"
	"encoding/json"
	"fmt"
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

func (app *application) getUserByApiKeyHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	//Send back user to client
	respondWithJSON(w, http.StatusOK, u)
}

type createFeedRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (app *application) createFeedHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	decoder := json.NewDecoder(r.Body)
	fmt.Println("here")
	var req createFeedRequest
	if err := decoder.Decode(&req); err != nil || req.Name == "" || req.Url == "" {
		respondWithError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	data, err := app.DB.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      req.Name,
		Url:       req.Url,
		UserID:    u.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithJSON(w, http.StatusCreated, data)
}

func (app *application) getAllFeedsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := app.DB.GetAllFeeds(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithJSON(w, http.StatusOK, data)

}
