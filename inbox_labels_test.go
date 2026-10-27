package resend

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testLabelId = "7f9c1d2e-4a6b-4c3d-8e15-9b0a7c6d5e34"

func TestListInboxLabels(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId+"/labels", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"object": "list",
			"data": [
				{"id": "%s", "name": "Urgent", "color": "#E93D82", "created_at": "2026-08-05T14:07:42.881Z"}
			]
		}`, testLabelId)
	})

	resp, err := client.Inboxes.Labels.List(testInboxId)
	assert.NoError(t, err)
	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, []InboxLabel{{
		Id:        testLabelId,
		Name:      "Urgent",
		Color:     "#E93D82",
		CreatedAt: "2026-08-05T14:07:42.881Z",
	}}, resp.Data)

	_, err = client.Inboxes.Labels.ListWithContext(context.Background(), "")
	assert.EqualError(t, err, "[ERROR]: Inbox ID is required")
}

func TestCreateInboxLabel(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId+"/labels", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assert.Equal(t, map[string]any{"name": "Urgent", "color": "#e93d82"}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{"object": "inbox_label", "id": "%s"}`, testLabelId)
	})

	resp, err := client.Inboxes.Labels.Create(testInboxId, &CreateInboxLabelRequest{Name: "Urgent", Color: "#e93d82"})
	assert.NoError(t, err)
	assert.Equal(t, "inbox_label", resp.Object)
	assert.Equal(t, testLabelId, resp.Id)
}

func TestCreateInboxLabelValidation(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Inboxes.Labels.Create("", &CreateInboxLabelRequest{Name: "Urgent"})
	assert.EqualError(t, err, "[ERROR]: Inbox ID is required")

	_, err = client.Inboxes.Labels.CreateWithContext(context.Background(), testInboxId, &CreateInboxLabelRequest{})
	assert.EqualError(t, err, "[ERROR]: Name is required")
}

func TestUpdateInboxLabel(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId+"/labels/"+testLabelId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		assert.Equal(t, map[string]any{"color": "#F76B15"}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "inbox_label", "id": "%s"}`, testLabelId)
	})

	resp, err := client.Inboxes.Labels.Update(testInboxId, testLabelId, &UpdateInboxLabelRequest{Color: "#F76B15"})
	assert.NoError(t, err)
	assert.Equal(t, testLabelId, resp.Id)

	_, err = client.Inboxes.Labels.UpdateWithContext(context.Background(), testInboxId, "", nil)
	assert.EqualError(t, err, "[ERROR]: Label ID is required")
}

func TestUpdateInboxLabelEmptyBody(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId+"/labels/"+testLabelId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		assert.Equal(t, map[string]any{}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "inbox_label", "id": "%s"}`, testLabelId)
	})

	_, err := client.Inboxes.Labels.Update(testInboxId, testLabelId, nil)
	assert.NoError(t, err)
}

func TestRemoveInboxLabel(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId+"/labels/"+testLabelId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "inbox_label", "id": "%s", "deleted": true}`, testLabelId)
	})

	resp, err := client.Inboxes.Labels.Remove(testInboxId, testLabelId)
	assert.NoError(t, err)
	assert.True(t, resp.Deleted)

	_, err = client.Inboxes.Labels.RemoveWithContext(context.Background(), "", testLabelId)
	assert.EqualError(t, err, "[ERROR]: Inbox ID is required")
}
