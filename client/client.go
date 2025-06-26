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
	return &Client{
		SMS:   sms.NewClient(cfg.BaseURL, cfg.SMSApiKey),
		Email: email.NewClient(cfg.BaseURL, cfg.EmailApiKey),
	}
}