package resend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testInboxId = "b3e2b2b6-3f0e-4c8e-9ad3-2f43a1e2c7f1"
)

func decodeBody(t *testing.T, r *http.Request) map[string]any {
	t.Helper()
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode request body: %v", err)
	}
	return body
}

func TestCreateInbox(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		body := decodeBody(t, r)
		assert.Equal(t, map[string]any{
			"email_address": "support@acme.com",
			"name":          "Customer Support",
			"from_name":     "Ada from Support",
			"forwarding":    true,
		}, body)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"object": "inbox", "id": "%s"}`, testInboxId)
	})

	resp, err := client.Inboxes.Create(&CreateInboxRequest{
		EmailAddress: "support@acme.com",
		Name:         "Customer Support",
		FromName:     "Ada from Support",
		Forwarding:   true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "inbox", resp.Object)
	assert.Equal(t, testInboxId, resp.Id)
}

func TestCreateInboxOmitsUnsetFields(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, map[string]any{"email_address": "support@acme.com"}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"object": "inbox", "id": "%s"}`, testInboxId)
	})

	_, err := client.Inboxes.CreateWithContext(context.Background(), &CreateInboxRequest{
		EmailAddress: "support@acme.com",
	})
	assert.NoError(t, err)
}

func TestCreateInboxMissingEmailAddress(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Inboxes.Create(&CreateInboxRequest{})
	assert.EqualError(t, err, "[ERROR]: EmailAddress is required")

	_, err = client.Inboxes.Create(nil)
	assert.EqualError(t, err, "[ERROR]: EmailAddress is required")
}

func TestListInboxes(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "cursor-id", r.URL.Query().Get("after"))

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"object": "list",
			"has_more": true,
			"data": [
				{
					"id": "%s",
					"name": "Customer Support",
					"email_address": "support@acme.com",
					"from_name": null,
					"unread": 3,
					"last_received": "2026-08-05T14:03:11.229Z"
				}
			]
		}`, testInboxId)
	})

	limit := 10
	after := "cursor-id"
	resp, err := client.Inboxes.List(&ListOptions{Limit: &limit, After: &after})
	assert.NoError(t, err)
	assert.Equal(t, "list", resp.Object)
	assert.True(t, resp.HasMore)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "Customer Support", *resp.Data[0].Name)
	assert.Nil(t, resp.Data[0].FromName)
	assert.Equal(t, 3, resp.Data[0].Unread)
	assert.Equal(t, "2026-08-05T14:03:11.229Z", *resp.Data[0].LastReceived)
}

func TestListInboxesWithoutOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes", func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"object": "list", "has_more": false, "data": []}`)
	})

	resp, err := client.Inboxes.ListWithContext(context.Background(), nil)
	assert.NoError(t, err)
	assert.Empty(t, resp.Data)
}

func TestGetInbox(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"object": "inbox",
			"id": "%s",
			"name": "Customer Support",
			"email_address": "support@acme.com",
			"domain_id": "d91cd9bd-1176-453e-8fc1-35364d380206",
			"receiving_address": null,
			"from_name": "Ada from Support",
			"unread": 3,
			"drafts": 2,
			"last_received": null,
			"created_at": "2026-08-01T09:12:44.512Z"
		}`, testInboxId)
	})

	resp, err := client.Inboxes.Get(testInboxId)
	assert.NoError(t, err)
	assert.Equal(t, "inbox", resp.Object)
	assert.Equal(t, "support@acme.com", resp.EmailAddress)
	assert.Equal(t, "d91cd9bd-1176-453e-8fc1-35364d380206", resp.DomainId)
	assert.Nil(t, resp.ReceivingAddress)
	assert.Equal(t, "Ada from Support", *resp.FromName)
	assert.Equal(t, 2, resp.Drafts)
	assert.Nil(t, resp.LastReceived)
	assert.Equal(t, "2026-08-01T09:12:44.512Z", resp.CreatedAt)
}

func TestGetInboxByEmail(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/support@acme.com", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "inbox", "id": "%s", "email_address": "support@acme.com"}`, testInboxId)
	})

	resp, err := client.Inboxes.GetWithContext(context.Background(), "support@acme.com")
	assert.NoError(t, err)
	assert.Equal(t, testInboxId, resp.Id)
}

func TestGetInboxMissingId(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Inboxes.Get("")
	assert.EqualError(t, err, "[ERROR]: Inbox ID or email is required")
}

func TestUpdateInbox(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		assert.Equal(t, map[string]any{"name": "Billing", "from_name": ""}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "inbox", "id": "%s"}`, testInboxId)
	})

	fromName := ""
	resp, err := client.Inboxes.Update(testInboxId, &UpdateInboxRequest{Name: "Billing", FromName: &fromName})
	assert.NoError(t, err)
	assert.Equal(t, "inbox", resp.Object)
	assert.Equal(t, testInboxId, resp.Id)
}

func TestUpdateInboxValidation(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Inboxes.Update("", &UpdateInboxRequest{Name: "Billing"})
	assert.EqualError(t, err, "[ERROR]: Inbox ID is required")

	_, err = client.Inboxes.UpdateWithContext(context.Background(), testInboxId, &UpdateInboxRequest{})
	assert.EqualError(t, err, "[ERROR]: At least one of Name or FromName is required")

	_, err = client.Inboxes.Update(testInboxId, nil)
	assert.EqualError(t, err, "[ERROR]: At least one of Name or FromName is required")
}

func TestRemoveInbox(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "inbox", "id": "%s", "deleted": true}`, testInboxId)
	})

	resp, err := client.Inboxes.Remove(testInboxId)
	assert.NoError(t, err)
	assert.Equal(t, testInboxId, resp.Id)
	assert.True(t, resp.Deleted)

	_, err = client.Inboxes.RemoveWithContext(context.Background(), "")
	assert.EqualError(t, err, "[ERROR]: Inbox ID is required")
}

func TestCreateInboxPlanLimitError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, `{"statusCode": 403, "name": "validation_error", "message": "You have reached the inbox limit of your plan"}`)
	})

	_, err := client.Inboxes.Create(&CreateInboxRequest{EmailAddress: "support@acme.com"})
	assert.EqualError(t, err, "[ERROR]: You have reached the inbox limit of your plan")
}
