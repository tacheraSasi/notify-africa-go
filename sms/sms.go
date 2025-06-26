package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultBaseURL = "https://api.notify.africa/v2"

// Client represents the SMS API client
// Now supports context for production readiness
type Client struct {
	BaseURL string
	Token   string
	client  *http.Client
}

// NewClient creates a new API client
func NewClient(token string) *Client {
	return &Client{
		BaseURL: defaultBaseURL,
		Token:   token,
		client:  &http.Client{},
	}
}

// SetBaseURL allows overriding the default BaseURL
func (c *Client) SetBaseURL(baseURL string) {
	if baseURL != "" {
		c.BaseURL = baseURL
	}
}

// Recipient represents a single SMS recipient
type Recipient struct {
	Number string `json:"number"`
}

// Payload represents the SMS request structure
type Payload struct {
	SenderID   int         `json:"sender_id"`
	Schedule   string      `json:"schedule"`
	SMS        string      `json:"sms"`
	Recipients []Recipient `json:"recipients"`
}

// SendResponse represents the API response
type SendResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// SendSMSWithContext sends an SMS message to multiple recipients with context
func (c *Client) SendSMSWithContext(ctx context.Context, senderID int, message string, recipients []string) (*SendResponse, error) {
	recipientList := make([]Recipient, len(recipients))
	for i, num := range recipients {
		recipientList[i] = Recipient{Number: num}
	}

	payload := Payload{
		SenderID:   senderID,
		Schedule:   "none",
		SMS:        message,
		Recipients: recipientList,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}

	url := fmt.Sprintf("%s/send-sms", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	var response SendResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing response JSON: %w", err)
	}

	return &response, nil
}

// SendSMS is a backward-compatible method that uses context.Background()
func (c *Client) SendSMS(senderID int, message string, recipients []string) (*SendResponse, error) {
	return c.SendSMSWithContext(context.Background(), senderID, message, recipients)
}
