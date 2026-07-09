package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

func oauthGrantsExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// List OAuth grants
	grants, err := client.OAuthGrants.ListWithContext(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("You have %d OAuth grants\n", len(grants.Data))

	if len(grants.Data) == 0 {
		return
	}

	// Revoke the first grant
	revoked, err := client.OAuthGrants.RevokeWithContext(ctx, grants.Data[0].Id)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Revoked OAuth grant %s (reason: %s)\n", revoked.Id, revoked.RevokedReason)
}
