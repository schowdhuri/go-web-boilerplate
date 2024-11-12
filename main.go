package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"viabl.ventures/gossr/internal/app"
	"viabl.ventures/gossr/internal/app/admin"
	"viabl.ventures/gossr/internal/app/home"
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
	homeContainer := home.NewHomeContainer(baseContainer)

	// Routes
	r.Route("/", homeContainer.Router.GetRoutes)
	r.Route("/admin", adminContainer.Router.GetRoutes)

	port := conf.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

	if isDev {
		var wg sync.WaitGroup
		wg.Add(3)
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
