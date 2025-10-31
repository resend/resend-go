package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

func sendEmailWithTemplateExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	fmt.Println("Creating a template with variables:")
	templateParams := &resend.CreateTemplateRequest{
		Name:    "user-welcome-template",
		Alias:   "welcome",
		From:    "onboarding@resend.dev",
		Subject: "Welcome to {{{companyName}}}, {{{userName}}}!",
		Html: `
		<html>
			<body>
				<h1>Hello {{{userName}}}!</h1>
				<p>Welcome to {{{companyName}}}. We're excited to have you on board.</p>
				<p>You currently have {{{messageCount}}} unread messages waiting for you.</p>
			</body>
		</html>
		`,
		Text: "Hello {{{userName}}}! Welcome to {{{companyName}}}.",
		Variables: []*resend.TemplateVariable{
			{
				Key:           "userName",
				Type:          resend.VariableTypeString,
				FallbackValue: "User",
			},
			{
				Key:           "companyName",
				Type:          resend.VariableTypeString,
				FallbackValue: "Our Company",
			},
			{
				Key:           "messageCount",
				Type:          resend.VariableTypeNumber,
				FallbackValue: 0,
			},
		},
	}

	template, err := client.Templates.CreateWithContext(ctx, templateParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✓ Created template: %s\n", template.Id)

	fmt.Println("\nPublishing template:")
	_, err = client.Templates.PublishWithContext(ctx, template.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✓ Published template: %s\n", template.Id)

	fmt.Println("\nSending email using template ID:")
	emailParams := &resend.SendEmailRequest{
		To: []string{"delivered@resend.dev"},
		Template: &resend.EmailTemplate{
			Id: template.Id,
			Variables: map[string]interface{}{
				"userName":     "Alice Johnson",
				"companyName":  "Acme Corporation",
				"messageCount": 12,
			},
		},
	}

	sent, err := client.Emails.SendWithContext(ctx, emailParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✓ Sent email: %s\n", sent.Id)

	fmt.Println("\nSending email using template alias:")
	emailParams2 := &resend.SendEmailRequest{
		To: []string{"delivered@resend.dev"},
		Template: &resend.EmailTemplate{
			Id: "welcome", // Using alias instead of ID
			Variables: map[string]interface{}{
				"userName":     "Bob Smith",
				"companyName":  "Tech Startup Inc",
				"messageCount": 3,
			},
		},
	}

	sent2, err := client.Emails.SendWithContext(ctx, emailParams2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✓ Sent email: %s\n", sent2.Id)

	fmt.Println("\nSending email with template and field overrides:")
	emailParams3 := &resend.SendEmailRequest{
		From:    "support@resend.dev", // Override template's From
		To:      []string{"delivered@resend.dev"},
		Subject: "Custom Subject Override", // Override template's Subject
		Bcc:     []string{"bcc@example.com"},
		ReplyTo: "noreply@resend.dev",
		Template: &resend.EmailTemplate{
			Id: template.Id,
			Variables: map[string]interface{}{
				"userName":     "Charlie Brown",
				"companyName":  "Example LLC",
				"messageCount": 7,
			},
		},
	}

	sent3, err := client.Emails.SendWithContext(ctx, emailParams3)
	if err != nil {
		panic(err)
	}
	fmt.Printf("✓ Sent email with overrides: %s\n", sent3.Id)

	fmt.Println("\nCleaning up template:")
	removeResp, err := client.Templates.RemoveWithContext(ctx, template.Id)
	if err != nil {
		fmt.Printf("Warning: Could not remove template %s: %v\n", template.Id, err)
	} else {
		fmt.Printf("✓ Removed template: %s\n", removeResp.Id)
	}

	fmt.Println("\nExample completed successfully!")
}
