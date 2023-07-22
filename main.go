package main

import (
	"fmt"
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
	homeTpl := views.Must(
		views.Parse(filepath.Join("templates", "home.html")),
	)

	contactTpl := views.Must(
		views.Parse(filepath.Join("templates", "contact.html")),
	)

	r.Use(middleware.Logger)

	r.Get("/", controllers.StaticHandler(homeTpl))
	r.Get("/contact", controllers.StaticHandler(contactTpl))
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not Found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
