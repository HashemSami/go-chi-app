package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func homeHendler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf-8")
	fmt.Fprint(w, "<h1>Home</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset-utf-8")
	fmt.Fprint(w, "<h1>Contact</h1>")
}

// func pathHandler(w http.ResponseWriter, r *http.Request) {
// 	switch r.URL.Path {
// 	case "/":
// 		homeHendler(w, r)
// 	case "/contact":
// 		contactHandler(w, r)
// 	}
// }

// implementing the handler inerface to add to the listen server
type Router struct{}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHendler(w, r)
	case "/contact":
		contactHandler(w, r)
	default:
		http.Error(w, "Page not Found", http.StatusNotFound)
	}
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
