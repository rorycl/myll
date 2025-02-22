package controller

import (
	"net/http"
	"net/url"
)

// UserNew shows a new user page
func UserNew(vf viewFunc) func(http.ResponseWriter, *http.Request) {
	rdr := makeRenderer("New User", vf)
	return func(w http.ResponseWriter, r *http.Request) {
		vals, uerr := url.ParseQuery(r.URL.RawQuery)
		ferr := r.ParseForm()
		data := struct {
			URLValues map[string][]string
			UErr      error
			Form      map[string][]string
			FErr      error
		}{
			URLValues: vals,
			UErr:      uerr,
			Form:      r.PostForm,
			FErr:      ferr,
		}
		rdr(w, data)
	}
}
