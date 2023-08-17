package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/HashemSami/go-chi-app/context"
	"github.com/HashemSami/go-chi-app/models"
)

type Users struct {
	Templates struct {
		SignUp         Template
		SignIn         Template
		ForgotPassword Template
		CheckYourEmail Template
		ResetPassword  Template
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
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

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	// FormValue can also get the query parameters from the URL
	// getting the same name attribute used in the html
	nu := models.NewUser{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
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
		// if the session doesn't already exists
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		//  TODO: long term, we should show a warning about not
		// being able to sign the user in

		// NOTE: we must return after the redirect
		// if we don't, below code will get get executed after
		// the redirect
		return
	}

	setCookie(w, CookieSession, session.Token)

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
	nu := models.NewUser{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	user, err := u.UserService.Authenticate(nu)
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
	setCookie(w, CookieSession, session.Token)

	http.Redirect(w, r, "/users/me", http.StatusFound)
}

// User required route
func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	// the context that will be set from the middleware
	user := context.User(r.Context())

	fmt.Fprintf(w, "Current user: %v\n", user.Email)
}

func (u Users) ProcessSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, CookieSession)
	if err != nil {
		// if the session already exists
		fmt.Println(err)
		// redirect to resigning to set the new cookies
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	err = u.SessionService.Delete(token)
	if err != nil {

		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	// delete the users cookie
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

// the page that will ask the user to provide the email that will
// receive the password reset
func (u Users) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

// getting the email value from the form, and process sending the forgot
// link to the user's email
func (u Users) ProcessForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		// TODO: Handle other cases in the future. for instance, if a user does not
		// exist with that email address.
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}

	// using the URL package to construct the reset URL query values
	vals := url.Values{
		"token": {pwReset.Token},
	}
	resetURL := "https://www.lenslocked.com/reset-pw?" + vals.Encode()

	err = u.EmailService.ForgotPassword(data.Email, resetURL)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	// redirect the user to the check your email page after sending the token
	// email to the user
	// note: the toke will not be rendered at this page, token will be only
	// sent to the users email
	u.Templates.CheckYourEmail.Execute(w, r, data)
}

func (u Users) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}

	data.Token = r.FormValue("token")
	u.Templates.ResetPassword.Execute(w, r, data)
}

func (u Users) ProcessResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}

	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")

	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		fmt.Println(err)
		// TODO: distinguish between types of errors
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// TODO: update the users password

	// sign in the user after resetting the password
	// any errors from this point onwards, should be redirected to the
	// sign in page

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}

	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

// ================================================================
// user middleware
type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get the cookie from the request
		token, err := readCookie(r, CookieSession)
		if err != nil {
			// proceed with the request and assume that the user is not signed in
			next.ServeHTTP(w, r)
			return
		}

		// get the user's data using the session token
		user, err := umw.SessionService.User(token)
		if err != nil {
			// if not able to bring the user data using the token
			fmt.Println(err)
			next.ServeHTTP(w, r)
			return
		}

		// if the user is found, set the context and set the request to the next handler
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		// update the request context
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// check the user if their present, if not, redirect the user to the
// sign in page
func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())

		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		// if the user is present, go to the next handler
		next.ServeHTTP(w, r)
	})
}
