package main

import (
	"encoding/json"
	"net/http"
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
