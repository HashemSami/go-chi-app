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
	"github.com/gorilla/csrf"
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
	signinTpl := views.Must(
		views.ParseFS(templates.FS, "signin.html", "tailwind.html"),
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
	sessionService := models.SessionService{
		DB: db,
	}

	// creating users controllers
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	usersC.Templates.SignUp = signupTpl
	usersC.Templates.SignIn = signinTpl

	r.Use(middleware.Logger)

	r.Get("/", controllers.StaticHandler(homeTpl))
	r.Get("/contact", controllers.StaticHandler(contactTpl))
	r.Get("/faq", controllers.FAQ(faqTpl))
	r.Get("/signup", usersC.SignUp)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Get("/users/me", usersC.CurrentUser)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not Found", http.StatusNotFound)
	})

	csrfKey := "kfjggbctiopwoidjipiuewdxhjksla"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		// TODO: FIX this before deploying
		// for HTTPS
		csrf.Secure(false),
	)

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", csrfMw(r))
}
