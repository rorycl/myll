package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
)

// View is the viewer struct containing the fs.FS fileSystem holding
// templates and from whose methods endpoint output is provided.
type View struct {
	templates     map[string]*template.Template
	inDevelopment bool
	fS            *fileSystem
}

// NewView creates a new View.
func NewView(fsName string, fS fs.FS, path string, inDevelopment bool) (*View, error) {
	fs, err := newFileSystem(fsName, fS, path, inDevelopment)
	if err != nil {
		return nil, fmt.Errorf("new view: cannot mount file system: %w", err)
	}
	v := &View{fS: fs, inDevelopment: inDevelopment}
	err = v.parseTemplates()
	if err != nil {
		return nil, fmt.Errorf("new view: cannot parse templates: %w", err)
	}
	return v, nil
}

// parseTemplates parses templates for endpoints from the embedded
// fileSystem and stores them in View.templates
func (v *View) parseTemplates() error {
	endpointToTpls := map[string][]string{
		"home":    []string{"page.html", "home.html"},
		"contact": []string{"page.html", "contact.html"},
		"faq":     []string{"page.html", "faq.html"},
		"404":     []string{"page.html", "404.html"},
	}
	v.templates = map[string]*template.Template{}
	for endpoint, pages := range endpointToTpls {
		t, err := template.ParseFS(v.fS.fS, pages...)
		if err != nil {
			return err
		}
		v.templates[endpoint] = t
	}
	return nil
}

// renderTemplate attempts to render a template t to ResponseWriter w
// using data, reporting an http.Error if the execution fails.
func (v *View) renderTemplate(w http.ResponseWriter, t *template.Template, data any) {
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

/*
func (v *View) render(endpoint string, w http.ResponseWriter, data any) {
	switch endpoint {
	case "home":
		v.Home(w, r)
	case "contact":
		v.Contact(w, r)
	case "faq":
		v.FAQ(w, r, data)
	default:
		v.NotFound(w, r)
	}
}
*/

// Home view
func (v *View) Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"URL": r.RequestURI}
	v.renderTemplate(w, v.templates["home"], data)
}

// Contact view
func (v *View) Contact(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"URL": r.RequestURI}
	v.renderTemplate(w, v.templates["contact"], data)
}

// FAQ view, which receives FAQ data.
func (v *View) FAQ(w http.ResponseWriter, r *http.Request, data any) {
	v.renderTemplate(w, v.templates["faq"], data)
}

// NotFound view for 404 errors.
func (v *View) NotFound(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"URL": r.RequestURI}
	w.WriteHeader(http.StatusNotFound)
	v.renderTemplate(w, v.templates["404"], data)
}
