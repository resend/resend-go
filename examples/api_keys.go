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

	fmt.Printf("%v", params)
	fmt.Println("GOING TO CREATE")

	resp, err := client.ApiKeys.Create(params)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Id)
	fmt.Println(resp.Token)

	// // Get
	// apiKey, err := client.Emails.Get(sent.Id)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%v\n", email)
}
