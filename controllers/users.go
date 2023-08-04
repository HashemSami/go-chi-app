package controllers

import (
	"fmt"
	"net/http"

	"github.com/HashemSami/go-chi-app/models"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	// this is getting values from the url query and add it
	// to the form as an initial data
	data.Email = r.FormValue("email")
	// we need a view to render
	u.Templates.New.Execute(w, data)
}

// func (u Users) Create(w http.ResponseWriter, r *http.Request) {
// 	// before calling post form, we must call
// 	// the parse form function first to be able to use the
// 	// post request as a form
// 	err := r.ParseForm()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// getting the same name attribute used in the html
// 	fmt.Fprint(w, "Email: ", r.PostForm.Get("email"))
// 	fmt.Fprint(w, "Password: ", r.PostForm.Get("password"))
// }

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
	u.Templates.SignIn.Execute(w, data)
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

	fmt.Fprintf(w, "User authenticated: %+v", user)

	// cookie:=http.Cookie(
	// )
}
