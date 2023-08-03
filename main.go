package main

import (
	"fmt"
	"net/http"

	"github.com/HashemSami/go-chi-app/controllers"
	"github.com/HashemSami/go-chi-app/models"
	"github.com/HashemSami/go-chi-app/templates"
	"github.com/HashemSami/go-chi-app/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	// parsing the html files before srving the app
	// to users
	homeTpl := views.Must(
		views.ParseFS(templates.FS, "home.html", "tailwind.html"),
	)
	contactTpl := views.Must(
		views.ParseFS(templates.FS, "contact.html", "tailwind.html"),
	)
	faqTpl := views.Must(
		views.ParseFS(templates.FS, "faq.html", "tailwind.html"),
	)
	signupTpl := views.Must(
		views.ParseFS(templates.FS, "signup.html", "tailwind.html"),
	)

	// get the DB connection
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// get the Services
	userService := models.UserService{
		DB: db,
	}

	// creating users controllers
	usersC := controllers.Users{
		UserService: &userService,
	}
	usersC.Templates.New = signupTpl

	r.Use(middleware.Logger)

	r.Get("/", controllers.StaticHandler(homeTpl))
	r.Get("/contact", controllers.StaticHandler(contactTpl))
	r.Get("/faq", controllers.FAQ(faqTpl))
	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not Found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
