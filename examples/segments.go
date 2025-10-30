package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v2"
)

func segmentsExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Create Segment params
	params := &resend.CreateSegmentRequest{
		Name: "New Segment",
	}

	segment, err := client.Segments.CreateWithContext(ctx, params)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created segment with entry id: " + segment.Id)

	// Get
	retrievedSegment, err := client.Segments.GetWithContext(ctx, "78b8d3bc-a55a-45a3-aee6-6ec0a5e13d7e")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nRetrieved segment: %v\n", retrievedSegment)

	// List
	segments, err := client.Segments.ListWithContext(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("You have %d segments in your project\n", len(segments.Data))

	// Remove
	removed, err := client.Segments.RemoveWithContext(ctx, segment.Id)
	if err != nil {
		panic(err)
	}
	println(removed.Deleted)
}
