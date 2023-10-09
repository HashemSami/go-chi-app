package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/HashemSami/go-chi-app/controllers"
	"github.com/HashemSami/go-chi-app/migrations"
	"github.com/HashemSami/go-chi-app/models"
	"github.com/HashemSami/go-chi-app/templates"
	"github.com/HashemSami/go-chi-app/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	// PSQL
	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	// put a message on the console if connection error
	if cfg.PSQL.Host == "" && cfg.PSQL.Port == "" {
		return cfg, fmt.Errorf("No PSQL config provided.")
	}
	// SMTP
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.UserName = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	// CSRF
	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	// if it has any value other than true. set to false
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"

	// Server
	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	err = run(cfg)

	if err != nil {
		panic(err)
	}
}

func run(cfg config) error {
	// get the DB connection
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		return err
	}
	defer db.Close()

	// setting the migration code
	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		return err
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
	forgotPasswordTpl := views.Must(
		views.ParseFS(templates.FS, "forgot_password.html", "tailwind.html"),
	)
	checkYourEmailTpl := views.Must(
		views.ParseFS(templates.FS, "check_your_email.html", "tailwind.html"),
	)
	resetPasswordTpl := views.Must(
		views.ParseFS(templates.FS, "reset_pw.html", "tailwind.html"),
	)
	newGalleryTpl := views.Must(
		views.ParseFS(templates.FS, "galleries/new.html", "tailwind.html"),
	)
	editGalleryTpl := views.Must(
		views.ParseFS(templates.FS, "galleries/edit.html", "tailwind.html"),
	)
	indexGalleryTpl := views.Must(
		views.ParseFS(templates.FS, "galleries/index.html", "tailwind.html"),
	)
	showGalleryTpl := views.Must(
		views.ParseFS(templates.FS, "galleries/show.html", "tailwind.html"),
	)

	// get the Services
	userService := &models.UserService{
		DB: db,
	}
	sessionService := &models.SessionService{
		DB: db,
	}
	passwordResetService := &models.PasswordResetService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)
	galleryService := &models.GalleryService{
		DB: db,
	}

	// creating users controllers
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: passwordResetService,
		EmailService:         emailService,
	}
	usersC.Templates.SignUp = signupTpl
	usersC.Templates.SignIn = signinTpl
	usersC.Templates.ForgotPassword = forgotPasswordTpl
	usersC.Templates.CheckYourEmail = checkYourEmailTpl
	usersC.Templates.ResetPassword = resetPasswordTpl

	galleriesC := controllers.Galleries{
		GalleryService: galleryService,
	}
	galleriesC.Templates.New = newGalleryTpl
	galleriesC.Templates.Edit = editGalleryTpl
	galleriesC.Templates.Index = indexGalleryTpl
	galleriesC.Templates.Show = showGalleryTpl

	// setting middleware
	umw := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Path("/"),
		// TODO: FIX this before deploying
		// for HTTPS
		csrf.Secure(cfg.CSRF.Secure),
	)

	// setup the handlers
	r := chi.NewRouter()
	// r2 := chi.NewRouter()
	// r2.Mount("/api", r)

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers

	})
	r.Use(cors.Handler)

	r.Use(middleware.Logger)
	// r.Use(umw.SetHeaders)
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
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)

	// provide a specific functionality to the current user
	// rout that will apply out user middleware
	r.Route("/users/me", func(r chi.Router) {
		// set the user middleware just for the routes that
		// start with this path
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	r.Route("/galleries", func(r chi.Router) {
		// group the routes that require the user
		// to be signed in
		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/", galleriesC.Index)
			r.Get("/new", galleriesC.New)
			r.Get("/{id}/edit", galleriesC.Edit)
			r.Post("/", galleriesC.Create)
			r.Post("/{id}", galleriesC.Update)
			r.Post("/{id}/delete", galleriesC.Delete)
			r.Post("/{id}/images", galleriesC.UploadImages)
			r.Post("/{id}/images/{filename}/delete", galleriesC.DeleteImage)
		})
		// routes that don't require a user to be signed in
		r.Get("/{id}", galleriesC.Show)
		r.Get("/{id}/images/{filename}", galleriesC.Image)
	})

	// serving the assets directory with the routes
	assetsHandler := http.FileServer(http.Dir("assets"))
	r.Get("/assets/*", http.StripPrefix("/assets", assetsHandler).ServeHTTP)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not Found", http.StatusNotFound)
	})

	// Start the server
	fmt.Printf("Starting the server on %s...\n", cfg.Server.Address)

	return http.ListenAndServe(cfg.Server.Address, r)
}
