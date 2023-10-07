package main

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
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
	// Health
	subrouter.Get("/readiness", app.readinessHandler)
	subrouter.Get("/err", app.errorHandler)
	//Users
	subrouter.Post("/users", app.createUserHandler)
	subrouter.Get("/users", app.middlewareAuth(app.getUserByApiKeyHandler))
	//Feeds
	subrouter.Post("/feeds", app.middlewareAuth(app.createFeedHandler))
	subrouter.Get("/feeds", app.getAllFeedsHandler)
	//Feed Follows
	subrouter.Get("/feed_follows", app.middlewareAuth(app.getAllFeedFollows))
	subrouter.Post("/feed_follows", app.middlewareAuth(app.createFeedFollowHandler))
	subrouter.Delete("/feed_follows/{id}", app.middlewareAuth(app.deleteFeedFollowHandler))

	router.Mount("/v1", subrouter)

	return router
}
