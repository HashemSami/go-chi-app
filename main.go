package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func executeTemplate(w http.ResponseWriter, filePath string) {
	w.Header().Set("Content-Type", "text/html; charset-utf-8")
	tpl, err := template.ParseFiles(filePath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There eas an error parsing the template.",
			http.StatusInternalServerError)
		return
	}

	templateErr := tpl.Execute(w, nil)
	if templateErr != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There eas an error executing the template.",
			http.StatusInternalServerError)
		return
	}
}

func homeHendler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "home.html")
	executeTemplate(w, tplPath)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tplPath := filepath.Join("templates", "contact.html")
	executeTemplate(w, tplPath)
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", homeHendler)
	r.Get("/contact", contactHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not Found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
