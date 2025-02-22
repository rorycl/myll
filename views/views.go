package views

// views.go sets out the general view funcs and 404/500 endpoints

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"regexp"
)

// tplFileRegexp matches files that we're interested in seeing changes
// to in the template directory
var tplFileRegexp *regexp.Regexp = regexp.MustCompile("(?i)^[a-z0-9].+html$")

// renderTemplate attempts to render a template t to ResponseWriter w
// using data, reporting an http.Error if the execution fails.
func renderTemplate(w http.ResponseWriter, t *template.Template, data any) {
	err := t.Execute(w, data)
	if err == nil {
		return
	}
	http.Error(
		w,
		fmt.Sprintf("template rendering error: <pre>%s</pre>", err),
		http.StatusInternalServerError,
	)
}

// notfound view is a 404 view
func NotFound(fS fs.FS) (func(w http.ResponseWriter, data any), error) {
	pages := []string{"404.html", "tailwind.html"}
	t, err := template.ParseFS(fS, pages...)
	if err != nil {
		return nil, fmt.Errorf("template parse error for %v: %w", pages, err)
	}
	return func(w http.ResponseWriter, data any) {
		renderTemplate(w, t, data)
	}, nil
}

// internalError view is a 500 view
func InternalError(fS fs.FS) (func(w http.ResponseWriter, data any), error) {
	pages := []string{"500.html", "tailwind.html"}
	t, err := template.ParseFS(fS, pages...)
	if err != nil {
		return nil, fmt.Errorf("template parse error for %v: %w", pages, err)
	}
	return func(w http.ResponseWriter, data any) {
		renderTemplate(w, t, data)
	}, nil
}
