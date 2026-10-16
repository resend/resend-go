package resend

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func threadLabelsPath() string {
	return "/inboxes/" + testInboxId + "/threads/" + testThreadId + "/labels"
}

func TestAddInboxThreadLabels(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(threadLabelsPath(), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assert.Equal(t, map[string]any{"labels": []any{testLabelId}}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "list", "data": [{"id": "%s", "name": "Urgent", "color": "#E93D82"}]}`, testLabelId)
	})

	resp, err := client.Inboxes.Threads.Labels.Add(testInboxId, testThreadId, &AddInboxThreadLabelsRequest{
		Labels: []string{testLabelId},
	})
	assert.NoError(t, err)
	assert.Equal(t, InboxThreadLabelsResponse{
		Object: "list",
		Data:   []InboxThreadLabel{{Id: testLabelId, Name: "Urgent", Color: "#E93D82"}},
	}, resp)
}

func TestAddInboxThreadLabelsValidation(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Inboxes.Threads.Labels.Add(testInboxId, testThreadId, &AddInboxThreadLabelsRequest{})
	assert.EqualError(t, err, "[ERROR]: At least one label ID is required")

	_, err = client.Inboxes.Threads.Labels.Add(testInboxId, testThreadId, nil)
	assert.EqualError(t, err, "[ERROR]: At least one label ID is required")

	_, err = client.Inboxes.Threads.Labels.Add("", testThreadId, nil)
	assert.EqualError(t, err, "[ERROR]: Inbox ID is required")

	_, err = client.Inboxes.Threads.Labels.Add(testInboxId, "", nil)
	assert.EqualError(t, err, "[ERROR]: Thread ID is required")
}

func TestSetInboxThreadLabels(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(threadLabelsPath(), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		assert.Equal(t, map[string]any{"labels": []any{testLabelId}}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"object": "list", "data": [{"id": "%s", "name": "Urgent", "color": "#E93D82"}]}`, testLabelId)
	})

	resp, err := client.Inboxes.Threads.Labels.Set(testInboxId, testThreadId, &SetInboxThreadLabelsRequest{
		Labels: []string{testLabelId},
	})
	assert.NoError(t, err)
	assert.Equal(t, testLabelId, resp.Data[0].Id)
}

func TestSetInboxThreadLabelsClearsWithNilOrEmpty(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(threadLabelsPath(), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		assert.Equal(t, map[string]any{"labels": []any{}}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"object": "list", "data": []}`)
	})

	resp, err := client.Inboxes.Threads.Labels.SetWithContext(context.Background(), testInboxId, testThreadId, &SetInboxThreadLabelsRequest{})
	assert.NoError(t, err)
	assert.Empty(t, resp.Data)

	_, err = client.Inboxes.Threads.Labels.Set(testInboxId, testThreadId, nil)
	assert.NoError(t, err)
}

func TestRemoveInboxThreadLabel(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(threadLabelsPath()+"/"+testLabelId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"object": "list", "data": []}`)
	})

	resp, err := client.Inboxes.Threads.Labels.Remove(testInboxId, testThreadId, testLabelId)
	assert.NoError(t, err)
	assert.Equal(t, "list", resp.Object)
	assert.Empty(t, resp.Data)

	_, err = client.Inboxes.Threads.Labels.RemoveWithContext(context.Background(), testInboxId, testThreadId, "")
	assert.EqualError(t, err, "[ERROR]: Label ID is required")
}

func TestInboxRoutesAcceptAnEmailAddress(t *testing.T) {
	setup()
	defer teardown()

	const address = "support@acme.com"
	var paths []string
	mux.HandleFunc("/inboxes/"+address+"/", func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{}`)
	})

	read := true
	_, err := client.Inboxes.Agent.Get(address)
	assert.NoError(t, err)
	_, err = client.Inboxes.Labels.List(address)
	assert.NoError(t, err)
	_, err = client.Inboxes.Threads.Update(address, testThreadId, &UpdateInboxThreadRequest{Read: &read})
	assert.NoError(t, err)
	_, err = client.Inboxes.Threads.Labels.Add(address, testThreadId, &AddInboxThreadLabelsRequest{Labels: []string{testLabelId}})
	assert.NoError(t, err)
	_, err = client.Inboxes.Threads.Emails.List(address, testThreadId, nil)
	assert.NoError(t, err)
	_, err = client.Inboxes.Drafts.Create(address, &CreateInboxDraftRequest{Text: "Hi"})
	assert.NoError(t, err)

	assert.Equal(t, []string{
		"GET /inboxes/" + address + "/agent",
		"GET /inboxes/" + address + "/labels",
		"PATCH /inboxes/" + address + "/threads/" + testThreadId,
		"POST /inboxes/" + address + "/threads/" + testThreadId + "/labels",
		"GET /inboxes/" + address + "/threads/" + testThreadId + "/emails",
		"POST /inboxes/" + address + "/drafts",
	}, paths)
}
