package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

func logsExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Retrieve a single log
	log, err := client.Logs.GetWithContext(ctx, "log_123abc")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Log ID: %s\n", log.Id)
	fmt.Printf("Endpoint: %s %s\n", log.Method, log.Endpoint)
	fmt.Printf("Response status: %d\n", log.ResponseStatus)
	fmt.Printf("Created at: %s\n", log.CreatedAt)

	// List logs
	logs, err := client.Logs.ListWithContext(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total logs: %d\n", len(logs.Data))
	fmt.Printf("Has more: %v\n", logs.HasMore)

	// List logs with pagination
	limit := 10
	logsPage, err := client.Logs.ListWithOptions(ctx, &resend.ListOptions{
		Limit: &limit,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Page logs: %d\n", len(logsPage.Data))

	if logsPage.HasMore {
		lastId := logsPage.Data[len(logsPage.Data)-1].Id
		nextPage, err := client.Logs.ListWithOptions(ctx, &resend.ListOptions{
			Limit: &limit,
			After: &lastId,
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("Next page logs: %d\n", len(nextPage.Data))
	}
}
