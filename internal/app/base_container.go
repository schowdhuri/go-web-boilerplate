package app

import (
	"gorm.io/gorm"
	"viabl.ventures/gossr/internal/config"
	database "viabl.ventures/gossr/internal/db"
	"viabl.ventures/gossr/internal/templates"
)

type BaseContainer struct {
	Config   *config.EnvVars
	DB       *gorm.DB
	Renderer *templates.Renderer
}

func NewBaseContainer(conf *config.EnvVars, renderer *templates.Renderer) *BaseContainer {
	c := &BaseContainer{
		Config:   conf,
		DB:       database.InitDB(conf),
		Renderer: renderer,
	}
	return c
}
