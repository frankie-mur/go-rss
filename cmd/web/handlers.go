package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/frankie-mur/go-rss/internal/database"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (app *application) readinessHandler(c echo.Context) error {
	data := map[string]string{
		"status": "OK",
	}
	return c.JSON(200, data)
}

func (app *application) createUserHandler(c echo.Context) error {
	/* Need to also check len(payload.Name) because default behavior
	does not error if field is not present */
	//	var req createUserRequest
	// if err := decodeJson(c.Request(), &req); err != nil || len(req.Name) == 0 {
	// 	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	// }
	//Using name we can create a new user
	name := c.FormValue("name")
	user, err := app.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	})
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, user)
}

func (app *application) getUserByApiKeyHandler(e echo.Context, u database.User) error {
	//Send back user to client
	return e.JSON(http.StatusOK, u)
}

type createFeedRequest struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type createFeedResponse struct {
	Feed        database.Feed       `json:"feed"`
	Feed_Follow database.FeedFollow `json:"feed_follow"`
}

func (app *application) createFeedHandler(e echo.Context, u database.User) error {
	var req createFeedRequest
	if err := decodeJson(e.Request(), &req); err != nil || len(req.Name) == 0 || len(req.Url) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	//Create a new feed
	feed, err := app.DB.CreateFeed(e.Request().Context(), database.CreateFeedParams{
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
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	//Create a feed follow
	feed_follow, err := app.DB.CreateFeedFollow(e.Request().Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    u.ID,
		FeedID:    feed.ID,
	})

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return e.JSON(http.StatusCreated, createFeedResponse{
		Feed:        feed,
		Feed_Follow: feed_follow,
	})
}

func (app *application) getAllFeedsHandler(e echo.Context) error {
	data, err := app.DB.GetAllFeeds(e.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, data)
}

type createFeedFollowRequest struct {
	Feed_id uuid.UUID `json:"feed_id"`
}

func (app *application) createFeedFollowHandler(e echo.Context, u database.User) error {
	var req createFeedFollowRequest
	if err := decodeJson(e.Request(), &req); err != nil || len(req.Feed_id) == 0 {
		fmt.Printf("failed with error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	data, err := app.DB.CreateFeedFollow(e.Request().Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    u.ID,
		FeedID:    req.Feed_id,
	})

	if err != nil {
		fmt.Printf("failed with error %v", err)
		//TODO: Need to handle error when feed_id does not match to a feed ID
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusCreated, data)
}

func (app *application) deleteFeedFollowHandler(e echo.Context, u database.User) error {
	param := e.Param("id")
	feedFollowID, err := uuid.Parse(param)
	if err != nil {
		return echo.ErrBadRequest
	}
	if len(param) == 0 {
		return echo.ErrBadRequest
	}
	fmt.Printf("Param %s\n", feedFollowID)
	feed_id, err := uuid.Parse(param)
	if err != nil {
		fmt.Printf("Error %v", err.Error())
		return echo.ErrBadRequest
	}
	err = app.DB.DeleteFeedFollow(e.Request().Context(), database.DeleteFeedFollowParams{
		FeedID: feed_id,
		UserID: u.ID,
	})
	if err != nil {
		fmt.Printf("failed with error %v", err)
		return echo.ErrInternalServerError
	}

	return e.String(http.StatusOK, "")
}

func (app *application) getAllFeedFollows(e echo.Context, u database.User) error {
	feed_follows, err := app.DB.GetAllFeedFollows(e.Request().Context(), u.ID)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, feed_follows)
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
	return e.Render(http.StatusOK, "posts", posts)
}
