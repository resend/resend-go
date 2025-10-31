package examples

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/resend/resend-go/v2"
)

func withPaginationExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// List domains with pagination
	fmt.Println("=== List domains with pagination ===")

	// List first 5 domains
	limit := 5
	domainsResp, err := client.Domains.ListWithOptions(ctx, &resend.ListOptions{
		Limit: &limit,
	})
	if err != nil {
		log.Printf("Error listing domains: %v", err)
	} else {
		fmt.Printf("Found %d domains (HasMore: %t)\n", len(domainsResp.Data), domainsResp.HasMore)
		for _, domain := range domainsResp.Data {
			fmt.Printf("  - %s (ID: %s)\n", domain.Name, domain.Id)
		}
	}

	// List API keys with cursor pagination
	fmt.Println("\n=== List API keys with cursor pagination ===")

	after := "550e8400-e29b-41d4-a716-446655440000" // Example UUID
	limit10 := 10
	apiKeysResp, err := client.ApiKeys.ListWithOptions(ctx, &resend.ListOptions{
		Limit: &limit10,
		After: &after,
	})
	if err != nil {
		log.Printf("Error listing API keys: %v", err)
	} else {
		fmt.Printf("Found %d API keys after cursor (HasMore: %t)\n", len(apiKeysResp.Data), apiKeysResp.HasMore)
		for _, apiKey := range apiKeysResp.Data {
			fmt.Printf("  - %s (ID: %s)\n", apiKey.Name, apiKey.Id)
		}
	}

	// List segments using before cursor
	fmt.Println("\n=== List segments with before cursor ===")

	before := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	limit3 := 3
	segmentsResp, err := client.Segments.ListWithOptions(ctx, &resend.ListOptions{
		Limit:  &limit3,
		Before: &before,
	})
	if err != nil {
		log.Printf("Error listing segments: %v", err)
	} else {
		fmt.Printf("Found %d segments before cursor (HasMore: %t)\n", len(segmentsResp.Data), segmentsResp.HasMore)
		for _, segment := range segmentsResp.Data {
			fmt.Printf("  - %s (ID: %s)\n", segment.Name, segment.Id)
		}
	}

	// List contacts in an audience with pagination
	fmt.Println("\n=== List contacts with pagination ===")

	audienceId := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	limit20 := 20
	contactsResp, err := client.Contacts.ListWithContext(ctx, &resend.ListContactsOptions{
		AudienceId: audienceId,
		Limit:      &limit20,
	})
	if err != nil {
		log.Printf("Error listing contacts: %v", err)
	} else {
		fmt.Printf("Found %d contacts in audience %s (HasMore: %t)\n", len(contactsResp.Data), audienceId, contactsResp.HasMore)
		for _, contact := range contactsResp.Data {
			fmt.Printf("  - %s (ID: %s)\n", contact.Email, contact.Id)
		}
	}

	// List broadcasts with pagination
	fmt.Println("\n=== List broadcasts with pagination ===")

	limit2 := 2
	broadcastsResp, err := client.Broadcasts.ListWithOptions(ctx, &resend.ListOptions{
		Limit: &limit2,
	})
	if err != nil {
		log.Printf("Error listing broadcasts: %v", err)
	} else {
		fmt.Printf("Found %d broadcasts (HasMore: %t)\n", len(broadcastsResp.Data), broadcastsResp.HasMore)
		for _, broadcast := range broadcastsResp.Data {
			fmt.Printf("  - %s (ID: %s)\n", broadcast.Name, broadcast.Id)
		}
	}
}
