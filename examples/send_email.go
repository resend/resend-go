package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resendlabs/resend-go"
)

func sendEmailExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Send
	params := &resend.SendEmailRequest{
		To:      []string{"delivered@resend.dev"},
		From:    "onboarding@resend.dev",
		Text:    "hello world",
		Subject: "Hello from Golang",
		Cc:      []string{"cc@example.com"},
		Bcc:     []string{"ccc@example.com"},
		ReplyTo: "to@example.com",
	}

	sent, err := client.Emails.SendWithContext(ctx, params)
	if err != nil {
		panic(err)
	}
	fmt.Println(sent.Id)

	// Get
	email, err := client.Emails.GetWithContext(ctx, sent.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", email)

}
