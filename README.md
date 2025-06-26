# notify-africa-go

Golang client for [notify.africa](https://notify.africa) CPaaS (SMS & Email).

## Installation

```
go get github.com/tacherasasi/notify-africa-go
```

## Environment Variables

Set your API keys as environment variables:

```
export SMS_APIKEY=your_sms_api_key
export EMAIL_APIKEY=your_email_api_key
```

## Usage

### Unified Client Example

```go
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
```

### Direct Package Usage

```go
import (
	"github.com/tacherasasi/notify-africa-go/sms"
	"github.com/tacherasasi/notify-africa-go/email"
)

smsClient := sms.NewClient(os.Getenv("SMS_APIKEY"))
smsClient.SetBaseURL("https://custom.url/v2") // Optional
// ...

emailClient := email.NewClient(os.Getenv("EMAIL_APIKEY"))
emailClient.SetBaseURL("https://custom.url/v2") // Optional
// ...
```

## License
MIT
