package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/frankie-mur/go-rss/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (app *application) readinessHandler(c echo.Context) error {
	data := map[string]string{
		"status": "OK",
	}
	return c.JSON(200, data)
}

type createUserRequest struct {
	Name string `json:"name"`
}

func (app *application) createUserHandler(c echo.Context) error {
	/* Need to also check len(payload.Name) because default behavior
	does not error if field is not present */
	var req createUserRequest
	if err := decodeJson(c.Request(), &req); err != nil || len(req.Name) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	//Using name we can create a new user
	user, err := app.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      req.Name,
	})
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

func (app *application) getUserByApiKeyHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	//Send back user to client
	respondWithJSON(w, http.StatusOK, u)
}

type createFeedRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type createFeedResponse struct {
	Feed        database.Feed       `json:"feed"`
	Feed_Follow database.FeedFollow `json:"feed_follow"`
}

func (app *application) createFeedHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	var req createFeedRequest
	if err := decodeJson(r, &req); err != nil || len(req.Name) == 0 || len(req.Url) == 0 {
		respondWithError(w, http.StatusBadRequest, "Bad Request")
		return
	}
	//Create a new feed
	feed, err := app.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      req.Name,
		Url:       req.Url,
		UserID:    u.ID,
	})
	if err != nil {
		fmt.Printf("failed with error %v", err)
		//TODO: handle duplicate key error (violates uniqueness)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	//Create a feed follow
	feed_follow, err := app.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    u.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		fmt.Printf("failed with error %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	respondWithJSON(w, http.StatusCreated, createFeedResponse{
		Feed:        feed,
		Feed_Follow: feed_follow,
	})
}

func (app *application) getAllFeedsHandler(w http.ResponseWriter, r *http.Request) {
	data, err := app.DB.GetAllFeeds(r.Context())
	if err != nil {
		fmt.Printf("failed with error %v", err)
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

	data, err := app.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
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

func (app *application) deleteFeedFollowHandler(w http.ResponseWriter, r *http.Request, u database.User) {
	param := chi.URLParam(r, "id")
	feedFollowID, err := uuid.Parse(param)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feed follow ID")
		return
	}
	// if len(param) == 0 {
	// 	respondWithError(w, http.StatusBadRequest, "Provide feed follow ID")
	// 	return
	// }
	fmt.Printf("Param %s\n", feedFollowID)
	feed_id, err := uuid.Parse(param)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feed ID")
		return
	}
	err = app.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		FeedID: feed_id,
		UserID: u.ID,
	})
	if err != nil {
		fmt.Printf("failed with error %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	respondWithJSON(w, http.StatusOK, "OK")
}

func (app *application) getAllFeedFollows(w http.ResponseWriter, r *http.Request, u database.User) {
	feed_follows, err := app.DB.GetAllFeedFollows(r.Context(), u.ID)
	if err != nil {
		fmt.Printf("failed with error %v", err)
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	respondWithJSON(w, http.StatusOK, feed_follows)
}

func (app *application) getPostsByUserHandler(e echo.Context) error {
	param := e.Param("limit")
	var limit int
	var err error
	fmt.Printf("Param: %v", param)
	if len(param) == 0 {
		limit = 15
	} else {
		limit, err = strconv.Atoi(param)
		if err != nil {
			echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	posts, err := app.DB.GetPostsByUserId(e.Request().Context(), database.GetPostsByUserIdParams{
		//Hard code in a UUID
		UserID: uuid.MustParse("8a2de18d-4813-431e-a038-38dac55e22d8"),
		Limit:  int32(limit),
	})
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return e.Render(http.StatusOK, "posts", posts)
}
