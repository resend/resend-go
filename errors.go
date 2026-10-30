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

// ResendError represents an API error other than a rate limit error.
// Use errors.As to inspect it. HTTP 429 responses return *RateLimitError instead.
type ResendError struct {
	// StatusCode is the HTTP response status, regardless of the response body.
	StatusCode int `json:"statusCode"`
	// Name is the API error name, when available.
	Name string `json:"name"`
	// Message is the API error message or a fallback for an unreadable body.
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *ResendError) Error() string {
	return "[ERROR]: " + e.Message
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
	ErrFailedToCreateBroadcastUpdateRequest     = errors.New("[ERROR]: Failed to create Broadcasts.Update request")
	ErrFailedToCreateBroadcastSendRequest       = errors.New("[ERROR]: Failed to create Broadcasts.Send request")
	ErrFailedToCreateBroadcastCreateRequest     = errors.New("[ERROR]: Failed to create Broadcasts.Create request")
	ErrFailedToCreateBroadcastCancelRequest     = errors.New("[ERROR]: Failed to create Broadcasts.Cancel request")
	ErrFailedToCreateBroadcastDuplicateRequest  = errors.New("[ERROR]: Failed to create Broadcasts.Duplicate request")
	ErrFailedToCreateBroadcastRecipientsRequest = errors.New("[ERROR]: Failed to create Broadcasts.Recipients request")
)

// ApiKeySvc errors
var (
	ErrFailedToCreateApiKeysCreateRequest = errors.New("[ERROR]: Failed to create ApiKeys.Create request")
	ErrFailedToCreateApiKeysListRequest   = errors.New("[ERROR]: Failed to create ApiKeys.List request")
	ErrFailedToCreateApiKeysRemoveRequest = errors.New("[ERROR]: Failed to create ApiKeys.Remove request")
	ErrFailedToCreateApiKeysUpdateRequest = errors.New("[ERROR]: Failed to create ApiKeys.Update request")
)

// EmailsSvc errors
var (
	ErrFailedToCreateUpdateEmailRequest           = errors.New("[ERROR]: Failed to create UpdateEmail request")
	ErrFailedToCreateEmailsSendRequest            = errors.New("[ERROR]: Failed to create SendEmail request")
	ErrFailedToCreateEmailsGetRequest             = errors.New("[ERROR]: Failed to create GetEmail request")
	ErrFailedToCreateEmailsListRequest            = errors.New("[ERROR]: Failed to create ListEmails request")
	ErrFailedToCreateEmailsGetAttachmentRequest   = errors.New("[ERROR]: Failed to create Emails.GetAttachment request")
	ErrFailedToCreateEmailsListAttachmentsRequest = errors.New("[ERROR]: Failed to create Emails.ListAttachments request")
	ErrFailedToCreateEmailsShareRequest           = errors.New("[ERROR]: Failed to create Emails.Share request")
	ErrFailedToCreateEmailsMetricsRequest         = errors.New("[ERROR]: Failed to create Emails.Metrics request")
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

// LogsSvc errors
var (
	ErrFailedToCreateLogsGetRequest  = errors.New("[ERROR]: Failed to create Logs.Get request")
	ErrFailedToCreateLogsListRequest = errors.New("[ERROR]: Failed to create Logs.List request")
)

// AutomationsSvc errors
var (
	ErrFailedToCreateAutomationCreateRequest    = errors.New("[ERROR]: Failed to create Automations.Create request")
	ErrFailedToCreateAutomationGetRequest       = errors.New("[ERROR]: Failed to create Automations.Get request")
	ErrFailedToCreateAutomationListRequest      = errors.New("[ERROR]: Failed to create Automations.List request")
	ErrFailedToCreateAutomationUpdateRequest    = errors.New("[ERROR]: Failed to create Automations.Update request")
	ErrFailedToCreateAutomationRemoveRequest    = errors.New("[ERROR]: Failed to create Automations.Remove request")
	ErrFailedToCreateAutomationDuplicateRequest = errors.New("[ERROR]: Failed to create Automations.Duplicate request")
	ErrFailedToCreateAutomationStopRequest      = errors.New("[ERROR]: Failed to create Automations.Stop request")
	ErrFailedToCreateAutomationListRunsRequest  = errors.New("[ERROR]: Failed to create Automations.ListRuns request")
	ErrFailedToCreateAutomationGetRunRequest    = errors.New("[ERROR]: Failed to create Automations.GetRun request")
)

// EventsSvc errors
var (
	ErrFailedToCreateEventCreateRequest = errors.New("[ERROR]: Failed to create Events.Create request")
	ErrFailedToCreateEventGetRequest    = errors.New("[ERROR]: Failed to create Events.Get request")
	ErrFailedToCreateEventListRequest   = errors.New("[ERROR]: Failed to create Events.List request")
	ErrFailedToCreateEventUpdateRequest = errors.New("[ERROR]: Failed to create Events.Update request")
	ErrFailedToCreateEventRemoveRequest = errors.New("[ERROR]: Failed to create Events.Remove request")
	ErrFailedToCreateEventSendRequest   = errors.New("[ERROR]: Failed to create Events.Send request")
)

// OAuthGrantsSvc errors
var (
	ErrFailedToCreateOAuthGrantsListRequest   = errors.New("[ERROR]: Failed to create OAuthGrants.List request")
	ErrFailedToCreateOAuthGrantsRevokeRequest = errors.New("[ERROR]: Failed to create OAuthGrants.Revoke request")
)

// UsageSvc errors
var (
	ErrFailedToCreateUsageGetRequest = errors.New("[ERROR]: Failed to create Usage.Get request")
)

// InboxesSvc errors
var (
	ErrFailedToCreateInboxesCreateRequest              = errors.New("[ERROR]: Failed to create Inboxes.Create request")
	ErrFailedToCreateInboxesListRequest                = errors.New("[ERROR]: Failed to create Inboxes.List request")
	ErrFailedToCreateInboxesGetRequest                 = errors.New("[ERROR]: Failed to create Inboxes.Get request")
	ErrFailedToCreateInboxesUpdateRequest              = errors.New("[ERROR]: Failed to create Inboxes.Update request")
	ErrFailedToCreateInboxesRemoveRequest              = errors.New("[ERROR]: Failed to create Inboxes.Remove request")
	ErrFailedToCreateInboxesAgentGetRequest            = errors.New("[ERROR]: Failed to create Inboxes.Agent.Get request")
	ErrFailedToCreateInboxesAgentUpdateRequest         = errors.New("[ERROR]: Failed to create Inboxes.Agent.Update request")
	ErrFailedToCreateInboxesLabelsListRequest          = errors.New("[ERROR]: Failed to create Inboxes.Labels.List request")
	ErrFailedToCreateInboxesLabelsCreateRequest        = errors.New("[ERROR]: Failed to create Inboxes.Labels.Create request")
	ErrFailedToCreateInboxesLabelsUpdateRequest        = errors.New("[ERROR]: Failed to create Inboxes.Labels.Update request")
	ErrFailedToCreateInboxesLabelsRemoveRequest        = errors.New("[ERROR]: Failed to create Inboxes.Labels.Remove request")
	ErrFailedToCreateInboxesThreadsListRequest         = errors.New("[ERROR]: Failed to create Inboxes.Threads.List request")
	ErrFailedToCreateInboxesThreadsSearchRequest       = errors.New("[ERROR]: Failed to create Inboxes.Threads.Search request")
	ErrFailedToCreateInboxesThreadsGetRequest          = errors.New("[ERROR]: Failed to create Inboxes.Threads.Get request")
	ErrFailedToCreateInboxesThreadsUpdateRequest       = errors.New("[ERROR]: Failed to create Inboxes.Threads.Update request")
	ErrFailedToCreateInboxesThreadsRemoveRequest       = errors.New("[ERROR]: Failed to create Inboxes.Threads.Remove request")
	ErrFailedToCreateInboxesThreadLabelsAddRequest     = errors.New("[ERROR]: Failed to create Inboxes.Threads.Labels.Add request")
	ErrFailedToCreateInboxesThreadLabelsSetRequest     = errors.New("[ERROR]: Failed to create Inboxes.Threads.Labels.Set request")
	ErrFailedToCreateInboxesThreadLabelsRemoveRequest  = errors.New("[ERROR]: Failed to create Inboxes.Threads.Labels.Remove request")
	ErrFailedToCreateInboxesThreadEmailsListRequest    = errors.New("[ERROR]: Failed to create Inboxes.Threads.Emails.List request")
	ErrFailedToCreateInboxesThreadEmailsGetRequest     = errors.New("[ERROR]: Failed to create Inboxes.Threads.Emails.Get request")
	ErrFailedToCreateInboxesThreadEmailsReplyRequest   = errors.New("[ERROR]: Failed to create Inboxes.Threads.Emails.Reply request")
	ErrFailedToCreateInboxesThreadEmailsForwardRequest = errors.New("[ERROR]: Failed to create Inboxes.Threads.Emails.Forward request")
	ErrFailedToCreateInboxesDraftsListRequest          = errors.New("[ERROR]: Failed to create Inboxes.Drafts.List request")
	ErrFailedToCreateInboxesDraftsCreateRequest        = errors.New("[ERROR]: Failed to create Inboxes.Drafts.Create request")
	ErrFailedToCreateInboxesDraftsGetRequest           = errors.New("[ERROR]: Failed to create Inboxes.Drafts.Get request")
	ErrFailedToCreateInboxesDraftsUpdateRequest        = errors.New("[ERROR]: Failed to create Inboxes.Drafts.Update request")
	ErrFailedToCreateInboxesDraftsRemoveRequest        = errors.New("[ERROR]: Failed to create Inboxes.Drafts.Remove request")
	ErrFailedToCreateInboxesDraftsSendRequest          = errors.New("[ERROR]: Failed to create Inboxes.Drafts.Send request")
)
