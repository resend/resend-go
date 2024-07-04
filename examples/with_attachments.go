package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func withAttachments() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	if apiKey == "" {
		panic("Api Key is missing")
	}

	// Read attachment file
	pwd, _ := os.Getwd()
	f, err := os.ReadFile(pwd + "/resources/invoice.pdf")
	if err != nil {
		panic(err)
	}

	client := resend.NewClient(apiKey)

	// Create attachments objects
	pdfAttachmentFromLocalFile := &resend.Attachment{
		Content:     f,
		Filename:    "invoice1.pdf",
		ContentType: "application/pdf",
	}

	pdfAttachmentFromRemotePath := &resend.Attachment{
		Path:        "https://github.com/resend/resend-go/raw/main/resources/invoice.pdf",
		Filename:    "invoice2.pdf",
		ContentType: "application/pdf",
	}

	params := &resend.SendEmailRequest{
		To:          []string{"delivered@resend.dev"},
		From:        "onboarding@resend.dev",
		Text:        "email with attachments !!",
		Html:        "<strong>email with attachments !!</strong>",
		Subject:     "Email with attachment",
		Attachments: []*resend.Attachment{pdfAttachmentFromLocalFile, pdfAttachmentFromRemotePath},
	}

	sent, err := client.Emails.SendWithContext(ctx, params)
	if err != nil {
		panic(err)
	}
	fmt.Println(sent.Id)
}
