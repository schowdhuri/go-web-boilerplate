package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"viabl.ventures/gossr/internal/config"
)

type EmailService struct {
	config *config.EnvVars
}

func NewEmailService(config *config.EnvVars) *EmailService {
	return &EmailService{config: config}
}

func (s *EmailService) SendMail(to, subject, htmlContent string) error {
	apiKey := s.config.BrevoAPIKey
	senderEmail := s.config.BrevoSender

	payload := map[string]interface{}{
		"sender":      map[string]string{"email": senderEmail, "name": "Contact Form"},
		"to":          []map[string]string{{"email": to}},
		"subject":     subject,
		"htmlContent": htmlContent,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.sendinblue.com/v3/smtp/email",
		bytes.NewBuffer(jsonPayload),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("API request failed: %s, body: %s", resp.Status, string(body))
	}

	return nil
}
