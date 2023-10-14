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

	//subrouter := chi.NewRouter()
	//Pages
	app.e.GET("/", app.indexHandler)
	// Health
	app.e.GET("/health", app.readinessHandler)
	//Users
	app.e.POST("/users", app.createUserHandler)
	//AUTH
	app.e.GET("/users", app.middlewareAuth(app.getUserByApiKeyHandler))
	//Feeds
	app.e.POST("/feeds", app.middlewareAuth(app.createFeedHandler))
	app.e.GET("/feeds", app.getAllFeedsHandler)
	// //Feed Follows
	app.e.GET("/feed_follows", app.middlewareAuth(app.getAllFeedFollows))
	app.e.POST("/feed_follows", app.middlewareAuth(app.createFeedFollowHandler))
	app.e.DELETE("/feed_follows/:id", app.middlewareAuth(app.deleteFeedFollowHandler))
	//Posts
	//Interim: Remove middleware
	app.e.GET("/posts/:limit", app.getPostsByUserHandler)

	return router
}
