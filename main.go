package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

var templates *template.Template

func loadTemplates(glob string) error {
	var err error
	templates, err = template.ParseGlob(glob)
	return err
}

func writeTemplate(page string, status int, data map[string]string, w http.ResponseWriter, r *http.Request) {
	if status != http.StatusOK {
		w.WriteHeader(status)
	}
	err := templates.ExecuteTemplate(w, page, map[string]string{"URL": r.RequestURI})
	if err == nil {
		return
	}
	http.Error(
		w,
		fmt.Sprintf("There was an error parsing the template: %s", err),
		http.StatusInternalServerError,
	)
}

func home(w http.ResponseWriter, r *http.Request) {
	writeTemplate("home.html", 200, map[string]string{"URL": r.RequestURI}, w, r)
}

func contact(w http.ResponseWriter, r *http.Request) {
	writeTemplate("contact.html", 200, map[string]string{"URL": r.RequestURI}, w, r)
}

func faq(w http.ResponseWriter, r *http.Request) {
	something := chi.URLParam(r, "something")
	fmt.Fprintf(w, "faq (with %s)", something)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	writeTemplate("404.html", 404, map[string]string{"URL": r.RequestURI}, w, r)
}

func main() {
	err := loadTemplates("templates/*")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)
	r.Get("/", home)
	// r.Get("/qaf(/{something})*", faq)
	r.Get("/faq", faq)
	r.Get("/faq/", faq)
	r.Get("/faq/{something}", faq)
	r.Get("/contact", contact)
	r.NotFound(notFound)
	fmt.Println("serving on :3000")
	http.ListenAndServe(":3000", r)
}
