package resend

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testThreadEmailId = "5b1a9f47-2c8d-4e6f-9a03-1d2e3f4a5b60"
	testSentEmailId   = "6a0c8e58-3d9e-4f7a-8b14-2e3f4a5b6c71"
)

const testInboxThreadEmailJSON = `{
	"id": "5b1a9f47-2c8d-4e6f-9a03-1d2e3f4a5b60",
	"direction": "inbound",
	"from": "Ada Lovelace <ada@example.org>",
	"to": ["support@acme.com"],
	"cc": [],
	"bcc": [],
	"reply_to": ["replies@example.org"],
	"subject": "Refund for order 1041",
	"message_id": "<1041@example.org>",
	"html": "<p>Could I get a refund for order 1041?</p>",
	"text": null,
	"attachments": [{"id": "9d4f2b81-6c3a-4e7d-8b12-0a5c6d7e8f90", "filename": "receipt.pdf", "size": 20841}],
	"read": false,
	"received_at": "2026-08-05T14:03:11.229Z"
}`

func threadEmailsPath(segments ...string) string {
	path := "/inboxes/" + testInboxId + "/threads/" + testThreadId + "/emails"
	for _, segment := range segments {
		path += "/" + segment
	}
	return path
}

func TestListInboxThreadEmails(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(threadEmailsPath(), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, "2", r.URL.Query().Get("limit"))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "list", "has_more": true, "data": [%s]}`, testInboxThreadEmailJSON)
	})

	limit := 2
	resp, err := client.Inboxes.Threads.Emails.List(testInboxId, testThreadId, &ListOptions{Limit: &limit})
	assert.NoError(t, err)
	assert.True(t, resp.HasMore)
	email := resp.Data[0]
	assert.Equal(t, testThreadEmailId, email.Id)
	assert.Equal(t, InboxThreadEmailDirectionInbound, email.Direction)
	assert.Equal(t, []string{"replies@example.org"}, email.ReplyTo)
	assert.Equal(t, "<1041@example.org>", *email.MessageId)
	assert.Nil(t, email.Text)
	assert.Equal(t, "receipt.pdf", *email.Attachments[0].Filename)
	assert.Equal(t, 20841, *email.Attachments[0].Size)

	_, err = client.Inboxes.Threads.Emails.ListWithContext(context.Background(), testInboxId, "", nil)
	assert.EqualError(t, err, "[ERROR]: Thread ID is required")
}

func TestGetInboxThreadEmail(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(threadEmailsPath(testThreadEmailId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, testInboxThreadEmailJSON)
	})

	resp, err := client.Inboxes.Threads.Emails.Get(testInboxId, testThreadId, testThreadEmailId)
	assert.NoError(t, err)
	assert.Equal(t, "Ada Lovelace <ada@example.org>", resp.From)
	assert.Equal(t, "<p>Could I get a refund for order 1041?</p>", *resp.Html)

	_, err = client.Inboxes.Threads.Emails.GetWithContext(context.Background(), testInboxId, testThreadId, "")
	assert.EqualError(t, err, "[ERROR]: Email ID is required")
}

func TestReplyInboxThreadEmail(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(threadEmailsPath(testThreadEmailId, "reply"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assert.Equal(t, map[string]any{
			"cc":        []any{"billing@acme.com"},
			"text":      "Refund issued for order 1041.",
			"subject":   "Re: Refund for order 1041",
			"reply_all": true,
			"attachments": []any{map[string]any{
				"content":  []any{float64(104), float64(101), float64(108), float64(108), float64(111)},
				"filename": "note.txt",
			}},
		}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"id": "%s"}`, testSentEmailId)
	})

	resp, err := client.Inboxes.Threads.Emails.Reply(testInboxId, testThreadId, testThreadEmailId, &ReplyInboxThreadEmailRequest{
		Cc:          []string{"billing@acme.com"},
		Text:        "Refund issued for order 1041.",
		Subject:     "Re: Refund for order 1041",
		ReplyAll:    true,
		Attachments: []*Attachment{{Content: []byte("hello"), Filename: "note.txt"}},
	})
	assert.NoError(t, err)
	assert.Equal(t, testSentEmailId, resp.Id)
}

func TestReplyInboxThreadEmailValidation(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Inboxes.Threads.Emails.Reply(testInboxId, testThreadId, testThreadEmailId, &ReplyInboxThreadEmailRequest{})
	assert.EqualError(t, err, "[ERROR]: At least one of Html or Text is required")

	_, err = client.Inboxes.Threads.Emails.Reply(testInboxId, testThreadId, testThreadEmailId, &ReplyInboxThreadEmailRequest{Html: " ", Text: "\n"})
	assert.EqualError(t, err, "[ERROR]: At least one of Html or Text is required")

	_, err = client.Inboxes.Threads.Emails.ReplyWithContext(context.Background(), testInboxId, testThreadId, testThreadEmailId, nil)
	assert.EqualError(t, err, "[ERROR]: At least one of Html or Text is required")

	_, err = client.Inboxes.Threads.Emails.Reply("", testThreadId, testThreadEmailId, &ReplyInboxThreadEmailRequest{Text: "Hi"})
	assert.EqualError(t, err, "[ERROR]: Inbox ID is required")
}

func TestForwardInboxThreadEmail(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(threadEmailsPath(testThreadEmailId, "forward"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assert.Equal(t, map[string]any{
			"to":   []any{"colleague@example.org"},
			"html": "<p>Flagging this refund request for you.</p>",
			"attachments": []any{map[string]any{
				"path":     "https://example.org/invoice.pdf",
				"filename": "invoice.pdf",
			}},
		}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"id": "%s"}`, testSentEmailId)
	})

	resp, err := client.Inboxes.Threads.Emails.Forward(testInboxId, testThreadId, testThreadEmailId, &ForwardInboxThreadEmailRequest{
		To:          []string{"colleague@example.org"},
		Html:        "<p>Flagging this refund request for you.</p>",
		Attachments: []*Attachment{{Path: "https://example.org/invoice.pdf", Filename: "invoice.pdf"}},
	})
	assert.NoError(t, err)
	assert.Equal(t, testSentEmailId, resp.Id)
}

func TestForwardInboxThreadEmailValidation(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Inboxes.Threads.Emails.Forward(testInboxId, testThreadId, testThreadEmailId, &ForwardInboxThreadEmailRequest{})
	assert.EqualError(t, err, "[ERROR]: To is required")

	_, err = client.Inboxes.Threads.Emails.ForwardWithContext(context.Background(), testInboxId, "", testThreadEmailId, nil)
	assert.EqualError(t, err, "[ERROR]: Thread ID is required")
}
