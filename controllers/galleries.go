package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/HashemSami/go-chi-app/context"
	"github.com/HashemSami/go-chi-app/models"
	"github.com/go-chi/chi/v5"
)

type Galleries struct {
	Templates struct {
		New  Template
		Edit Template
	}
	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}

	data.Title = r.FormValue("title")
	g.Templates.New.Execute(w, r, data)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID int
		Title  string
	}

	data.UserID = context.User(r.Context()).ID
	data.Title = r.FormValue("title")

	gallery, err := g.GalleryService.Create(data.Title, data.UserID)
	if err != nil {
		g.Templates.New.Execute(w, r, data, err)
		return
	}

	// redirect the user to the edit the gallery
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)

	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}

	gallery, err := g.GalleryService.ByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Gallery not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// making sure that the user is editing the gallery belongs to the him/her
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "you are not authorized to use this gallery", http.StatusForbidden)
		return
	}

	var data struct {
		ID    int
		Title string
	}
	data.ID = gallery.ID
	data.Title = gallery.Title

	g.Templates.Edit.Execute(w, r, data)
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusNotFound)
		return
	}

	gallery, err := g.GalleryService.ByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Gallery not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// making sure that the user is editing the gallery belongs to the him/her
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "you are not authorized to use this gallery", http.StatusForbidden)
		return
	}

	gallery.Title = r.FormValue("title")
	err = g.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}

	// redirect the user after editing the galleries to the edit page
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)

	http.Redirect(w, r, editPath, http.StatusFound)
}
