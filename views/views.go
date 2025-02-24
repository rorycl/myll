package views

// views.go sets out the general view funcs and 404/500 endpoints

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"regexp"
)

// View is a generic view struct containing the fs.FS fileSystem holding
// templates and from whose methods endpoint output is provided.
type View struct {
	templates     map[string]*template.Template
	inDevelopment bool
	fS            *fileSystem
	done          chan bool
	UpdateChan    <-chan bool
}

// NewView makes a View
func NewView(fsName string, fS fs.FS, path string, inDevelopment bool) (*View, error) {
	fs, err := newFileSystem(fsName, fS, path, inDevelopment)
	if err != nil {
		return nil, fmt.Errorf("new view: cannot mount file system: %w", err)
	}
	v := &View{fS: fs, inDevelopment: inDevelopment}
	if err != nil {
		return nil, fmt.Errorf("new view: cannot parse templates: %w", err)
	}
	if inDevelopment { // make template monitor
		v.done = make(chan bool)
		v.UpdateChan, err = watchDir(path, tplFileRegexp, v.done)
		if err != nil {
			return nil, fmt.Errorf("new view: could not monitor templates: %w", err)
		}
	}
	return v, nil
}

// tplFileRegexp matches files that we're interested in seeing changes
// to in the template directory
var tplFileRegexp *regexp.Regexp = regexp.MustCompile("(?i)^[a-z0-9].+html$")

// genericViewMaker makes a generic view
func genericViewMaker(fS fs.FS, pages ...string) (func(w http.ResponseWriter, data any), error) {
	t, err := template.ParseFS(fS, pages...)
	if err != nil {
		return nil, fmt.Errorf("template parse error for %v: %w", pages, err)
	}
	return func(w http.ResponseWriter, data any) {
		renderTemplate(w, t, data)
	}, nil
}

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
