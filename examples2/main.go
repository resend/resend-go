package main

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func main() {
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
}
