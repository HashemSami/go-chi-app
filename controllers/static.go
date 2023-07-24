package controllers

import (
	"html/template"
	"net/http"

	"github.com/HashemSami/go-chi-app/views"
)

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func FAQ(tpl views.Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "What are you?",
			Answer:   "Human",
		}, {
			Question: "Is there a free version?",
			Answer:   "Yes, look at <a href=\"http://google.com\">Here</a>",
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}
