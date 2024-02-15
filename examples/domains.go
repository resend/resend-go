package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func domainsExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Create Domain params
	params := &resend.CreateDomainRequest{
		Name: "exampledomain.com",
	}

	domain, err := client.Domains.CreateWithContext(ctx, params)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created Domain entry id: " + domain.Id)
	fmt.Println("Status: " + domain.Status)

	for _, record := range domain.Records {
		fmt.Printf("%v\n", record)
	}

	// Get
	retrievedDomain, err := client.Domains.GetWithContext(ctx, domain.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Retrieved domain: %v", retrievedDomain)

	updateDomainParams := &resend.UpdateDomainRequest{
		OpenTracking:  true,
		ClickTracking: true,
	}

	updated, err := client.Domains.UpdateWithContext(ctx, domain.Id, updateDomainParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", updated)

	// List
	domains, err := client.Domains.ListWithContext(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("You have %d domains in your project\n", len(domains.Data))

	// Verify
	verified, err := client.Domains.VerifyWithContext(ctx, domain.Id)
	if err != nil {
		panic(err)
	}
	if verified {
		println("verified domain id: " + domain.Id)
	} else {
		println("could not verify domain id: " + domain.Id)
	}

	// Remove
	removed, err := client.Domains.RemoveWithContext(ctx, domain.Id)
	if err != nil {
		panic(err)
	}
	if removed {
		println("removed domain id: " + domain.Id)
	} else {
		println("could not remove domain id: " + domain.Id)
	}
}
