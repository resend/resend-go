package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func inlineAttachmentExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	if apiKey == "" {
		panic("Api Key is missing")
	}

	client := resend.NewClient(apiKey)

	// Create attachments objects
	attachment := &resend.Attachment{
		Path:      "https://resend.com/static/brand/resend-wordmark-black.png",
		Filename:  "resend-wordmark-black.png",
		ContentId: "my-test-image",
	}

	params := &resend.SendEmailRequest{
		To:          []string{"delivered@resend.dev"},
		From:        "onboarding@resend.dev",
		Text:        "email with inline content attachment",
		Html:        "<p>This is an email with an <img width=100 height=40 src=\"cid:my-test-image\" /> embed image</p>",
		Subject:     "Email with inline attachment",
		Attachments: []*resend.Attachment{attachment},
	}

	sent, err := client.Emails.SendWithContext(ctx, params)
	if err != nil {
		panic(err)
	}
	fmt.Println(sent.Id)
}
