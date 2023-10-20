package main

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

type PageData struct {
	Flash           string
	IsAuthenticated bool
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
		//	Flash:           app.session.PopString(e.Request().Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(e),
	}
}

func (app *application) isAuthenticated(e echo.Context) bool {
	isAuthenticated := app.session.GetBool(e.Request().Context(), "IsAuthenticated")
	if !isAuthenticated {
		return false
	}
	return isAuthenticated
}
