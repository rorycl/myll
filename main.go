package main

import (
	"fmt"
	"net/http"
	"os"

	"mylenslocked/views"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

// global: to replace
var writeTemplate views.WriteTplFunc

func home(w http.ResponseWriter, r *http.Request) {
	writeTemplate(w, "home.html", map[string]string{"URL": r.RequestURI})
}

func contact(w http.ResponseWriter, r *http.Request) {
	writeTemplate(w, "contact.html", map[string]string{"URL": r.RequestURI})
}

func faq(w http.ResponseWriter, r *http.Request) {
	something := chi.URLParam(r, "something")
	fmt.Fprintf(w, "faq (with %s)", something)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	writeTemplate(w, "404.html", map[string]string{"URL": r.RequestURI})
}

func main() {
	var err error
	writeTemplate, err = views.LoadTemplates("views/templates/*html")
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
