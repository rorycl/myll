package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
)

// public.go are the public pages

// PublicView is the viewer struct containing the fs.FS fileSystem holding
// templates and from whose methods endpoint output is provided.
type PublicView struct {
	templates     map[string]*template.Template
	inDevelopment bool
	fS            *fileSystem
	done          chan bool
	UpdateChan    <-chan bool
}

// NewView creates a new View.
func NewPublicView(fsName string, fS fs.FS, path string, inDevelopment bool) (*PublicView, error) {
	fs, err := newFileSystem(fsName, fS, path, inDevelopment)
	if err != nil {
		return nil, fmt.Errorf("new view: cannot mount file system: %w", err)
	}
	v := &PublicView{fS: fs, inDevelopment: inDevelopment}
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

// home view
func (pv *PublicView) Home() (func(w http.ResponseWriter, data any), error) {
	pages := []string{"home.html", "tailwind.html"}
	t, err := template.ParseFS(pv.fS.fS, pages...)
	if err != nil {
		return nil, fmt.Errorf("template parse error for %v: %w", pages, err)
	}
	return func(w http.ResponseWriter, data any) {
		renderTemplate(w, t, data)
	}, nil
}

// contact view
func (pv *PublicView) Contact() (func(w http.ResponseWriter, data any), error) {
	pages := []string{"contact.html", "tailwind.html"}
	t, err := template.ParseFS(pv.fS.fS, pages...)
	if err != nil {
		return nil, fmt.Errorf("template parse error for %v: %w", pages, err)
	}
	return func(w http.ResponseWriter, data any) {
		renderTemplate(w, t, data)
	}, nil
}

// faq view
func (pv *PublicView) FAQ() (func(w http.ResponseWriter, data any), error) {
	pages := []string{"faq.html", "tailwind.html"}
	t, err := template.ParseFS(pv.fS.fS, pages...)
	if err != nil {
		return nil, fmt.Errorf("template parse error for %v: %w", pages, err)
	}
	return func(w http.ResponseWriter, data any) {
		renderTemplate(w, t, data)
	}, nil
}

// signup view
func (pv *PublicView) Signup() (func(w http.ResponseWriter, data any), error) {
	pages := []string{"signup.html", "tailwind.html"}
	t, err := template.ParseFS(pv.fS.fS, pages...)
	if err != nil {
		return nil, fmt.Errorf("template parse error for %v: %w", pages, err)
	}
	return func(w http.ResponseWriter, data any) {
		renderTemplate(w, t, data)
	}, nil
}

// NotFound uses the generic NotFound endpoint
func (pv *PublicView) NotFound() (func(w http.ResponseWriter, data any), error) {
	return NotFound(pv.fS.fS)
}

// InternalError uses the generic Internal Error
func (pv *PublicView) InternalError() (func(w http.ResponseWriter, data any), error) {
	return InternalError(pv.fS.fS)
}
