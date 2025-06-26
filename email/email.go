package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type (
	Respond struct {
		Message string `json:"message"`
		Success bool   `json:"success"`
	}

	EmailPayload struct {
		Sender     string   `json:"sender"`
		Subject    string   `json:"subject"`
		Body       string   `json:"body"`
		Recipients []string `json:"recipients"`
	}
)

// EmailEndpoint sends an email using notify.africa
func EmailEndpoint(sender, subject, body string, recipients []string) error {
	url := "https://api.notify.africa/v2/send-email"

	apiKey := os.Getenv("EMAIL_APIKEY")
	if apiKey == "" {
		return fmt.Errorf("EMAIL_APIKEY environment variable is not set")
	}

	payload := EmailPayload{
		Sender:     sender,
		Subject:    subject,
		Body:       body,
		Recipients: recipients,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Printf("Status: %s\nResponse: %s\n", resp.Status, string(respBody))

	return nil
}
