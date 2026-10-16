package resend

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testThreadId = "4d8e2a1c-9b3f-4c6d-8a21-3e5f7c9c0d12"

const testInboxThreadJSON = `{
	"id": "4d8e2a1c-9b3f-4c6d-8a21-3e5f7c9c0d12",
	"subject": "Refund for order 1041",
	"from": "ada@example.org",
	"to": ["support@example.com"],
	"cc": [],
	"bcc": [],
	"labels": [{"id": "7f9c1d2e-4a6b-4c3d-8e15-9b0a7c6d5e34", "name": "Urgent", "color": "#E93D82"}],
	"email_count": 2,
	"has_attachment": true,
	"has_draft": false,
	"read": false,
	"received_at": "2026-08-05T14:03:11.229Z",
	"folder": "inbox"`

func TestListInboxThreads(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId+"/threads", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		query := r.URL.Query()
		assert.Equal(t, "inbox,sent", query.Get("folders"))
		assert.Equal(t, testLabelId, query.Get("labels"))
		assert.Equal(t, "false", query.Get("read"))
		assert.Equal(t, "5", query.Get("limit"))
		assert.Equal(t, testThreadId, query.Get("before"))

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "list", "has_more": true, "data": [%s}]}`, testInboxThreadJSON)
	})

	read := false
	limit := 5
	before := testThreadId
	resp, err := client.Inboxes.Threads.List(testInboxId, &ListInboxThreadsOptions{
		Folders: []InboxFolder{InboxFolderInbox, InboxFolderSent},
		Labels:  []string{testLabelId},
		Read:    &read,
		Limit:   &limit,
		Before:  &before,
	})
	assert.NoError(t, err)
	assert.True(t, resp.HasMore)
	assert.Len(t, resp.Data, 1)
	thread := resp.Data[0]
	assert.Equal(t, testThreadId, thread.Id)
	assert.Equal(t, "Refund for order 1041", *thread.Subject)
	assert.Equal(t, "ada@example.org", *thread.From)
	assert.Equal(t, []string{"support@example.com"}, thread.To)
	assert.Equal(t, "#E93D82", thread.Labels[0].Color)
	assert.Equal(t, 2, thread.EmailCount)
	assert.True(t, thread.HasAttachment)
	assert.Equal(t, InboxFolderInbox, thread.Folder)
}

func TestListInboxThreadsWithoutOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId+"/threads", func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"object": "list", "has_more": false, "data": []}`)
	})

	_, err := client.Inboxes.Threads.ListWithContext(context.Background(), testInboxId, nil)
	assert.NoError(t, err)

	_, err = client.Inboxes.Threads.List("", nil)
	assert.EqualError(t, err, "[ERROR]: Inbox ID is required")
}

func TestSearchInboxThreads(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId+"/threads/search", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		query := r.URL.Query()
		assert.Equal(t, `"refund" -draft`, query.Get("query"))
		assert.Equal(t, "isabella@example.com,carolina@example.com", query.Get("from"))
		assert.Equal(t, "support@acme.com", query.Get("to"))
		assert.Equal(t, "billing@acme.com", query.Get("cc"))
		assert.Equal(t, "records@acme.com", query.Get("bcc"))
		assert.Equal(t, "true", query.Get("has_attachment"))
		assert.Equal(t, "2026-09-01", query.Get("start_date"))
		assert.Equal(t, "2026-09-30", query.Get("end_date"))
		assert.Equal(t, "archive", query.Get("folders"))
		assert.Equal(t, testLabelId, query.Get("labels"))
		assert.Equal(t, "true", query.Get("read"))
		assert.Equal(t, "20", query.Get("limit"))
		assert.Equal(t, testThreadId, query.Get("after"))

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "list", "has_more": false, "data": [%s,
			"matched_email_id": "e8a3f1c6-2b9d-4e5a-a7c4-9d1b6f3e2a80"
		}]}`, testInboxThreadJSON)
	})

	hasAttachment := true
	read := true
	limit := 20
	after := testThreadId
	resp, err := client.Inboxes.Threads.Search(testInboxId, &SearchInboxThreadsOptions{
		Query:         `"refund" -draft`,
		From:          []string{"isabella@example.com", "carolina@example.com"},
		To:            []string{"support@acme.com"},
		Cc:            []string{"billing@acme.com"},
		Bcc:           []string{"records@acme.com"},
		HasAttachment: &hasAttachment,
		StartDate:     "2026-09-01",
		EndDate:       "2026-09-30",
		Folders:       []InboxFolder{InboxFolderArchive},
		Labels:        []string{testLabelId},
		Read:          &read,
		Limit:         &limit,
		After:         &after,
	})
	assert.NoError(t, err)
	assert.Len(t, resp.Data, 1)
	result := resp.Data[0]
	assert.Equal(t, testThreadId, result.Id)
	assert.Equal(t, "e8a3f1c6-2b9d-4e5a-a7c4-9d1b6f3e2a80", *result.MatchedEmailId)
}

func TestSearchInboxThreadsNullMatchedEmail(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId+"/threads/search", func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "list", "has_more": false, "data": [%s, "matched_email_id": null}]}`, testInboxThreadJSON)
	})

	resp, err := client.Inboxes.Threads.SearchWithContext(context.Background(), testInboxId, nil)
	assert.NoError(t, err)
	assert.Nil(t, resp.Data[0].MatchedEmailId)
}

func TestGetInboxThread(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId+"/threads/"+testThreadId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"object": "inbox_thread",
			"id": "%s",
			"subject": null,
			"folder": "archive",
			"labels": [],
			"read": true
		}`, testThreadId)
	})

	resp, err := client.Inboxes.Threads.Get(testInboxId, testThreadId)
	assert.NoError(t, err)
	assert.Equal(t, "inbox_thread", resp.Object)
	assert.Nil(t, resp.Subject)
	assert.Equal(t, InboxFolderArchive, resp.Folder)
	assert.True(t, resp.Read)

	_, err = client.Inboxes.Threads.GetWithContext(context.Background(), testInboxId, "")
	assert.EqualError(t, err, "[ERROR]: Thread ID is required")
}

func TestUpdateInboxThread(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId+"/threads/"+testThreadId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		assert.Equal(t, map[string]any{
			"read":   false,
			"folder": "archive",
		}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"object": "inbox_thread",
			"id": "%s",
			"subject": "Refund for order 1041",
			"folder": "archive",
			"labels": [{"id": "%s", "name": "Urgent", "color": "#E93D82"}],
			"read": false
		}`, testThreadId, testLabelId)
	})

	read := false
	resp, err := client.Inboxes.Threads.Update(testInboxId, testThreadId, &UpdateInboxThreadRequest{
		Read:   &read,
		Folder: InboxMoveThreadFolderArchive,
	})
	assert.NoError(t, err)
	assert.Equal(t, InboxFolderArchive, resp.Folder)
	assert.Equal(t, testLabelId, resp.Labels[0].Id)
}

func TestUpdateInboxThreadValidation(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Inboxes.Threads.Update(testInboxId, testThreadId, &UpdateInboxThreadRequest{})
	assert.EqualError(t, err, "[ERROR]: At least one of Read or Folder is required")

	_, err = client.Inboxes.Threads.Update(testInboxId, testThreadId, nil)
	assert.EqualError(t, err, "[ERROR]: At least one of Read or Folder is required")

	_, err = client.Inboxes.Threads.Update("", testThreadId, nil)
	assert.EqualError(t, err, "[ERROR]: Inbox ID is required")
}

func TestRemoveInboxThread(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId+"/threads/"+testThreadId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "inbox_thread", "id": "%s", "deleted": true}`, testThreadId)
	})

	resp, err := client.Inboxes.Threads.Remove(testInboxId, testThreadId)
	assert.NoError(t, err)
	assert.Equal(t, "inbox_thread", resp.Object)
	assert.True(t, resp.Deleted)

	_, err = client.Inboxes.Threads.RemoveWithContext(context.Background(), testInboxId, "")
	assert.EqualError(t, err, "[ERROR]: Thread ID is required")
}
