package resend

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAutomation(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/automations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"object":"automation","id":"aut_123"}`)
	})

	resp, err := client.Automations.Create(&CreateAutomationRequest{
		Name: "Welcome Flow",
		Steps: []AutomationStep{
			{Key: "trigger_1", Type: AutomationStepTypeTrigger, Config: map[string]any{"event_name": "user.created"}},
			{Key: "send_1", Type: AutomationStepTypeSendEmail, Config: map[string]any{"template": map[string]any{"id": "tpl_abc"}}},
		},
		Connections: []AutomationConnection{
			{From: "trigger_1", To: "send_1"},
		},
	})
	if err != nil {
		t.Errorf("Automations.Create returned error: %v", err)
	}
	assert.Equal(t, "automation", resp.Object)
	assert.Equal(t, "aut_123", resp.Id)
}

func TestGetAutomation(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/automations/aut_123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "automation",
			"id": "aut_123",
			"name": "Welcome Flow",
			"status": "enabled",
			"created_at": "2026-04-01T00:00:00Z",
			"updated_at": "2026-04-01T00:00:00Z",
			"steps": [
				{"key": "trigger_1", "type": "trigger", "config": {"event_name": "user.created"}},
				{"key": "send_1", "type": "send_email", "config": {"template": {"id": "tpl_abc"}}}
			],
			"connections": [
				{"from": "trigger_1", "to": "send_1"}
			]
		}`)
	})

	resp, err := client.Automations.Get("aut_123")
	if err != nil {
		t.Errorf("Automations.Get returned error: %v", err)
	}
	assert.Equal(t, "automation", resp.Object)
	assert.Equal(t, "aut_123", resp.Id)
	assert.Equal(t, "Welcome Flow", resp.Name)
	assert.Equal(t, AutomationStatusEnabled, resp.Status)
	assert.Equal(t, 2, len(resp.Steps))
	assert.Equal(t, "trigger_1", resp.Steps[0].Key)
	assert.Equal(t, AutomationStepTypeTrigger, resp.Steps[0].Type)
	assert.Equal(t, "send_1", resp.Steps[1].Key)
	assert.Equal(t, 1, len(resp.Connections))
}

func TestListAutomations(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/automations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "list",
			"has_more": false,
			"data": [
				{"id": "aut_123", "name": "Welcome Flow", "status": "enabled", "created_at": "2026-04-01T00:00:00Z", "updated_at": "2026-04-01T00:00:00Z"},
				{"id": "aut_456", "name": "Onboarding", "status": "disabled", "created_at": "2026-04-02T00:00:00Z", "updated_at": "2026-04-02T00:00:00Z"}
			]
		}`)
	})

	resp, err := client.Automations.List()
	if err != nil {
		t.Errorf("Automations.List returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, false, resp.HasMore)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "aut_123", resp.Data[0].Id)
	assert.Equal(t, "Welcome Flow", resp.Data[0].Name)
	assert.Equal(t, AutomationStatusEnabled, resp.Data[0].Status)
	assert.Equal(t, "aut_456", resp.Data[1].Id)
	assert.Equal(t, AutomationStatusDisabled, resp.Data[1].Status)
}

func TestListAutomationsWithOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/automations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, "enabled", r.URL.Query().Get("status"))
		assert.Equal(t, "5", r.URL.Query().Get("limit"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "list",
			"has_more": true,
			"data": [
				{"id": "aut_123", "name": "Welcome Flow", "status": "enabled", "created_at": "2026-04-01T00:00:00Z", "updated_at": "2026-04-01T00:00:00Z"}
			]
		}`)
	})

	status := AutomationStatusEnabled
	limit := 5
	resp, err := client.Automations.ListWithOptions(context.Background(), &ListAutomationsOptions{
		Status: &status,
		Limit:  &limit,
	})
	if err != nil {
		t.Errorf("Automations.ListWithOptions returned error: %v", err)
	}
	assert.Equal(t, true, resp.HasMore)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "aut_123", resp.Data[0].Id)
}

func TestUpdateAutomation(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/automations/aut_123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"object":"automation","id":"aut_123"}`)
	})

	resp, err := client.Automations.Update("aut_123", &UpdateAutomationRequest{
		Status: AutomationStatusEnabled,
	})
	if err != nil {
		t.Errorf("Automations.Update returned error: %v", err)
	}
	assert.Equal(t, "automation", resp.Object)
	assert.Equal(t, "aut_123", resp.Id)
}

func TestRemoveAutomation(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/automations/aut_123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"object":"automation","id":"aut_123","deleted":true}`)
	})

	resp, err := client.Automations.Remove("aut_123")
	if err != nil {
		t.Errorf("Automations.Remove returned error: %v", err)
	}
	assert.Equal(t, "automation", resp.Object)
	assert.Equal(t, "aut_123", resp.Id)
	assert.Equal(t, true, resp.Deleted)
}

func TestStopAutomation(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/automations/aut_123/stop", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"object":"automation","id":"aut_123","status":"disabled"}`)
	})

	resp, err := client.Automations.Stop("aut_123")
	if err != nil {
		t.Errorf("Automations.Stop returned error: %v", err)
	}
	assert.Equal(t, "automation", resp.Object)
	assert.Equal(t, "aut_123", resp.Id)
	assert.Equal(t, "disabled", resp.Status)
}

func TestListAutomationRuns(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/automations/aut_123/runs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "list",
			"has_more": false,
			"data": [
				{"id": "run_1", "status": "completed", "started_at": "2026-04-01T00:00:00Z", "completed_at": "2026-04-01T00:01:00Z", "created_at": "2026-04-01T00:00:00Z"},
				{"id": "run_2", "status": "running", "started_at": "2026-04-02T00:00:00Z", "completed_at": null, "created_at": "2026-04-02T00:00:00Z"}
			]
		}`)
	})

	resp, err := client.Automations.ListRuns("aut_123")
	if err != nil {
		t.Errorf("Automations.ListRuns returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, false, resp.HasMore)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "run_1", resp.Data[0].Id)
	assert.Equal(t, AutomationRunStatusCompleted, resp.Data[0].Status)
	assert.NotNil(t, resp.Data[0].StartedAt)
	assert.NotNil(t, resp.Data[0].CompletedAt)
	assert.Equal(t, "run_2", resp.Data[1].Id)
	assert.Equal(t, AutomationRunStatusRunning, resp.Data[1].Status)
	assert.Nil(t, resp.Data[1].CompletedAt)
}

func TestListAutomationRunsWithOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/automations/aut_123/runs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, "completed,failed", r.URL.Query().Get("status"))
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "list",
			"has_more": false,
			"data": [
				{"id": "run_1", "status": "completed", "started_at": "2026-04-01T00:00:00Z", "completed_at": "2026-04-01T00:01:00Z", "created_at": "2026-04-01T00:00:00Z"}
			]
		}`)
	})

	limit := 10
	resp, err := client.Automations.ListRunsWithContext(context.Background(), "aut_123", &ListAutomationRunsOptions{
		Status: []AutomationRunStatus{AutomationRunStatusCompleted, AutomationRunStatusFailed},
		Limit:  &limit,
	})
	if err != nil {
		t.Errorf("Automations.ListRunsWithContext returned error: %v", err)
	}
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "run_1", resp.Data[0].Id)
}

func TestCreateAutomationWithDelayAndWaitForEvent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/automations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"object":"automation","id":"aut_456"}`)
	})

	resp, err := client.Automations.Create(&CreateAutomationRequest{
		Name: "Onboarding Flow",
		Steps: []AutomationStep{
			{Key: "trigger_1", Type: AutomationStepTypeTrigger, Config: map[string]any{
				"event_name": "user.created",
			}},
			{Key: "delay_1", Type: AutomationStepTypeDelay, Config: map[string]any{
				"duration": "30 minutes",
			}},
			{Key: "wait_1", Type: AutomationStepTypeWaitForEvent, Config: map[string]any{
				"event_name": "user.verified",
				"timeout":    "1 hour",
			}},
			{Key: "send_1", Type: AutomationStepTypeSendEmail, Config: map[string]any{
				"template": map[string]any{"id": "tpl_abc"},
			}},
		},
		Connections: []AutomationConnection{
			{From: "trigger_1", To: "delay_1"},
			{From: "delay_1", To: "wait_1"},
			{From: "wait_1", To: "send_1", Type: AutomationConnectionTypeEventReceived},
			{From: "wait_1", To: "send_1", Type: AutomationConnectionTypeTimeout},
		},
	})
	if err != nil {
		t.Errorf("Automations.Create returned error: %v", err)
	}
	assert.Equal(t, "automation", resp.Object)
	assert.Equal(t, "aut_456", resp.Id)
}

func TestGetAutomationStepResponseKeys(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/automations/aut_456", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "automation",
			"id": "aut_456",
			"name": "Onboarding Flow",
			"status": "disabled",
			"created_at": "2026-04-10T00:00:00Z",
			"updated_at": "2026-04-10T00:00:00Z",
			"steps": [
				{"key": "trigger_1", "type": "trigger", "config": {"event_name": "user.created"}},
				{"key": "delay_1", "type": "delay", "config": {"duration": "30 minutes"}},
				{"key": "wait_1", "type": "wait_for_event", "config": {"event_name": "user.verified", "timeout": "1 hour"}},
				{"key": "send_1", "type": "send_email", "config": {"template": {"id": "tpl_abc"}}}
			],
			"connections": [
				{"from": "trigger_1", "to": "delay_1"},
				{"from": "delay_1", "to": "wait_1"},
				{"from": "wait_1", "to": "send_1", "type": "event_received"},
				{"from": "wait_1", "to": "send_1", "type": "timeout"}
			]
		}`)
	})

	resp, err := client.Automations.Get("aut_456")
	if err != nil {
		t.Errorf("Automations.Get returned error: %v", err)
	}
	assert.Equal(t, 4, len(resp.Steps))

	assert.Equal(t, "trigger_1", resp.Steps[0].Key)
	assert.Equal(t, AutomationStepTypeTrigger, resp.Steps[0].Type)

	assert.Equal(t, "delay_1", resp.Steps[1].Key)
	assert.Equal(t, AutomationStepTypeDelay, resp.Steps[1].Type)
	assert.Equal(t, "30 minutes", resp.Steps[1].Config["duration"])

	assert.Equal(t, "wait_1", resp.Steps[2].Key)
	assert.Equal(t, AutomationStepTypeWaitForEvent, resp.Steps[2].Type)
	assert.Equal(t, "user.verified", resp.Steps[2].Config["event_name"])
	assert.Equal(t, "1 hour", resp.Steps[2].Config["timeout"])

	assert.Equal(t, "send_1", resp.Steps[3].Key)
	assert.Equal(t, AutomationStepTypeSendEmail, resp.Steps[3].Type)

	assert.Equal(t, 4, len(resp.Connections))
	assert.Equal(t, AutomationConnectionTypeEventReceived, resp.Connections[2].Type)
	assert.Equal(t, AutomationConnectionTypeTimeout, resp.Connections[3].Type)
}

func TestGetAutomationRun(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/automations/aut_123/runs/run_1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "automation_run",
			"id": "run_1",
			"status": "completed",
			"started_at": "2026-04-01T00:00:00Z",
			"completed_at": "2026-04-01T00:01:00Z",
			"created_at": "2026-04-01T00:00:00Z",
			"steps": [
				{"key": "trigger_1", "type": "trigger", "status": "completed", "started_at": "2026-04-01T00:00:00Z", "completed_at": "2026-04-01T00:00:01Z", "output": null, "error": null, "created_at": "2026-04-01T00:00:00Z"},
				{"key": "send_1", "type": "send_email", "status": "completed", "started_at": "2026-04-01T00:00:01Z", "completed_at": "2026-04-01T00:01:00Z", "output": null, "error": null, "created_at": "2026-04-01T00:00:01Z"}
			]
		}`)
	})

	resp, err := client.Automations.GetRun("aut_123", "run_1")
	if err != nil {
		t.Errorf("Automations.GetRun returned error: %v", err)
	}
	assert.Equal(t, "automation_run", resp.Object)
	assert.Equal(t, "run_1", resp.Id)
	assert.Equal(t, AutomationRunStatusCompleted, resp.Status)
	assert.NotNil(t, resp.StartedAt)
	assert.NotNil(t, resp.CompletedAt)
	assert.Equal(t, 2, len(resp.Steps))
	assert.Equal(t, "trigger_1", resp.Steps[0].Key)
	assert.Equal(t, AutomationStepTypeTrigger, resp.Steps[0].Type)
	assert.Equal(t, "completed", resp.Steps[0].Status)
	assert.Equal(t, "send_1", resp.Steps[1].Key)
	assert.Equal(t, AutomationStepTypeSendEmail, resp.Steps[1].Type)
}
