package main

import (
	"context"
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
	/* Need to also check len(payload.Name) because default behavior
	does not error if field is not present */
	var req createUserRequest
	if err := decodeJson(r, &req); err != nil || len(req.Name) == 0 {
		respondWithError(w, http.StatusBadRequest, "Bad request")
		return
	}
	//Using name we can create a new user
	user, err := app.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      req.Name,
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
	var req createFeedRequest
	if err := decodeJson(r, &req); err != nil || len(req.Name) == 0 || len(req.Url) == 0 {
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
		//TODO: handle duplicate key error (violates uniqueness)
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

type createFeedFollowRequest struct {
	Feed_id uuid.UUID `json:"feed_id"`
}

func (app *application) createFeedFollowHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	var req createFeedFollowRequest
	if err := decodeJson(r, &req); err != nil || len(req.Feed_id) == 0 {
		fmt.Printf("failed with error %v", err)
		respondWithError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	data, err := app.DB.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    u.ID,
		FeedID:    req.Feed_id,
	})
	if err != nil {
		fmt.Printf("failed with error %v", err)
		//TODO: Need to handle error when feed_id does not match to a feed ID
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithJSON(w, http.StatusCreated, data)
}
