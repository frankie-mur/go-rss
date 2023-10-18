package main

import session "github.com/spazzymoto/echo-scs-session"

func (app *application) routes() {
	//Init our session group to attach middleware that use sessions
	sessionRoutes := app.e.Group("/session")
	sessionRoutes.Use(session.LoadAndSave(app.session))
	// Health
	app.e.GET("/health", app.readinessHandler)
	//Pages
	app.e.GET("/", app.indexHandler)
	app.e.GET("/signup", app.signupHandler)
	app.e.GET("/login", app.loginHandler)
	//Users
	app.e.POST("/users/signup", app.createUserHandler)
	sessionRoutes.POST("/users/login", app.loginUserHandler)
	//AUTH
	sessionRoutes.GET("/users", app.middlewareAuth(app.getUserByApiKeyHandler))
	//Feeds
	sessionRoutes.POST("/feeds", app.middlewareAuth(app.createFeedHandler))
	app.e.GET("/feeds", app.getAllFeedsHandler)
	//Feed Follows
	sessionRoutes.GET("/feed_follows", app.middlewareAuth(app.getAllFeedFollows))
	sessionRoutes.POST("/feed_follows", app.middlewareAuth(app.createFeedFollowHandler))
	sessionRoutes.DELETE("/feed_follows/:id", app.middlewareAuth(app.deleteFeedFollowHandler))
	//Posts
	sessionRoutes.GET("/posts", app.middlewareAuth(app.getPostsByUserHandler))

}
