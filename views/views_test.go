package views

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// tpls defined in fs_test.go

func TestViewGeneral(t *testing.T) {

	thisFs, err := newFileSystem("test", nil, "templates", true)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		input        map[string]string
		thisFunc     func(fs.FS) (func(w http.ResponseWriter, data any), error)
		expectedBody string
	}{

		{
			input:        map[string]string{"URL": "/notfound"},
			thisFunc:     NotFound,
			expectedBody: "/notfound",
		},
		{
			input:        map[string]string{"inDevelopment": "true", "problem": "definitely"},
			thisFunc:     InternalError,
			expectedBody: "definitely",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			w := httptest.NewRecorder()
			writerFunc, err := tt.thisFunc(thisFs.fS)
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
