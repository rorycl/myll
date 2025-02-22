package views

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// tpls defined in fs_test.go

func TestViewPublic(t *testing.T) {

	pv, err := NewPublicView("test", nil, "templates", true)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		input        map[string]string
		thisFunc     func() (func(w http.ResponseWriter, data any), error)
		expectedBody string
	}{

		{
			input:        map[string]string{"URL": "/"},
			thisFunc:     pv.Home,
			expectedBody: "Welcome",
		},
		{
			input:        map[string]string{"URL": "/contact"},
			thisFunc:     pv.Contact,
			expectedBody: "contact",
		},
		{
			input:        map[string]string{"URL": "/faq"},
			thisFunc:     pv.FAQ,
			expectedBody: "FAQ",
		},
		{
			input:        map[string]string{"URL": "/signup"},
			thisFunc:     pv.Signup,
			expectedBody: "Start sharing",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			w := httptest.NewRecorder()
			writerFunc, err := tt.thisFunc()
			if err != nil {
				t.Fatal(err)
			}
			writerFunc(w, tt.input)
			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(body), tt.expectedBody) {
				fmt.Println(string(body))
				t.Errorf("no '%s' found in body", tt.expectedBody)
			}
		})
	}
}
