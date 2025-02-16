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
	questions := []struct {
		Question string
		Answer   string
	}{
		{
			Question: "Is there a free version?",
			Answer:   "Yes! We offer a free trial for 30 days on any paid plans.",
		},
		{
			Question: "What are your support hours?",
			Answer:   "We have support staff answering emails 24/7, though response times may be a bit slower on weekends.",
		},
		{
			Question: "How do I contact support?",
			Answer:   `Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a>`,
		},
	}
	writeTemplate(w, "faq.html", questions)
}

func faq2(w http.ResponseWriter, r *http.Request) {
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
