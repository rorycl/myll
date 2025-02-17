package controller

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type renderFunc func(w http.ResponseWriter, data any, pages ...string)

type controller struct {
	render renderFunc
}

func newController(f fs.FS, rf renderFunc) *controller {
	return &controller{f, rf}
}

func (c *controller) home(w http.ResponseWriter, r *http.Request) {
	c.render(w, map[string]string{"URL": r.RequestURI}, "page.html", "home.html")
}

func (c *controller) contact(w http.ResponseWriter, r *http.Request) {
	c.render(w, map[string]string{"URL": r.RequestURI}, "page.html", "contact.html")
}

func (c *controller) faq(w http.ResponseWriter, r *http.Request) {
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
	c.render(w, questions, "page.html", "faq.html")
}

func (c *controller) faq2(w http.ResponseWriter, r *http.Request) {
	something := chi.URLParam(r, "something")
	fmt.Fprintf(w, "faq (with %s)", something)
}

func (c *controller) notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	c.render(w, questions, "page.html", "404.html")
}

func Serve(rf renderFunc) {

	c := newController(rf)

	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)

	// routes
	r.Get("/", c.home)
	r.Get("/faq", c.faq)
	r.Get("/faq/", c.faq)
	r.Get("/faq/{something}", c.faq)
	r.Get("/contact", c.contact)
	r.NotFound(c.notFound)
	fmt.Println("serving on :3000")
	http.ListenAndServe(":3000", r)
}
