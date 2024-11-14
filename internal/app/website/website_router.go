package home

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"viabl.ventures/gossr/internal/templates"
)

type WebsiteRouter struct {
	renderer *templates.Renderer
}

func NewWebsiteRouter(renderer *templates.Renderer) *WebsiteRouter {
	return &WebsiteRouter{renderer}
}

func (router *WebsiteRouter) GetRoutes(r chi.Router) {
	r.Get("/", router.homePageView)
}

func (router *WebsiteRouter) homePageView(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title": "Home Page",
	}
	router.renderer.RenderTemplate(w, "home.html", data)
}
