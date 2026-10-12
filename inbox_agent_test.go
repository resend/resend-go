package resend

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInboxAgent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId+"/agent", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"object": "inbox_agent",
			"instructions": "Answer refund questions yourself.",
			"tone": null,
			"enabled_actions": ["draft_reply", "add_labels"]
		}`)
	})

	resp, err := client.Inboxes.Agent.Get(testInboxId)
	assert.NoError(t, err)
	assert.Equal(t, "inbox_agent", resp.Object)
	assert.Equal(t, "Answer refund questions yourself.", *resp.Instructions)
	assert.Nil(t, resp.Tone)
	assert.Equal(t, []InboxAgentAction{InboxAgentActionDraftReply, InboxAgentActionAddLabels}, resp.EnabledActions)

	_, err = client.Inboxes.Agent.GetWithContext(context.Background(), "")
	assert.EqualError(t, err, "[ERROR]: Inbox ID is required")
}

func TestUpdateInboxAgent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId+"/agent", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		assert.Equal(t, map[string]any{
			"instructions":    "Escalate legal threats.",
			"tone":            nil,
			"enabled_actions": []any{"mark_as_spam"},
		}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"object": "inbox_agent", "id": "a1f84a4e-6f2b-4f0a-9c1d-8a2e5b3c7d90"}`)
	})

	instructions := "Escalate legal threats."
	params := &UpdateInboxAgentRequest{
		Instructions:   &instructions,
		EnabledActions: []InboxAgentAction{InboxAgentActionMarkAsSpam},
	}
	params.ClearTone()

	resp, err := client.Inboxes.Agent.Update(testInboxId, params)
	assert.NoError(t, err)
	assert.Equal(t, "inbox_agent", resp.Object)
	assert.Equal(t, "a1f84a4e-6f2b-4f0a-9c1d-8a2e5b3c7d90", resp.Id)
}

func TestUpdateInboxAgentEmptyActionsDisablesAll(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/inboxes/"+testInboxId+"/agent", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, map[string]any{"enabled_actions": []any{}}, decodeBody(t, r))
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"object": "inbox_agent", "id": "a1f84a4e-6f2b-4f0a-9c1d-8a2e5b3c7d90"}`)
	})

	_, err := client.Inboxes.Agent.UpdateWithContext(context.Background(), testInboxId, &UpdateInboxAgentRequest{
		EnabledActions: []InboxAgentAction{},
	})
	assert.NoError(t, err)
}

func TestUpdateInboxAgentValidation(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Inboxes.Agent.Update("", &UpdateInboxAgentRequest{})
	assert.EqualError(t, err, "[ERROR]: Inbox ID is required")

	_, err = client.Inboxes.Agent.Update(testInboxId, &UpdateInboxAgentRequest{})
	assert.EqualError(t, err, "[ERROR]: At least one of Instructions, Tone or EnabledActions is required")

	_, err = client.Inboxes.Agent.Update(testInboxId, nil)
	assert.EqualError(t, err, "[ERROR]: At least one of Instructions, Tone or EnabledActions is required")
}
