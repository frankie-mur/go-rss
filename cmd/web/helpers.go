package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)

	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	errMsg := map[string]string{
		"error": msg,
	}
	err := respondWithJSON(w, code, errMsg)
	if err != nil {
		return err
	}
	return nil
}

func validateAuthHeader(authHeader string) (string, error) {
	// Split the header to check for the "Bearer" format
	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return "", errors.New("invalid auth header")
	}
	return splitToken[1], nil // The actual token
}

func decodeJson(r *http.Request, payload interface{}) error {
	return json.NewDecoder(r.Body).Decode(payload)
}
