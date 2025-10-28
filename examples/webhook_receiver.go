package examples

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/resend/resend-go/v2"
)

// Demonstrates how to receive and verify webhooks
// This example creates an HTTP server that listens for webhook POST requests
// and verifies them using HMAC-SHA256 signature validation
func webhookReceiverExample() {
	webhookSecret := "whsec_1234567890abcdefghijklmnopqrstuvwxyz"
	apiKey := "re_1234567890abcdefghijklmnopqrstuvwxyz"

	client := resend.NewClient(apiKey)

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read the raw body (must be raw for signature verification)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Extract Svix headers
		headers := resend.WebhookHeaders{
			Id:        r.Header.Get("svix-id"),
			Timestamp: r.Header.Get("svix-timestamp"),
			Signature: r.Header.Get("svix-signature"),
		}

		// Verify the webhook
		err = client.Webhooks.Verify(&resend.VerifyWebhookOptions{
			Payload:       string(body),
			Headers:       headers,
			WebhookSecret: webhookSecret,
		})

		if err != nil {
			log.Printf("Webhook verification failed: %v", err)
			http.Error(w, "Webhook verification failed", http.StatusBadRequest)
			return
		}

		// Parse the verified payload
		var payload map[string]interface{}
		if err := json.Unmarshal(body, &payload); err != nil {
			log.Printf("Error parsing JSON: %v", err)
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		log.Printf("✓ Webhook verified successfully!")
		log.Printf("Event Type: %v", payload["type"])
		log.Printf("Payload: %v", payload)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]bool{"success": true})
	})

	port := ":5000"
	fmt.Printf("🚀 Webhook receiver listening on http://localhost%s/webhook\n", port)
	fmt.Println("Send a POST request with Resend webhook headers to test verification")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
