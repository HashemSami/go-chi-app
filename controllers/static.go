package controllers

import (
	"html/template"
	"net/http"
)

func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func FAQ(tpl Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "What are you?",
			Answer:   "Human",
		}, {
			Question: "Is there a free version?",
			Answer:   "Yes, look at <a class=\"underline\" href=\"http://google.com\">Here</a>",
		}, {
			Question: "new question?",
			Answer:   "new answer",
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}
