package examples

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func globalContactsExample() {
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// First, define custom properties (required before using them on contacts)
	propertiesToCreate := []struct {
		key string
		typ string
	}{
		{"tier", "string"},
		{"role", "string"},
		{"signup_source", "string"},
		{"active", "string"},
		{"age", "string"},
		{"source", "string"},
	}

	for _, prop := range propertiesToCreate {
		_, err := client.Contacts.Properties.Create(&resend.CreateContactPropertyRequest{
			Key:  prop.key,
			Type: prop.typ,
		})
		if err != nil {
			// Property might already exist, that's okay
			fmt.Printf("Note: Could not create property '%s': %v\n", prop.key, err)
		} else {
			fmt.Printf("Created property: %s\n", prop.key)
		}
	}

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
	// First create a segment
	segment, err := client.Segments.Create(&resend.CreateSegmentRequest{
		Name: "Example Segment",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created segment: %s (ID: %s)\n", segment.Name, segment.Id)

	// Create a global contact
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

	// Add the contact to the segment using the Contacts.Segments API
	addToSegment := &resend.AddContactSegmentRequest{
		ContactId: contactResp.Id,
		SegmentId: segment.Id,
	}

	_, err = client.Contacts.Segments.Add(addToSegment)
	if err != nil {
		panic(err)
	}
	fmt.Println("Added contact to segment")

	// List all segments for this global contact
	listSegments := &resend.ListContactSegmentsRequest{
		ContactId: contactResp.Id,
	}

	segmentsList, err := client.Contacts.Segments.List(listSegments)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Contact is in %d segment(s):\n", len(segmentsList.Data))
	for _, seg := range segmentsList.Data {
		fmt.Printf("  - %s (ID: %s)\n", seg.Name, seg.Id)
	}

	// Clean up: Remove the segment
	removedSegment, err := client.Segments.Remove(segment.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Segment deleted: %v\n", removedSegment.Deleted)
}
