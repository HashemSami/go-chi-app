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
	UserService    *models.UserService
	SessionService *models.SessionService
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

	// creating session token in the sessions database
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		// if the session already exists
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		//  TODO: long term, we should show a warning about not
		// being able to sign the user in

		// NOTE: we must return after the redirect
		// if we don't, below code will get get executed after
		// the redirect
		return
	}

	coockie := http.Cookie{
		Name:     "session",
		Value:    session.Token,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, &coockie)

	http.Redirect(w, r, "/users/me", http.StatusFound)
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

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		// this will execute when the used successfully signed in but
		// something went wrong with the session creation
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// if nothing went wrong, set the user cookies
	cookie := http.Cookie{
		Name:  "session",
		Value: session.Token,
		Path:  "/",
		// this will prevent js to access the cookie
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/users/me", http.StatusFound)
}

// validating the user request byt verifying the token
// taken from the headers cookies
func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	// get the cookie from the request
	tokenCookie, err := r.Cookie("session")
	if err != nil {
		// if the session already exists
		fmt.Println(err)
		// redirect to resigning to set the new cookies
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	// get the user's data using the session token
	user, err := u.SessionService.User(tokenCookie.Value)
	if err != nil {
		// if not able to bring the user data using the token
		fmt.Println(err)
		// redirect to resigning to set the new cookies
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	fmt.Fprintf(w, "Current user: %s\n", user.Email)
}
