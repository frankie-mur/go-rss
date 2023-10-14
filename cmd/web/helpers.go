package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

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
