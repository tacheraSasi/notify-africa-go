package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultBaseURL = "https://api.notify.africa"

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
type SendSinglePayload struct {
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
	SenderID    string `json:"sender_id"`
}

// SendBatchPayload represents the payload for sending batch SMS
type SendBatchPayload struct {
	PhoneNumbers []string `json:"phone_numbers"`
	Message      string   `json:"message"`
	SenderID     string   `json:"sender_id"`
}

// SendSingleResponse represents the response from sending a single SMS
type SendSingleResponse struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Path      string `json:"path"`
	Data      struct {
		MessageID string `json:"messageId"`
		Status    string `json:"status"`
	} `json:"data"`
}

// SendBatchResponse represents the response from sending batch SMS
type SendBatchResponse struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Path      string `json:"path"`
	Data      struct {
		MessageCount     int `json:"messageCount"`
		CreditsDeducted  int `json:"creditsDeducted"`
		RemainingBalance int `json:"remainingBalance"`
	} `json:"data"`
}


// MessageStatusResponse represents the response from checking message status
type MessageStatusResponse struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Path      string `json:"path"`
	Data      struct {
		MessageID   string  `json:"messageId"`
		Status      string  `json:"status"`
		SentAt      *string `json:"sentAt"`
		DeliveredAt *string `json:"deliveredAt"`
	} `json:"data"`
}

// SendSingleSMS sends a single SMS message
func (c *Client) SendSingleSMS(ctx context.Context, phoneNumber, message, senderID string) (*SendSingleResponse, error) {
	payload := SendSinglePayload{
		PhoneNumber: phoneNumber,
		Message:     message,
		SenderID:    senderID,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}
	url := fmt.Sprintf("%s/api/v1/api/messages/send", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
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
	var response SendSingleResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing response JSON: %w", err)
	}
	return &response, nil
}

// SendBatchSMS sends batch SMS messages
func (c *Client) SendBatchSMS(ctx context.Context, phoneNumbers []string, message, senderID string) (*SendBatchResponse, error) {
	payload := SendBatchPayload{
		PhoneNumbers: phoneNumbers,
		Message:      message,
		SenderID:     senderID,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}
	url := fmt.Sprintf("%s/api/v1/api/messages/batch", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
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
	var response SendBatchResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing response JSON: %w", err)
	}
	return &response, nil
}

// CheckMessageStatus checks the status of a sent message
func (c *Client) CheckMessageStatus(ctx context.Context, messageID string) (*MessageStatusResponse, error) {
	url := fmt.Sprintf("%s/api/v1/api/messages/status/%s", c.BaseURL, messageID)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
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
	var response MessageStatusResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing response JSON: %w", err)
	}
	return &response, nil
}
