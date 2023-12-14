package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resendlabs/resend-go/v2"
)

// func contactsExample() {
func main() {

	audienceId := "709d076c-2bb1-4be6-94ed-3f8f32622db6"

	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Create Contact params
	params := &resend.CreateContactRequest{
		Email: "hi2@example.com",
	}

	contact, err := client.Contacts.CreateWithContext(ctx, audienceId, params)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created contact with entry id: " + contact.Id)

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

	// Remove
	removed, err := client.Contacts.RemoveWithContext(ctx, audienceId, contact.Id)
	if err != nil {
		panic(err)
	}
	println(removed.Deleted)
}
