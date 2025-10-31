package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

func contactSegmentsExample() {
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Example 1: Add a contact to a segment by contact ID
	addParams := &resend.AddContactSegmentRequest{
		SegmentId: "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
		ContactId: "479e3145-dd38-476b-932c-529ceb705947",
	}

	addResp, err := client.Contacts.Segments.Add(addParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Contact added to segment: %s\n", addResp.Id)

	// Example 2: Add a contact to a segment by email
	addByEmailParams := &resend.AddContactSegmentRequest{
		SegmentId: "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
		Email:     "user@example.com",
	}

	addByEmailResp, err := client.Contacts.Segments.Add(addByEmailParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Contact added to segment by email: %s\n", addByEmailResp.Id)

	// Example 3: List all segments for a contact by contact ID
	listParams := &resend.ListContactSegmentsRequest{
		ContactId: "479e3145-dd38-476b-932c-529ceb705947",
	}

	listResp, err := client.Contacts.Segments.List(listParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d segments for contact\n", len(listResp.Data))
	for _, segment := range listResp.Data {
		fmt.Printf("  - %s (ID: %s)\n", segment.Name, segment.Id)
	}

	// Example 4: List all segments for a contact by email with pagination
	limit := 10
	listWithPaginationParams := &resend.ListContactSegmentsRequest{
		Email: "user@example.com",
	}
	options := &resend.ListOptions{
		Limit: &limit,
	}

	paginatedResp, err := client.Contacts.Segments.ListWithOptions(
		context.Background(),
		listWithPaginationParams,
		options,
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Found %d segments (has more: %v)\n", len(paginatedResp.Data), paginatedResp.HasMore)

	// Example 5: Remove a contact from a segment by contact ID
	removeParams := &resend.RemoveContactSegmentRequest{
		SegmentId: "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
		ContactId: "479e3145-dd38-476b-932c-529ceb705947",
	}

	removeResp, err := client.Contacts.Segments.Remove(removeParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Contact removed from segment: %v\n", removeResp.Deleted)

	// Example 6: Remove a contact from a segment by email
	removeByEmailParams := &resend.RemoveContactSegmentRequest{
		SegmentId: "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
		Email:     "user@example.com",
	}

	removeByEmailResp, err := client.Contacts.Segments.Remove(removeByEmailParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Contact removed from segment by email: %v\n", removeByEmailResp.Deleted)
}
