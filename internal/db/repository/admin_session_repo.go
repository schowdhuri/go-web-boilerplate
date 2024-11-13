package repository

import (
	"time"

	"gorm.io/gorm"
	"viabl.ventures/gossr/internal/db/models"
)

type AdminSessionRepository struct {
	db *gorm.DB
}

func NewAdminSessionRepository(db *gorm.DB) *AdminSessionRepository {
	return &AdminSessionRepository{
		db: db,
	}
}

func (r *AdminSessionRepository) Create(sessionID string, adminUserID uint, lifetime time.Duration) (*models.AdminSession, error) {
	session := &models.AdminSession{
		SessionID:   sessionID,
		AdminUserID: adminUserID,
		ExpiresAt:   time.Now().Add(lifetime),
	}
	err := r.db.Create(session).Error
	return session, err
}

func (r *AdminSessionRepository) FindBySessionID(sessionID string) (*models.AdminSession, error) {
	var session models.AdminSession
	// query the database for the session and check if it hasnt expired
	err := r.db.Where("session_id = ? AND expires_at > ?", sessionID, time.Now()).First(&session).Error
	return &session, err
}

func (r *AdminSessionRepository) FindAllForUser(adminUserID uint) ([]*models.AdminSession, error) {
	var sessions []*models.AdminSession
	// query the database for all sessions for a user
	err := r.db.Where("admin_user_id = ?", adminUserID).Find(&sessions).Error
	return sessions, err
}

func (r *AdminSessionRepository) UpdateExpiry(sessionID string, lifetime time.Duration) error {
	// update the session expiration time
	return r.db.Model(&models.AdminSession{}).Where("session_id = ?", sessionID).Update("expires_at", time.Now().Add(lifetime)).Error
}

func (r *AdminSessionRepository) Delete(sessionID string) error {
	// delete the session
	return r.db.Where("session_id = ?", sessionID).Delete(&models.AdminSession{}).Error
}

func (r *AdminSessionRepository) DeleteForUser(adminUserID uint) error {
	// delete all sessions for a user
	return r.db.Where("admin_user_id = ?", adminUserID).Delete(&models.AdminSession{}).Error
}

func (r *AdminSessionRepository) DeleteBySessionID(sessionID string) error {
	// delete the session
	return r.db.Where("session_id = ?", sessionID).Delete(&models.AdminSession{}).Error
}

func (r *AdminSessionRepository) DeleteExpiredSessions() error {
	// delete all expired sessions
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.AdminSession{}).Error
}
