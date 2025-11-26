# notify-africa-go

Golang client for [notify.africa](https://notify.africa) CPaaS (SMS & Email).

## Installation

```
go get github.com/tacherasasi/notify-africa-go
```

## Environment Variables

Set your API keys and sender ID as environment variables:

```
export SMS_APIKEY=your_sms_api_key
export SMS_SENDER_ID=your_sender_id
export SMS_URL=https://api.notify.africa/api/v1/api/messages/batch
export EMAIL_APIKEY=your_email_api_key
```

## Usage

### Send Test SMS Example

```go
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
	}
	c := client.NewClient(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Optionally override SMS base URL from env
	smsURL := os.Getenv("SMS_URL")
	if smsURL != "" {
		c.SMS.SetBaseURL(smsURL[:len(smsURL)-len("/api/v1/api/messages/batch")])
	}

	// Send test SMS to 255686477074
	senderID := os.Getenv("SMS_SENDER_ID")
	batchResp, err := c.SMS.SendBatchSMS(ctx, []string{"255686477074"}, "Test SMS from Notify Africa Go!", senderID)
	if err != nil {
		log.Fatalf("Test SMS error: %v", err)
	}
	log.Printf("Test SMS sent! Status: %d, Message: %s, Count: %d, Credits: %d, Balance: %d", batchResp.Status, batchResp.Message, batchResp.Data.MessageCount, batchResp.Data.CreditsDeducted, batchResp.Data.RemainingBalance)
}
```

### Direct Package Usage

```go
import (
	"github.com/tacherasasi/notify-africa-go/sms"
	"github.com/tacherasasi/notify-africa-go/email"
)

smsClient := sms.NewClient(os.Getenv("SMS_APIKEY"))
smsClient.SetBaseURL("https://api.notify.africa") // Optional
// ...

emailClient := email.NewClient(os.Getenv("EMAIL_APIKEY"))
emailClient.SetBaseURL("https://api.notify.africa/v2") // Optional
// ...
```

## License

MIT
