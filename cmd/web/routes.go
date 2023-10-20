package main

func (app *application) routes() {
	// Health
	app.e.GET("/health", app.readinessHandler)
	//Pages
	app.e.GET("/", app.indexHandler)
	app.e.GET("/signup", app.signupHandler)
	app.e.GET("/login", app.loginHandler)
	//Users
	app.e.POST("/users/signup", app.createUserHandler)
	app.e.POST("/users/login", app.loginUserHandler)
	app.e.POST("/users/logout", app.logoutUserHandler)
	//Feeds
	app.e.POST("/feeds", app.middlewareAuth(app.createFeedHandler))

	app.e.GET("/feeds", app.getAllFeedsHandler)
	//Feed Follows
	app.e.GET("/feed_follows", app.middlewareAuth(app.getAllFeedFollows))
	app.e.POST("/feed_follows", app.middlewareAuth(app.createFeedFollowHandler))
	app.e.DELETE("/feed_follows/:id", app.middlewareAuth(app.deleteFeedFollowHandler))
	//Posts
	app.e.GET("/posts", app.middlewareAuth(app.getPostsByUserHandler))
}
