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

func TestGetReceivedEmailAttachment(t *testing.T) {
	setup()
	defer teardown()

	emailID := "4ef9a417-02e9-4d39-ad75-9611e0fcc33c"
	attachmentID := "2a0c9ce0-3112-4728-976e-47ddcd16a318"

	mux.HandleFunc(fmt.Sprintf("/emails/receiving/%s/attachments/%s", emailID, attachmentID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "attachment",
			"data": {
				"id": "2a0c9ce0-3112-4728-976e-47ddcd16a318",
				"filename": "avatar.png",
				"content_type": "image/png",
				"content_disposition": "inline",
				"content_id": "img001",
				"download_url": "https://inbound-cdn.resend.com/4ef9a417-02e9-4d39-ad75-9611e0fcc33c/attachments/2a0c9ce0-3112-4728-976e-47ddcd16a318?some-params=example&signature=sig-123",
				"expires_at": "2025-10-17T14:29:41.521Z"
			}
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Receiving.GetAttachment(emailID, attachmentID)
	if err != nil {
		t.Errorf("Receiving.GetAttachment returned error: %v", err)
	}
	assert.Equal(t, "2a0c9ce0-3112-4728-976e-47ddcd16a318", resp.Id)
	assert.Equal(t, "avatar.png", resp.Filename)
	assert.Equal(t, "image/png", resp.ContentType)
	assert.Equal(t, "inline", resp.ContentDisposition)
	assert.Equal(t, "img001", resp.ContentId)
	assert.Equal(t, "https://inbound-cdn.resend.com/4ef9a417-02e9-4d39-ad75-9611e0fcc33c/attachments/2a0c9ce0-3112-4728-976e-47ddcd16a318?some-params=example&signature=sig-123", resp.DownloadUrl)
	assert.Equal(t, "2025-10-17T14:29:41.521Z", resp.ExpiresAt)
}

func TestGetReceivedEmailAttachmentWithContext(t *testing.T) {
	setup()
	defer teardown()

	emailID := "test-email-id"
	attachmentID := "test-attachment-id"

	mux.HandleFunc(fmt.Sprintf("/emails/receiving/%s/attachments/%s", emailID, attachmentID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "attachment",
			"data": {
				"id": "test-attachment-id",
				"filename": "document.pdf",
				"content_type": "application/pdf",
				"content_disposition": "attachment",
				"content_id": "doc001",
				"download_url": "https://inbound-cdn.resend.com/test-email-id/attachments/test-attachment-id",
				"expires_at": "2025-10-18T12:00:00.000Z"
			}
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	resp, err := client.Receiving.GetAttachmentWithContext(ctx, emailID, attachmentID)
	if err != nil {
		t.Errorf("Receiving.GetAttachmentWithContext returned error: %v", err)
	}
	assert.Equal(t, "test-attachment-id", resp.Id)
	assert.Equal(t, "document.pdf", resp.Filename)
	assert.Equal(t, "application/pdf", resp.ContentType)
	assert.Equal(t, "attachment", resp.ContentDisposition)
	assert.Equal(t, "doc001", resp.ContentId)
	assert.NotEmpty(t, resp.DownloadUrl)
	assert.NotEmpty(t, resp.ExpiresAt)
}

func TestListReceivedEmailAttachments(t *testing.T) {
	setup()
	defer teardown()

	emailID := "4ef9a417-02e9-4d39-ad75-9611e0fcc33c"

	mux.HandleFunc(fmt.Sprintf("/emails/receiving/%s/attachments", emailID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &ListReceivedEmailAttachmentsResponse{
			Object:  "list",
			HasMore: false,
			Data: []ReceivedEmailAttachment{
				{
					Id:                 "2a0c9ce0-3112-4728-976e-47ddcd16a318",
					Filename:           "avatar.png",
					ContentType:        "image/png",
					ContentDisposition: "inline",
					ContentId:          "img001",
					DownloadUrl:        "https://inbound-cdn.resend.com/4ef9a417-02e9-4d39-ad75-9611e0fcc33c/attachments/2a0c9ce0-3112-4728-976e-47ddcd16a318?some-params=example&signature=sig-123",
					ExpiresAt:          "2025-10-17T14:29:41.521Z",
				},
				{
					Id:                 "3b1d0df1-4223-5839-a87f-58eedd17b429",
					Filename:           "document.pdf",
					ContentType:        "application/pdf",
					ContentDisposition: "attachment",
					ContentId:          "doc001",
					DownloadUrl:        "https://inbound-cdn.resend.com/4ef9a417-02e9-4d39-ad75-9611e0fcc33c/attachments/3b1d0df1-4223-5839-a87f-58eedd17b429?some-params=example&signature=sig-456",
					ExpiresAt:          "2025-10-17T14:29:41.521Z",
				},
			},
		}

		if err := json.NewEncoder(w).Encode(ret); err != nil {
			panic(err)
		}
	})

	resp, err := client.Receiving.ListAttachments(emailID)
	if err != nil {
		t.Errorf("Receiving.ListAttachments returned error: %v", err)
	}

	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, false, resp.HasMore)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "2a0c9ce0-3112-4728-976e-47ddcd16a318", resp.Data[0].Id)
	assert.Equal(t, "avatar.png", resp.Data[0].Filename)
	assert.Equal(t, "image/png", resp.Data[0].ContentType)
	assert.Equal(t, "inline", resp.Data[0].ContentDisposition)
	assert.Equal(t, "img001", resp.Data[0].ContentId)
	assert.NotEmpty(t, resp.Data[0].DownloadUrl)
	assert.NotEmpty(t, resp.Data[0].ExpiresAt)
	assert.Equal(t, "3b1d0df1-4223-5839-a87f-58eedd17b429", resp.Data[1].Id)
}

func TestListReceivedEmailAttachmentsWithParameters(t *testing.T) {
	setup()
	defer teardown()

	emailID := "test-email-id"

	mux.HandleFunc(fmt.Sprintf("/emails/receiving/%s/attachments", emailID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		// Verify query parameters
		query := r.URL.Query()
		assert.Equal(t, "5", query.Get("limit"))
		assert.Equal(t, "cursor456", query.Get("after"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &ListReceivedEmailAttachmentsResponse{
			Object:  "list",
			HasMore: true,
			Data: []ReceivedEmailAttachment{
				{
					Id:                 "attachment-1",
					Filename:           "file1.jpg",
					ContentType:        "image/jpeg",
					ContentDisposition: "attachment",
					ContentId:          "img1",
					DownloadUrl:        "https://example.com/file1.jpg",
					ExpiresAt:          "2025-10-18T12:00:00.000Z",
				},
			},
		}

		if err := json.NewEncoder(w).Encode(ret); err != nil {
			panic(err)
		}
	})

	limit := 5
	after := "cursor456"
	options := &ListOptions{
		Limit: &limit,
		After: &after,
	}

	resp, err := client.Receiving.ListAttachmentsWithOptions(context.Background(), emailID, options)
	if err != nil {
		t.Errorf("Receiving.ListAttachmentsWithOptions returned error: %v", err)
	}

	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, true, resp.HasMore)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "attachment-1", resp.Data[0].Id)
}

func TestListReceivedEmailAttachmentsEmpty(t *testing.T) {
	setup()
	defer teardown()

	emailID := "email-no-attachments"

	mux.HandleFunc(fmt.Sprintf("/emails/receiving/%s/attachments", emailID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &ListReceivedEmailAttachmentsResponse{
			Object:  "list",
			HasMore: false,
			Data:    []ReceivedEmailAttachment{},
		}

		if err := json.NewEncoder(w).Encode(ret); err != nil {
			panic(err)
		}
	})

	resp, err := client.Receiving.ListAttachments(emailID)
	if err != nil {
		t.Errorf("Receiving.ListAttachments returned error: %v", err)
	}

	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, false, resp.HasMore)
	assert.Equal(t, 0, len(resp.Data))
}

func TestListReceivedEmailAttachmentsWithContext(t *testing.T) {
	setup()
	defer teardown()

	emailID := "context-test-id"

	mux.HandleFunc(fmt.Sprintf("/emails/receiving/%s/attachments", emailID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &ListReceivedEmailAttachmentsResponse{
			Object:  "list",
			HasMore: false,
			Data: []ReceivedEmailAttachment{
				{
					Id:                 "ctx-att-1",
					Filename:           "context-test.txt",
					ContentType:        "text/plain",
					ContentDisposition: "attachment",
					ContentId:          "txt1",
					DownloadUrl:        "https://example.com/context-test.txt",
					ExpiresAt:          "2025-10-18T12:00:00.000Z",
				},
			},
		}

		if err := json.NewEncoder(w).Encode(ret); err != nil {
			panic(err)
		}
	})

	ctx := context.Background()
	resp, err := client.Receiving.ListAttachmentsWithContext(ctx, emailID)
	if err != nil {
		t.Errorf("Receiving.ListAttachmentsWithContext returned error: %v", err)
	}

	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "ctx-att-1", resp.Data[0].Id)
	assert.Equal(t, "context-test.txt", resp.Data[0].Filename)
}
