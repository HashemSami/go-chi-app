```go
func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	// before calling post form, we must call
	// the parse form function first to be able to use the
	// post request as a form
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// getting the same name attribute used in the html
	fmt.Fprint(w, "Email: ", r.PostForm.Get("email"))
	fmt.Fprint(w, "Password: ", r.PostForm.Get("password"))
}

```
