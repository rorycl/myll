package controller

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

// global: to replace
var writeTemplate func(w http.ResponseWriter, page string, data any)

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

func Serve(fn func(w http.ResponseWriter, page string, data any)) {

	writeTemplate = fn

	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)

	// routes
	r.Get("/", home)
	r.Get("/faq", faq)
	r.Get("/faq/", faq)
	r.Get("/faq/{something}", faq)
	r.Get("/contact", contact)
	r.NotFound(notFound)
	fmt.Println("serving on :3000")
	http.ListenAndServe(":3000", r)
}
