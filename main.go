package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "home")
}

func faq(w http.ResponseWriter, r *http.Request) {
	something := chi.URLParam(r, "something")
	fmt.Fprintf(w, "faq (with %s)", something)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not found", http.StatusNotFound)
}

func main() {
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
