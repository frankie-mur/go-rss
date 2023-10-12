package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (app *application) indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("ui/html/pages/index.go.tmpl"))
	if err := tmpl.Execute(w, r); err != nil {
		fmt.Println("Could not execute template", err)
	}
}
