package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func templatesExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Create a simple template without variables
	template, err := client.Templates.CreateWithContext(ctx, &resend.CreateTemplateRequest{
		Name: "welcome-email",
		Html: "<strong>Welcome to our service!</strong>",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created template: %s\n", template.Id)

	// Create a template with variables
	// IMPORTANT: All variables used in the HTML (e.g., {{{NAME}}}, {{{AGE}}})
	// MUST be declared in the Variables array, or the API will return an error:
	// "Variable 'NAME' is used in the template but not defined in the variables list"
	templateWithVars, err := client.Templates.Create(&resend.CreateTemplateRequest{
		Name:    "user-notification",
		From:    "notifications@example.com",
		Subject: "Hello {{{NAME}}}",
		Html:    "<strong>Hey, {{{NAME}}}, you are {{{AGE}}} years old.</strong>",
		Variables: []*resend.TemplateVariable{
			{
				Key:           "NAME",
				Type:          resend.VariableTypeString,
				FallbackValue: "user",
			},
			{
				Key:           "AGE",
				Type:          resend.VariableTypeNumber,
				FallbackValue: 25,
			},
			{
				Key:  "OPTIONAL_VARIABLE",
				Type: resend.VariableTypeString,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created template with variables: %s\n", templateWithVars.Id)

	// Create a template with all optional fields
	fullTemplate, err := client.Templates.Create(&resend.CreateTemplateRequest{
		Name:    "full-template-example",
		Alias:   "full-example",
		From:    "Team <team@example.com>",
		Subject: "Important Update for {{{NAME}}}",
		ReplyTo: []string{"support@example.com", "help@example.com"},
		Html:    "<h1>Hello {{{NAME}}}</h1><p>You have {{{COUNT}}} new messages.</p>",
		Text:    "Hello {{{NAME}}}\nYou have {{{COUNT}}} new messages.",
		Variables: []*resend.TemplateVariable{
			{
				Key:           "NAME",
				Type:          resend.VariableTypeString,
				FallbackValue: "Guest",
			},
			{
				Key:           "COUNT",
				Type:          resend.VariableTypeNumber,
				FallbackValue: 0,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created full template: %s (object: %s)\n", fullTemplate.Id, fullTemplate.Object)

	// Create a template with boolean, list, and object variables
	// IMPORTANT: 'object' and 'list' types REQUIRE a fallback_value
	// - For 'list' type: Must contain at least one item (cannot be empty array)
	// - For 'object' type: Must be a valid object
	advancedTemplate, err := client.Templates.Create(&resend.CreateTemplateRequest{
		Name: "advanced-template",
		Html: "<div>{{{IS_PREMIUM}}} - {{{ITEMS}}} - {{{USER}}}</div>",
		Variables: []*resend.TemplateVariable{
			{
				Key:           "IS_PREMIUM",
				Type:          resend.VariableTypeBoolean,
				FallbackValue: false,
			},
			{
				Key:           "ITEMS",
				Type:          resend.VariableTypeList,
				FallbackValue: []interface{}{"default-item"}, // Must have at least one item
			},
			{
				Key:           "USER",
				Type:          resend.VariableTypeObject,
				FallbackValue: map[string]interface{}{"name": "Guest", "id": 0}, // Object with default values
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created advanced template: %s\n", advancedTemplate.Id)

	// Get a template by ID
	retrievedTemplate, err := client.Templates.GetWithContext(ctx, template.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nRetrieved template by ID:\n")
	fmt.Printf("  Id: %s\n", retrievedTemplate.Id)
	fmt.Printf("  Name: %s\n", retrievedTemplate.Name)
	fmt.Printf("  Alias: %s\n", retrievedTemplate.Alias)
	fmt.Printf("  Status: %s\n", retrievedTemplate.Status)
	fmt.Printf("  CreatedAt: %s\n", retrievedTemplate.CreatedAt)
	fmt.Printf("  UpdatedAt: %s\n", retrievedTemplate.UpdatedAt)
	fmt.Printf("  PublishedAt: %s\n", retrievedTemplate.PublishedAt)
	fmt.Printf("  From: %s\n", retrievedTemplate.From)
	fmt.Printf("  Subject: %s\n", retrievedTemplate.Subject)
	fmt.Printf("  Html length: %d\n", len(retrievedTemplate.Html))
	fmt.Printf("  Text length: %d\n", len(retrievedTemplate.Text))
	fmt.Printf("  Variables count: %d\n", len(retrievedTemplate.Variables))

	// Get a template by alias (if alias was set)
	if fullTemplate.Id != "" {
		// Note: In real usage, you would use the alias like "full-example"
		// but to avoid duplicate key errors, we're using ID here
		templateByAlias, err := client.Templates.Get(fullTemplate.Id)
		if err != nil {
			panic(err)
		}
		fmt.Printf("\nRetrieved template: %s (Status: %s)\n", templateByAlias.Name, templateByAlias.Status)
	}

	// Update a template
	updatedTemplate, err := client.Templates.UpdateWithContext(ctx, template.Id, &resend.UpdateTemplateRequest{
		Name:    "welcome-email-updated",
		Html:    "<strong>Welcome to our updated service, {{{NAME}}}!</strong>",
		Subject: "Welcome {{{NAME}}}!",
		Variables: []*resend.TemplateVariable{
			{
				Key:           "NAME",
				Type:          resend.VariableTypeString,
				FallbackValue: "Guest",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nUpdated template: %s\n", updatedTemplate.Id)

	// Verify the update by getting the template again
	verifyTemplate, err := client.Templates.Get(updatedTemplate.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Verified updated template:\n")
	fmt.Printf("  Name: %s\n", verifyTemplate.Name)
	fmt.Printf("  Subject: %s\n", verifyTemplate.Subject)
	fmt.Printf("  Variables count: %d\n", len(verifyTemplate.Variables))

	// Publish a template
	// Note: Templates must be published before they can be used to send emails
	publishedTemplate, err := client.Templates.PublishWithContext(ctx, template.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nPublished template: %s\n", publishedTemplate.Id)

	// Verify the template is now published
	publishedCheck, err := client.Templates.Get(publishedTemplate.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Template status after publish: %s\n", publishedCheck.Status)
	fmt.Printf("PublishedAt: %s\n", publishedCheck.PublishedAt)

	// Duplicate a template
	// Note: This creates a new template as a copy of the original
	duplicatedTemplate, err := client.Templates.DuplicateWithContext(ctx, fullTemplate.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nDuplicated template: %s\n", duplicatedTemplate.Id)

	// Verify the duplicated template was created
	// The duplicate will be a separate template with its own ID
	duplicateCheck, err := client.Templates.Get(duplicatedTemplate.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Duplicated template details:\n")
	fmt.Printf("  Id: %s\n", duplicateCheck.Id)
	fmt.Printf("  Name: %s\n", duplicateCheck.Name)
	fmt.Printf("  Status: %s\n", duplicateCheck.Status)
	fmt.Printf("  CreatedAt: %s\n", duplicateCheck.CreatedAt)

	// Remove a template
	// Note: This permanently deletes the template
	removedTemplate, err := client.Templates.RemoveWithContext(ctx, duplicatedTemplate.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nRemoved template:\n")
	fmt.Printf("  Id: %s\n", removedTemplate.Id)
	fmt.Printf("  Object: %s\n", removedTemplate.Object)
	fmt.Printf("  Deleted: %t\n", removedTemplate.Deleted)
}
