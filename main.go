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

	// Send a single SMS
	singleResp, err := c.SMS.SendSingleSMS(ctx, "2557654321", "Hello from Notify Africa!", "137")
	if err != nil {
		log.Fatalf("Single SMS error: %v", err)
	}
	log.Printf("Single SMS sent! Status: %d, Message: %s, MessageID: %s", singleResp.Status, singleResp.Message, singleResp.Data.MessageID)

	// Send batch SMS
	batchResp, err := c.SMS.SendBatchSMS(ctx, []string{"255763765548", "255689737839"}, "Batch test message", "137")
	if err != nil {
		log.Fatalf("Batch SMS error: %v", err)
	}
	log.Printf("Batch SMS sent! Status: %d, Message: %s, Count: %d, Credits: %d, Balance: %d", batchResp.Status, batchResp.Message, batchResp.Data.MessageCount, batchResp.Data.CreditsDeducted, batchResp.Data.RemainingBalance)

	// Check message status (using the MessageID from singleResp)
	statusResp, err := c.SMS.CheckMessageStatus(ctx, singleResp.Data.MessageID)
	if err != nil {
		log.Fatalf("Check status error: %v", err)
	}
	log.Printf("Message Status: %s, DeliveredAt: %v", statusResp.Data.Status, statusResp.Data.DeliveredAt)

	// Send Email (unchanged)
	emailResp, err := c.Email.SendEmailWithContext(ctx, "noreply@yourdomain.com", "Subject", "Body text", []string{"user@example.com"})
	if err != nil {
		log.Fatalf("Email error: %v", err)
	}
	log.Printf("Email sent! Message: %s, Success: %v", emailResp.Message, emailResp.Success)
}
