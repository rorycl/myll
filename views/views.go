package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
)

// WriteTplFunc is the function signature for a template writing
// function returned from LoadTemplates
type WriteTplFunc func(w http.ResponseWriter, page string, data any)

// Load Templates is a closure for loading a group of templates either
// from disk (if inDevelopment) or an fS embedded file system at path
// using the glob pattern pat.
func LoadTemplates(fS fs.FS, path, pat string, inDevelopment bool) (WriteTplFunc, error) {

	f, err := newFileSystem("path", fS, path, inDevelopment)
	if err != nil {
		return nil, fmt.Errorf("error loading templates: %w", err)
	}
	templates, err := template.ParseFS(f.fS, pat)
	if err != nil {
		return nil, fmt.Errorf("error parsing templates: %w", err)
	}

	return func(w http.ResponseWriter, page string, data any) {
		err := templates.ExecuteTemplate(w, page, data)
		if err == nil {
			return
		}
		http.Error(
			w,
			fmt.Sprintf("There was a template error: %s", err),
			http.StatusInternalServerError,
		)
	}, nil
}
