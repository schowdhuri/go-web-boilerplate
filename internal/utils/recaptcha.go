package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"viabl.ventures/gossr/internal/config"
)

const (
	recaptchaThreshold = 0.75
	recaptchaVerifyURL = "https://www.google.com/recaptcha/api/siteverify"
)

type recaptchaService struct {
	config *config.EnvVars
}

func NewRecaptchaService(config *config.EnvVars) *recaptchaService {
	return &recaptchaService{config}
}

func (s *recaptchaService) VerifyRecaptcha(token string, version int) (bool, error) {
	params := url.Values{}
	params.Set("secret", s.config.RecaptchaSecret)
	params.Set("response", token)

	resp, err := http.PostForm(recaptchaVerifyURL, params)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var result struct {
		Success bool    `json:"success"`
		Score   float64 `json:"score"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return false, err
	}

	if !result.Success {
		return false, fmt.Errorf("reCAPTCHA verification failed: %v", result)
	}

	if version == 3 && result.Score < recaptchaThreshold {
		log.Default().Printf("reCAPTCHA score too low: %v", result)
		return false, nil
	}

	return true, nil
}
