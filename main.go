package main

import (
	"fmt"
	"net/http"
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
	var router Router
	// http.Handler - interface with the ServeHTTP method
	// http.HandlerFunc - a function type accepts same args as ServeHTTP method,
	// and also implments http.Handler

	// http.Handle("/", http.Handler)
	// http.HandleFunc("/", pathHandler)

	fmt.Println("Startin gthe server on port :3000")
	http.ListenAndServe(":3000", router)
}
