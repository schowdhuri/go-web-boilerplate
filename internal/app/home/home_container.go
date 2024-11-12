package home

import (
	"viabl.ventures/gossr/internal/app"
)

type HomeContainer struct {
	app.BaseContainer
	Router *HomeRouter
}

func NewHomeContainer(bc *app.BaseContainer) *HomeContainer {
	c := &HomeContainer{
		BaseContainer: *bc,
	}
	c.Router = NewHomeRouter(c.BaseContainer.Renderer)
	return c
}
