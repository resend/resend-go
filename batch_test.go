package resend

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatchSendEmail(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/batch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &BatchEmailResponse{
			Data: []SendEmailResponse{
				{Id: "1"},
				{Id: "2"},
			},
		}
		err := json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	req := []*SendEmailRequest{
		{
			To: []string{"d@e.com"},
		},
		{
			To: []string{"d@e.com"},
		},
	}
	resp, err := client.Batch.Send(req)
	if err != nil {
		t.Errorf("BatchEmail.Send returned error: %v", err)
	}
	assert.Equal(t, resp.Data[0].Id, "1")
	assert.Equal(t, resp.Data[1].Id, "2")
	// Verify Errors field is nil (backward compatibility)
	assert.Nil(t, resp.Errors)
}

func TestBatchSendEmailWithErrors(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/batch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Simulate a permissive mode response with both successful and failed emails
		ret := &BatchEmailResponse{
			Data: []SendEmailResponse{
				{Id: "success-1"},
				{Id: "success-2"},
			},
			Errors: []BatchError{
				{Index: 2, Message: "The `to` field is missing."},
				{Index: 3, Message: "Invalid email address."},
			},
		}
		err := json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	req := []*SendEmailRequest{
		{To: []string{"valid1@example.com"}},
		{To: []string{"valid2@example.com"}},
		{To: []string{}},                // Missing 'to' field
		{To: []string{"invalid-email"}}, // Invalid email
	}

	resp, err := client.Batch.Send(req)
	if err != nil {
		t.Errorf("BatchEmail.Send returned error: %v", err)
	}

	// Verify successful responses
	assert.Equal(t, len(resp.Data), 2)
	assert.Equal(t, resp.Data[0].Id, "success-1")
	assert.Equal(t, resp.Data[1].Id, "success-2")

	// Verify error responses
	assert.NotNil(t, resp.Errors)
	assert.Equal(t, len(resp.Errors), 2)
	assert.Equal(t, resp.Errors[0].Index, 2)
	assert.Equal(t, resp.Errors[0].Message, "The `to` field is missing.")
	assert.Equal(t, resp.Errors[1].Index, 3)
	assert.Equal(t, resp.Errors[1].Message, "Invalid email address.")
}

func TestBatchSendEmailWithTagsAndScheduledAt(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/batch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		// Decode the request body to verify tags and scheduled_at are serialized per email
		var body []map[string]any
		err := json.NewDecoder(r.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Len(t, body, 2)

		// First email: relative scheduledAt + tags
		assert.Equal(t, "in 1 hour", body[0]["scheduled_at"])
		tags0, ok := body[0]["tags"].([]any)
		assert.True(t, ok)
		assert.Len(t, tags0, 1)
		tag0 := tags0[0].(map[string]any)
		assert.Equal(t, "category", tag0["name"])
		assert.Equal(t, "welcome", tag0["value"])

		// Second email: ISO 8601 scheduledAt + multiple tags
		assert.Equal(t, "2024-08-05T11:52:01.858Z", body[1]["scheduled_at"])
		tags1, ok := body[1]["tags"].([]any)
		assert.True(t, ok)
		assert.Len(t, tags1, 2)
		tag1a := tags1[0].(map[string]any)
		assert.Equal(t, "category", tag1a["name"])
		assert.Equal(t, "confirm", tag1a["value"])
		tag1b := tags1[1].(map[string]any)
		assert.Equal(t, "env", tag1b["name"])
		assert.Equal(t, "prod", tag1b["value"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &BatchEmailResponse{
			Data: []SendEmailResponse{
				{Id: "1"},
				{Id: "2"},
			},
		}
		err = json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	req := []*SendEmailRequest{
		{
			To:          []string{"d@e.com"},
			ScheduledAt: "in 1 hour",
			Tags: []Tag{
				{Name: "category", Value: "welcome"},
			},
		},
		{
			To:          []string{"d@e.com"},
			ScheduledAt: "2024-08-05T11:52:01.858Z",
			Tags: []Tag{
				{Name: "category", Value: "confirm"},
				{Name: "env", Value: "prod"},
			},
		},
	}

	resp, err := client.Batch.Send(req)
	if err != nil {
		t.Errorf("BatchEmail.Send returned error: %v", err)
	}
	assert.Equal(t, resp.Data[0].Id, "1")
	assert.Equal(t, resp.Data[1].Id, "2")
}

func TestBatchSendWithOptionsEmail(t *testing.T) {
	setup()
	defer teardown()
	ctx := context.Background()

	mux.HandleFunc("/emails/batch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		// Verify Idempotency-Key header is set
		assert.Equal(t, "1234567890", r.Header.Get("Idempotency-Key"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &BatchEmailResponse{
			Data: []SendEmailResponse{
				{Id: "1"},
				{Id: "2"},
			},
		}
		err := json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	req := []*SendEmailRequest{
		{
			To: []string{"d@e.com"},
		},
		{
			To: []string{"d@e.com"},
		},
	}

	options := &BatchSendEmailOptions{
		IdempotencyKey: "1234567890",
	}

	resp, err := client.Batch.SendWithOptions(ctx, req, options)
	if err != nil {
		t.Errorf("BatchEmail.SendWithOptions returned error: %v", err)
	}
	assert.Equal(t, resp.Data[0].Id, "1")
	assert.Equal(t, resp.Data[1].Id, "2")
	// Verify Errors field is nil (backward compatibility)
	assert.Nil(t, resp.Errors)
}

func TestBatchSendWithValidationMode(t *testing.T) {
	setup()
	defer teardown()
	ctx := context.Background()

	mux.HandleFunc("/emails/batch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		// Verify x-batch-validation header is set to permissive
		assert.Equal(t, "permissive", r.Header.Get("x-batch-validation"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Simulate permissive mode response with errors
		ret := &BatchEmailResponse{
			Data: []SendEmailResponse{
				{Id: "success-1"},
			},
			Errors: []BatchError{
				{Index: 1, Message: "Invalid email format"},
			},
		}
		err := json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	req := []*SendEmailRequest{
		{To: []string{"valid@example.com"}},
		{To: []string{"invalid"}},
	}

	options := &BatchSendEmailOptions{
		BatchValidation: "permissive",
	}

	resp, err := client.Batch.SendWithOptions(ctx, req, options)
	if err != nil {
		t.Errorf("BatchEmail.SendWithOptions returned error: %v", err)
	}

	// Verify successful response
	assert.Equal(t, len(resp.Data), 1)
	assert.Equal(t, resp.Data[0].Id, "success-1")

	// Verify error response
	assert.NotNil(t, resp.Errors)
	assert.Equal(t, len(resp.Errors), 1)
	assert.Equal(t, resp.Errors[0].Index, 1)
	assert.Equal(t, resp.Errors[0].Message, "Invalid email format")
}

func TestBatchSendWithStrictValidation(t *testing.T) {
	setup()
	defer teardown()
	ctx := context.Background()

	mux.HandleFunc("/emails/batch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		// Verify x-batch-validation header is set to strict
		assert.Equal(t, "strict", r.Header.Get("x-batch-validation"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &BatchEmailResponse{
			Data: []SendEmailResponse{
				{Id: "1"},
				{Id: "2"},
			},
		}
		err := json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	req := []*SendEmailRequest{
		{To: []string{"test1@example.com"}},
		{To: []string{"test2@example.com"}},
	}

	options := &BatchSendEmailOptions{
		BatchValidation: "strict",
	}

	resp, err := client.Batch.SendWithOptions(ctx, req, options)
	if err != nil {
		t.Errorf("BatchEmail.SendWithOptions returned error: %v", err)
	}

	assert.Equal(t, len(resp.Data), 2)
	assert.Nil(t, resp.Errors)
}

func TestBatchSendWithInvalidValidationMode(t *testing.T) {
	setup()
	defer teardown()
	ctx := context.Background()

	req := []*SendEmailRequest{
		{To: []string{"test@example.com"}},
	}

	// Test with invalid validation mode
	options := &BatchSendEmailOptions{
		BatchValidation: "invalid-mode",
	}

	resp, err := client.Batch.SendWithOptions(ctx, req, options)

	// Should return an error for invalid validation mode
	assert.NotNil(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "BatchValidation must be either BatchValidationStrict or BatchValidationPermissive")
}
