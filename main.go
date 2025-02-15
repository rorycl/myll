package main

import (
	"fmt"
	"html/template"
	"log"
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

func home(w http.ResponseWriter, r *http.Request) {
	// log.Printf("%#v\n", templates.DefinedTemplates())
	err := templates.ExecuteTemplate(w, "home.html", map[string]string{"URL": r.RequestURI})
	if err == nil {
		return
	}
	log.Println(err)
	http.Error(
		w,
		fmt.Sprintf("There was an error parsing the template: %s", err),
		http.StatusInternalServerError,
	)
}

func faq(w http.ResponseWriter, r *http.Request) {
	something := chi.URLParam(r, "something")
	fmt.Fprintf(w, "faq (with %s)", something)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not found", http.StatusNotFound)
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
	r.NotFound(notFound)
	fmt.Println("serving on :3000")
	http.ListenAndServe(":3000", r)
}
