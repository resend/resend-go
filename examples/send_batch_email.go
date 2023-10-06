package examples

import (
	"fmt"
	"os"

	"github.com/resendlabs/resend-go"
)

func sendEmailExample() {

	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Send
	var batchEmails = []*resend.SendEmailRequest{
		&resend.SendEmailRequest{
			To:      []string{"delivered@resend.dev"},
			Text:    "hey",
			Subject: "Hello from emails",
		},
		&resend.SendEmailRequest{
			From:    "onboarding@resend.dev",
			Text:    "hellooo",
			Subject: "Hello from batch emails 2",
		},
	}

	sent, err := client.Batch.Send(batchEmails)
	if err != nil {
		panic(err)
	}
	fmt.Println(sent.Data)
}
