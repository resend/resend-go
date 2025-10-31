package examples

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func globalContactsExample() {
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Example 1: Create a global contact with custom properties
	// Global contacts don't require an audience_id and support custom properties
	createParams := &resend.CreateContactRequest{
		Email:     "user@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Properties: map[string]interface{}{
			"tier":          "premium",
			"role":          "admin",
			"signup_source": "website",
			"active":        "true",
			"age":           "30",
		},
	}

	created, err := client.Contacts.Create(createParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created global contact: %s\n", created.Id)

	// Example 2: List all global contacts
	// Omit AudienceId to list global contacts
	contacts, err := client.Contacts.List(&resend.ListContactsOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d global contacts\n", len(contacts.Data))
	for _, contact := range contacts.Data {
		fmt.Printf("  - %s %s (%s)\n", contact.FirstName, contact.LastName, contact.Email)
		if contact.Properties != nil {
			fmt.Printf("    Properties: %+v\n", contact.Properties)
		}
	}

	// Example 3: Get a specific global contact
	// Omit AudienceId to get a global contact
	contact, err := client.Contacts.Get(&resend.GetContactOptions{
		Id: created.Id,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Retrieved contact: %s (%s)\n", contact.Email, contact.Id)
	if contact.Properties != nil {
		fmt.Printf("Properties: %+v\n", contact.Properties)
	}

	// Example 4: Update a global contact with new properties
	updateParams := &resend.UpdateContactRequest{
		Id:        created.Id,
		FirstName: "Jane",
		Properties: map[string]interface{}{
			"tier":   "enterprise",
			"role":   "owner",
			"active": "true", // Use string, not boolean
		},
	}

	updated, err := client.Contacts.Update(updateParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Updated contact: %s\n", updated.Data.Email)
	fmt.Printf("New properties: %+v\n", updated.Data.Properties)

	// Example 5: Remove a global contact
	// Omit AudienceId to remove a global contact
	removed, err := client.Contacts.Remove(&resend.RemoveContactOptions{
		Id: created.Id,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Removed contact: %s (deleted: %v)\n", removed.Id, removed.Deleted)

	// Example 6: Create a global contact and add to segments
	// First create a global contact
	globalContact := &resend.CreateContactRequest{
		Email:     "segmented@example.com",
		FirstName: "Segmented",
		LastName:  "User",
		Properties: map[string]interface{}{
			"source": "landing_page",
		},
	}

	contactResp, err := client.Contacts.Create(globalContact)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created global contact: %s\n", contactResp.Id)

	// Then add to a segment using the Contacts.Segments API
	addToSegment := &resend.AddContactSegmentRequest{
		ContactId: contactResp.Id,
		SegmentId: "segment_id_here",
	}

	_, err = client.Contacts.Segments.Add(addToSegment)
	if err != nil {
		fmt.Printf("Note: Could not add to segment (may not exist): %v\n", err)
	} else {
		fmt.Println("Added contact to segment")
	}

	// List all segments for this global contact
	listSegments := &resend.ListContactSegmentsRequest{
		ContactId: contactResp.Id,
	}

	segments, err := client.Contacts.Segments.List(listSegments)
	if err != nil {
		fmt.Printf("Note: Could not list segments: %v\n", err)
	} else {
		fmt.Printf("Contact is in %d segments\n", len(segments.Data))
	}
}
