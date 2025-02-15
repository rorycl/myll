package views

import (
	"fmt"
	"html/template"
	"net/http"
)

type WriteTplFunc func(w http.ResponseWriter, page string, data any)

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
