package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func contactsExample() {

	audienceId := "ca4e37c5-a82a-4199-a3b8-bf912a6472aa"

	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")
	contactEmail := "hi@example.com"

	client := resend.NewClient(apiKey)

	// Create Contact params
	params := &resend.CreateContactRequest{
		Email:        contactEmail,
		AudienceId:   audienceId,
		FirstName:    "Steve",
		LastName:     "Woz",
		Unsubscribed: true,
	}

	contact, err := client.Contacts.CreateWithContext(ctx, params)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created contact with entry id: " + contact.Id)

	// Update
	updateParams := &resend.UpdateContactRequest{
		AudienceId: audienceId,
		Id:         "88ffbe62-9bd6-4a39-9ddc-4d51053e172a",
		FirstName:  "new Updated First Name",
		LastName:   "new Updated Last Name",
	}

	// Set unsubscribed to false
	updateParams.SetUnsubscribed(false)

	_, err = client.Contacts.UpdateWithContext(ctx, updateParams)
	if err != nil {
		panic(err)
	}

	// Get by ID
	retrievedContact, err := client.Contacts.GetWithContext(ctx, audienceId, contact.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nRetrieved contact by ID: %v\n", retrievedContact)

	// Get by email
	retrievedByEmail, err := client.Contacts.GetWithContext(ctx, audienceId, contactEmail)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nRetrieved contact by email: %v\n", retrievedByEmail)

	// List
	contacts, err := client.Contacts.ListWithContext(ctx, audienceId)
	if err != nil {
		panic(err)
	}
	fmt.Printf("You have %d contacts in your audience\n", len(contacts.Data))
	for _, c := range contacts.Data {
		fmt.Printf("%v\n", c)
	}

	// Remove by id
	removed, err := client.Contacts.RemoveWithContext(ctx, audienceId, contact.Id)

	// Remove by email
	// removed, err = client.Contacts.RemoveWithContext(ctx, audienceId, "hi@example.com")
	if err != nil {
		panic(err)
	}
	println(removed.Deleted)
}
