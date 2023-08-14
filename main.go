package main

import (
	"fmt"
	"net/http"

	"github.com/HashemSami/go-chi-app/controllers"
	"github.com/HashemSami/go-chi-app/migrations"
	"github.com/HashemSami/go-chi-app/models"
	"github.com/HashemSami/go-chi-app/templates"
	"github.com/HashemSami/go-chi-app/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
)

func main() {
	// get the DB connection
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// setting the migration code
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// parsing the html files before serving the app
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

	// setting middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfKey := "kfjggbctiopwoidjipiuewdxhjksla"
	csrfMw := csrf.Protect(
		[]byte(csrfKey),
		csrf.Path("/"),
		// TODO: FIX this before deploying
		// for HTTPS
		csrf.Secure(false),
	)

	// setup the handlers
	r := chi.NewRouter()
	// r2 := chi.NewRouter()
	// r2.Mount("/api", r)

	r.Use(middleware.Logger)
	r.Use(csrfMw)
	r.Use(umw.SetUser)

	// setting our routes
	r.Get("/", controllers.StaticHandler(homeTpl))
	r.Get("/contact", controllers.StaticHandler(contactTpl))
	r.Get("/faq", controllers.FAQ(faqTpl))
	r.Get("/signup", usersC.SignUp)
	r.Post("/signup", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/signout", usersC.ProcessSignOut)

	// provide a spicific functionality to the current user
	// rout that will apply out user middleware
	r.Route("/users/me", func(r chi.Router) {
		// set the user middleware just for the routes that
		// start with this path
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not Found", http.StatusNotFound)
	})

	// Start the server
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
