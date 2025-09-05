package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func sendEmailExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Send params
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
	fmt.Printf("Sent basic email: %s\n", sent.Id)

	// Sending with IdempotencyKey
	options := &resend.SendEmailOptions{
		IdempotencyKey: "unique-idempotency-key",
	}

	sent, err = client.Emails.SendWithOptions(ctx, params, options)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Sent email with idempotency key: %s\n", sent.Id)

	// Get Email
	email, err := client.Emails.GetWithContext(ctx, sent.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", email)

	// List emails
	fmt.Println("\nListing recent emails:")
	listResp, err := client.Emails.ListWithContext(ctx, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d emails\n", len(listResp.Data))
	fmt.Printf("Has more emails: %v\n", listResp.HasMore)

	for i, email := range listResp.Data {
		if i < 5 { // Show first 5 emails
			fmt.Printf("  - ID: %s, Subject: %s, To: %v\n",
				email.Id, email.Subject, email.To)
		}
	}

	// List emails with pagination
	fmt.Println("\nListing emails with limit:")
	paginatedResp, err := client.Emails.ListWithContext(ctx, &resend.ListEmailsRequest{
		Limit: 3,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d emails (limited to 3)\n", len(paginatedResp.Data))

	// Example of cursor-based pagination
	if paginatedResp.HasMore && len(paginatedResp.Data) > 0 {
		lastEmailID := paginatedResp.Data[len(paginatedResp.Data)-1].Id
		fmt.Printf("\nFetching next page after email ID: %s\n", lastEmailID)

		nextPage, err := client.Emails.ListWithContext(ctx, &resend.ListEmailsRequest{
			Limit: 3,
			After: lastEmailID,
		})
		if err != nil {
			panic(err)
		}

		fmt.Printf("Found %d more emails in next page\n", len(nextPage.Data))
	}
}
