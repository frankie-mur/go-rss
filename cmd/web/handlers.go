package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/frankie-mur/go-rss/internal/email_generator"
	"github.com/frankie-mur/go-rss/internal/validator"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/frankie-mur/go-rss/internal/database"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (app *application) readinessHandler(e echo.Context) error {
	data := map[string]string{
		"status": "OK",
	}
	return e.JSON(200, data)
}

func (app *application) createUserHandler(e echo.Context) error {
	//Using name we can create a new user
	name := e.FormValue("name")
	email := e.FormValue("email")
	if !validator.Matches(email, validator.EmailRX) {
		return echo.ErrBadRequest
	}
	_, err := validator.VerifyEmail(email)
	if err != nil {
		fmt.Printf("email verfication failed with error%v\n", err)
		return echo.ErrBadRequest
	}
	// User email is valid, we can create a new user
	user, err := app.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Email: sql.NullString{
			String: email,
			Valid:  true,
		},
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return echo.NewHTTPError(http.StatusBadRequest, "User already exists")
		}
		return echo.ErrInternalServerError
	}
	log.Printf("successfully created user %v", user.Name)
	//Send welcome email
	err = email_generator.SendEmail(email, email_generator.GenerateWelcomeEmail(name))
	if err != nil {
		fmt.Printf("error sending welcome email: %v\n", err)
		return echo.ErrInternalServerError
	}
	app.session.Put(e.Request().Context(), "flash", "Successfully Created Account!")
	app.session.Put(e.Request().Context(), "isAuthenticated", true)
	app.session.Put(e.Request().Context(), "authenticatedUserID", fmt.Sprintf("%s", user.ID))

	return e.Redirect(http.StatusSeeOther, "/")
}

func (app *application) loginUserHandler(e echo.Context) error {
	email := e.FormValue("email")
	if !validator.Matches(email, validator.EmailRX) {
		return echo.ErrBadRequest
	}
	user, err := app.DB.GetUserByEmail(e.Request().Context(), sql.NullString{
		String: email,
		Valid:  true,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	log.Printf("Succesfully fetched user")
	// Use the RenewToken() method on the current session to change the session
	// ID. It's good practice to generate a new session ID when the
	// authentication state or privilege levels changes for the user (e.g. login
	// and logout operations).
	err = app.session.RenewToken(e.Request().Context())
	if err != nil {
		return echo.ErrInternalServerError
	}

	app.session.Put(e.Request().Context(), "flash", "Successfully Logged In!")
	app.session.Put(e.Request().Context(), "isAuthenticated", true)
	app.session.Put(e.Request().Context(), "authenticatedUserID", fmt.Sprintf("%s", user.ID))

	return e.Redirect(http.StatusFound, "/")
}

func (app *application) logoutUserHandler(e echo.Context) error {
	err := app.session.RenewToken(e.Request().Context())
	if err != nil {
		return echo.ErrInternalServerError
	}

	app.session.Remove(e.Request().Context(), "isAuthenticated")
	app.session.Remove(e.Request().Context(), "authenticatedUserID")
	app.session.Put(e.Request().Context(), "flash", "Successfully Logged Out!")

	return e.Redirect(http.StatusSeeOther, "/")
}

func (app *application) getUserByApiKeyHandler(e echo.Context, u database.User) error {
	//Send back user to client
	return e.JSON(http.StatusOK, u)
}

type createFeedResponse struct {
	Feed        database.Feed       `json:"feed"`
	Feed_Follow database.FeedFollow `json:"feed_follow"`
}

func (app *application) createFeedHandler(e echo.Context, u database.User) error {
	//var req createFeedRequest
	//if err := decodeJson(e.Request(), &req); err != nil || len(req.Name) == 0 || len(req.Url) == 0 {
	//	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	//}
	fmt.Println("In create feed")
	name := e.FormValue("name")
	url := e.FormValue("url")
	//Validate form data
	if len(name) == 0 || len(url) == 0 || !validator.Matches(url, validator.RssUrlRX) {
		return echo.ErrBadRequest
	}
	//Create a new feed
	feed, err := app.DB.CreateFeed(e.Request().Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
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
	fmt.Println("created feed")
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
	feedFollows, err := app.DB.GetAllFeedFollows(e.Request().Context(), u.ID)
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, feedFollows)
}

func (app *application) getPostsByUserHandler(e echo.Context, u database.User) error {
	limitParam := e.QueryParam("limit")
	//sortedParam := e.QueryParam("sorted")
	var limit int
	var err error
	if len(limitParam) == 0 {
		//Default to query for 15 posts
		limit = 15
	} else {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	posts, err := app.DB.GetPostsByUserId(e.Request().Context(), database.GetPostsByUserIdParams{
		UserID: u.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	//Group posts by name
	groupedPosts := make(map[string][]database.GetPostsByUserIdRow)
	for _, post := range posts {
		groupedPosts[post.Name.String] = append(groupedPosts[post.Name.String], post)
	}
	return e.Render(http.StatusOK, "posts", groupedPosts)
}
