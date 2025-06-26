package client

import (
	"github.com/tacherasasi/notify-africa-go/email"
	"github.com/tacherasasi/notify-africa-go/sms"
)

type Client struct {
	SMS   *sms.Client
	Email *email.Client
}

type Config struct {
	BaseURL     string
	SMSApiKey   string
	EmailApiKey string
}

func NewClient(cfg Config) *Client {
	smsClient := sms.NewClient(cfg.SMSApiKey)
	emailClient := email.NewClient(cfg.EmailApiKey)
	if cfg.BaseURL != "" {
		smsClient.SetBaseURL(cfg.BaseURL)
		emailClient.SetBaseURL(cfg.BaseURL)
	}
	return &Client{
		SMS:   smsClient,
		Email: emailClient,
	}
}
