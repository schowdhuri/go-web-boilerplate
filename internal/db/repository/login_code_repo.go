package repository

import (
	"time"

	"gorm.io/gorm"
	"viabl.ventures/gossr/internal/db/models"
)

type LoginCodeRepository struct {
	db *gorm.DB
}

var lifetime = time.Minute * 5

func NewLoginCodeRepository(db *gorm.DB) *LoginCodeRepository {
	return &LoginCodeRepository{
		db: db,
	}
}

func (r *LoginCodeRepository) Create(code string, adminUserID uint) (*models.LoginCode, error) {
	loginCode := &models.LoginCode{
		Code:        code,
		AdminUserID: adminUserID,
		ExpiresAt:   time.Now().Add(lifetime),
	}
	err := r.db.Create(loginCode).Error
	return loginCode, err
}

func (r *LoginCodeRepository) FindByCode(email, code string) (*models.LoginCode, error) {
	var loginCode models.LoginCode
	// query the database for the login code and check if it hasnt expired
	err := r.db.Where("email = ? AND code = ? AND expires_at > ?", email, code, time.Now()).First(&loginCode).Error
	return &loginCode, err
}

func (r *LoginCodeRepository) Delete(id uint) error {
	// delete the login code that has been used
	return r.db.Delete(&models.LoginCode{}, id).Error
}

func (r *LoginCodeRepository) DeleteExpired(duration string) error {
	// delete all login codes that have expired
	return r.db.Where("expires_at < ?", time.Now().Add(-lifetime)).Delete(&models.LoginCode{}).Error
}
