package admin

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"viabl.ventures/gossr/internal/db/models"
	"viabl.ventures/gossr/internal/db/repository"
)

type AuthService struct {
	adminUserRepo  *repository.AdminUserRepository
	loginCodeRepo  *repository.LoginCodeRepository
	sessionService *SessionService
}

func NewSigninService(adminUserRepo *repository.AdminUserRepository, loginCodeRepo *repository.LoginCodeRepository, sessionService *SessionService) *AuthService {
	return &AuthService{adminUserRepo, loginCodeRepo, sessionService}
}

func (s *AuthService) GenerateCode(email string) (*models.LoginCode, error) {
	// Find the admin user by email
	adminUser, err := s.adminUserRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	// Generate random code
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	code := hex.EncodeToString(bytes)

	loginCode, err := s.loginCodeRepo.Create(
		code,
		adminUser.ID,
	)
	if err != nil {
		return nil, err
	}

	return loginCode, nil
}

func (s *AuthService) ValidateAndUseCode(email, code string) (*models.LoginCode, error) {
	adminUser, err := s.adminUserRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email")
	}

	loginCode, err := s.loginCodeRepo.FindByCode(code, adminUser.ID)
	if err != nil {
		return nil, err
	}

	if loginCode == nil {
		return nil, errors.New("invalid code")
	}

	// Delete the login code after it has been used
	err = s.loginCodeRepo.Delete(loginCode.ID)
	if err != nil {
		return nil, err
	}

	return loginCode, nil
}

func (s *AuthService) CreateSessionCookie(email string) (*http.Cookie, error) {
	// Find the admin user by email
	adminUser, err := s.adminUserRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email")
	}
	return s.sessionService.CreateSessionCookie(adminUser.ID)
}

func (s *AuthService) ValidateSession(sessionID string) (*http.Cookie, error) {
	session, err := s.sessionService.ValidateSession(sessionID)
	if err != nil {
		return nil, errors.New("invalid session")
	}

	// Update session expiration time if remaining session lifetime is less than 15 minutes
	if time.Until(session.ExpiresAt) < 15*time.Minute {
		return s.sessionService.RenewSession(session)
	}
	return nil, nil
}

func (s *AuthService) DeleteSession(sessionID string) error {
	return s.sessionService.DeleteSession(sessionID)
}

func (s *AuthService) CleanupExpiredCodes() error {
	// Delete codes older than 24 hours
	return s.loginCodeRepo.DeleteExpired("24 hours")
}
