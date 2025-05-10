package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func sendBatchEmails() {

	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Batch Send
	var batchEmails = []*resend.SendEmailRequest{
		{
			To:      []string{"delivered@resend.dev"},
			From:    "onboarding@resend.dev",
			Text:    "hey",
			Subject: "Hello from emails",
		},
		{
			To:      []string{"delivered@resend.dev"},
			From:    "onboarding@resend.dev",
			Text:    "hellooo",
			Subject: "Hello from batch emails 2",
		},
	}

	// Regular send without options
	sent, err := client.Batch.SendWithContext(ctx, batchEmails)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sent without options")
	fmt.Println(sent.Data)

	// Send with options: IdempotencyKey
	options := &resend.BatchSendEmailOptions{
		IdempotencyKey: "68656c6c6f2d776f726c64",
	}
	sent, err = client.Batch.SendWithOptions(ctx, batchEmails, options)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sent with idempotency key")
	fmt.Println(sent.Data)
}
