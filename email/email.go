package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	Client struct {
		BaseURL string
		Token   string
		client  *http.Client
	}
)

func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		client:  &http.Client{},
	}
}

// SendEmailWithContext sends an email using notify.africa with context
func (c *Client) SendEmailWithContext(ctx context.Context, sender, subject, body string, recipients []string) (*Respond, error) {
	url := c.BaseURL + "/send-email"

	payload := EmailPayload{
		Sender:     sender,
		Subject:    subject,
		Body:       body,
		Recipients: recipients,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(respBody))
	}

	var response Respond
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("error parsing response JSON: %w", err)
	}

	return &response, nil
}

// EmailEndpoint is a backward-compatible function that uses context.Background()
func (c *Client) SendEmail(sender, subject, body string, recipients []string) (*Respond, error) {
	return c.SendEmailWithContext(context.Background(), sender, subject, body, recipients)
}

// For backward compatibility, keep the old EmailEndpoint function, but recommend using the client.
func EmailEndpoint(sender, subject, body string, recipients []string) error {
	client := NewClient("https://api.notify.africa/v2", "") // Token should be set via env or config
	_, err := client.SendEmail(sender, subject, body, recipients)
	return err
}
