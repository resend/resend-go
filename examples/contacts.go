package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func contactsExample() {

	audienceId := "ca4e37c5-a82a-4199-a3b8-bf912a6472aa"

	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")
	contactEmail := "hi@example.com"

	client := resend.NewClient(apiKey)

	// Create a topic first (for demonstrating topic subscriptions later)
	topicParams := &resend.CreateTopicRequest{
		Name:                "Product Updates",
		DefaultSubscription: resend.DefaultSubscriptionOptOut,
		Description:         "Latest product updates and announcements",
	}

	topic, err := client.Topics.CreateWithContext(ctx, topicParams)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created topic with ID: " + topic.Id)

	// Create Contact params
	params := &resend.CreateContactRequest{
		Email:        contactEmail,
		AudienceId:   audienceId,
		FirstName:    "Steve",
		LastName:     "Woz",
		Unsubscribed: true,
	}

	contact, err := client.Contacts.CreateWithContext(ctx, params)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created contact with entry id: " + contact.Id)

	// Update
	updateParams := &resend.UpdateContactRequest{
		AudienceId: audienceId,
		Id:         "88ffbe62-9bd6-4a39-9ddc-4d51053e172a",
		FirstName:  "new Updated First Name",
		LastName:   "new Updated Last Name",
	}

	// Set unsubscribed to false
	updateParams.SetUnsubscribed(false)

	_, err = client.Contacts.UpdateWithContext(ctx, updateParams)
	if err != nil {
		panic(err)
	}

	// Get by ID
	retrievedContact, err := client.Contacts.GetWithContext(ctx, &resend.GetContactOptions{
		AudienceId: audienceId,
		Id:         contact.Id,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nRetrieved contact by ID: %v\n", retrievedContact)

	// Get by email
	retrievedByEmail, err := client.Contacts.GetWithContext(ctx, &resend.GetContactOptions{
		AudienceId: audienceId,
		Id:         contactEmail,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nRetrieved contact by email: %v\n", retrievedByEmail)

	// List
	contacts, err := client.Contacts.ListWithContext(ctx, &resend.ListContactsOptions{
		AudienceId: audienceId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("You have %d contacts in your audience\n", len(contacts.Data))
	for _, c := range contacts.Data {
		fmt.Printf("%v\n", c)
	}

	// ========== Contact Topics ==========

	// Retrieve contact topics by ID
	topics, err := client.Contacts.Topics.ListWithContext(ctx, contact.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nContact has %d topic subscriptions:\n", len(topics.Data))
	for _, topic := range topics.Data {
		fmt.Printf("  - %s (%s): %s\n", topic.Name, topic.Id, topic.Subscription)
	}

	// Retrieve contact topics by email
	topicsByEmail, err := client.Contacts.Topics.ListWithContext(ctx, contactEmail)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nContact (by email) has %d topic subscriptions\n", len(topicsByEmail.Data))

	// Update topic subscriptions by contact ID
	updateTopicsParams := &resend.UpdateContactTopicsRequest{
		Id: contact.Id,
		Topics: []resend.TopicSubscriptionUpdate{
			{
				Id:           topic.Id,
				Subscription: "opt_in",
			},
		},
	}

	updatedTopics, err := client.Contacts.Topics.UpdateWithContext(ctx, updateTopicsParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nUpdated topic subscriptions for contact: %s\n", updatedTopics.Id)

	// Retrieve updated topics to verify
	updatedTopicsList, err := client.Contacts.Topics.ListWithContext(ctx, contact.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Contact now has %d topic subscription(s):\n", len(updatedTopicsList.Data))
	for _, t := range updatedTopicsList.Data {
		fmt.Printf("  - %s: %s\n", t.Name, t.Subscription)
	}

	// You can also update by email instead of ID
	updateTopicsByEmailParams := &resend.UpdateContactTopicsRequest{
		Email: contactEmail,
		Topics: []resend.TopicSubscriptionUpdate{
			{
				Id:           topic.Id,
				Subscription: "opt_out",
			},
		},
	}

	updatedTopicsByEmail, err := client.Contacts.Topics.UpdateWithContext(ctx, updateTopicsByEmailParams)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nUpdated topic subscriptions for contact (by email): %s\n", updatedTopicsByEmail.Id)

	// ====================================

	// Remove by id
	removed, err := client.Contacts.RemoveWithContext(ctx, &resend.RemoveContactOptions{
		AudienceId: audienceId,
		Id:         contact.Id,
	})

	// Remove by email
	// removed, err := client.Contacts.RemoveWithContext(ctx, &resend.RemoveContactOptions{
	//   AudienceId: audienceId,
	//   Id:         "hi@example.com",
	// })
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nContact deleted: %v\n", removed.Deleted)

	// Clean up: Remove the topic we created
	removedTopic, err := client.Topics.RemoveWithContext(ctx, topic.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Topic deleted: %v\n", removedTopic.Deleted)
}
