package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

func webhooksExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Create Webhook
	createParams := &resend.CreateWebhookRequest{
		Endpoint: "https://webhook.example.com/handler",
		Events:   []string{"email.sent", "email.delivered", "email.bounced"},
	}

	created, err := client.Webhooks.CreateWithContext(ctx, createParams)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created Webhook ID: " + created.Id)
	fmt.Println("Signing Secret: " + created.SigningSecret)

	// Get Webhook
	webhook, err := client.Webhooks.GetWithContext(ctx, created.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Webhook Status: %s\n", webhook.Status)
	fmt.Printf("Webhook Endpoint: %s\n", webhook.Endpoint)
	fmt.Printf("Webhook Events: %v\n", webhook.Events)

	// Update Webhook
	newEndpoint := "https://new-webhook.example.com/handler"
	newStatus := "disabled"
	updateParams := &resend.UpdateWebhookRequest{
		Endpoint: &newEndpoint,
		Events:   []string{"email.sent", "email.delivered"},
		Status:   &newStatus,
	}

	updated, err := client.Webhooks.UpdateWithContext(ctx, created.Id, updateParams)
	if err != nil {
		panic(err)
	}
	fmt.Println("Updated Webhook ID: " + updated.Id)

	// List Webhooks
	webhooks, err := client.Webhooks.ListWithContext(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("You have %d webhooks in your project\n", len(webhooks.Data))
	for i, wh := range webhooks.Data {
		fmt.Printf("  [%d] ID: %s, Status: %s, Endpoint: %s\n", i+1, wh.Id, wh.Status, wh.Endpoint)
	}

	// List Webhooks with Pagination
	limit := 10
	listOptions := &resend.ListOptions{
		Limit: &limit,
	}
	paginatedWebhooks, err := client.Webhooks.ListWithOptions(ctx, listOptions)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Fetched %d webhooks (has_more: %v)\n", len(paginatedWebhooks.Data), paginatedWebhooks.HasMore)

	// Delete Webhook
	deleted, err := client.Webhooks.RemoveWithContext(ctx, created.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted Webhook ID: %s (deleted: %v)\n", deleted.Id, deleted.Deleted)
}
