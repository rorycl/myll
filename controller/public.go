package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// home resolves the / and /home endpoints
func Home(vf viewFunc) func(http.ResponseWriter, *http.Request) {
	rdr := makeRenderer("Home", vf)
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{"URL": r.RequestURI}
		rdr(w, data)
	}
}

// contact resolves the /contact endpoint
func Contact(vf viewFunc) func(http.ResponseWriter, *http.Request) {
	rdr := makeRenderer("Contact", vf)
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{"URL": r.RequestURI}
		rdr(w, data)
	}
}

// faq resolves the /faq endpoint
func FAQ(vf viewFunc) func(http.ResponseWriter, *http.Request) {
	rdr := makeRenderer("FAQ", vf)
	return func(w http.ResponseWriter, r *http.Request) {
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
		rdr(w, map[string]any{"Questions": questions, "Params": etc})
	}
}

// faq2 is an endpoint for messing around with routing, url params etc.
/*
func faq2() func(http.ResponseWriter, *http.Request) {
	 rdr := pv
	 return func(w http.ResponseWriter, r *http.Request) {
	etc := chi.URLParam(r, "etc")
	fmt.Fprintf(w, "faq (with %s)", etc)
}
*/

// signup resolves the /signup endpoint
func Signup(vf viewFunc) func(http.ResponseWriter, *http.Request) {
	rdr := makeRenderer("Signup", vf)
	return func(w http.ResponseWriter, r *http.Request) {
		// validate user
		// data := map[string]any{
		// 	"CSRFField": csrf.TemplateField(r),
		// }
		data := struct{}
		rdr(w, data)
	}
}

func NotFound(vf viewFunc) func(http.ResponseWriter, *http.Request) {
	rdr := makeRenderer("NotFound", vf)
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		data := map[string]string{}
		rdr(w, data)
	}
}

func InternalError(vf viewFunc) func(http.ResponseWriter, *http.Request) {
	rdr := makeRenderer("InternalError", vf)
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		data := map[string]string{}
		rdr(w, data)
	}
}
