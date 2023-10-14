package main

import (
	"database/sql"
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
	DB *database.Queries
	e  *echo.Echo
}

func main() {
	godotenv.Load(".env")

	//port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	datastore := database.New(db)
	e := echo.New()

	app := &application{
		DB: datastore,
		e:  e,
	}

	app.routes()

	go initScraper(datastore, 10, time.Hour)

	// Little bit of middlewares for housekeeping
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	NewTemplateRenderer(e, "ui/html/*/*.tmpl")

	// srv := &http.Server{
	// 	Addr:    fmt.Sprintf("localhost:%s", port),
	// 	Handler: app.routes(),
	// }
	//	fmt.Printf("Starting server on addr %s\n", srv.Addr)
	//	err = srv.ListenAndServe()
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
