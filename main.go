package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"viabl.ventures/gossr/internal/app"
	"viabl.ventures/gossr/internal/app/admin"
	home "viabl.ventures/gossr/internal/app/website"
	"viabl.ventures/gossr/internal/assets"
	"viabl.ventures/gossr/internal/config"
	"viabl.ventures/gossr/internal/middleware"
	"viabl.ventures/gossr/internal/templates"
)

func main() {
	conf := config.NewConfig()
	isDev := conf.GoEnv == "development"
	assetPipe := assets.NewAssetPipeline(isDev)
	renderer := templates.NewRenderer(assetPipe, isDev)

	r := chi.NewRouter()

	// Middleware
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.CompressionMiddleware)
	r.Use(middleware.CacheControlMiddleware)

	// Static files
	fileServer := http.FileServer(http.Dir("dist"))
	r.Handle("/dist/*", http.StripPrefix("/dist/", fileServer))

	// initialize app containers
	baseContainer := app.NewBaseContainer(conf, renderer)
	adminContainer := admin.NewAdminContainer(baseContainer)
	websiteContainer := home.NewWebsiteContainer(baseContainer)

	// Routes
	r.Route("/", websiteContainer.Router.GetRoutes)
	r.Route("/admin", adminContainer.Router.GetRoutes)

	port := conf.Port
	if port == "" {
		port = "8080"
	}

	var wg sync.WaitGroup
	if isDev {
		wg.Add(4)
	} else {
		wg.Add(1)
	}

	go func() {
		defer wg.Done()
		log.Printf("Server starting on :%s\n", port)
		log.Fatal(http.ListenAndServe(":"+port, r))
	}()

	if isDev {
		fmt.Println("Starting asset pipeline...")
		go func() {
			defer wg.Done()
			assets.BuildJs(true)
		}()
		go func() {
			defer wg.Done()
			assets.BuildCss(true)
		}()
		go func() {
			defer wg.Done()
			assets.CopyPublicAssets(true)
		}()

		wg.Wait()
	}
}
