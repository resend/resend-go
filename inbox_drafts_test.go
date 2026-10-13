package resend

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testDraftId = "c3a1e8b4-2d5f-4a7c-9e10-6b8d7f5a4c32"

func draftsPath(segments ...string) string {
	path := "/inboxes/" + testInboxId + "/drafts"
	for _, segment := range segments {
		path += "/" + segment
	}
	return path
}

func TestListInboxDrafts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(draftsPath(), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, testDraftId, r.URL.Query().Get("after"))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"object": "list",
			"has_more": false,
			"data": [{
				"id": "%s",
				"type": "reply",
				"to": null,
				"cc": [],
				"bcc": [],
				"subject": null,
				"snippet": "Refund issued for order 1041.",
				"thread_id": "%s",
				"reply_to_email_id": "%s",
				"updated_at": "2026-08-05T14:12:04.110Z",
				"attachments": []
			}]
		}`, testDraftId, testThreadId, testThreadEmailId)
	})

	after := testDraftId
	resp, err := client.Inboxes.Drafts.List(testInboxId, &ListOptions{After: &after})
	assert.NoError(t, err)
	draft := resp.Data[0]
	assert.Equal(t, InboxDraftTypeReply, draft.Type)
	assert.Nil(t, draft.To)
	assert.Equal(t, "Refund issued for order 1041.", *draft.Snippet)
	assert.Equal(t, testThreadId, *draft.ThreadId)
	assert.Empty(t, draft.Attachments)

	_, err = client.Inboxes.Drafts.ListWithContext(context.Background(), "", nil)
	assert.EqualError(t, err, "[ERROR]: Inbox ID is required")
}

func TestCreateInboxDraft(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(draftsPath(), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assert.Equal(t, map[string]any{
			"to":                []any{},
			"text":              "Refund issued for order 1041.",
			"thread_id":         testThreadId,
			"reply_to_email_id": testThreadEmailId,
		}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"object": "inbox_draft", "id": "%s"}`, testDraftId)
	})

	resp, err := client.Inboxes.Drafts.Create(testInboxId, &CreateInboxDraftRequest{
		To:             []string{},
		Text:           "Refund issued for order 1041.",
		ThreadId:       testThreadId,
		ReplyToEmailId: testThreadEmailId,
	})
	assert.NoError(t, err)
	assert.Equal(t, "inbox_draft", resp.Object)
	assert.Equal(t, testDraftId, resp.Id)
}

func TestCreateInboxDraftOmitsNilRecipients(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(draftsPath(), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, map[string]any{
			"to":      []any{"ada@example.org"},
			"subject": "Refund for order 1041",
			"html":    "<p>Refund issued.</p>",
		}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "inbox_draft", "id": "%s"}`, testDraftId)
	})

	_, err := client.Inboxes.Drafts.CreateWithContext(context.Background(), testInboxId, &CreateInboxDraftRequest{
		To:      []string{"ada@example.org"},
		Subject: "Refund for order 1041",
		Html:    "<p>Refund issued.</p>",
	})
	assert.NoError(t, err)
}

func TestCreateInboxDraftValidation(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Inboxes.Drafts.Create(testInboxId, &CreateInboxDraftRequest{Text: "Hi", ThreadId: testThreadId})
	assert.EqualError(t, err, "[ERROR]: ThreadId and ReplyToEmailId must be provided together")

	_, err = client.Inboxes.Drafts.Create(testInboxId, &CreateInboxDraftRequest{To: []string{}, Subject: " ", ThreadId: testThreadId, ReplyToEmailId: testThreadEmailId})
	assert.EqualError(t, err, "[ERROR]: At least one of To, Cc, Bcc, Subject, Text, Html or Attachments must be non-empty")

	_, err = client.Inboxes.Drafts.Create(testInboxId, nil)
	assert.EqualError(t, err, "[ERROR]: Draft params are required")

	_, err = client.Inboxes.Drafts.Create("", &CreateInboxDraftRequest{Text: "Hi"})
	assert.EqualError(t, err, "[ERROR]: Inbox ID is required")
}

func TestGetInboxDraft(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(draftsPath(testDraftId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"object": "inbox_draft",
			"id": "%s",
			"type": "standalone",
			"to": ["ada@example.org"],
			"cc": [],
			"bcc": [],
			"subject": "Refund for order 1041",
			"html": null,
			"text": "Refund issued for order 1041.",
			"thread_id": null,
			"reply_to_email_id": null,
			"email_id": null,
			"created_at": "2026-08-05T14:12:04.110Z",
			"updated_at": "2026-08-05T14:12:04.110Z",
			"attachments": [{"id": "5c1e9d2a-7b3f-4e8a-9c6d-2f0a1b3c4d5e", "filename": "receipt.pdf", "content_type": "application/pdf", "size": 20841}]
		}`, testDraftId)
	})

	resp, err := client.Inboxes.Drafts.Get(testInboxId, testDraftId)
	assert.NoError(t, err)
	assert.Equal(t, InboxDraftTypeStandalone, resp.Type)
	assert.Equal(t, []string{"ada@example.org"}, resp.To)
	assert.Nil(t, resp.Html)
	assert.Equal(t, "Refund issued for order 1041.", *resp.Text)
	assert.Nil(t, resp.ThreadId)
	assert.Nil(t, resp.EmailId)
	assert.Equal(t, []InboxDraftAttachment{{
		Id:          "5c1e9d2a-7b3f-4e8a-9c6d-2f0a1b3c4d5e",
		Filename:    "receipt.pdf",
		ContentType: "application/pdf",
		Size:        20841,
	}}, resp.Attachments)

	_, err = client.Inboxes.Drafts.GetWithContext(context.Background(), testInboxId, "")
	assert.EqualError(t, err, "[ERROR]: Draft ID is required")
}

func TestUpdateInboxDraft(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(draftsPath(testDraftId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		assert.Equal(t, map[string]any{
			"cc":      []any{},
			"subject": "Refund issued for order 1041",
			"html":    nil,
			"text":    "Updated body.",
		}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "inbox_draft", "id": "%s", "subject": "Refund issued for order 1041", "html": null, "text": "Updated body."}`, testDraftId)
	})

	subject := "Refund issued for order 1041"
	text := "Updated body."
	params := &UpdateInboxDraftRequest{
		Cc:      []string{},
		Subject: &subject,
		Text:    &text,
	}
	params.ClearHtml()

	resp, err := client.Inboxes.Drafts.Update(testInboxId, testDraftId, params)
	assert.NoError(t, err)
	assert.Equal(t, testDraftId, resp.Id)
	assert.Equal(t, "Updated body.", *resp.Text)
	assert.Nil(t, resp.Html)
}

func TestUpdateInboxDraftClearOnly(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(draftsPath(testDraftId), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, map[string]any{"subject": nil, "text": nil}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "inbox_draft", "id": "%s"}`, testDraftId)
	})

	params := &UpdateInboxDraftRequest{}
	params.ClearSubject()
	params.ClearText()

	_, err := client.Inboxes.Drafts.UpdateWithContext(context.Background(), testInboxId, testDraftId, params)
	assert.NoError(t, err)
}

func TestUpdateInboxDraftValidation(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Inboxes.Drafts.Update(testInboxId, testDraftId, &UpdateInboxDraftRequest{})
	assert.EqualError(t, err, "[ERROR]: At least one of To, Cc, Bcc, Subject, Text, Html or Attachments is required")

	_, err = client.Inboxes.Drafts.Update(testInboxId, testDraftId, nil)
	assert.EqualError(t, err, "[ERROR]: At least one of To, Cc, Bcc, Subject, Text, Html or Attachments is required")
}

func TestRemoveInboxDraft(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(draftsPath(testDraftId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "inbox_draft", "id": "%s", "deleted": true}`, testDraftId)
	})

	resp, err := client.Inboxes.Drafts.Remove(testInboxId, testDraftId)
	assert.NoError(t, err)
	assert.True(t, resp.Deleted)

	_, err = client.Inboxes.Drafts.RemoveWithContext(context.Background(), "", testDraftId)
	assert.EqualError(t, err, "[ERROR]: Inbox ID is required")
}

func TestSendInboxDraft(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(draftsPath(testDraftId, "send"), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assert.Equal(t, int64(0), r.ContentLength)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"id": "%s"}`, testSentEmailId)
	})

	resp, err := client.Inboxes.Drafts.Send(testInboxId, testDraftId)
	assert.NoError(t, err)
	assert.Equal(t, testSentEmailId, resp.Id)

	_, err = client.Inboxes.Drafts.SendWithContext(context.Background(), testInboxId, "")
	assert.EqualError(t, err, "[ERROR]: Draft ID is required")
}

func TestCreateInboxDraftWithOnlyAttachments(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(draftsPath(), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assert.Equal(t, map[string]any{
			"attachments": []any{map[string]any{
				"content":      []any{float64(104), float64(105)},
				"filename":     "note.txt",
				"content_type": "text/plain",
			}},
		}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"object": "inbox_draft", "id": "%s"}`, testDraftId)
	})

	_, err := client.Inboxes.Drafts.Create(testInboxId, &CreateInboxDraftRequest{
		Attachments: []*Attachment{{Content: []byte("hi"), Filename: "note.txt", ContentType: "text/plain"}},
	})
	assert.NoError(t, err)
}

func TestUpdateInboxDraftAttachments(t *testing.T) {
	setup()
	defer teardown()

	var bodies []map[string]any
	mux.HandleFunc(draftsPath(testDraftId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		bodies = append(bodies, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "inbox_draft", "id": "%s", "attachments": []}`, testDraftId)
	})

	_, err := client.Inboxes.Drafts.Update(testInboxId, testDraftId, &UpdateInboxDraftRequest{
		Attachments: []*Attachment{{Content: []byte("hi"), Filename: "note.txt", ContentType: "text/plain"}},
	})
	assert.NoError(t, err)

	resp, err := client.Inboxes.Drafts.Update(testInboxId, testDraftId, &UpdateInboxDraftRequest{
		Attachments: []*Attachment{},
	})
	assert.NoError(t, err)
	assert.Empty(t, resp.Attachments)

	assert.Equal(t, []map[string]any{
		{"attachments": []any{map[string]any{
			"content":      []any{float64(104), float64(105)},
			"filename":     "note.txt",
			"content_type": "text/plain",
		}}},
		{"attachments": []any{}},
	}, bodies)
}
