package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/frankie-mur/go-rss/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type application struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	dbURL := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	app := &application{
		DB: database.New(db),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", port),
		Handler: app.routes(),
	}
	fmt.Printf("Starting server on addr %s\n", srv.Addr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}

}
