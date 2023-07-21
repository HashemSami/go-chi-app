package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func Parse(filePath string) (Template, error) {
	tpl, err := template.ParseFiles(filePath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{htmlTpl: tpl}, nil
}

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset-utf-8")

	templateErr := t.htmlTpl.Execute(w, data)

	if templateErr != nil {
		log.Printf("executing template: %v", templateErr)
		http.Error(w, "There eas an error executing the template.",
			http.StatusInternalServerError)
		return
	}
}
