package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
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

	// List topics with pagination
	// By default, the API will return the most recent 20 topics
	// You can use limit, after, or before parameters to control pagination
	limit := 2
	listResponse, err := client.Topics.ListWithContext(ctx, &resend.ListOptions{
		Limit: &limit,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nListed topics:\n")
	fmt.Printf("  Object: %s\n", listResponse.Object)
	fmt.Printf("  HasMore: %t\n", listResponse.HasMore)
	fmt.Printf("  Topics count: %d\n", len(listResponse.Data))
	for i, t := range listResponse.Data {
		fmt.Printf("\n  Topic %d:\n", i+1)
		fmt.Printf("    Id: %s\n", t.Id)
		fmt.Printf("    Name: %s\n", t.Name)
		fmt.Printf("    Description: %s\n", t.Description)
		fmt.Printf("    DefaultSubscription: %s\n", t.DefaultSubscription)
		fmt.Printf("    CreatedAt: %s\n", t.CreatedAt)
	}

	// List topics with pagination using after parameter
	if listResponse.HasMore && len(listResponse.Data) > 0 {
		lastId := listResponse.Data[len(listResponse.Data)-1].Id
		nextPage, err := client.Topics.List(&resend.ListOptions{
			Limit: &limit,
			After: &lastId,
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("\nNext page topics count: %d\n", len(nextPage.Data))
	}

	// Update a topic
	// Note: default_subscription cannot be changed after creation
	updatedTopic, err := client.Topics.UpdateWithContext(ctx, topic.Id, &resend.UpdateTopicRequest{
		Name:        "Weekly Newsletter - Updated",
		Description: "Updated weekly newsletter for our valued subscribers",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nUpdated topic: %s\n", updatedTopic.Id)

	// Verify the update by getting the topic again
	verifyTopic, err := client.Topics.Get(updatedTopic.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Verified updated topic:\n")
	fmt.Printf("  Name: %s\n", verifyTopic.Name)
	fmt.Printf("  Description: %s\n", verifyTopic.Description)

	// Remove a topic
	// Note: This permanently deletes the topic
	removedTopic, err := client.Topics.RemoveWithContext(ctx, optOutTopic.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nRemoved topic:\n")
	fmt.Printf("  Object: %s\n", removedTopic.Object)
	fmt.Printf("  Id: %s\n", removedTopic.Id)
	fmt.Printf("  Deleted: %t\n", removedTopic.Deleted)
}
