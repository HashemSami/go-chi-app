package controllers

import (
	"fmt"
	"net/http"

	"github.com/HashemSami/go-chi-app/models"
)

type Users struct {
	Templates struct {
		SignUp Template
		SignIn Template
	}
	UserService *models.UserService
}

func (u Users) SignUp(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	// this is getting values from the url query and add it
	// to the form as an initial data
	data.Email = r.FormValue("email")
	// we need a view to render
	u.Templates.SignUp.Execute(w, r, data)
}

// another version
func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	// FormValue can also get the query parameters from the URL
	// getting the same name attribute used in the html
	email := r.FormValue("email")
	password := r.FormValue("password")

	nu := models.NewUser{
		Email:    email,
		Password: password,
	}

	user, err := u.UserService.Create(nu)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "User created: %+v", user)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	// this is getting values from the url query and add it
	// to the form as an initial data
	data.Email = r.FormValue("email")
	// we need a view to render
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) ProcessSignIn(w http.ResponseWriter, r *http.Request) {
	data := models.NewUser{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	user, err := u.UserService.Authenticate(data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// if nothing went wrong, set the user cookies
	cookie := http.Cookie{
		Name:  "email",
		Value: user.Email,
		Path:  "/",
		// this will prevent js to access the cookie
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	fmt.Fprintf(w, "User authenticated: %+v", user)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	// get the cookie from the request
	email, err := r.Cookie("email")
	if err != nil {
		fmt.Fprint(w, "The email cookie could not be read.")
		return
	}

	fmt.Fprintf(w, "Email cookie: %s\n", email.Value)
	fmt.Fprintf(w, "Headers: %+v\n", r.Header)
}
