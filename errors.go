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
