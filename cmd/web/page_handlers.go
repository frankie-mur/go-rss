package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (app *application) indexHandler(e echo.Context) error {
	data := app.newPageData(e)
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
