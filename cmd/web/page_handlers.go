package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (app *application) indexHandler(e echo.Context) error {
	data := app.newPageData(e)
	fmt.Print(data)
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
