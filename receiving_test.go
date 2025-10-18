package resend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetReceivedEmail(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/receiving/8136d3fb-0439-4b09-b939-b8436a3524b6", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "inbound",
			"id": "8136d3fb-0439-4b09-b939-b8436a3524b6",
			"to": ["delivered@resend.dev"],
			"from": "Acme <onboarding@resend.dev>",
			"created_at": "2023-04-03T22:13:42.674981+00:00",
			"subject": "Hello World",
			"html": "Congrats on sending your <strong>first email</strong>!",
			"text": "Congrats on sending your first email!",
			"bcc": [],
			"cc": ["cc@example.com"],
			"reply_to": ["reply@example.com"],
			"headers": {
				"X-Custom-Header": "value"
			},
			"attachments": [
				{
					"id": "2a0c9ce0-3112-4728-976e-47ddcd16a318",
					"filename": "avatar.png",
					"content_type": "image/png",
					"content_disposition": "inline",
					"content_id": "img001"
				}
			]
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Receiving.Get("8136d3fb-0439-4b09-b939-b8436a3524b6")
	if err != nil {
		t.Errorf("Receiving.Get returned error: %v", err)
	}
	assert.Equal(t, "8136d3fb-0439-4b09-b939-b8436a3524b6", resp.Id)
	assert.Equal(t, "inbound", resp.Object)
	assert.Equal(t, "Acme <onboarding@resend.dev>", resp.From)
	assert.Equal(t, "Hello World", resp.Subject)
	assert.Equal(t, "Congrats on sending your <strong>first email</strong>!", resp.Html)
	assert.Equal(t, "Congrats on sending your first email!", resp.Text)
	assert.Equal(t, 1, len(resp.To))
	assert.Equal(t, "delivered@resend.dev", resp.To[0])
	assert.Equal(t, 1, len(resp.Cc))
	assert.Equal(t, "cc@example.com", resp.Cc[0])
	assert.Equal(t, 1, len(resp.ReplyTo))
	assert.Equal(t, "reply@example.com", resp.ReplyTo[0])
	assert.Equal(t, 1, len(resp.Headers))
	assert.Equal(t, "value", resp.Headers["X-Custom-Header"])
	assert.Equal(t, 1, len(resp.Attachments))
	assert.Equal(t, "2a0c9ce0-3112-4728-976e-47ddcd16a318", resp.Attachments[0].Id)
	assert.Equal(t, "avatar.png", resp.Attachments[0].Filename)
	assert.Equal(t, "image/png", resp.Attachments[0].ContentType)
	assert.Equal(t, "inline", resp.Attachments[0].ContentDisposition)
	assert.Equal(t, "img001", resp.Attachments[0].ContentId)
}

func TestGetReceivedEmailWithNullFields(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/receiving/null-fields-id", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "inbound",
			"id": "null-fields-id",
			"to": ["delivered@resend.dev"],
			"from": "sender@example.com",
			"created_at": "2023-04-03T22:13:42.674981+00:00",
			"subject": "Test Subject",
			"html": "",
			"text": "",
			"bcc": [],
			"cc": [],
			"reply_to": [],
			"headers": {},
			"attachments": []
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Receiving.Get("null-fields-id")
	if err != nil {
		t.Errorf("Receiving.Get returned error: %v", err)
	}
	assert.Equal(t, "null-fields-id", resp.Id)
	assert.Equal(t, "", resp.Html)
	assert.Equal(t, "", resp.Text)
	assert.Equal(t, 0, len(resp.Bcc))
	assert.Equal(t, 0, len(resp.Cc))
	assert.Equal(t, 0, len(resp.ReplyTo))
	assert.Equal(t, 0, len(resp.Headers))
	assert.Equal(t, 0, len(resp.Attachments))
}

func TestListReceivedEmails(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/receiving", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &ListReceivedEmailsResponse{
			Object:  "list",
			HasMore: true,
			Data: []ListReceivedEmail{
				{
					Id:        "1",
					To:        []string{"recipient@example.com"},
					From:      "sender@example.com",
					CreatedAt: "2024-01-01T00:00:00Z",
					Subject:   "Test Email 1",
					Bcc:       []string{},
					Cc:        []string{"cc@example.com"},
					ReplyTo:   []string{},
					Attachments: []ReceivedAttachment{
						{
							Id:                 "att1",
							Filename:           "file.pdf",
							ContentType:        "application/pdf",
							ContentDisposition: "attachment",
							ContentId:          "cid1",
						},
					},
				},
				{
					Id:          "2",
					To:          []string{"recipient2@example.com"},
					From:        "sender2@example.com",
					CreatedAt:   "2024-01-02T00:00:00Z",
					Subject:     "Test Email 2",
					Bcc:         []string{},
					Cc:          []string{},
					ReplyTo:     []string{},
					Attachments: []ReceivedAttachment{},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(ret); err != nil {
			panic(err)
		}
	})

	resp, err := client.Receiving.List()
	if err != nil {
		t.Errorf("Receiving.List returned error: %v", err)
	}

	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, true, resp.HasMore)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "1", resp.Data[0].Id)
	assert.Equal(t, "sender@example.com", resp.Data[0].From)
	assert.Equal(t, "Test Email 1", resp.Data[0].Subject)
	assert.Equal(t, 1, len(resp.Data[0].Attachments))
	assert.Equal(t, "att1", resp.Data[0].Attachments[0].Id)
	assert.Equal(t, "2", resp.Data[1].Id)
}

func TestListReceivedEmailsWithParameters(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/receiving", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		// Verify query parameters
		query := r.URL.Query()
		assert.Equal(t, "10", query.Get("limit"))
		assert.Equal(t, "cursor123", query.Get("after"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &ListReceivedEmailsResponse{
			Object:  "list",
			HasMore: false,
			Data: []ListReceivedEmail{
				{
					Id:          "1",
					To:          []string{"recipient@example.com"},
					From:        "sender@example.com",
					CreatedAt:   "2024-01-01T00:00:00Z",
					Subject:     "Test Email 1",
					Bcc:         []string{},
					Cc:          []string{},
					ReplyTo:     []string{},
					Attachments: []ReceivedAttachment{},
				},
			},
		}

		if err := json.NewEncoder(w).Encode(ret); err != nil {
			panic(err)
		}
	})

	limit := 10
	after := "cursor123"
	options := &ListOptions{
		Limit: &limit,
		After: &after,
	}

	resp, err := client.Receiving.ListWithOptions(context.Background(), options)
	if err != nil {
		t.Errorf("Receiving.ListWithOptions returned error: %v", err)
	}

	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, false, resp.HasMore)
	assert.Equal(t, 1, len(resp.Data))
}

func TestListReceivedEmailsEmpty(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/receiving", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &ListReceivedEmailsResponse{
			Object:  "list",
			HasMore: false,
			Data:    []ListReceivedEmail{},
		}

		if err := json.NewEncoder(w).Encode(ret); err != nil {
			panic(err)
		}
	})

	resp, err := client.Receiving.List()
	if err != nil {
		t.Errorf("Receiving.List returned error: %v", err)
	}

	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, false, resp.HasMore)
	assert.Equal(t, 0, len(resp.Data))
}
