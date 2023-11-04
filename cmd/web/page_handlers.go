package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *application) indexHandler(e echo.Context) error {
	data := app.newPageData(e)
	if data.IsAuthenticated {
		//Load feeds and feed_follows
		err := app.getPosts(e, &data)
		if err != nil {
			return err
		}

	}
	//fmt.Printf("--Page data is %v\n\n", data)
	//Return must be logged in to view feeds
	return e.Render(http.StatusOK, "index", data)
}

func (app *application) signupHandler(e echo.Context) error {
	data := app.newPageData(e)
	return e.Render(http.StatusOK, "signup", data)
}

func (app *application) loginHandler(e echo.Context) error {
	data := app.newPageData(e)
	return e.Render(http.StatusOK, "login", data)
}
