package views

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
)

// user.go are the pages for user-related actions

// UserView is the viewer struct containing the fs.FS fileSystem holding
// templates and from whose methods endpoint output is provided.
type UserView struct {
	templates     map[string]*template.Template
	inDevelopment bool
	fS            *fileSystem
	done          chan bool
	UpdateChan    <-chan bool
}

// NewUserView creates a new View.
func NewUserView(fsName string, fS fs.FS, path string, inDevelopment bool) (*UserView, error) {
	fs, err := newFileSystem(fsName, fS, path, inDevelopment)
	if err != nil {
		return nil, fmt.Errorf("new view: cannot mount file system: %w", err)
	}
	v := &UserView{fS: fs, inDevelopment: inDevelopment}
	if err != nil {
		return nil, fmt.Errorf("new view: cannot parse templates: %w", err)
	}
	if inDevelopment { // make template monitor
		v.done = make(chan bool)
		v.UpdateChan, err = watchDir(path, tplFileRegexp, v.done)
		if err != nil {
			return nil, fmt.Errorf("new view: could not monitor templates: %w", err)
		}
	}
	return v, nil
}

// NewUser shows a new user after succesful creation
func (pv *UserView) NewUser() (func(w http.ResponseWriter, data any), error) {
	return genericViewMaker(pv.fS.fS, "user_new.html", "tailwind.html")
}
