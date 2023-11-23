package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resendlabs/resend-go/v2"
)

func apiKeysExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Create API Key
	params := &resend.CreateApiKeyRequest{
		Name: "nice api key",
	}

	resp, err := client.ApiKeys.CreateWithContext(ctx, params)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created API Key id: " + resp.Id)
	fmt.Println("Token: " + resp.Token)

	// List
	apiKeys, err := client.ApiKeys.ListWithContext(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("You have %d api keys in your project\n", len(apiKeys.Data))

	// Delete
	_, err = client.ApiKeys.RemoveWithContext(ctx, resp.Id)
	if err != nil {
		panic(err)
	}
	println("deleted api key id: " + resp.Id)
}
