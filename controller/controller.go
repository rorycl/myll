package controller

// general funcs and types for controllers

import (
	"fmt"
	"net/http"
)

// makeRenderer mounts the renderer in a controller closure
func makeRenderer(endpoint string, vf viewFunc) func(http.ResponseWriter, any) {
	r, err := vf.render()
	if err != nil {
		panic(fmt.Sprintf("view for %s could not be mounted: %v", endpoint, err))
	}
	return r
}

// viewFunc makes an abstrace type of any view function
type viewFunc func() (func(w http.ResponseWriter, data any), error)

// render adds a generic "render" function to a viewFunc
func (v viewFunc) render() (func(w http.ResponseWriter, data any), error) {
	return v()
}

// renderer is an interface that describes a decorated viewFunc
// (currently unused)
type renderer interface {
	render() (func(w http.ResponseWriter, data any), error)
}
