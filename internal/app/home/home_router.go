package home

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"viabl.ventures/gossr/internal/templates"
)

type HomeRouter struct {
	renderer *templates.Renderer
}

func NewHomeRouter(renderer *templates.Renderer) *HomeRouter {
	return &HomeRouter{renderer}
}

func (router *HomeRouter) GetRoutes(r chi.Router) {
	r.Get("/", homePageView(router.renderer))
}

func homePageView(renderer *templates.Renderer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title": "Home Page",
		}
		renderer.RenderTemplate(w, "home.html", data)
	}
}
