package admin

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"viabl.ventures/gossr/internal/db/models"
	"viabl.ventures/gossr/internal/db/repository"
)

var SESSION_LIFETIME = 1 * time.Hour

type SessionService struct {
	sessionRepo *repository.AdminSessionRepository
}

func NewSessionService(sessionRepo *repository.AdminSessionRepository) *SessionService {
	return &SessionService{sessionRepo}
}

func (s *SessionService) CreateSessionCookie(userID uint) (*http.Cookie, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, err
	}

	// Create session in database
	s.sessionRepo.Create(sessionID, userID, SESSION_LIFETIME)

	cookie := http.Cookie{
		Name:     "session",
		Value:    sessionID,
		HttpOnly: true, // Prevent JavaScript access
		Secure:   true, // Only send over HTTPS
		Path:     "/",
		Expires:  time.Now().Add(SESSION_LIFETIME),
	}
	return &cookie, nil
}

func (s *SessionService) ValidateSession(sessionID string) (*models.AdminSession, error) {
	session, err := s.sessionRepo.FindBySessionID(sessionID)
	if err != nil {
		return nil, errors.New("invalid session")
	}
	// REVIEW: any more session validation checks?

	return session, nil
}

func (s *SessionService) RenewSession(session *models.AdminSession) (*http.Cookie, error) {
	if time.Until(session.ExpiresAt) <= 0 {
		return nil, errors.New("session has expired")
	}

	// Update session last access time if remaining session lifetime is less than 15 minutes
	if time.Until(session.ExpiresAt) < 15*time.Minute {
		s.sessionRepo.UpdateExpiry(session.SessionID, SESSION_LIFETIME)

		// create a new session cookie with updated expiry
		cookie := http.Cookie{
			Name:     "session",
			Value:    session.SessionID,
			HttpOnly: true, // Prevent JavaScript access
			Secure:   true, // Only send over HTTPS
			Expires:  time.Now().Add(SESSION_LIFETIME),
		}
		return &cookie, nil
	}

	// no new cookie, no errors
	return nil, nil
}

func (s *SessionService) DeleteSession(sessionID string) error {
	return s.sessionRepo.Delete(sessionID)
}

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
