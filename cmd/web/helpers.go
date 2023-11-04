package main

import (
	"encoding/json"
	"errors"
	"github.com/frankie-mur/go-rss/internal/database"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type PageData struct {
	Flash           string
	IsAuthenticated bool
	Posts           map[string][]database.GetPostsByUserIdRow
	FeedFollows     []database.FeedFollow
}

//func validateAuthHeader(authHeader string) (string, error) {
//	// Split the header to check for the "Bearer" format
//	splitToken := strings.Split(authHeader, "Bearer ")
//	if len(splitToken) != 2 {
//		return "", errors.New("invalid auth header")
//	}
//	return splitToken[1], nil // The actual token
//}

func decodeJson(r *http.Request, payload interface{}) error {
	return json.NewDecoder(r.Body).Decode(payload)
}

func (app *application) newPageData(e echo.Context) PageData {
	return PageData{
		Flash:           app.session.PopString(e.Request().Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(e),
	}
}

func (app *application) isAuthenticated(e echo.Context) bool {
	isAuthenticated := app.session.GetBool(e.Request().Context(), "isAuthenticated")
	if !isAuthenticated {
		return false
	}
	return isAuthenticated
}

func (app *application) getUserId(e echo.Context) (*uuid.UUID, error) {
	userId := app.session.GetString(e.Request().Context(), "authenticatedUserID")
	if len(userId) == 0 {
		return nil, errors.New("no authenticated user ID in session")
	}
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("failed to parse user ID from session")
	}
	return &userUUID, nil
}

func (app *application) getPosts(e echo.Context, data *PageData) error {
	if !data.IsAuthenticated {
		return errors.New("user must be authenticated to get feeds")
	}
	userUUID, err := app.getUserId(e)
	if err != nil {
		return errors.New("failed to get user ID from")
	}
	// Get feeds from the database
	posts, err := app.DB.GetPostsByUserId(e.Request().Context(), database.GetPostsByUserIdParams{
		UserID: *userUUID,
		Limit:  int32(30),
	})
	//Group posts by name
	groupedPosts := make(map[string][]database.GetPostsByUserIdRow)
	for _, post := range posts {
		groupedPosts[post.Name.String] = append(groupedPosts[post.Name.String], post)
	}
	//fmt.Printf("Got posts %v\n", groupedPosts)
	data.Posts = groupedPosts
	return nil
}
