package examples

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/resend/resend-go/v4"
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
		var payload resend.WebhookPayload
		if err := json.Unmarshal(body, &payload); err != nil {
			log.Printf("Error parsing JSON: %v", err)
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		log.Printf("✓ Webhook verified successfully!")
		log.Printf("Event Type: %v", payload.Type)
		switch t := payload.Data.(type) {
		case *resend.WebhookPayloadContactCreated, *resend.WebhookPayloadContactDeleted, *resend.WebhookPayloadContactUpdated:
			fmt.Printf("%s: ID: %s, Audience ID: %s, Segment IDs: %v\n", payload.Type, t.ID, t.AudienceID, t.SegmentIDs)
		case *resend.WebhookPayloadContactTopicsUpdated:
			fmt.Printf("%s: Topic IDs: ", payload.Type)
			for _, topic := range t.Topics {
				fmt.Printf("%s ", topic.ID)
			}
			fmt.Println()
		case *resend.WebhookPayloadDomainCreated, *resend.WebhookPayloadDomainDeleted, *resend.WebhookPayloadDomainUpdated:
			fmt.Printf("%s: ID: %s\n", payload.Type, t.ID)
		case *resend.WebhookPayloadEmailBounced:
			fmt.Printf("%s: Broadcast ID: %s, Email ID: %s, Message ID: %s, Template ID: %s\n", payload.Type, t.BroadcastID, t.EmailID, t.MessageID, t.TemplateID)
		case *resend.WebhookPayloadEmailClicked:
			fmt.Printf("%s: Broadcast ID: %s, Email ID: %s, Message ID: %s, Template ID: %s\n", payload.Type, t.BroadcastID, t.EmailID, t.MessageID, t.TemplateID)
		case *resend.WebhookPayloadEmailComplained, *resend.WebhookPayloadEmailDelivered, *resend.WebhookPayloadEmailDeliveryDelayed:
			fmt.Printf("%s: Broadcast ID: %s, Email ID: %s, Message ID: %s, Template ID: %s\n", payload.Type, t.BroadcastID, t.EmailID, t.MessageID, t.TemplateID)
		case *resend.WebhookPayloadEmailFailed:
			fmt.Printf("%s: Broadcast ID: %s, Email ID: %s, Template ID: %s\n", payload.Type, t.BroadcastID, t.EmailID, t.TemplateID)
		case *resend.WebhookPayloadEmailOpened:
			fmt.Printf("%s: Broadcast ID: %s, Email ID: %s, Message ID: %s, Template ID: %s\n", payload.Type, t.BroadcastID, t.EmailID, t.MessageID, t.TemplateID)
		case *resend.WebhookPayloadEmailReceived:
			fmt.Printf("%s: Email ID: %s, Message ID: %s\n", payload.Type, t.EmailID, t.MessageID)
		case *resend.WebhookPayloadEmailScheduled:
			fmt.Printf("%s: Broadcast ID: %s, Email ID: %s, Template ID: %s\n", payload.Type, t.BroadcastID, t.EmailID, t.TemplateID)
		case *resend.WebhookPayloadEmailSent:
			fmt.Printf("%s: Broadcast ID: %s, Email ID: %s, Message ID: %s, Template ID: %s\n", payload.Type, t.BroadcastID, t.EmailID, t.MessageID, t.TemplateID)
		case *resend.WebhookPayloadEmailSuppressed:
			fmt.Printf("%s: Broadcast ID: %s, Email ID: %s, Template ID: %s\n", payload.Type, t.BroadcastID, t.EmailID, t.TemplateID)
		case *resend.WebhookPayloadInboxCreated, *resend.WebhookPayloadInboxDeleted:
			fmt.Printf("%s: Inbox ID: %s\n", payload.Type, t.InboxID)
		case *resend.WebhookPayloadInboxDraftCreated, *resend.WebhookPayloadInboxDraftUpdated, *resend.WebhookPayloadInboxEmailReceived, *resend.WebhookPayloadInboxEmailSent, *resend.WebhookPayloadInboxThreadCreated:
			fmt.Printf("%s: Inbox ID: %s, Thread ID: %s, Draft ID: %s\n", payload.Type, t.InboxID, t.ThreadID, t.DraftID)
		case *resend.WebhookPayloadInboxDraftDeleted:
			fmt.Printf("%s: Inbox ID: %s, Thread ID: %s, Draft ID: %s\n", payload.Type, t.InboxID, t.ThreadID, t.DraftID)
		case *resend.WebhookPayloadInboxDraftSent:
			fmt.Printf("%s: Inbox ID: %s, Thread ID: %s, Email ID: %s, Draft ID: %s\n", payload.Type, t.InboxID, t.ThreadID, t.EmailID, t.DraftID)
		case *resend.WebhookPayloadInboxThreadAssigned:
			fmt.Printf("%s: Inbox ID: %s, Thread ID: %s\n", payload.Type, t.InboxID, t.ThreadID)
		case *resend.WebhookPayloadInboxThreadFolderUpdated:
			fmt.Printf("%s: Inbox ID: %s, Thread ID: %s\n", payload.Type, t.InboxID, t.ThreadID)
		case *resend.WebhookPayloadInboxThreadLabelsUpdated:
			fmt.Printf("%s: Inbox ID: %s, Thread ID: %s\n", payload.Type, t.InboxID, t.ThreadID)
		case *resend.WebhookPayloadInboxThreadUnassigned:
			fmt.Printf("%s: Inbox ID: %s, Thread ID: %s\n", payload.Type, t.InboxID, t.ThreadID)
		case *resend.WebhookPayloadInboxUpdated, *resend.WebhookPayloadInboxUpdated:
			fmt.Printf("%s: Inbox ID: %s\n", payload.Type, t.InboxID)
		case *resend.WebhookPayloadSuppressionAdded:
			fmt.Printf("%s: ID: %s, Source ID: %s\n", payload.Type, t.ID, t.SourceID)
		case *resend.WebhookPayloadSuppressionRemoved:
			fmt.Printf("%s: ID: %s, Source ID: %v\n", payload.Type, t.ID, t.SourceID)
		case *resend.WebhookPayloadTopicCreated, *resend.WebhookPayloadTopicDeleted, *resend.WebhookPayloadTopicUpdated:
			fmt.Printf("%s: ID: %s\n", payload.Type, t.ID)
		default:
			fmt.Printf("%s: Unknown type: %T\n", payload.Type, t)
		}

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
