package repository

import (
	"gorm.io/gorm"
	"viabl.ventures/gossr/internal/db/models"
)

type AdminUserRepository struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) *AdminUserRepository {
	return &AdminUserRepository{db}
}

func (r *AdminUserRepository) FindByEmail(email string) (*models.AdminUser, error) {
	var user models.AdminUser
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *AdminUserRepository) FindByID(id uint) (*models.AdminUser, error) {
	var user models.AdminUser
	err := r.db.First(&user, id).Error
	return &user, err
}
