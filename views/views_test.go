package views

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestViews(t *testing.T) {

	tplFunc, err := LoadTemplates("templates/*html")
	if err != nil {
		t.Fatal(err)
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
}
