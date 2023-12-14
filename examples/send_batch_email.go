package examples

import (
	"fmt"
	"os"

	"github.com/resendlabs/resend-go/v2"
)

func sendBatchEmails() {

	// ctx := context.TODO()
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

	sent, err := client.Batch.Send(batchEmails)
	// sent, err := client.Batch.SendWithContext(ctx, batchEmails)
	if err != nil {
		panic(err)
	}
	fmt.Println(sent.Data)
}
