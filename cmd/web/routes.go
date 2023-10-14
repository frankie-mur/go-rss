package main

func (app *application) routes() {
	// Health
	app.e.GET("/health", app.readinessHandler)
	//Pages
	app.e.GET("/", app.indexHandler)
	app.e.GET("/signup", app.signupHandler)
	//Users
	app.e.POST("/users", app.createUserHandler)
	//AUTH
	app.e.GET("/users", app.middlewareAuth(app.getUserByApiKeyHandler))
	//Feeds
	app.e.POST("/feeds", app.middlewareAuth(app.createFeedHandler))
	app.e.GET("/feeds", app.getAllFeedsHandler)
	//Feed Follows
	app.e.GET("/feed_follows", app.middlewareAuth(app.getAllFeedFollows))
	app.e.POST("/feed_follows", app.middlewareAuth(app.createFeedFollowHandler))
	app.e.DELETE("/feed_follows/:id", app.middlewareAuth(app.deleteFeedFollowHandler))
	//Posts
	//Interim: Remove middleware
	app.e.GET("/posts/:limit", app.getPostsByUserHandler)

}
