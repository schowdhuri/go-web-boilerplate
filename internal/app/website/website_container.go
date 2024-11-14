package home

import (
	"viabl.ventures/gossr/internal/app"
)

type WebsiteContainer struct {
	app.BaseContainer
	Router *WebsiteRouter
}

func NewWebsiteContainer(bc *app.BaseContainer) *WebsiteContainer {
	c := &WebsiteContainer{
		BaseContainer: *bc,
	}
	c.Router = NewWebsiteRouter(c.BaseContainer.Renderer)
	return c
}
