package resend

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEvent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"object":"event","id":"evt_123"}`)
	})

	resp, err := client.Events.Create(&CreateEventRequest{
		Name: "user.created",
		Schema: map[string]string{
			"plan": EventSchemaTypeString,
			"age":  EventSchemaTypeNumber,
		},
	})
	if err != nil {
		t.Fatalf("Events.Create returned error: %v", err)
	}
	assert.Equal(t, "event", resp.Object)
	assert.Equal(t, "evt_123", resp.Id)
}

func TestCreateEventNoSchema(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"object":"event","id":"evt_456"}`)
	})

	resp, err := client.Events.Create(&CreateEventRequest{
		Name: "order.placed",
	})
	if err != nil {
		t.Fatalf("Events.Create returned error: %v", err)
	}
	assert.Equal(t, "evt_456", resp.Id)
}

func TestGetEvent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/events/evt_123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "event",
			"id": "evt_123",
			"name": "user.created",
			"schema": {"plan": "string", "age": "number"},
			"created_at": "2026-04-01T00:00:00Z",
			"updated_at": null
		}`)
	})

	resp, err := client.Events.Get("evt_123")
	if err != nil {
		t.Fatalf("Events.Get returned error: %v", err)
	}
	assert.Equal(t, "event", resp.Object)
	assert.Equal(t, "evt_123", resp.Id)
	assert.Equal(t, "user.created", resp.Name)
	assert.Equal(t, "string", resp.Schema["plan"])
	assert.Equal(t, "number", resp.Schema["age"])
	assert.Nil(t, resp.UpdatedAt)
}

func TestGetEventByName(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/events/user.created", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "event",
			"id": "evt_123",
			"name": "user.created",
			"schema": null,
			"created_at": "2026-04-01T00:00:00Z",
			"updated_at": "2026-04-02T00:00:00Z"
		}`)
	})

	resp, err := client.Events.Get("user.created")
	if err != nil {
		t.Fatalf("Events.Get returned error: %v", err)
	}
	assert.Equal(t, "evt_123", resp.Id)
	assert.Equal(t, "user.created", resp.Name)
	assert.NotNil(t, resp.UpdatedAt)
}

func TestListEvents(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "list",
			"has_more": false,
			"data": [
				{"id": "evt_123", "name": "user.created", "schema": {"plan": "string"}, "created_at": "2026-04-01T00:00:00Z", "updated_at": null},
				{"id": "evt_456", "name": "order.placed", "schema": null, "created_at": "2026-04-02T00:00:00Z", "updated_at": null}
			]
		}`)
	})

	resp, err := client.Events.List()
	if err != nil {
		t.Fatalf("Events.List returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, false, resp.HasMore)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "evt_123", resp.Data[0].Id)
	assert.Equal(t, "user.created", resp.Data[0].Name)
	assert.Equal(t, "evt_456", resp.Data[1].Id)
}

func TestListEventsWithOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, "5", r.URL.Query().Get("limit"))
		assert.Equal(t, "evt_000", r.URL.Query().Get("after"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "list",
			"has_more": true,
			"data": [
				{"id": "evt_123", "name": "user.created", "schema": null, "created_at": "2026-04-01T00:00:00Z", "updated_at": null}
			]
		}`)
	})

	limit := 5
	after := "evt_000"
	resp, err := client.Events.ListWithOptions(context.Background(), &ListOptions{
		Limit: &limit,
		After: &after,
	})
	if err != nil {
		t.Fatalf("Events.ListWithOptions returned error: %v", err)
	}
	assert.Equal(t, true, resp.HasMore)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "evt_123", resp.Data[0].Id)
}

func TestUpdateEvent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/events/evt_123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"object":"event","id":"evt_123"}`)
	})

	resp, err := client.Events.Update("evt_123", &UpdateEventRequest{
		Schema: map[string]string{
			"plan":     EventSchemaTypeString,
			"verified": EventSchemaTypeBoolean,
		},
	})
	if err != nil {
		t.Fatalf("Events.Update returned error: %v", err)
	}
	assert.Equal(t, "event", resp.Object)
	assert.Equal(t, "evt_123", resp.Id)
}

func TestRemoveEvent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/events/evt_123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"object":"event","id":"evt_123","deleted":true}`)
	})

	resp, err := client.Events.Remove("evt_123")
	if err != nil {
		t.Fatalf("Events.Remove returned error: %v", err)
	}
	assert.Equal(t, "event", resp.Object)
	assert.Equal(t, "evt_123", resp.Id)
	assert.Equal(t, true, resp.Deleted)
}

func TestSendEvent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/events/send", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, `{"object":"event","event":"user.created"}`)
	})

	resp, err := client.Events.Send(&SendEventRequest{
		Event: "user.created",
		Email: "user@example.com",
		Payload: map[string]any{
			"plan": "pro",
		},
	})
	if err != nil {
		t.Fatalf("Events.Send returned error: %v", err)
	}
	assert.Equal(t, "event", resp.Object)
	assert.Equal(t, "user.created", resp.Event)
}

func TestSendEventWithContactId(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/events/send", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, `{"object":"event","event":"order.placed"}`)
	})

	resp, err := client.Events.Send(&SendEventRequest{
		Event:     "order.placed",
		ContactId: "c1b2c3d4-e5f6-7890-abcd-ef1234567890",
	})
	if err != nil {
		t.Fatalf("Events.Send returned error: %v", err)
	}
	assert.Equal(t, "order.placed", resp.Event)
}
