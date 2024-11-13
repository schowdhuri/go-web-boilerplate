package admin

import (
	"viabl.ventures/gossr/internal/app"
	"viabl.ventures/gossr/internal/db/repository"
	"viabl.ventures/gossr/internal/utils"
)

type AdminContainer struct {
	app.BaseContainer
	emailService   *utils.EmailService
	authService    *AuthService
	sessionService *SessionService
	Router         *AdminRouter
}

func NewAdminContainer(bc *app.BaseContainer) *AdminContainer {
	c := &AdminContainer{
		BaseContainer: *bc,
		emailService:  utils.NewEmailService(bc.Config),
	}
	signinCodeRepo := repository.NewLoginCodeRepository(c.BaseContainer.DB)
	adminRepo := repository.NewAdminUserRepository(c.BaseContainer.DB)
	sessionRepo := repository.NewAdminSessionRepository(c.BaseContainer.DB)
	c.sessionService = NewSessionService(sessionRepo)
	c.authService = NewSigninService(adminRepo, signinCodeRepo, c.sessionService)
	c.Router = NewAdminRouter(c.BaseContainer.Config, c.BaseContainer.Renderer, c.emailService, c.authService)
	return c
}
