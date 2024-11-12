package models

import "time"

type LoginCode struct {
	Base
	Code        string    `gorm:"not null" json:"code"`
	AdminUser   AdminUser `gorm:"foreignKey:AdminUserID" json:"-"`
	AdminUserID uint      `json:"admin_user_id"`
	ExpiresAt   time.Time `json:"expires_at"`
}
