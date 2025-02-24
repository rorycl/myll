package views

import (
	"io/fs"
	"net/http"
)

// public.go are the public pages

// PublicView is a View type
type PublicView View

// NewPublicView initialises a new PublicView
func NewPublicView(fsName string, fS fs.FS, path string, inDevelopment bool) (*PublicView, error) {
	v, err := NewView(fsName, fS, path, inDevelopment)
	pv := PublicView(*v)
	return &pv, err
}

// NewView creates a new View.

// home view
func (pv *PublicView) Home() (func(w http.ResponseWriter, data any), error) {
	return genericViewMaker(pv.fS.fS, "home.html", "tailwind.html")
}

// contact view
func (pv *PublicView) Contact() (func(w http.ResponseWriter, data any), error) {
	return genericViewMaker(pv.fS.fS, "contact.html", "tailwind.html")
}

// faq view
func (pv *PublicView) FAQ() (func(w http.ResponseWriter, data any), error) {
	return genericViewMaker(pv.fS.fS, "faq.html", "tailwind.html")
}

// signup view
func (pv *PublicView) Signup() (func(w http.ResponseWriter, data any), error) {
	return genericViewMaker(pv.fS.fS, "signup.html", "tailwind.html")
}

// NotFound uses the generic NotFound endpoint
func (pv *PublicView) NotFound() (func(w http.ResponseWriter, data any), error) {
	return NotFound(pv.fS.fS)
}

// InternalError uses the generic Internal Error
func (pv *PublicView) InternalError() (func(w http.ResponseWriter, data any), error) {
	return InternalError(pv.fS.fS)
}
