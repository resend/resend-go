package resend

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateWebhook(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "webhook",
			"id": "4dd369bc-aa82-4ff3-97de-514ae3000ee0",
			"signing_secret": "whsec_xxxxxxxxxx"
		}`
		fmt.Fprintf(w, ret)
	})

	req := &CreateWebhookRequest{
		Endpoint: "https://webhook.example.com/handler",
		Events:   []string{"email.sent", "email.delivered", "email.bounced"},
	}
	resp, err := client.Webhooks.Create(req)
	if err != nil {
		t.Errorf("Webhooks.Create returned error: %v", err)
	}
	assert.Equal(t, "webhook", resp.Object)
	assert.Equal(t, "4dd369bc-aa82-4ff3-97de-514ae3000ee0", resp.Id)
	assert.Equal(t, "whsec_xxxxxxxxxx", resp.SigningSecret)
}

func TestCreateWebhookWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "webhook",
			"id": "test-webhook-id",
			"signing_secret": "whsec_test_secret"
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	req := &CreateWebhookRequest{
		Endpoint: "https://webhook.example.com/handler",
		Events:   []string{"email.sent"},
	}
	resp, err := client.Webhooks.CreateWithContext(ctx, req)
	if err != nil {
		t.Errorf("Webhooks.CreateWithContext returned error: %v", err)
	}
	assert.Equal(t, "webhook", resp.Object)
	assert.Equal(t, "test-webhook-id", resp.Id)
	assert.Equal(t, "whsec_test_secret", resp.SigningSecret)
}

func TestGetWebhook(t *testing.T) {
	setup()
	defer teardown()

	webhookId := "4dd369bc-aa82-4ff3-97de-514ae3000ee0"

	mux.HandleFunc(fmt.Sprintf("/webhooks/%s", webhookId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "webhook",
			"id": "4dd369bc-aa82-4ff3-97de-514ae3000ee0",
			"created_at": "2023-08-22T15:28:00.000Z",
			"status": "enabled",
			"endpoint": "https://webhook.example.com/handler",
			"events": ["email.sent", "email.received"],
			"signing_secret": "whsec_xxxxxxxxxx"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Webhooks.Get(webhookId)
	if err != nil {
		t.Errorf("Webhooks.Get returned error: %v", err)
	}
	assert.Equal(t, "webhook", resp.Object)
	assert.Equal(t, "4dd369bc-aa82-4ff3-97de-514ae3000ee0", resp.Id)
	assert.Equal(t, "2023-08-22T15:28:00.000Z", resp.CreatedAt)
	assert.Equal(t, "enabled", resp.Status)
	assert.Equal(t, "https://webhook.example.com/handler", resp.Endpoint)
	assert.Equal(t, 2, len(resp.Events))
	assert.Equal(t, "email.sent", resp.Events[0])
	assert.Equal(t, "email.received", resp.Events[1])
	assert.Equal(t, "whsec_xxxxxxxxxx", resp.SigningSecret)
}

func TestGetWebhookWithContext(t *testing.T) {
	setup()
	defer teardown()

	webhookId := "test-webhook-id"

	mux.HandleFunc(fmt.Sprintf("/webhooks/%s", webhookId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "webhook",
			"id": "test-webhook-id",
			"created_at": "2024-01-01T00:00:00.000Z",
			"status": "enabled",
			"endpoint": "https://test.example.com/webhook",
			"events": ["email.delivered"],
			"signing_secret": "whsec_test_secret"
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	resp, err := client.Webhooks.GetWithContext(ctx, webhookId)
	if err != nil {
		t.Errorf("Webhooks.GetWithContext returned error: %v", err)
	}
	assert.Equal(t, "webhook", resp.Object)
	assert.Equal(t, "test-webhook-id", resp.Id)
	assert.Equal(t, "2024-01-01T00:00:00.000Z", resp.CreatedAt)
	assert.Equal(t, "enabled", resp.Status)
	assert.Equal(t, "https://test.example.com/webhook", resp.Endpoint)
	assert.Equal(t, 1, len(resp.Events))
	assert.Equal(t, "email.delivered", resp.Events[0])
	assert.Equal(t, "whsec_test_secret", resp.SigningSecret)
}

func TestUpdateWebhook(t *testing.T) {
	setup()
	defer teardown()

	webhookId := "430eed87-632a-4ea6-90db-0aace67ec228"

	mux.HandleFunc(fmt.Sprintf("/webhooks/%s", webhookId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "webhook",
			"id": "430eed87-632a-4ea6-90db-0aace67ec228"
		}`
		fmt.Fprintf(w, ret)
	})

	endpoint := "https://new-webhook.example.com/handler"
	status := "enabled"
	req := &UpdateWebhookRequest{
		Endpoint: &endpoint,
		Events:   []string{"email.sent", "email.delivered"},
		Status:   &status,
	}
	resp, err := client.Webhooks.Update(webhookId, req)
	if err != nil {
		t.Errorf("Webhooks.Update returned error: %v", err)
	}
	assert.Equal(t, "webhook", resp.Object)
	assert.Equal(t, "430eed87-632a-4ea6-90db-0aace67ec228", resp.Id)
}

func TestUpdateWebhookWithContext(t *testing.T) {
	setup()
	defer teardown()

	webhookId := "test-update-id"

	mux.HandleFunc(fmt.Sprintf("/webhooks/%s", webhookId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "webhook",
			"id": "test-update-id"
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	status := "disabled"
	req := &UpdateWebhookRequest{
		Status: &status,
	}
	resp, err := client.Webhooks.UpdateWithContext(ctx, webhookId, req)
	if err != nil {
		t.Errorf("Webhooks.UpdateWithContext returned error: %v", err)
	}
	assert.Equal(t, "webhook", resp.Object)
	assert.Equal(t, "test-update-id", resp.Id)
}

func TestListWebhooks(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
			"has_more": false,
			"data": [
				{
					"id": "7ab123cd-ef45-6789-abcd-ef0123456789",
					"created_at": "2023-09-10T10:15:30.000Z",
					"status": "disabled",
					"endpoint": "https://first-webhook.example.com/handler",
					"events": ["email.delivered", "email.bounced"]
				},
				{
					"id": "4dd369bc-aa82-4ff3-97de-514ae3000ee0",
					"created_at": "2023-08-22T15:28:00.000Z",
					"status": "enabled",
					"endpoint": "https://second-webhook.example.com/receive",
					"events": ["email.received"]
				}
			]
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Webhooks.List()
	if err != nil {
		t.Errorf("Webhooks.List returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, false, resp.HasMore)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "7ab123cd-ef45-6789-abcd-ef0123456789", resp.Data[0].Id)
	assert.Equal(t, "disabled", resp.Data[0].Status)
	assert.Equal(t, "https://first-webhook.example.com/handler", resp.Data[0].Endpoint)
	assert.Equal(t, 2, len(resp.Data[0].Events))
	assert.Equal(t, "4dd369bc-aa82-4ff3-97de-514ae3000ee0", resp.Data[1].Id)
	assert.Equal(t, "enabled", resp.Data[1].Status)
}

func TestListWebhooksWithOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/webhooks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		// Verify query parameters
		query := r.URL.Query()
		assert.Equal(t, "10", query.Get("limit"))
		assert.Equal(t, "cursor123", query.Get("after"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
			"has_more": true,
			"data": [
				{
					"id": "webhook-1",
					"created_at": "2024-01-01T00:00:00.000Z",
					"status": "enabled",
					"endpoint": "https://test.example.com/webhook",
					"events": ["email.sent"]
				}
			]
		}`
		fmt.Fprintf(w, ret)
	})

	limit := 10
	after := "cursor123"
	options := &ListOptions{
		Limit: &limit,
		After: &after,
	}

	resp, err := client.Webhooks.ListWithOptions(context.Background(), options)
	if err != nil {
		t.Errorf("Webhooks.ListWithOptions returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, true, resp.HasMore)
	assert.Equal(t, 1, len(resp.Data))
}

func TestRemoveWebhook(t *testing.T) {
	setup()
	defer teardown()

	webhookId := "4dd369bc-aa82-4ff3-97de-514ae3000ee0"

	mux.HandleFunc(fmt.Sprintf("/webhooks/%s", webhookId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "webhook",
			"id": "4dd369bc-aa82-4ff3-97de-514ae3000ee0",
			"deleted": true
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Webhooks.Remove(webhookId)
	if err != nil {
		t.Errorf("Webhooks.Remove returned error: %v", err)
	}
	assert.Equal(t, "webhook", resp.Object)
	assert.Equal(t, "4dd369bc-aa82-4ff3-97de-514ae3000ee0", resp.Id)
	assert.Equal(t, true, resp.Deleted)
}

func TestRemoveWebhookWithContext(t *testing.T) {
	setup()
	defer teardown()

	webhookId := "test-delete-id"

	mux.HandleFunc(fmt.Sprintf("/webhooks/%s", webhookId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "webhook",
			"id": "test-delete-id",
			"deleted": true
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	resp, err := client.Webhooks.RemoveWithContext(ctx, webhookId)
	if err != nil {
		t.Errorf("Webhooks.RemoveWithContext returned error: %v", err)
	}
	assert.Equal(t, "webhook", resp.Object)
	assert.Equal(t, "test-delete-id", resp.Id)
	assert.Equal(t, true, resp.Deleted)
}
