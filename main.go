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
	writeTemplate, err := views.LoadTemplates(tpls, "views/templates", "*html", false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	controller.Serve(writeTemplate)
}
