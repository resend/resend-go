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
	return e.message
}

// ErrRateLimit is a sentinel error for rate limit detection with errors.Is
var ErrRateLimit = errors.New("rate limit exceeded")

// ErrResendAPI is the base error for API responses with messages.
var (
	ErrResendAPI     = errors.New("[ERROR]")
	ErrResendUnknown = errors.New("[ERROR]: Unknown Error")
)

// ResendError represents a message-based error returned by the API.
type ResendError struct { //nolint:revive
	Message string
}

func (e *ResendError) Error() string {
	return ErrResendAPI.Error() + ": " + e.Message
}

// Is implements errors.Is support for detecting generic API errors.
func (e *ResendError) Is(target error) bool {
	return errors.Is(target, ErrResendAPI)
}

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
	return errors.Is(target, ErrRateLimit)
}

// BroadcastsSvc errors
var (
	ErrFailedToCreateBroadcastUpdateRequest = errors.New("[ERROR]: Failed to create Broadcasts.Update request")
	ErrFailedToCreateBroadcastSendRequest   = errors.New("[ERROR]: Failed to create Broadcasts.Send request")
	ErrFailedToCreateBroadcastCreateRequest = errors.New("[ERROR]: Failed to create Broadcasts.Create request")
	ErrFailedToCreateBroadcastGetRequest    = errors.New("[ERROR]: Failed to create Broadcast.Get request")
	ErrFailedToCreateBroadcastRemoveRequest = errors.New("[ERROR]: Failed to create Broadcast.Remove request")
	ErrFailedToCreateBroadcastsListRequest  = errors.New("[ERROR]: Failed to create Broadcasts.List request")
	ErrBroadcastSegmentOrAudienceRequired   = errors.New("[ERROR]: Either SegmentId or AudienceId must be provided")
	ErrBroadcastFromRequired                = errors.New("[ERROR]: From cannot be empty")
	ErrBroadcastSubjectRequired             = errors.New("[ERROR]: Subject cannot be empty")
	ErrBroadcastIDRequired                  = errors.New("[ERROR]: BroadcastId cannot be empty")
	ErrBroadcastIDRequiredLowercase         = errors.New("[ERROR]: broadcastId cannot be empty")
)

// ApiKeySvc errors
var (
	ErrFailedToCreateApiKeysCreateRequest = errors.New("[ERROR]: Failed to create ApiKeys.Create request") //nolint:revive
	ErrFailedToCreateApiKeysListRequest   = errors.New("[ERROR]: Failed to create ApiKeys.List request")   //nolint:revive
	ErrFailedToCreateApiKeysRemoveRequest = errors.New("[ERROR]: Failed to create ApiKeys.Remove request") //nolint:revive
)

// EmailsSvc errors
var (
	ErrFailedToCreateUpdateEmailRequest           = errors.New("[ERROR]: Failed to create UpdateEmail request")
	ErrFailedToCreateEmailsSendRequest            = errors.New("[ERROR]: Failed to create SendEmail request")
	ErrFailedToCreateEmailsGetRequest             = errors.New("[ERROR]: Failed to create GetEmail request")
	ErrFailedToCreateEmailsListRequest            = errors.New("[ERROR]: Failed to create ListEmails request")
	ErrFailedToCreateEmailsGetAttachmentRequest   = errors.New("[ERROR]: Failed to create Emails.GetAttachment request")
	ErrFailedToCreateEmailsListAttachmentsRequest = errors.New("[ERROR]: Failed to create Emails.ListAttachments request")
)

// TemplatesSvc errors
var (
	ErrFailedToCreateTemplateCreateRequest    = errors.New("[ERROR]: Failed to create Templates.Create request")
	ErrFailedToCreateTemplateGetRequest       = errors.New("[ERROR]: Failed to create Templates.Get request")
	ErrFailedToCreateTemplateListRequest      = errors.New("[ERROR]: Failed to create Templates.List request")
	ErrFailedToCreateTemplateUpdateRequest    = errors.New("[ERROR]: Failed to create Templates.Update request")
	ErrFailedToCreateTemplatePublishRequest   = errors.New("[ERROR]: Failed to create Templates.Publish request")
	ErrFailedToCreateTemplateDuplicateRequest = errors.New("[ERROR]: Failed to create Templates.Duplicate request")
	ErrFailedToCreateTemplateRemoveRequest    = errors.New("[ERROR]: Failed to create Templates.Remove request")
)

// ReceivingSvc errors
var (
	ErrFailedToCreateReceivingGetRequest             = errors.New("[ERROR]: Failed to create Receiving.Get request")
	ErrFailedToCreateReceivingListRequest            = errors.New("[ERROR]: Failed to create Receiving.List request")
	ErrFailedToCreateReceivingGetAttachmentRequest   = errors.New("[ERROR]: Failed to create Receiving.GetAttachment request")
	ErrFailedToCreateReceivingListAttachmentsRequest = errors.New("[ERROR]: Failed to create Receiving.ListAttachments request")
)

// TopicsSvc errors
var (
	ErrFailedToCreateTopicCreateRequest = errors.New("[ERROR]: Failed to create Topics.Create request")
	ErrFailedToCreateTopicGetRequest    = errors.New("[ERROR]: Failed to create Topics.Get request")
	ErrFailedToCreateTopicListRequest   = errors.New("[ERROR]: Failed to create Topics.List request")
	ErrFailedToCreateTopicUpdateRequest = errors.New("[ERROR]: Failed to create Topics.Update request")
	ErrFailedToCreateTopicRemoveRequest = errors.New("[ERROR]: Failed to create Topics.Remove request")
)

// BatchSvc errors
var (
	ErrFailedToCreateBatchEmailRequest = errors.New("[ERROR]: Failed to create BatchEmail request")
	ErrInvalidBatchValidation          = errors.New("[ERROR]: BatchValidation must be either BatchValidationStrict or BatchValidationPermissive")
)

// SegmentsSvc errors
var (
	ErrFailedToCreateSegmentsCreateRequest = errors.New("[ERROR]: Failed to create Segments.Create request")
	ErrFailedToCreateSegmentsListRequest   = errors.New("[ERROR]: Failed to create Segments.List request")
	ErrFailedToCreateSegmentRemoveRequest  = errors.New("[ERROR]: Failed to create Segment.Remove request")
	ErrFailedToCreateSegmentGetRequest     = errors.New("[ERROR]: Failed to create Segment.Get request")
)

// ContactSegmentsSvc errors
var (
	ErrContactSegmentIDRequired                   = errors.New("[ERROR]: SegmentId is required")
	ErrContactSegmentContactIDOrEmailRequired     = errors.New("[ERROR]: Either ContactId or Email must be provided")
	ErrFailedToCreateContactSegmentsAddRequest    = errors.New("[ERROR]: Failed to create ContactSegments.Add request")
	ErrFailedToCreateContactSegmentsRemoveRequest = errors.New("[ERROR]: Failed to create ContactSegments.Remove request")
	ErrFailedToCreateContactSegmentsListRequest   = errors.New("[ERROR]: Failed to create ContactSegments.List request")
)

// ContactPropertiesSvc errors
var (
	ErrContactPropertyKeyMissing                    = errors.New("[ERROR]: Key is missing")
	ErrContactPropertyTypeMissing                   = errors.New("[ERROR]: Type is missing")
	ErrContactPropertyIDMissing                     = errors.New("[ERROR]: ID is missing")
	ErrFailedToCreateContactPropertiesCreateRequest = errors.New("[ERROR]: Failed to create ContactProperties.Create request")
	ErrFailedToCreateContactPropertiesListRequest   = errors.New("[ERROR]: Failed to create ContactProperties.List request")
	ErrFailedToCreateContactPropertiesGetRequest    = errors.New("[ERROR]: Failed to create ContactProperties.Get request")
	ErrFailedToCreateContactPropertiesUpdateRequest = errors.New("[ERROR]: Failed to create ContactProperties.Update request")
	ErrFailedToCreateContactPropertiesRemoveRequest = errors.New("[ERROR]: Failed to create ContactProperties.Remove request")
)

// ContactTopicsSvc errors
var (
	ErrContactTopicsContactIDOrEmailMissing     = errors.New("[ERROR]: Contact ID or email is missing")
	ErrContactTopicsArrayEmpty                  = errors.New("[ERROR]: Topics array is empty")
	ErrFailedToCreateContactTopicsListRequest   = errors.New("[ERROR]: Failed to create ContactTopics.List request")
	ErrFailedToCreateContactTopicsUpdateRequest = errors.New("[ERROR]: Failed to create ContactTopics.Update request")
)

// ContactsSvc errors
var (
	ErrFailedToCreateContactsCreateRequest = errors.New("[ERROR]: Failed to create Contacts.Create request")
	ErrFailedToCreateContactsListRequest   = errors.New("[ERROR]: Failed to create Contacts.List request")
	ErrContactIDRequired                   = errors.New("[ERROR]: Id is required")
	ErrFailedToCreateContactRemoveRequest  = errors.New("[ERROR]: Failed to create Contact.Remove request")
	ErrFailedToCreateContactGetRequest     = errors.New("[ERROR]: Failed to create Contact.Get request")
	ErrFailedToCreateContactsUpdateRequest = errors.New("[ERROR]: Failed to create Contacts.Update request")
)

// DomainsSvc errors
var (
	ErrFailedToCreateDomainsUpdateRequest = errors.New("[ERROR]: Failed to create Domains.Update request")
	ErrFailedToCreateDomainsCreateRequest = errors.New("[ERROR]: Failed to create Domains.Create request")
	ErrFailedToCreateDomainsVerifyRequest = errors.New("[ERROR]: Failed to create Domains.Verify request")
	ErrFailedToCreateDomainsListRequest   = errors.New("[ERROR]: Failed to create Domains.List request")
	ErrFailedToCreateDomainsRemoveRequest = errors.New("[ERROR]: Failed to create Domains.Remove request")
	ErrFailedToCreateDomainsGetRequest    = errors.New("[ERROR]: Failed to create Domains.Get request")
)

// WebhooksSvc errors
var (
	ErrWebhookOptionsNil                = errors.New("options cannot be nil")
	ErrWebhookPayloadEmpty              = errors.New("payload cannot be empty")
	ErrWebhookSecretEmpty               = errors.New("webhook secret cannot be empty")
	ErrWebhookHeaderIDRequired          = errors.New("svix-id header is required")
	ErrWebhookHeaderTimestampRequired   = errors.New("svix-timestamp header is required")
	ErrWebhookHeaderSignatureRequired   = errors.New("svix-signature header is required")
	ErrWebhookSignatureNotFound         = errors.New("no matching signature found")
	ErrWebhookTimestampOutsideTolerance = errors.New("timestamp outside tolerance window")
)
