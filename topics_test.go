package resend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTopic(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify it
		var req CreateTopicRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		assert.Equal(t, "Weekly Newsletter", req.Name)
		assert.Equal(t, DefaultSubscriptionOptIn, req.DefaultSubscription)

		ret := `
		{
			"id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Topics.Create(&CreateTopicRequest{
		Name:                "Weekly Newsletter",
		DefaultSubscription: DefaultSubscriptionOptIn,
	})
	if err != nil {
		t.Errorf("Topics.Create returned error: %v", err)
	}
	assert.Equal(t, "b6d24b8e-af0b-4c3c-be0c-359bbd97381e", resp.Id)
}

func TestCreateTopicWithOptOut(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify it
		var req CreateTopicRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		assert.Equal(t, "Product Updates", req.Name)
		assert.Equal(t, DefaultSubscriptionOptOut, req.DefaultSubscription)

		ret := `
		{
			"id": "opt-out-topic-id"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Topics.Create(&CreateTopicRequest{
		Name:                "Product Updates",
		DefaultSubscription: DefaultSubscriptionOptOut,
	})
	if err != nil {
		t.Errorf("Topics.Create returned error: %v", err)
	}
	assert.Equal(t, "opt-out-topic-id", resp.Id)
}

func TestCreateTopicWithDescription(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify it
		var req CreateTopicRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		assert.Equal(t, "Monthly Summary", req.Name)
		assert.Equal(t, DefaultSubscriptionOptIn, req.DefaultSubscription)
		assert.Equal(t, "Monthly summary of your account activity", req.Description)

		ret := `
		{
			"id": "topic-with-description-id"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Topics.Create(&CreateTopicRequest{
		Name:                "Monthly Summary",
		DefaultSubscription: DefaultSubscriptionOptIn,
		Description:         "Monthly summary of your account activity",
	})
	if err != nil {
		t.Errorf("Topics.Create returned error: %v", err)
	}
	assert.Equal(t, "topic-with-description-id", resp.Id)
}

func TestCreateTopicWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "context-topic-id"
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	resp, err := client.Topics.CreateWithContext(ctx, &CreateTopicRequest{
		Name:                "Context Topic",
		DefaultSubscription: DefaultSubscriptionOptIn,
	})
	if err != nil {
		t.Errorf("Topics.CreateWithContext returned error: %v", err)
	}
	assert.Equal(t, "context-topic-id", resp.Id)
}

func TestGetTopic(t *testing.T) {
	setup()
	defer teardown()

	topicID := "b6d24b8e-af0b-4c3c-be0c-359bbd97381e"

	mux.HandleFunc(fmt.Sprintf("/topics/%s", topicID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
			"name": "Weekly Newsletter",
			"description": "Weekly newsletter for our subscribers",
			"default_subscription": "opt_in",
			"created_at": "2023-04-08T00:11:13.110779+00:00"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Topics.Get(topicID)
	if err != nil {
		t.Errorf("Topics.Get returned error: %v", err)
	}
	assert.Equal(t, "b6d24b8e-af0b-4c3c-be0c-359bbd97381e", resp.Id)
	assert.Equal(t, "Weekly Newsletter", resp.Name)
	assert.Equal(t, "Weekly newsletter for our subscribers", resp.Description)
	assert.Equal(t, DefaultSubscriptionOptIn, resp.DefaultSubscription)
	assert.Equal(t, "2023-04-08T00:11:13.110779+00:00", resp.CreatedAt)
}

func TestGetTopicWithOptOut(t *testing.T) {
	setup()
	defer teardown()

	topicID := "opt-out-topic-id"

	mux.HandleFunc(fmt.Sprintf("/topics/%s", topicID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "opt-out-topic-id",
			"name": "Product Updates",
			"description": "",
			"default_subscription": "opt_out",
			"created_at": "2023-04-08T00:11:13.110779+00:00"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Topics.Get(topicID)
	if err != nil {
		t.Errorf("Topics.Get returned error: %v", err)
	}
	assert.Equal(t, "opt-out-topic-id", resp.Id)
	assert.Equal(t, "Product Updates", resp.Name)
	assert.Equal(t, "", resp.Description)
	assert.Equal(t, DefaultSubscriptionOptOut, resp.DefaultSubscription)
	assert.Equal(t, "2023-04-08T00:11:13.110779+00:00", resp.CreatedAt)
}

func TestGetTopicWithContext(t *testing.T) {
	setup()
	defer teardown()

	topicID := "context-topic-id"

	mux.HandleFunc(fmt.Sprintf("/topics/%s", topicID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "context-topic-id",
			"name": "Context Topic",
			"description": "Test topic with context",
			"default_subscription": "opt_in",
			"created_at": "2023-04-08T00:11:13.110779+00:00"
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	resp, err := client.Topics.GetWithContext(ctx, topicID)
	if err != nil {
		t.Errorf("Topics.GetWithContext returned error: %v", err)
	}
	assert.Equal(t, "context-topic-id", resp.Id)
	assert.Equal(t, "Context Topic", resp.Name)
	assert.Equal(t, "Test topic with context", resp.Description)
	assert.Equal(t, DefaultSubscriptionOptIn, resp.DefaultSubscription)
	assert.Equal(t, "2023-04-08T00:11:13.110779+00:00", resp.CreatedAt)
}
