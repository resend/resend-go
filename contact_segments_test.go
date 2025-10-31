package resend

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContactSegmentsAddWithContactId(t *testing.T) {
	setup()
	defer teardown()

	contactId := "479e3145-dd38-476b-932c-529ceb705947"
	segmentId := "b6d24b8e-af0b-4c3c-be0c-359bbd97381e"

	mux.HandleFunc("/contacts/"+contactId+"/segments/"+segmentId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"id": "c4a7e1e0-5f1f-4b8d-9e3a-7c2b1d8e9f0a",
			"object": "contact_segment"
		}`)
	})

	params := &AddContactSegmentRequest{
		SegmentId: segmentId,
		ContactId: contactId,
	}

	resp, err := client.Contacts.Segments.Add(params)
	if err != nil {
		t.Fatalf("failed to add contact to segment: %v", err)
	}

	assert.Equal(t, "c4a7e1e0-5f1f-4b8d-9e3a-7c2b1d8e9f0a", resp.Id)
	assert.Equal(t, "contact_segment", resp.Object)
}

func TestContactSegmentsAddWithEmail(t *testing.T) {
	setup()
	defer teardown()

	email := "user@example.com"
	segmentId := "b6d24b8e-af0b-4c3c-be0c-359bbd97381e"

	mux.HandleFunc("/contacts/"+email+"/segments/"+segmentId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"id": "c4a7e1e0-5f1f-4b8d-9e3a-7c2b1d8e9f0a",
			"object": "contact_segment"
		}`)
	})

	params := &AddContactSegmentRequest{
		SegmentId: segmentId,
		Email:     email,
	}

	resp, err := client.Contacts.Segments.Add(params)
	if err != nil {
		t.Fatalf("failed to add contact to segment: %v", err)
	}

	assert.NotEmpty(t, resp.Id)
}

func TestContactSegmentsAddValidation(t *testing.T) {
	// Test missing segment_id
	_, err := client.Contacts.Segments.Add(&AddContactSegmentRequest{
		ContactId: "contact_123",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SegmentId is required")

	// Test missing both ContactId and Email
	_, err = client.Contacts.Segments.Add(&AddContactSegmentRequest{
		SegmentId: "segment_123",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Either ContactId or Email must be provided")
}

func TestContactSegmentsRemoveWithContactId(t *testing.T) {
	setup()
	defer teardown()

	contactId := "479e3145-dd38-476b-932c-529ceb705947"
	segmentId := "b6d24b8e-af0b-4c3c-be0c-359bbd97381e"

	mux.HandleFunc("/contacts/"+contactId+"/segments/"+segmentId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"id": "c4a7e1e0-5f1f-4b8d-9e3a-7c2b1d8e9f0a",
			"object": "contact_segment",
			"deleted": true
		}`)
	})

	params := &RemoveContactSegmentRequest{
		SegmentId: segmentId,
		ContactId: contactId,
	}

	resp, err := client.Contacts.Segments.Remove(params)
	if err != nil {
		t.Fatalf("failed to remove contact from segment: %v", err)
	}

	assert.True(t, resp.Deleted)
	assert.Equal(t, "contact_segment", resp.Object)
}

func TestContactSegmentsListWithContactId(t *testing.T) {
	setup()
	defer teardown()

	contactId := "479e3145-dd38-476b-932c-529ceb705947"

	mux.HandleFunc("/contacts/"+contactId+"/segments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "list",
			"data": [
				{
					"id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
					"name": "Test Segment",
					"object": "segment",
					"created_at": "2023-01-01T00:00:00.000Z"
				}
			],
			"has_more": false
		}`)
	})

	params := &ListContactSegmentsRequest{
		ContactId: contactId,
	}

	resp, err := client.Contacts.Segments.List(params)
	if err != nil {
		t.Fatalf("failed to list contact segments: %v", err)
	}

	assert.Equal(t, "list", resp.Object)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "Test Segment", resp.Data[0].Name)
	assert.False(t, resp.HasMore)
}

func TestContactSegmentsListWithEmail(t *testing.T) {
	setup()
	defer teardown()

	email := "user@example.com"

	mux.HandleFunc("/contacts/"+email+"/segments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "list",
			"data": [],
			"has_more": false
		}`)
	})

	params := &ListContactSegmentsRequest{
		Email: email,
	}

	resp, err := client.Contacts.Segments.List(params)
	if err != nil {
		t.Fatalf("failed to list contact segments: %v", err)
	}

	assert.Equal(t, "list", resp.Object)
	assert.Len(t, resp.Data, 0)
}

func TestContactSegmentsListWithPagination(t *testing.T) {
	setup()
	defer teardown()

	contactId := "479e3145-dd38-476b-932c-529ceb705947"

	mux.HandleFunc("/contacts/"+contactId+"/segments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		query := r.URL.Query()
		assert.Equal(t, "10", query.Get("limit"))
		assert.Equal(t, "cursor_123", query.Get("after"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "list",
			"data": [],
			"has_more": true
		}`)
	})

	limit := 10
	after := "cursor_123"
	params := &ListContactSegmentsRequest{
		ContactId: contactId,
	}

	options := &ListOptions{
		Limit: &limit,
		After: &after,
	}

	resp, err := client.Contacts.Segments.ListWithOptions(context.Background(), params, options)
	if err != nil {
		t.Fatalf("failed to list contact segments with options: %v", err)
	}

	assert.True(t, resp.HasMore)
}

func TestContactSegmentsListValidation(t *testing.T) {
	// Test missing both ContactId and Email
	_, err := client.Contacts.Segments.List(&ListContactSegmentsRequest{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Either ContactId or Email must be provided")
}
