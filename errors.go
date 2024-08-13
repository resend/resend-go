package resend

import "errors"

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
