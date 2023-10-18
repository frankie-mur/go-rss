package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"

	"github.com/frankie-mur/go-rss/internal/database"
	"github.com/labstack/echo/v4"
)

type authedHandler func(e echo.Context, u database.User) error

func (app *application) middlewareAuth(handler authedHandler) echo.HandlerFunc {
	return func(e echo.Context) error {
		userId := app.session.GetString(e.Request().Context(), "authenticatedUserID")
		//userId is default value if not present
		if len(userId) == 0 {
			fmt.Printf("Failed with error %v", "no User Id")
			return echo.ErrUnauthorized
		}
		//Ensure the userId is valid UUID
		userUIID, err := uuid.Parse(userId)
		if err != nil {
			fmt.Printf("Failed with error %v", err.Error())
			return echo.ErrUnauthorized
		}
		//check if apikey matches to a user
		user, err := app.DB.GetUserByID(e.Request().Context(), userUIID)
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

//// SessionMiddleware Wrap the scs LoadAndSave middleware to make compatible with echo middleware
//func (app *application) SessionMiddleware(next echo.HandlerFunc) echo.MiddlewareFunc {
//	return echo.WrapMiddleware(app.session.LoadAndSave)
//}
