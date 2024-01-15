package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func contactsExample() {

	audienceId := "78b8d3bc-a55a-45a3-aee6-6ec0a5e13d7e"

	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Create Contact params
	params := &resend.CreateContactRequest{
		Email:      "hi@example.com",
		AudienceId: audienceId,
		FirstName:  "Steve",
		LastName:   "Woz",
	}

	contact, err := client.Contacts.CreateWithContext(ctx, params)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created contact with entry id: " + contact.Id)

	// Update
	updateParams := &resend.UpdateContactRequest{
		AudienceId:   audienceId,
		Id:           contact.Id,
		FirstName:    "Updated First Name",
		LastName:     "Updated Last Name",
		Unsubscribed: true,
	}
	_, err = client.Contacts.UpdateWithContext(ctx, updateParams)
	if err != nil {
		panic(err)
	}

	// Get
	retrievedContact, err := client.Contacts.GetWithContext(ctx, audienceId, contact.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nRetrieved contact: %v\n", retrievedContact)

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
	// removed, err := client.Contacts.RemoveWithContext(ctx, audienceId, "hi@example.com")
	if err != nil {
		panic(err)
	}
	println(removed.Deleted)
}
