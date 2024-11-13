package models

import "time"

type AdminSession struct {
	Base
	SessionID   string    `gorm:"uniqueIndex not null" json:"session_id"`
	AdminUser   AdminUser `gorm:"foreignKey:AdminUserID" json:"-"`
	AdminUserID uint      `json:"admin_user_id"`
	ExpiresAt   time.Time `json:"expires_at"`
}
