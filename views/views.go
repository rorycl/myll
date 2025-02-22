package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"regexp"
)

// View is the viewer struct containing the fs.FS fileSystem holding
// templates and from whose methods endpoint output is provided.
type View struct {
	templates     map[string]*template.Template
	inDevelopment bool
	fS            *fileSystem
}

var fileMatcher *regexp.Regexp = regexp.MustCompile("(?i)^[a-z0-9].+html")

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

	// detect updates to files
	if inDevelopment {
		done := make(chan bool)
		updateEvents, err := watchDir(path, fileMatcher, done)
		if err != nil {
			return v, fmt.Errorf("notify event watcher error: %w", err)
		}
		go func() {
			for range updateEvents {
				log.Println("template change detected, reloading...")
				_ = v.parseTemplates()
			}
		}()
	}
	return v, nil
}

// parseTemplates parses templates for endpoints from the embedded
// fileSystem and stores them in View.templates
func (v *View) parseTemplates() error {
	endpointToTpls := map[string][]string{
		"home":    []string{"home.html", "tailwind.html"},
		"contact": []string{"contact.html", "tailwind.html"},
		"faq":     []string{"faq.html", "tailwind.html"},
		"404":     []string{"404.html", "tailwind.html"},
		"signup":  []string{"signup.html", "tailwind.html"},
		"usernew": []string{"usernew.html", "tailwind.html"},
	}
	v.templates = map[string]*template.Template{}
	for endpoint, pages := range endpointToTpls {
		t, err := template.ParseFS(v.fS.fS, pages...)
		if err != nil {
			return fmt.Errorf("template parse error for %v: %w", pages, err)
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

// Render is a router for view template rendering
func (v *View) Render(endpoint string, w http.ResponseWriter, data any) {
	switch endpoint {
	case "home":
		v.home(w, data)
	case "contact":
		v.contact(w, data)
	case "faq":
		v.faq(w, data)
	case "signup":
		v.signup(w, data)
	case "usernew":
		v.usernew(w, data)
	default:
		v.notFound(w, data)
	}
}

// home view
func (v *View) home(w http.ResponseWriter, data any) {
	v.renderTemplate(w, v.templates["home"], data)
}

// contact view
func (v *View) contact(w http.ResponseWriter, data any) {
	v.renderTemplate(w, v.templates["contact"], data)
}

// faq view, which receives faq data.
func (v *View) faq(w http.ResponseWriter, data any) {
	v.renderTemplate(w, v.templates["faq"], data)
}

// signup view, which receives signup data.
func (v *View) signup(w http.ResponseWriter, data any) {
	v.renderTemplate(w, v.templates["signup"], data)
}

// usernew view, to make a new user
func (v *View) usernew(w http.ResponseWriter, data any) {
	v.renderTemplate(w, v.templates["usernew"], data)
}

// notfound view for 404 errors.
func (v *View) notFound(w http.ResponseWriter, data any) {
	w.WriteHeader(http.StatusNotFound)
	v.renderTemplate(w, v.templates["404"], data)
}
