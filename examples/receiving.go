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
	email, err := client.Emails.Receiving.GetWithContext(ctx, "8136d3fb-0439-4b09-b939-b8436a3524b6")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Retrieved received email: %s\n", email.Subject)
	fmt.Printf("From: %s\n", email.From)
	fmt.Printf("To: %v\n", email.To)
	fmt.Printf("Has %d attachments\n", len(email.Attachments))

	// List received emails
	emails, err := client.Emails.Receiving.ListWithContext(ctx)
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
	paginatedEmails, err := client.Emails.Receiving.ListWithOptions(ctx, listOptions)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nPaginated list returned %d emails\n", len(paginatedEmails.Data))

	// Get email with attachments
	emailWithAttachments, err := client.Emails.Receiving.GetWithContext(ctx, "006e2796-ff6a-4436-91ad-0429e600bf8a")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nEmail '%s' has %d attachments\n", emailWithAttachments.Subject, len(emailWithAttachments.Attachments))

	// Get each attachment's full details including download URLs
	for _, att := range emailWithAttachments.Attachments {
		attachment, err := client.Emails.Receiving.GetAttachmentWithContext(ctx, emailWithAttachments.Id, att.Id)
		if err != nil {
			panic(err)
		}
		fmt.Printf("\nAttachment #%s:\n", att.Id)
		fmt.Printf("  Filename: %s\n", attachment.Filename)
		fmt.Printf("  Content Type: %s\n", attachment.ContentType)
		fmt.Printf("  Download URL: %s\n", attachment.DownloadUrl)
		fmt.Printf("  Expires At: %s\n", attachment.ExpiresAt)
	}

	// List all attachments for a received email
	attachmentsList, err := client.Emails.Receiving.ListAttachmentsWithContext(ctx, "006e2796-ff6a-4436-91ad-0429e600bf8a")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nEmail has %d attachments\n", len(attachmentsList.Data))
	fmt.Printf("Has more: %v\n", attachmentsList.HasMore)

	// List attachments with pagination
	attachmentsLimit := 5
	attachmentsOptions := &resend.ListOptions{
		Limit: &attachmentsLimit,
	}
	paginatedAttachments, err := client.Emails.Receiving.ListAttachmentsWithOptions(ctx, "006e2796-ff6a-4436-91ad-0429e600bf8a", attachmentsOptions)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nPaginated attachments list returned %d attachments\n", len(paginatedAttachments.Data))
	for _, att := range paginatedAttachments.Data {
		fmt.Printf("  - %s (%s)\n", att.Filename, att.ContentType)
	}
}
