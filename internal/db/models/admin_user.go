package models

type AdminUser struct {
	Base
	Email string `gorm:"uniqueIndex;not null" json:"email"`
}
