package examples

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/resend/resend-go/v2"
)

func handleRateLimitExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		To:      []string{"delivered@resend.dev"},
		From:    "onboarding@resend.dev",
		Text:    "hello world",
		Subject: "Hello from Golang",
	}

	sent, err := client.Emails.SendWithContext(ctx, params)
	if err != nil {
		// Check if it's a rate limit error using errors.Is
		if errors.Is(err, resend.ErrRateLimit) {
			fmt.Println("Rate limit exceeded!")

			// Extract detailed rate limit information
			var rateLimitErr *resend.RateLimitError
			if errors.As(err, &rateLimitErr) {
				fmt.Printf("Message: %s\n", rateLimitErr.Message)
				fmt.Printf("Limit: %s requests\n", rateLimitErr.Limit)
				fmt.Printf("Remaining: %s requests\n", rateLimitErr.Remaining)
				fmt.Printf("Reset in: %s seconds\n", rateLimitErr.Reset)
				fmt.Printf("Retry after: %s seconds\n", rateLimitErr.RetryAfter)

				// Implement retry logic
				if retryAfter, err := strconv.Atoi(rateLimitErr.RetryAfter); err == nil {
					fmt.Printf("Waiting %d seconds before retry...\n", retryAfter)
					time.Sleep(time.Duration(retryAfter) * time.Second)

					// Retry the request
					sent, err = client.Emails.SendWithContext(ctx, params)
					if err != nil {
						panic(err)
					}
					fmt.Printf("Successfully sent after retry: %s\n", sent.Id)
				}
			}
			return
		}

		// Handle other errors
		panic(err)
	}

	fmt.Printf("Email sent successfully: %s\n", sent.Id)
}
