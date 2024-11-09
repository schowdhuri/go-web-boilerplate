package handlers

import (
	"net/http"

	"viabl.ventures/gossr/internal/templates"
)

func CreateHomePageHandler(renderer *templates.Renderer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title": "Home Page",
		}
		renderer.RenderTemplate(w, "home.html", data)
	}
}
