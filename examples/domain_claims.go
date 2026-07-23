package examples

import (
	"context"
	"fmt"
	"os"

	"github.com/resend/resend-go/v3"
)

func domainClaimsExample() {
	ctx := context.TODO()
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	// Start a claim for a domain that another Resend account has already verified.
	claimParams := &resend.CreateDomainClaimRequest{
		Name:   "example.com",
		Region: "us-east-1",
	}

	claim, err := client.DomainClaims.CreateWithContext(ctx, claimParams)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created domain claim id: " + claim.Id)
	fmt.Println("Status: " + claim.Status)
	if claim.Record != nil {
		fmt.Printf("Add this TXT record to prove ownership: %s = %s\n", claim.Record.Name, claim.Record.Value)
	}

	// Get: poll the claim until the TXT record has been added and verification can run.
	retrievedClaim, err := client.DomainClaims.GetWithContext(ctx, claim.DomainId)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Retrieved domain claim: %v\n", retrievedClaim)

	// Verify: trigger asynchronous DNS verification and ownership transfer.
	verifiedClaim, err := client.DomainClaims.VerifyWithContext(ctx, claim.DomainId)
	if err != nil {
		panic(err)
	}
	fmt.Println("Verification triggered, claim status: " + verifiedClaim.Status)
}
