package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client represents the SMS API client
type Client struct {
	BaseURL string
	Token   string
}

// NewClient creates a new API client
func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
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

// SendSMS sends an SMS message to multiple recipients
func (c *Client) SendSMS(senderID int, message string, recipients []string) (*SendResponse, error) {
	// Prepare recipients list
	recipientList := make([]Recipient, len(recipients))
	for i, num := range recipients {
		recipientList[i] = Recipient{Number: num}
	}

	// Create payload
	payload := Payload{
		SenderID:   senderID,
		Schedule:   "none",
		SMS:        message,
		Recipients: recipientList,
	}

	// Marshal payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}

	// Create request
	url := fmt.Sprintf("%s/send-sms", c.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Check for non-success status
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var response SendResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing response JSON: %w", err)
	}

	return &response, nil
}