package controller

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServe(t *testing.T) {

	// https://www.cloudbees.com/blog/testing-http-handlers-go
	// w := httptest.NewResponse()

	tests := []struct {
		h      http.HandlerFunc
		status int
		want   string
	}{
		{
			h:      home,
			status: http.StatusOK,
			want:   "url: /",
		},
		{
			h:      contact,
			status: http.StatusOK,
			want:   "page: contact.html",
		},
		{
			h:      notFound,
			status: http.StatusNotFound,
			want:   "page: 404.html",
		},
	}

	wt := func(w http.ResponseWriter, r *http.Response) {
		fmt.Fprintf(w, "url: %s\n", r.Request.URL)
	}
	writeTemplate = wt

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {

			req, err := http.NewRequest("GET", "whatever", nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			tt.h.ServeHTTP(rr, req) // tt.h is already coerced to a HandlerFunc

			if got, want := rr.Code, tt.status; got != want {
				t.Errorf("status code: got %d want %d", got, want)
			}

			if !strings.Contains(rr.Body.String(), tt.want) {
				t.Errorf("unexpected body: got %v want %v",
					rr.Body.String(), tt.want)
			}
		})
	}
}
