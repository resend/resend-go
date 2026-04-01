package resend

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLog(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/logs/log_123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `{
			"object": "log",
			"id": "log_123",
			"created_at": "2026-03-30 13:43:54.622865+00",
			"endpoint": "/emails",
			"method": "POST",
			"response_status": 200,
			"user_agent": "resend-go/3.2.0",
			"request_body": {"from": "hello@example.com"},
			"response_body": {"id": "email_abc"}
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Logs.Get("log_123")
	if err != nil {
		t.Errorf("Logs.Get returned error: %v", err)
	}
	assert.Equal(t, "log", resp.Object)
	assert.Equal(t, "log_123", resp.Id)
	assert.Equal(t, "2026-03-30 13:43:54.622865+00", resp.CreatedAt)
	assert.Equal(t, "/emails", resp.Endpoint)
	assert.Equal(t, "POST", resp.Method)
	assert.Equal(t, 200, resp.ResponseStatus)
	assert.NotNil(t, resp.UserAgent)
	assert.Equal(t, "resend-go/3.2.0", *resp.UserAgent)
	assert.NotNil(t, resp.RequestBody)
	assert.NotNil(t, resp.ResponseBody)
}

func TestGetLogNullUserAgent(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/logs/log_456", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `{
			"object": "log",
			"id": "log_456",
			"created_at": "2026-03-30 14:00:00.000000+00",
			"endpoint": "/api-keys",
			"method": "GET",
			"response_status": 200,
			"user_agent": null
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Logs.Get("log_456")
	if err != nil {
		t.Errorf("Logs.Get returned error: %v", err)
	}
	assert.Equal(t, "log_456", resp.Id)
	assert.Nil(t, resp.UserAgent)
}

func TestListLogs(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `{
			"object": "list",
			"has_more": false,
			"data": [
				{
					"id": "log_111",
					"created_at": "2026-03-30 13:43:54.622865+00",
					"endpoint": "/emails",
					"method": "POST",
					"response_status": 200,
					"user_agent": "resend-go/3.2.0"
				},
				{
					"id": "log_222",
					"created_at": "2026-03-29 10:00:00.000000+00",
					"endpoint": "/domains",
					"method": "GET",
					"response_status": 200,
					"user_agent": null
				}
			]
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Logs.List()
	if err != nil {
		t.Errorf("Logs.List returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, false, resp.HasMore)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "log_111", resp.Data[0].Id)
	assert.Equal(t, "/emails", resp.Data[0].Endpoint)
	assert.Equal(t, "POST", resp.Data[0].Method)
	assert.Equal(t, 200, resp.Data[0].ResponseStatus)
	assert.NotNil(t, resp.Data[0].UserAgent)
	assert.Equal(t, "resend-go/3.2.0", *resp.Data[0].UserAgent)
	assert.Equal(t, "log_222", resp.Data[1].Id)
	assert.Nil(t, resp.Data[1].UserAgent)
}

func TestListLogsWithOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		assert.Equal(t, "5", r.URL.Query().Get("limit"))
		assert.Equal(t, "log_000", r.URL.Query().Get("after"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `{
			"object": "list",
			"has_more": true,
			"data": [
				{
					"id": "log_333",
					"created_at": "2026-03-28 09:00:00.000000+00",
					"endpoint": "/emails",
					"method": "POST",
					"response_status": 422,
					"user_agent": "resend-node/6.0.3"
				}
			]
		}`
		fmt.Fprintf(w, ret)
	})

	limit := 5
	after := "log_000"
	resp, err := client.Logs.ListWithOptions(context.Background(), &ListOptions{
		Limit: &limit,
		After: &after,
	})
	if err != nil {
		t.Errorf("Logs.ListWithOptions returned error: %v", err)
	}
	assert.Equal(t, true, resp.HasMore)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "log_333", resp.Data[0].Id)
	assert.Equal(t, 422, resp.Data[0].ResponseStatus)
}
