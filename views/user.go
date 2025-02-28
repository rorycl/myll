package views

import (
	"io/fs"
	"net/http"
)

// user.go are the pages for user-related actions

// UserView is a View type
type UserView View

// NewUserView initialises a new UserView
func NewUserView(fsName string, fS fs.FS, path string, inDevelopment bool) (*UserView, error) {
	v, err := NewView(fsName, fS, path, inDevelopment)
	uv := UserView(*v)
	return &uv, err
}

// NewUserViewFromView initialises a new UserView from an existing view
// which has the same mount point.
func NewUserViewFromView(v *View) *UserView {
	uv := UserView(*v)
	return &uv
}

// NewUser shows a new user after succesful creation
func (pv *UserView) NewUser() (func(w http.ResponseWriter, data any), error) {
	return genericViewMaker(pv.fS.fS, "user_new.html", "tailwind.html")
}
