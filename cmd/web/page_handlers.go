package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) indexHandler(e echo.Context) error {
	return e.Render(http.StatusOK, "index", nil)
}

func (app *application) signupHandler(e echo.Context) error {
	return e.Render(http.StatusOK, "signup", nil)
}
