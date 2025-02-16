package views

import (
	"fmt"
	"html/template"
	"net/http"
)

// WriteTplFunc is the function signature for a template writing
// function returned from LoadTemplates
type WriteTplFunc func(w http.ResponseWriter, page string, data any)

// Load Templates is a closure for loading a group of templates
func LoadTemplates(globPath string) (WriteTplFunc, error) {
	templates, err := template.ParseGlob(globPath)
	if err != nil {
		return nil, err
	}
	return func(w http.ResponseWriter, page string, data any) {
		err := templates.ExecuteTemplate(w, page, data)
		if err == nil {
			return
		}
		http.Error(
			w,
			fmt.Sprintf("There was an error parsing the template: %s", err),
			http.StatusInternalServerError,
		)
	}, nil
}
