package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/frankie-mur/go-rss/internal/database"
	"github.com/labstack/echo/v4"
)

type authedHandler func(e echo.Context, u database.User) error

func (app *application) middlewareAuth(handler authedHandler) echo.HandlerFunc {
	return func(e echo.Context) error {
		authHeader := e.Request().Header.Get("Authorization")
		if authHeader == "" {
			echo.NewHTTPError(http.StatusUnauthorized, errors.New("invalid authorization header"))
		}
		//Validate authorization header is in correct format
		apikey, err := validateAuthHeader(authHeader)
		if err != nil {
			echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		//check if apikey matches to a user
		user, err := app.DB.GetUserByApiKey(e.Request().Context(), apikey)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			} else {
				echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
		}
		fmt.Println("User is authenticated")
		handler(e, user)
		return nil
	}
}
