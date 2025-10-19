package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func topicsExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Create a topic with opt_in default subscription
	// Note: default_subscription cannot be changed after creation
	topic, err := client.Topics.CreateWithContext(ctx, &resend.CreateTopicRequest{
		Name:                "Weekly Newsletter",
		DefaultSubscription: resend.DefaultSubscriptionOptIn,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created topic: %s\n", topic.Id)

	// Create a topic with opt_out default subscription
	optOutTopic, err := client.Topics.Create(&resend.CreateTopicRequest{
		Name:                "Product Updates",
		DefaultSubscription: resend.DefaultSubscriptionOptOut,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created opt-out topic: %s\n", optOutTopic.Id)

	// Create a topic with description
	// Note: name max length is 50 characters, description max length is 200 characters
	topicWithDescription, err := client.Topics.Create(&resend.CreateTopicRequest{
		Name:                "Monthly Summary",
		DefaultSubscription: resend.DefaultSubscriptionOptIn,
		Description:         "Monthly summary of your account activity and updates",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created topic with description: %s\n", topicWithDescription.Id)

	// Get a topic by ID
	retrievedTopic, err := client.Topics.GetWithContext(ctx, topic.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nRetrieved topic:\n")
	fmt.Printf("  Id: %s\n", retrievedTopic.Id)
	fmt.Printf("  Name: %s\n", retrievedTopic.Name)
	fmt.Printf("  Description: %s\n", retrievedTopic.Description)
	fmt.Printf("  DefaultSubscription: %s\n", retrievedTopic.DefaultSubscription)
	fmt.Printf("  CreatedAt: %s\n", retrievedTopic.CreatedAt)
}
