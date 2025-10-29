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
