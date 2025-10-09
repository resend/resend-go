package resend

import (
	"errors"
	"fmt"
)

// MissingRequiredFieldsError is used when a required field is missing before making an API request
type MissingRequiredFieldsError struct {
	message string
}

func (e *MissingRequiredFieldsError) Error() string {
	return fmt.Sprintf("%s", e.message)
}

// ErrRateLimit is a sentinel error for rate limit detection with errors.Is
var ErrRateLimit = errors.New("rate limit exceeded")

// RateLimitError represents a rate limit error with metadata from response headers
type RateLimitError struct {
	// Message is the error message from the API
	Message string

	// Limit is the maximum number of requests allowed in the current window (raw header value)
	Limit string

	// Remaining is the number of requests remaining in the current window (raw header value)
	Remaining string

	// Reset is the time when the rate limit will reset in seconds (raw header value)
	Reset string

	// RetryAfter is the recommended wait time before retrying in seconds (raw header value)
	RetryAfter string
}

// Error implements the error interface
func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limit exceeded: %s (limit: %s, remaining: %s, reset: %s, retry after: %s)",
		e.Message, e.Limit, e.Remaining, e.Reset, e.RetryAfter)
}

// Is implements errors.Is support for detecting rate limit errors
func (e *RateLimitError) Is(target error) bool {
	return target == ErrRateLimit
}

// BroadcastsSvc errors
var (
	ErrFailedToCreateBroadcastUpdateRequest = errors.New("[ERROR]: Failed to create Broadcasts.Update request")
	ErrFailedToCreateBroadcastSendRequest   = errors.New("[ERROR]: Failed to create Broadcasts.Send request")
	ErrFailedToCreateBroadcastCreateRequest = errors.New("[ERROR]: Failed to create Broadcasts.Create request")
)

// ApiKeySvc errors
var (
	ErrFailedToCreateApiKeysCreateRequest = errors.New("[ERROR]: Failed to create ApiKeys.Create request")
	ErrFailedToCreateApiKeysListRequest   = errors.New("[ERROR]: Failed to create ApiKeys.List request")
	ErrFailedToCreateApiKeysRemoveRequest = errors.New("[ERROR]: Failed to create ApiKeys.Remove request")
)

// EmailsSvc errors
var (
	ErrFailedToCreateUpdateEmailRequest = errors.New("[ERROR]: Failed to create UpdateEmail request")
	ErrFailedToCreateEmailsSendRequest  = errors.New("[ERROR]: Failed to create SendEmail request")
	ErrFailedToCreateEmailsGetRequest   = errors.New("[ERROR]: Failed to create GetEmail request")
	ErrFailedToCreateEmailsListRequest  = errors.New("[ERROR]: Failed to create ListEmails request")
)
