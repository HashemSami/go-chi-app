package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/HashemSami/go-chi-app/controllers"
	"github.com/HashemSami/go-chi-app/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	// parsing the html files before srving the app
	// to users
	homeTpl, err := views.Parse(filepath.Join("templates", "home.html"))
	if err != nil {
		log.Printf("parsing template: %v", err)
		return
	}

	contactTpl, err := views.Parse(filepath.Join("templates", "contact.html"))
	if err != nil {
		log.Printf("parsing template: %v", err)
		return
	}

	r.Use(middleware.Logger)

	r.Get("/", controllers.StaticHandler(homeTpl))
	r.Get("/contact", controllers.StaticHandler(contactTpl))
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not Found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
