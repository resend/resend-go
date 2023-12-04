package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resendlabs/resend-go/v2"
)

func audiencesExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Create Audience params
	params := &resend.CreateAudienceRequest{
		Name: "New Audience",
	}

	audience, err := client.Audiences.CreateWithContext(ctx, params)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created audience with entry id: " + audience.Id)

	// Get
	retrievedAudience, err := client.Audiences.GetWithContext(ctx, "78b8d3bc-a55a-45a3-aee6-6ec0a5e13d7e")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nRetrieved audience: %v\n", retrievedAudience)

	// List
	audiences, err := client.Audiences.ListWithContext(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("You have %d audiences in your project\n", len(audiences.Data))

	// Remove
	removed, err := client.Audiences.RemoveWithContext(ctx, audience.Id)
	if err != nil {
		panic(err)
	}
	println(removed.Deleted)
}
