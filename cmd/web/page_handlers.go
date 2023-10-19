package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) indexHandler(e echo.Context) error {
	return e.Render(http.StatusOK, "index", app.pageData)
}

func (app *application) signupHandler(e echo.Context) error {
	return e.Render(http.StatusOK, "signup", app.pageData)
}

func (app *application) loginHandler(e echo.Context) error {
	return e.Render(http.StatusOK, "login", app.pageData)
}
