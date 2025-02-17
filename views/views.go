package views

import (
	"fmt"
	"io/fs"
	"net/http"
)

// WriteTplFunc is the function signature for a template writing
// function returned from LoadTemplates
type WriteTplFunc func(w http.ResponseWriter, data any, pages ...string)

// Load Templates is a closure for loading a group of templates either
// from disk (if inDevelopment) or an fS embedded file system at path
// using the glob pattern pat.
func LoadTemplates(fS fs.FS, path, pat string, inDevelopment bool) (WriteTplFunc, error) {

	f, err := newFileSystem("path", fS, path, inDevelopment)
	if err != nil {
		return nil, fmt.Errorf("error loading templates: %w", err)
	}

	return func(w http.ResponseWriter, data any, pages ...string) {
		t, err := templates.ParseFS(f.fS, pages)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Template : %s", err),
				http.StatusInternalServerError,
			)
			return
		}
		err := t.Execute(w, data)
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
