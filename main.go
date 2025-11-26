package main

import (
	"context"
	"log"
	"os"
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

	// Send test SMS to 255686477074 using batch endpoint
	smsURL := os.Getenv("SMS_URL")
	smsSenderID := os.Getenv("SMS_SENDER_ID")
	if smsURL != "" {
		c.SMS.SetBaseURL(smsURL[:len(smsURL)-len("/api/v1/api/messages/batch")]) // Remove endpoint, keep base
	}
	batchResp, err := c.SMS.SendBatchSMS(ctx, []string{"255686477074"}, "Test SMS from Notify Africa Go!", smsSenderID)
	if err != nil {
		log.Fatalf("Test SMS error: %v", err)
	}
	log.Printf("Test SMS sent! Status: %d, Message: %s, Count: %d, Credits: %d, Balance: %d", batchResp.Status, batchResp.Message, batchResp.Data.MessageCount, batchResp.Data.CreditsDeducted, batchResp.Data.RemainingBalance)

	// Send Email (unchanged)
	emailResp, err := c.Email.SendEmailWithContext(ctx, "noreply@yourdomain.com", "Subject", "Body text", []string{"user@example.com"})
	if err != nil {
		log.Fatalf("Email error: %v", err)
	}
	log.Printf("Email sent! Message: %s, Success: %v", emailResp.Message, emailResp.Success)
}
