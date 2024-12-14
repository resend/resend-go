package resend

import "errors"

// BroadcastsSvc errors
var (
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
)
