package main

import (
	"embed"
	"fmt"
	"os"

	"mylenslocked/controller"
	"mylenslocked/views"
)

//go:embed views/templates/*html
var tpls embed.FS

func main() {
	viewer, err := views.NewView("viewfs", tpls, "views/templates", true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	controller.Serve(viewer)
}
