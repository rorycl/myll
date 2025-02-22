package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"mylenslocked/controller"
	"mylenslocked/views"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

//go:embed views/templates/*html
var tpls embed.FS

var inDevelopment bool = true

func main() {
	viewerPublic, err := views.NewPublicView("public", nil, "views/templates", inDevelopment)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viewerUser, err := views.NewUserView("user", nil, "views/templates", inDevelopment)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Logger)
	// r.Use(middleware.Recoverer)

	// attach public routes
	publicRoutes := func(pv *views.PublicView) {
		r.Get("/", controller.Home(pv.Home))
		r.Get("/home", controller.Home(pv.Home))
		r.Get("/faq", controller.FAQ(pv.FAQ))
		r.Get("/faq/", controller.FAQ(pv.FAQ))
		r.Get("/faq/{etc}", controller.FAQ(pv.FAQ))
		r.Get("/contact", controller.Contact(pv.Contact))
		r.Get("/signup", controller.Signup(pv.Signup))
		r.NotFound(controller.NotFound(pv.NotFound))
	}

	// attach user routes
	userRoutes := func(uv *views.UserView) {
		r.Post("/user", controller.UserNew(uv.NewUser))
	}

	publicRoutes(viewerPublic)
	userRoutes(viewerUser)

	if inDevelopment {
		go func() {
			// ok to do this since these two viewers are the same
			for range viewerPublic.UpdateChan {
				log.Println("template change detected...reloading")
				publicRoutes(viewerPublic)
				userRoutes(viewerUser)
			}
		}()
	}

	fmt.Println("serving on :3000")
	http.ListenAndServe(":3000", r)

}
