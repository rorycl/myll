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
	render func(endpoint string, w http.ResponseWriter, data any)
}

// newController makes a new controller
func newController(v *views.View) *controller {
	return &controller{render: v.Render}
}

// routes adds the routes for the controller
func (c *controller) routes() *chi.Mux {
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)

	// routes
	r.Get("/", c.home)
	r.Get("/home", c.home)
	r.Get("/faq", c.faq)
	r.Get("/faq/", c.faq)
	r.Get("/faq/{etc}", c.faq)
	r.Get("/contact", c.contact)
	r.Get("/signup", c.signup)
	r.NotFound(c.notFound)

	return r
}

// home resolves the / and /home endpoints
func (c *controller) home(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"URL": r.RequestURI}
	c.render("home", w, data)
}

// contact resolves the /contact endpoint
func (c *controller) contact(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"URL": r.RequestURI}
	c.render("home", w, data)
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
	etc := chi.URLParam(r, "etc")
	c.render("faq", w, map[string]any{"Questions": questions, "Params": etc})
}

// faq2 is an endpoint for messing around with routing, url params etc.
func (c *controller) faq2(w http.ResponseWriter, r *http.Request) {
	etc := chi.URLParam(r, "etc")
	fmt.Fprintf(w, "faq (with %s)", etc)
}

// signup resolves the /signup endpoint
func (c *controller) signup(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{}
	c.render("signup", w, data)
}

// notFound is a 404 endpoint
func (c *controller) notFound(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"URL": r.RequestURI}
	c.render("notfound", w, data)
}

// Serve serves the urls and routes them to the associated endpoints
func Serve(viewer *views.View) {
	c := newController(viewer)
	r := c.routes()
	fmt.Println("serving on :3000")
	http.ListenAndServe(":3000", r)
}
