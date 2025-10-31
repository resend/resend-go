package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

func contactPropertiesExample() {

	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Create Contact Property with number type
	createParams := &resend.CreateContactPropertyRequest{
		Key:           "age",
		Type:          "number",
		FallbackValue: 0,
	}

	property, err := client.ContactProperties.CreateWithContext(ctx, createParams)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created contact property with ID: " + property.Id)

	// Create another property with string type
	createStringProperty := &resend.CreateContactPropertyRequest{
		Key:           "country",
		Type:          "string",
		FallbackValue: "US",
	}

	stringProperty, err := client.ContactProperties.CreateWithContext(ctx, createStringProperty)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created contact property with ID: " + stringProperty.Id)

	// Get by ID
	retrievedProperty, err := client.ContactProperties.GetWithContext(ctx, property.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nRetrieved contact property by ID: %v\n", retrievedProperty)

	// List all contact properties
	properties, err := client.ContactProperties.ListWithContext(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nYou have %d contact properties\n", len(properties.Data))
	for _, p := range properties.Data {
		fmt.Printf("Property: %s (type: %s, fallback: %v)\n", p.Key, p.Type, p.FallbackValue)
	}

	// Update contact property
	updateParams := &resend.UpdateContactPropertyRequest{
		Id:            property.Id,
		FallbackValue: 25,
	}

	updated, err := client.ContactProperties.UpdateWithContext(ctx, updateParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nUpdated contact property: %s\n", updated.Id)

	// Remove contact property
	removed, err := client.ContactProperties.RemoveWithContext(ctx, property.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted contact property: %s (deleted: %v)\n", removed.Id, removed.Deleted)
}
