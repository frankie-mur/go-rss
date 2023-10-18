package main

import (
	"database/sql"
	"fmt"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/frankie-mur/go-rss/internal/database"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type application struct {
	DB       *database.Queries
	e        *echo.Echo
	session  *scs.SessionManager
	pageData map[string]interface{}
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		return
	}

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	dbURL := os.Getenv("DB_URL")
	//Initialize postgresDB
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//Initialize our queries (setup by sqlc)
	datastore := database.New(db)
	//Our echo instance
	e := echo.New()
	//Our Session manager connecting to our postgres db
	//For persistent storage
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = postgresstore.New(db)
	session.Cookie.Secure = true
	//Declare and assign our application struct
	app := &application{
		DB:      datastore,
		e:       e,
		session: session,
	}
	//Declare our routes
	app.routes()
	//Send off our scrapper in a go routine
	//This will scrape all RSS feeds stored in db every hour
	go initScraper(datastore, 10, time.Hour)

	//bit of middlewares for housekeeping
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use()

	//Render all of our templates
	NewTemplateRenderer(e, "ui/html/*/*.tmpl")
	//Start the server!
	if err := e.Start(fmt.Sprintf(":%s", port)); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
