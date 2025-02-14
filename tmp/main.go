package main

import (
	"fmt"
	"net/http"
)

func v1() {
	http.HandleFunc("/", home)
	http.HandleFunc("/faq", faq)
	fmt.Println("starting on port 3000")
	http.ListenAndServe(":3000", nil)
}

func v2() {
	s := http.NewServeMux()
	s.HandleFunc("/", home)
	s.HandleFunc("/faq", faq)
	fmt.Println("starting on port 3000")
	http.ListenAndServe(":3000", s)
}

func v3() {
	pathHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			home(w, r)
		case "/faq":
			faq(w, r)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}
	var s http.HandlerFunc
	s = pathHandler
	fmt.Println("starting on port 3000")
	http.ListenAndServe(":3000", s)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "home")
}

func faq(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "faq")
}

func main() {
	v3()
}
