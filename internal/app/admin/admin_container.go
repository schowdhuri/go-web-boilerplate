package admin

import (
	"viabl.ventures/gossr/internal/app"
	"viabl.ventures/gossr/internal/db/repository"
	"viabl.ventures/gossr/internal/utils"
)

type AdminContainer struct {
	app.BaseContainer
	emailService  *utils.EmailService
	signinService *SigninService
	Router        *AdminRouter
}

func NewAdminContainer(bc *app.BaseContainer) *AdminContainer {
	c := &AdminContainer{
		BaseContainer: *bc,
		emailService:  utils.NewEmailService(bc.Config),
	}
	signinCodeRepo := repository.NewLoginCodeRepository(c.BaseContainer.DB)
	c.signinService = NewSigninService(signinCodeRepo)
	c.Router = NewAdminRouter(c.BaseContainer.Renderer, c.emailService, c.signinService)
	return c
}
