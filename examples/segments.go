package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

func segmentsExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Create a segment
	params := &resend.CreateSegmentRequest{
		Name: "Premium Users",
	}

	segment, err := client.Segments.CreateWithContext(ctx, params)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created segment: %s (ID: %s)\n", segment.Name, segment.Id)

	// Create a global contact
	contactParams := &resend.CreateContactRequest{
		Email:     "premium.user@example.com",
		FirstName: "Premium",
		LastName:  "User",
	}

	contact, err := client.Contacts.CreateWithContext(ctx, contactParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created contact: %s\n", contact.Id)

	// Add the contact to the segment
	addToSegment := &resend.AddContactSegmentRequest{
		ContactId: contact.Id,
		SegmentId: segment.Id,
	}

	_, err = client.Contacts.Segments.AddWithContext(ctx, addToSegment)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Added contact to segment\n")

	// List all segments for this contact
	contactSegments, err := client.Contacts.Segments.ListWithContext(ctx, &resend.ListContactSegmentsRequest{
		ContactId: contact.Id,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Contact is in %d segment(s):\n", len(contactSegments.Data))
	for _, seg := range contactSegments.Data {
		fmt.Printf("  - %s (ID: %s)\n", seg.Name, seg.Id)
	}

	// Get segment details
	retrievedSegment, err := client.Segments.GetWithContext(ctx, segment.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nRetrieved segment: %s\n", retrievedSegment.Name)

	// List all segments
	segments, err := client.Segments.ListWithContext(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("You have %d segments in your project\n", len(segments.Data))

	// Clean up: Remove the contact from the segment
	removeFromSegment, err := client.Contacts.Segments.RemoveWithContext(ctx, &resend.RemoveContactSegmentRequest{
		ContactId: contact.Id,
		SegmentId: segment.Id,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nRemoved contact from segment: %v\n", removeFromSegment.Deleted)

	// Clean up: Remove the contact
	removedContact, err := client.Contacts.RemoveWithContext(ctx, &resend.RemoveContactOptions{
		Id: contact.Id,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted contact: %v\n", removedContact.Deleted)

	// Clean up: Remove the segment
	removedSegment, err := client.Segments.RemoveWithContext(ctx, segment.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted segment: %v\n", removedSegment.Deleted)
}
