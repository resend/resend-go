package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func receivingExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Get a single received email
	email, err := client.Receiving.GetWithContext(ctx, "8136d3fb-0439-4b09-b939-b8436a3524b6")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Retrieved received email: %s\n", email.Subject)
	fmt.Printf("From: %s\n", email.From)
	fmt.Printf("To: %v\n", email.To)
	fmt.Printf("Has %d attachments\n", len(email.Attachments))

	// List received emails
	emails, err := client.Receiving.ListWithContext(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nYou have %d received emails\n", len(emails.Data))
	fmt.Printf("Has more: %v\n", emails.HasMore)

	// List with pagination
	limit := 10
	listOptions := &resend.ListOptions{
		Limit: &limit,
	}
	paginatedEmails, err := client.Receiving.ListWithOptions(ctx, listOptions)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nPaginated list returned %d emails\n", len(paginatedEmails.Data))
}
