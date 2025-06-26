package main

import (
	"log"
	"os"
	"context"
	"time"

	"github.com/tacherasasi/notify-africa-go/client"
)

func main() {
	cfg := client.Config{
		SMSApiKey:   os.Getenv("SMS_APIKEY"),
		EmailApiKey: os.Getenv("EMAIL_APIKEY"),
		// Optionally override BaseURL:
		// BaseURL: "https://custom.url/v2",
	}
	c := client.NewClient(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Send SMS
	smsResp, err := c.SMS.SendSMSWithContext(ctx, 1, "Hello from Notify Africa!", []string{"2557654321"})
	if err != nil {
		log.Fatalf("SMS error: %v", err)
	}
	log.Printf("SMS sent! Status: %d, Message: %s", smsResp.Status, smsResp.Message)

	// Send Email
	emailResp, err := c.Email.SendEmailWithContext(ctx, "noreply@yourdomain.com", "Subject", "Body text", []string{"user@example.com"})
	if err != nil {
		log.Fatalf("Email error: %v", err)
	}
	log.Printf("Email sent! Message: %s, Success: %v", emailResp.Message, emailResp.Success)
}