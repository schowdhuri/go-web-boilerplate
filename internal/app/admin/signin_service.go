package admin

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"viabl.ventures/gossr/internal/db/models"
	"viabl.ventures/gossr/internal/db/repository"
)

type SigninService struct {
	repo *repository.LoginCodeRepository
}

func NewSigninService(repo *repository.LoginCodeRepository) *SigninService {
	return &SigninService{
		repo: repo,
	}
}

func (s *SigninService) GenerateCode(email string) (*models.LoginCode, error) {
	// Find the admin user by email
	adminRepo := repository.NewAdminUserRepository()
	adminUser, err := adminRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	// Generate random code
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	code := hex.EncodeToString(bytes)

	loginCode, err := s.repo.Create(
		code,
		adminUser.ID,
	)
	if err != nil {
		return nil, err
	}

	return loginCode, nil
}

func (s *SigninService) ValidateAndUseCode(email, code string) (*models.LoginCode, error) {
	loginCode, err := s.repo.FindByCode(email, code)
	if err != nil {
		return nil, err
	}

	if loginCode == nil {
		return nil, errors.New("invalid code")
	}

	err = s.repo.Delete(loginCode.ID)
	if err != nil {
		return nil, err
	}

	return loginCode, nil
}

func (s *SigninService) CleanupExpiredCodes() error {
	// Delete codes older than 24 hours
	return s.repo.DeleteExpired("24 hours")
}
