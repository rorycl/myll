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
		h      func(viewFunc) func(http.ResponseWriter, *http.Request)
		url    string
		status int
		want   string
	}{
		{
			h:      Home,
			url:    "/",
			status: http.StatusOK,
			want:   "content: /",
		},
		{
			h:      Contact,
			url:    "/contact",
			status: http.StatusOK,
			want:   "content: /contact",
		},
		{
			h:      FAQ,
			url:    "/faq",
			status: http.StatusOK,
			want:   "content: /faq",
		},
		{
			h:      NotFound,
			url:    "/404",
			status: http.StatusNotFound,
			want:   "content: /404",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {

			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			wr := httptest.NewRecorder()

			fakeView := func() (func(w http.ResponseWriter, data any), error) {
				return func(w http.ResponseWriter, data any) {
					fmt.Fprintf(w, "content: %s\n", req.URL)
				}, nil
			}

			var endpoint http.HandlerFunc
			endpoint = tt.h(fakeView)
			if err != nil {
				t.Fatalf("could not initialise endpoint: %v", err)
			}
			endpoint.ServeHTTP(wr, req) // tt.endpoint is already coerced to a HandlerFunc

			if got, want := wr.Code, tt.status; got != want {
				t.Errorf("status code: got %d want %d", got, want)
			}

			if !strings.Contains(wr.Body.String(), tt.want) {
				t.Errorf("unexpected body: got %v want %v",
					wr.Body.String(), tt.want)
			}
		})
	}
}
