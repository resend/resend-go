package main

import (
	"fmt"
	"os"

	"github.com/drish/resend-go"
)

func main() {

	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Create API Key
	params := &resend.CreateApiKeyRequest{
		Name: "nice api key",
	}

	resp, err := client.ApiKeys.Create(params)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created API Key id: " + resp.Id)
	fmt.Println("Token: " + resp.Token)

	// List
	apiKeys, err := client.ApiKeys.List()
	if err != nil {
		panic(err)
	}
	fmt.Printf("You have %d api keys in your project\n", len(apiKeys.Data))

	// Delete
	client.ApiKeys.Delete(resp.Id)
	println("deleted api key id: " + resp.Id)
}
