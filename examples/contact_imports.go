package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

func contactImportsExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	file, err := os.ReadFile("contacts.csv")
	if err != nil {
		panic(err)
	}

	createParams := &resend.CreateContactImportRequest{
		File:       file,
		Filename:   "contacts.csv",
		OnConflict: "upsert",
		ColumnMap: map[string]any{
			"email":      "Email",
			"first_name": "First Name",
			"last_name":  "Last Name",
			"properties": map[string]any{
				"plan": map[string]any{
					"column": "Plan",
					"type":   "string",
				},
			},
		},
		Segments: []string{"78e7a5c6-9a91-4c63-9d1f-3b9c0b5b9ab6"},
	}

	created, err := client.Contacts.Imports.CreateWithContext(ctx, createParams)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created contact import with ID: " + created.Id)

	// Get a contact import by ID
	imported, err := client.Contacts.Imports.GetWithContext(ctx, created.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Contact import status: %s\n", imported.Status)

	// List contact imports
	limit := 10
	list, err := client.Contacts.Imports.ListWithContext(ctx, &resend.ListContactImportsOptions{
		Status: string(resend.ContactImportStatusCompleted),
		Limit:  &limit,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("You have %d completed contact import(s)\n", len(list.Data))
}
