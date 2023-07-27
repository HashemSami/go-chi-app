package controllers

import (
	"fmt"
	"net/http"
)

type Users struct {
	Templates struct {
		New Template
	}
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
	// getting the same name attribute used in the html
	fmt.Fprint(w, "Email: ", r.FormValue("email"))
	fmt.Fprint(w, "Password: ", r.FormValue("password"))
}
