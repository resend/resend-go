package resend

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailSentEventUnmarshal(t *testing.T) {
	payload := `{
		"type": "email.sent",
		"created_at": "2026-02-22T23:41:12.126Z",
		"data": {
			"broadcast_id": "8b146471-e88e-4322-86af-016cd36fd216",
			"created_at": "2026-02-22T23:41:11.894719+00:00",
			"email_id": "56761188-7520-42d8-8898-ff6fc54ce618",
			"message_id": "<111-222-333@email.example.com>",
			"from": "Acme <onboarding@resend.dev>",
			"to": ["delivered@resend.dev"],
			"subject": "Sending this example",
			"template_id": "43f68331-0622-4e15-8202-246a0388854b",
			"tags": {
				"category": "confirm_email"
			}
		}
	}`

	var event EmailSentEvent
	err := json.Unmarshal([]byte(payload), &event)
	assert.NoError(t, err)
	assert.Equal(t, EventEmailSent, event.Type)
	assert.Equal(t, "56761188-7520-42d8-8898-ff6fc54ce618", event.Data.EmailId)
	assert.Equal(t, "<111-222-333@email.example.com>", event.Data.MessageId)
	assert.Equal(t, "confirm_email", event.Data.Tags["category"])
}

func TestEmailBouncedEventUnmarshal(t *testing.T) {
	payload := `{
		"type": "email.bounced",
		"created_at": "2026-11-22T23:41:12.126Z",
		"data": {
			"created_at": "2026-11-22T23:41:11.894719+00:00",
			"email_id": "56761188-7520-42d8-8898-ff6fc54ce618",
			"message_id": "<111-222-333@email.example.com>",
			"from": "Acme <onboarding@resend.dev>",
			"to": ["delivered@resend.dev"],
			"subject": "Sending this example",
			"bounce": {
				"message": "The recipient's email address is on the suppression list.",
				"subType": "Suppressed",
				"type": "Permanent"
			}
		}
	}`

	var event EmailBouncedEvent
	err := json.Unmarshal([]byte(payload), &event)
	assert.NoError(t, err)
	assert.Equal(t, EventEmailBounced, event.Type)
	assert.Equal(t, "<111-222-333@email.example.com>", event.Data.MessageId)
	assert.Equal(t, "Suppressed", event.Data.Bounce.SubType)
}

func TestEmailClickedEventUnmarshal(t *testing.T) {
	payload := `{
		"type": "email.clicked",
		"created_at": "2026-11-22T23:41:12.126Z",
		"data": {
			"created_at": "2026-11-22T23:41:11.894719+00:00",
			"email_id": "56761188-7520-42d8-8898-ff6fc54ce618",
			"message_id": "<111-222-333@email.example.com>",
			"from": "Acme <onboarding@resend.dev>",
			"to": ["delivered@resend.dev"],
			"subject": "Sending this example",
			"click": {
				"ipAddress": "122.115.53.11",
				"link": "https://resend.com",
				"timestamp": "2026-11-24T05:00:57.163Z",
				"userAgent": "Mozilla/5.0"
			}
		}
	}`

	var event EmailClickedEvent
	err := json.Unmarshal([]byte(payload), &event)
	assert.NoError(t, err)
	assert.Equal(t, EventEmailClicked, event.Type)
	assert.Equal(t, "<111-222-333@email.example.com>", event.Data.MessageId)
	assert.Equal(t, "https://resend.com", event.Data.Click.Link)
}

func TestEmailReceivedEventUnmarshal(t *testing.T) {
	payload := `{
		"type": "email.received",
		"created_at": "2026-02-22T23:41:12.126Z",
		"data": {
			"email_id": "56761188-7520-42d8-8898-ff6fc54ce618",
			"created_at": "2026-02-22T23:41:11.894719+00:00",
			"from": "onboarding@resend.dev",
			"to": ["delivered@resend.dev"],
			"bcc": [],
			"cc": [],
			"received_for": ["forwarded@example.com"],
			"message_id": "<111-222-333@email.example.com>",
			"subject": "Sending this example",
			"attachments": [
				{
					"id": "2a0c9ce0-3112-4728-976e-47ddcd16a318",
					"filename": "avatar.png",
					"content_type": "image/png",
					"content_disposition": "inline",
					"content_id": "img001"
				}
			]
		}
	}`

	var event EmailReceivedEvent
	err := json.Unmarshal([]byte(payload), &event)
	assert.NoError(t, err)
	assert.Equal(t, EventEmailReceived, event.Type)
	assert.Equal(t, "<111-222-333@email.example.com>", event.Data.MessageId)
	assert.Equal(t, "forwarded@example.com", event.Data.ReceivedFor[0])
}
