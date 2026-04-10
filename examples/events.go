package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

func eventsExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Create an event definition
	event, err := client.Events.CreateWithContext(ctx, &resend.CreateEventRequest{
		Name: "user.created",
		Schema: map[string]string{
			"plan":     resend.EventSchemaTypeString,
			"age":      resend.EventSchemaTypeNumber,
			"verified": resend.EventSchemaTypeBoolean,
			"joined_at": resend.EventSchemaTypeDate,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created event id: %s\n", event.Id)

	// Get an event by ID
	retrieved, err := client.Events.GetWithContext(ctx, event.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Event name: %s\n", retrieved.Name)
	fmt.Printf("Event schema: %v\n", retrieved.Schema)

	// Get an event by name
	byName, err := client.Events.GetWithContext(ctx, "user.created")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Event by name id: %s\n", byName.Id)

	// List events
	events, err := client.Events.ListWithContext(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total events: %d\n", len(events.Data))

	// List events with pagination
	limit := 10
	eventsPage, err := client.Events.ListWithOptions(ctx, &resend.ListOptions{
		Limit: &limit,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Page events: %d\n", len(eventsPage.Data))

	if eventsPage.HasMore {
		lastId := eventsPage.Data[len(eventsPage.Data)-1].Id
		nextPage, err := client.Events.ListWithOptions(ctx, &resend.ListOptions{
			Limit: &limit,
			After: &lastId,
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Next page events: %d\n", len(nextPage.Data))
	}

	// Update an event's schema
	updated, err := client.Events.UpdateWithContext(ctx, event.Id, &resend.UpdateEventRequest{
		Schema: map[string]string{
			"plan":     resend.EventSchemaTypeString,
			"verified": resend.EventSchemaTypeBoolean,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Updated event id: %s\n", updated.Id)

	// Send an event by email address
	sentByEmail, err := client.Events.SendWithContext(ctx, &resend.SendEventRequest{
		Event: "user.created",
		Email: "user@example.com",
		Payload: map[string]any{
			"plan":     "pro",
			"verified": true,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Sent event: %s\n", sentByEmail.Event)

	// Create a contact to use its ID for sending an event
	contact, err := client.Contacts.CreateWithContext(ctx, &resend.CreateContactRequest{
		Email:     "user@example.com",
		FirstName: "John",
		LastName:  "Doe",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created contact id: %s\n", contact.Id)

	// Send an event by contact ID
	sentByContact, err := client.Events.SendWithContext(ctx, &resend.SendEventRequest{
		Event:     "user.created",
		ContactId: contact.Id,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Sent event: %s\n", sentByContact.Event)

	// Delete an event
	deleted, err := client.Events.RemoveWithContext(ctx, event.Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted event id: %s, deleted: %v\n", deleted.Id, deleted.Deleted)
}
