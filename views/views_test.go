package views

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

// tpls defined in fs_test.go

func TestViews(t *testing.T) {

	tests := []struct {
		fS            fs.FS
		path, pat     string
		inDevelopment bool
		ok            bool
	}{
		{
			fS:            tpls,
			path:          "templates",
			pat:           "*html",
			inDevelopment: false,
			ok:            true,
		},
		{
			fS:            tpls,
			path:          "templates",
			pat:           "abcd",
			inDevelopment: true,
			ok:            false,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			tplFunc, err := LoadTemplates(tt.fS, tt.path, tt.pat, tt.inDevelopment)
			if !(tt.ok == (err == nil)) {
				t.Fatalf("ok %t and err %s", tt.ok, err)
			}
			if !tt.ok {
				return
			}
			w := httptest.NewRecorder()
			tplFunc(w, "home", map[string]string{"URL": "test"})
			body, err := ioutil.ReadAll(w.Body)
			if err != nil {
				t.Fatal(err)
			}
			if strings.Contains(string(body), "test") {
				t.Errorf("no test in body")
			}

		})
	}
}
