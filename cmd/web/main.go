package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	subrouter := chi.NewRouter()
	router.Mount("/v1", subrouter)

	subrouter.Get("/readiness", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"status": "OK",
		}
		err := respondWithJSON(w, 200, data)
		if err != nil {
			log.Fatal(err)
		}
	})

	subrouter.Get("/err", func(w http.ResponseWriter, r *http.Request) {
		err := respondWithError(w, 500, "Internal Server Error")
		if err != nil {
			log.Fatal(err)
		}
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", port),
		Handler: router,
	}
	fmt.Printf("Starting server on addr %s", srv.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}

}
