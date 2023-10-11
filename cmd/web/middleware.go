package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/frankie-mur/go-rss/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (app *application) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		//Validate authorization header is in correct format
		apikey, err := validateAuthHeader(authHeader)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		//check if apikey matches to a user
		user, err := app.DB.GetUserByApiKey(r.Context(), apikey)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, http.StatusUnauthorized, "Unauthorized")
				return
			} else {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		fmt.Println("User is authenticated")
		handler(w, r, user)
	}
}
