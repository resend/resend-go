package resend

import (
	"context"
	"errors"
	"net/http"
)

// BatchValidationMode represents the validation mode for batch emails
type BatchValidationMode string

const (
	// BatchValidationStrict only sends the batch if all emails are valid
	BatchValidationStrict BatchValidationMode = "strict"
	// BatchValidationPermissive processes all emails, allowing partial success
	BatchValidationPermissive BatchValidationMode = "permissive"
)

// IsValid checks if the BatchValidationMode has a valid value
func (b BatchValidationMode) IsValid() bool {
	return b == BatchValidationStrict || b == BatchValidationPermissive
}

// String returns the string representation of the BatchValidationMode
func (b BatchValidationMode) String() string {
	return string(b)
}

// BatchSendEmailOptions is the additional options struct for the Batch.SendEmail call.
type BatchSendEmailOptions struct {
	IdempotencyKey string `json:"idempotency_key,omitempty"`
	// BatchValidation controls the validation behavior for batch emails.
	// Can be BatchValidationStrict (default) or BatchValidationPermissive.
	// - BatchValidationStrict: Only sends the batch if all emails are valid
	// - BatchValidationPermissive: Processes all emails, allowing partial success
	BatchValidation BatchValidationMode `json:"-"`
}

// GetIdempotencyKey returns the idempotency key for the batch send email request.
func (o BatchSendEmailOptions) GetIdempotencyKey() string {
	return o.IdempotencyKey
}

// GetBatchValidation returns the batch validation mode for the batch send email request.
func (o BatchSendEmailOptions) GetBatchValidation() string {
	return o.BatchValidation.String()
}

// BatchError represents an error for a specific email in a batch request
// when using permissive validation mode.
type BatchError struct {
	Index   int    `json:"index"`
	Message string `json:"message"`
}

// BatchEmailResponse is the response from the BatchSendEmail call.
// see https://resend.com/docs/api-reference/emails/send-batch-emails
type BatchEmailResponse struct {
	Data   []SendEmailResponse `json:"data"`
	Errors []BatchError        `json:"errors,omitempty"`
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
	// Validate BatchValidation field if provided
	if options != nil && options.BatchValidation != "" {
		if !options.BatchValidation.IsValid() {
			return nil, errors.New("[ERROR]: BatchValidation must be either BatchValidationStrict or BatchValidationPermissive")
		}
	}

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
