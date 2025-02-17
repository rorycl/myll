package controller

import (
	"fmt"
	"net/http"

	"mylenslocked/views"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

// controller is the struct off which url endpoints are routed to views
// in viewer.
type controller struct {
	viewer *views.View
	// render func(endpoint string, w http.ResponseWriter, r *http.Response, data any)
}

// newController makes a new controller
func newController(v *views.View) *controller {
	// return &controller{viewer: v, render: v.render}
	return &controller{viewer: v}
}

// home resolves the / and /home endpoints
func (c *controller) home(w http.ResponseWriter, r *http.Request) {
	// c.render("home", w, r, nil)
	c.viewer.Home(w, r)
}

// contact resolves the /contact endpoint
func (c *controller) contact(w http.ResponseWriter, r *http.Request) {
	c.viewer.Contact(w, r)
}

// faq resolves the /faq endpoint
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
	c.viewer.FAQ(w, r, questions)
}

// faq2 is an endpoint for messing around with routing, url params
func (c *controller) faq2(w http.ResponseWriter, r *http.Request) {
	something := chi.URLParam(r, "something")
	fmt.Fprintf(w, "faq (with %s)", something)
}

// notFound is a 404 endpoint
func (c *controller) notFound(w http.ResponseWriter, r *http.Request) {
	c.viewer.NotFound(w, r)
}

// Serve serves the urls and routes them to the associated endpoints
func Serve(viewer *views.View) {

	c := newController(viewer)
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)

	// routes
	r.Get("/", c.home)
	r.Get("/home", c.home)
	r.Get("/faq", c.faq)
	r.Get("/faq/", c.faq)
	r.Get("/faq/{something}", c.faq)
	r.Get("/contact", c.contact)
	r.NotFound(c.notFound)
	fmt.Println("serving on :3000")
	http.ListenAndServe(":3000", r)
}
