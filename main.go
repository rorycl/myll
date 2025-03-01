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
	"github.com/gorilla/csrf"
)

//go:embed views/templates/*html
var tpls embed.FS

var inDevelopment bool = true

func main() {

	csrfToken := os.Getenv("MYLL_CSRF_TOKEN")
	if csrfToken == "" || len(csrfToken) != 32 {
		fmt.Println("MYLL_CSRF_TOKEN is not set or invalid, exiting")
		os.Exit(1)
	}
	dbURL := os.Getenv("MYLL_DB_URL")
	if dbURL == "" {
		fmt.Println("MYLL_DB_URL is not set, exiting")
		os.Exit(1)
	}

	viewer, err := views.NewView("public", nil, "views/templates", inDevelopment)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viewerPublic := views.NewPublicViewFromView(viewer)
	viewerUser := views.NewUserViewFromView(viewer)

	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Logger)

	csrfMiddleware := csrf.Protect(
		[]byte(csrfToken),
		csrf.MaxAge(0),
		csrf.Secure(true),
		csrf.HttpOnly(true),
		csrf.CookieName("myll"),
		csrf.SameSite(csrf.SameSiteStrictMode),
	)
	if inDevelopment { // override csrf if in development
		csrfMiddleware = csrf.Protect(
			[]byte(csrfToken),
			csrf.MaxAge(0),
			csrf.Secure(false), // allow non-secure
			csrf.HttpOnly(true),
			csrf.CookieName("myll"),
			csrf.SameSite(csrf.SameSiteStrictMode),
		)
	}
	r.Use(csrfMiddleware)
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
			// views share same mount point
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
