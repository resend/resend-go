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

	topicId := "b6d24b8e-af0b-4c3c-be0c-359bbd97381e"

	mux.HandleFunc(fmt.Sprintf("/topics/%s", topicId), func(w http.ResponseWriter, r *http.Request) {
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

	resp, err := client.Topics.Get(topicId)
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

	topicId := "opt-out-topic-id"

	mux.HandleFunc(fmt.Sprintf("/topics/%s", topicId), func(w http.ResponseWriter, r *http.Request) {
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

	resp, err := client.Topics.Get(topicId)
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

	topicId := "context-topic-id"

	mux.HandleFunc(fmt.Sprintf("/topics/%s", topicId), func(w http.ResponseWriter, r *http.Request) {
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
	resp, err := client.Topics.GetWithContext(ctx, topicId)
	if err != nil {
		t.Errorf("Topics.GetWithContext returned error: %v", err)
	}
	assert.Equal(t, "context-topic-id", resp.Id)
	assert.Equal(t, "Context Topic", resp.Name)
	assert.Equal(t, "Test topic with context", resp.Description)
	assert.Equal(t, DefaultSubscriptionOptIn, resp.DefaultSubscription)
	assert.Equal(t, "2023-04-08T00:11:13.110779+00:00", resp.CreatedAt)
}

func TestListTopics(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "2", query.Get("limit"))

		ret := `
		{
			"object": "list",
			"has_more": false,
			"data": [
				{
					"id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
					"name": "Weekly Newsletter",
					"description": "Weekly newsletter for our subscribers",
					"default_subscription": "opt_in",
					"created_at": "2023-04-08T00:11:13.110779+00:00"
				},
				{
					"id": "topic-2-id",
					"name": "Product Updates",
					"description": "",
					"default_subscription": "opt_out",
					"created_at": "2023-04-09T00:11:13.110779+00:00"
				}
			]
		}`
		fmt.Fprintf(w, ret)
	})

	limit := 2
	resp, err := client.Topics.List(&ListOptions{
		Limit: &limit,
	})
	if err != nil {
		t.Errorf("Topics.List returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.False(t, resp.HasMore)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "b6d24b8e-af0b-4c3c-be0c-359bbd97381e", resp.Data[0].Id)
	assert.Equal(t, "Weekly Newsletter", resp.Data[0].Name)
	assert.Equal(t, "Weekly newsletter for our subscribers", resp.Data[0].Description)
	assert.Equal(t, DefaultSubscriptionOptIn, resp.Data[0].DefaultSubscription)
	assert.Equal(t, "topic-2-id", resp.Data[1].Id)
	assert.Equal(t, "Product Updates", resp.Data[1].Name)
	assert.Equal(t, DefaultSubscriptionOptOut, resp.Data[1].DefaultSubscription)
}

func TestListTopicsWithAfter(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "2", query.Get("limit"))
		assert.Equal(t, "topic-1-id", query.Get("after"))

		ret := `
		{
			"object": "list",
			"has_more": true,
			"data": [
				{
					"id": "topic-2-id",
					"name": "Next Topic",
					"description": "Next topic description",
					"default_subscription": "opt_in",
					"created_at": "2023-04-10T00:11:13.110779+00:00"
				}
			]
		}`
		fmt.Fprintf(w, ret)
	})

	limit := 2
	after := "topic-1-id"
	resp, err := client.Topics.List(&ListOptions{
		Limit: &limit,
		After: &after,
	})
	if err != nil {
		t.Errorf("Topics.List returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.True(t, resp.HasMore)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "topic-2-id", resp.Data[0].Id)
}

func TestListTopicsWithBefore(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "2", query.Get("limit"))
		assert.Equal(t, "topic-3-id", query.Get("before"))

		ret := `
		{
			"object": "list",
			"has_more": false,
			"data": [
				{
					"id": "topic-1-id",
					"name": "Previous Topic",
					"description": "Previous topic description",
					"default_subscription": "opt_out",
					"created_at": "2023-04-07T00:11:13.110779+00:00"
				}
			]
		}`
		fmt.Fprintf(w, ret)
	})

	limit := 2
	before := "topic-3-id"
	resp, err := client.Topics.List(&ListOptions{
		Limit:  &limit,
		Before: &before,
	})
	if err != nil {
		t.Errorf("Topics.List returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.False(t, resp.HasMore)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "topic-1-id", resp.Data[0].Id)
}

func TestListTopicsWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
			"has_more": false,
			"data": [
				{
					"id": "context-topic-id",
					"name": "Context Topic",
					"description": "Context topic description",
					"default_subscription": "opt_in",
					"created_at": "2023-04-08T00:11:13.110779+00:00"
				}
			]
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	resp, err := client.Topics.ListWithContext(ctx, &ListOptions{})
	if err != nil {
		t.Errorf("Topics.ListWithContext returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.False(t, resp.HasMore)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "context-topic-id", resp.Data[0].Id)
}

func TestListTopicsWithoutOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Check that there are no query parameters
		query := r.URL.Query()
		assert.Equal(t, "", query.Get("limit"))
		assert.Equal(t, "", query.Get("after"))
		assert.Equal(t, "", query.Get("before"))

		ret := `
		{
			"object": "list",
			"has_more": false,
			"data": [
				{
					"id": "topic-1-id",
					"name": "Topic 1",
					"description": "First topic",
					"default_subscription": "opt_in",
					"created_at": "2023-04-08T00:11:13.110779+00:00"
				},
				{
					"id": "topic-2-id",
					"name": "Topic 2",
					"description": "Second topic",
					"default_subscription": "opt_out",
					"created_at": "2023-04-09T00:11:13.110779+00:00"
				}
			]
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Topics.List(nil)
	if err != nil {
		t.Errorf("Topics.List returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.False(t, resp.HasMore)
	assert.Equal(t, 2, len(resp.Data))
}

func TestUpdateTopic(t *testing.T) {
	setup()
	defer teardown()

	topicId := "b6d24b8e-af0b-4c3c-be0c-359bbd97381e"

	mux.HandleFunc(fmt.Sprintf("/topics/%s", topicId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify it
		var req UpdateTopicRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		assert.Equal(t, "Weekly Newsletter - Updated", req.Name)
		assert.Equal(t, "Updated description", req.Description)

		ret := `
		{
			"id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Topics.Update(topicId, &UpdateTopicRequest{
		Name:        "Weekly Newsletter - Updated",
		Description: "Updated description",
	})
	if err != nil {
		t.Errorf("Topics.Update returned error: %v", err)
	}
	assert.Equal(t, "b6d24b8e-af0b-4c3c-be0c-359bbd97381e", resp.Id)
}

func TestUpdateTopicNameOnly(t *testing.T) {
	setup()
	defer teardown()

	topicId := "topic-name-only-id"

	mux.HandleFunc(fmt.Sprintf("/topics/%s", topicId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify it
		var req UpdateTopicRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		assert.Equal(t, "New Name", req.Name)
		assert.Equal(t, "", req.Description)

		ret := `
		{
			"id": "topic-name-only-id"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Topics.Update(topicId, &UpdateTopicRequest{
		Name: "New Name",
	})
	if err != nil {
		t.Errorf("Topics.Update returned error: %v", err)
	}
	assert.Equal(t, "topic-name-only-id", resp.Id)
}

func TestUpdateTopicDescriptionOnly(t *testing.T) {
	setup()
	defer teardown()

	topicId := "topic-description-only-id"

	mux.HandleFunc(fmt.Sprintf("/topics/%s", topicId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify it
		var req UpdateTopicRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		assert.Equal(t, "", req.Name)
		assert.Equal(t, "New description only", req.Description)

		ret := `
		{
			"id": "topic-description-only-id"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Topics.Update(topicId, &UpdateTopicRequest{
		Description: "New description only",
	})
	if err != nil {
		t.Errorf("Topics.Update returned error: %v", err)
	}
	assert.Equal(t, "topic-description-only-id", resp.Id)
}

func TestUpdateTopicWithContext(t *testing.T) {
	setup()
	defer teardown()

	topicId := "context-update-id"

	mux.HandleFunc(fmt.Sprintf("/topics/%s", topicId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "context-update-id"
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	resp, err := client.Topics.UpdateWithContext(ctx, topicId, &UpdateTopicRequest{
		Name:        "Context Update",
		Description: "Updated with context",
	})
	if err != nil {
		t.Errorf("Topics.UpdateWithContext returned error: %v", err)
	}
	assert.Equal(t, "context-update-id", resp.Id)
}

func TestRemoveTopic(t *testing.T) {
	setup()
	defer teardown()

	topicId := "b6d24b8e-af0b-4c3c-be0c-359bbd97381e"

	mux.HandleFunc(fmt.Sprintf("/topics/%s", topicId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "topic",
			"id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
			"deleted": true
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Topics.Remove(topicId)
	if err != nil {
		t.Errorf("Topics.Remove returned error: %v", err)
	}
	assert.Equal(t, "topic", resp.Object)
	assert.Equal(t, "b6d24b8e-af0b-4c3c-be0c-359bbd97381e", resp.Id)
	assert.True(t, resp.Deleted)
}

func TestRemoveTopicWithContext(t *testing.T) {
	setup()
	defer teardown()

	topicId := "context-remove-id"

	mux.HandleFunc(fmt.Sprintf("/topics/%s", topicId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "topic",
			"id": "context-remove-id",
			"deleted": true
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	resp, err := client.Topics.RemoveWithContext(ctx, topicId)
	if err != nil {
		t.Errorf("Topics.RemoveWithContext returned error: %v", err)
	}
	assert.Equal(t, "topic", resp.Object)
	assert.Equal(t, "context-remove-id", resp.Id)
	assert.True(t, resp.Deleted)
}
