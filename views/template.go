package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
)

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}

	return t
}

// parsing function to embed the html file while building the app binaries
func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	// creating an empty template first so we can add
	// our custom template functions before parsing the HTML
	// so it can identify the function written inside the HTML
	tpl := template.New(patterns[0])

	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
		},
	)

	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{htmlTpl: tpl}, nil
}

// func Parse(filePath string) (Template, error) {
// 	tpl, err := template.ParseFiles(filePath)
// 	if err != nil {
// 		return Template{}, fmt.Errorf("parsing template: %w", err)
// 	}

// 	return Template{htmlTpl: tpl}, nil
// }

type Template struct {
	htmlTpl *template.Template
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	// to avoid multiple request executing the same template
	// will clone the template every time the execute function is called
	// so every user will their own version of the template
	tpl, err := t.htmlTpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error executing the template.",
			http.StatusInternalServerError)
		return
	}

	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset-utf-8")

	// putting all the HTML inside a buffer fist before executing will
	// make sure that all the data is rendered
	// can cause issue if we have a big HTML size
	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	// err = tpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template: %v", err)
		http.Error(w, "There was an error executing the template.",
			http.StatusInternalServerError)
		return
	}

	io.Copy(w, &buf)
}
