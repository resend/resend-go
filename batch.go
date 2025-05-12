package resend

import (
	"context"
	"errors"
	"net/http"
)

// BatchSendEmailOptions is the additional options struct for the Batch.SendEmail call.
type BatchSendEmailOptions struct {
	IdempotencyKey string `json:"idempotency_key,omitempty"`
}

// GetIdempotencyKey returns the idempotency key for the batch send email request.
func (o BatchSendEmailOptions) GetIdempotencyKey() string {
	return o.IdempotencyKey
}

// BatchEmailResponse is the response from the BatchSendEmail call.
// see https://resend.com/docs/api-reference/emails/send-batch-emails
type BatchEmailResponse struct {
	Data []SendEmailResponse `json:"data"`
}

type BatchSvc interface {
	Send([]*SendEmailRequest) (*BatchEmailResponse, error)
	SendWithContext(ctx context.Context, params []*SendEmailRequest) (*BatchEmailResponse, error)
	SendWithOptions(ctx context.Context, params []*SendEmailRequest, options *BatchSendEmailOptions) (*BatchEmailResponse, error)
}

type BatchSvcImpl struct {
	client *Client
}

// Send send a batch of emails
// https://resend.com/docs/api-reference/emails/send-batch-emails
func (s *BatchSvcImpl) Send(params []*SendEmailRequest) (*BatchEmailResponse, error) {
	return s.SendWithContext(context.Background(), params)
}

// SendWithContext is the same as Send but accepts a ctx as argument
func (s *BatchSvcImpl) SendWithContext(ctx context.Context, params []*SendEmailRequest) (*BatchEmailResponse, error) {
	path := "emails/batch"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return nil, errors.New("[ERROR]: Failed to create BatchEmail request")
	}

	// Build response recipient obj
	batchSendEmailResponse := new(BatchEmailResponse)

	// Send Request
	_, err = s.client.Perform(req, batchSendEmailResponse)

	if err != nil {
		return nil, err
	}

	return batchSendEmailResponse, nil
}

// SendWithOptions is the same as Send but accepts a ctx and options as arguments
func (s *BatchSvcImpl) SendWithOptions(ctx context.Context, params []*SendEmailRequest, options *BatchSendEmailOptions) (*BatchEmailResponse, error) {
	path := "emails/batch"

	// Prepare request
	req, err := s.client.NewRequestWithOptions(ctx, http.MethodPost, path, params, options)
	if err != nil {
		return nil, errors.New("[ERROR]: Failed to create BatchEmail request")
	}

	// Build response recipient obj
	batchSendEmailResponse := new(BatchEmailResponse)

	// Send Request
	_, err = s.client.Perform(req, batchSendEmailResponse)

	if err != nil {
		return nil, err
	}

	return batchSendEmailResponse, nil
}
