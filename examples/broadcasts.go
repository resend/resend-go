package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func broadcastsExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Create Broadcast
	params := &resend.CreateBroadcastRequest{
		AudienceId: "78b8d3bc-a55a-45a3-aee6-6ec0a5e13d7e",
		From:       "onboarding@resend.dev",
		Html:       "<html><body><h1>Hello, world!</h1></body></html>",
		Name:       "Test Broadcast",
		Subject:    "Hello, world!",
	}

	broadcast, err := client.Broadcasts.CreateWithContext(ctx, params)

	if err != nil {
		panic(err)
	}
	fmt.Println("Created broadcast with entry id: " + broadcast.Id)

	// Get Broadcast
	retrievedBroadcast, err := client.Broadcasts.GetWithContext(ctx, broadcast.Id)
	if err != nil {
		panic(err)
	}

	fmt.Println("ID: " + retrievedBroadcast.Id)
	fmt.Println("Name: " + retrievedBroadcast.Name)
	fmt.Println("Audience ID: " + retrievedBroadcast.AudienceId)
	fmt.Println("From: " + retrievedBroadcast.From)
	fmt.Println("Subject: " + retrievedBroadcast.Subject)
	fmt.Println("Preview Text: " + retrievedBroadcast.PreviewText)
	fmt.Println("Status: " + retrievedBroadcast.Status)
	fmt.Println("Created At: " + retrievedBroadcast.CreatedAt)
	fmt.Println("Scheduled At: " + retrievedBroadcast.ScheduledAt)
	fmt.Println("Sent At: " + retrievedBroadcast.SentAt)

	// Send Broadcast
	sendParams := &resend.SendBroadcastRequest{
		BroadcastId: broadcast.Id,
	}

	sendResponse, err := client.Broadcasts.SendWithContext(ctx, sendParams)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sent broadcast with entry id: " + sendResponse.Id)

	// List Broadcasts

	listResponse, err := client.Broadcasts.ListWithContext(ctx)
	if err != nil {
		panic(err)
	}

	for _, b := range listResponse.Data {
		fmt.Println("ID: " + b.Id)
		fmt.Println("Name: " + b.Name)
		fmt.Println("Audience ID: " + b.AudienceId)
		fmt.Println("From: " + b.From)
		fmt.Println("Subject: " + b.Subject)
		fmt.Println("Preview Text: " + b.PreviewText)
		fmt.Println("Status: " + b.Status)
		fmt.Println("Created At: " + b.CreatedAt)
		fmt.Println("Scheduled At: " + b.ScheduledAt)
		fmt.Println("Sent At: " + b.SentAt)
	}

	// Remove Broadcast (Only Draft Broadcasts can be deleted)
	// removeResponse, err := client.Broadcasts.RemoveWithContext(ctx, broadcast.Id)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Deleted broadcast with entry id: " + removeResponse.Id)
}
