package resend

import (
	"reflect"
	"time"
)

type WebhookPayload struct {
	CreatedAt time.Time          `json:"created_at"`
	Type      string             `json:"type"`
	Data      WebhookPayloadData `json:"data"`
}

type WebhookPayloadData interface {
	IsWebhookPayloadData()
}

// Each event has its own concrete data type. Shared nested objects are modeled below.

type WebhookDomainCapabilities struct {
	Sending   string `json:"sending"`
	Receiving string `json:"receiving"`
}

type WebhookDomainRecord struct {
	Record   string `json:"record"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	TTL      string `json:"ttl"`
	Status   string `json:"status"`
	Value    string `json:"value"`
	Priority int    `json:"priority,omitempty"`
}

type Bounce struct {
	Message        string   `json:"message"`
	SubType        string   `json:"subType"`
	Type           string   `json:"type"`
	DiagnosticCode []string `json:"diagnosticCode,omitempty"`
}

type WebhookEmailClick struct {
	IPAddress string    `json:"ipAddress"`
	Link      string    `json:"link"`
	Timestamp time.Time `json:"timestamp"`
	UserAgent string    `json:"userAgent"`
}

type WebhookEmailFailure struct {
	Reason string `json:"reason"`
}

type WebhookEmailAttachment struct {
	ID                 string `json:"id"`
	Filename           string `json:"filename"`
	ContentType        string `json:"content_type"`
	ContentDisposition string `json:"content_disposition"`
	ContentID          string `json:"content_id"`
}

type WebhookEmailSuppression struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type WebhookInbox struct {
	Object           string     `json:"object"`
	ID               string     `json:"id"`
	Name             string     `json:"name"`
	EmailAddress     string     `json:"email_address"`
	DomainID         string     `json:"domain_id"`
	ReceivingAddress *string    `json:"receiving_address"`
	FromName         string     `json:"from_name"`
	Unread           int        `json:"unread"`
	Drafts           int        `json:"drafts"`
	LastReceived     *time.Time `json:"last_received"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type WebhookInboxThread struct {
	Object  string              `json:"object"`
	ID      string              `json:"id"`
	Subject string              `json:"subject"`
	Folder  string              `json:"folder"`
	Labels  []WebhookInboxLabel `json:"labels"`
	Read    bool                `json:"read"`
}

type WebhookInboxDraft struct {
	Object         string                   `json:"object"`
	ID             string                   `json:"id"`
	Type           string                   `json:"type"`
	To             []string                 `json:"to"`
	CC             []string                 `json:"cc"`
	BCC            []string                 `json:"bcc"`
	Subject        string                   `json:"subject"`
	ThreadID       string                   `json:"thread_id"`
	ReplyToEmailID string                   `json:"reply_to_email_id"`
	EmailID        *string                  `json:"email_id"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
	Attachments    []WebhookInboxAttachment `json:"attachments"`
}

type WebhookInboxEmail struct {
	ID          string                   `json:"id"`
	Direction   string                   `json:"direction"`
	From        string                   `json:"from"`
	To          []string                 `json:"to"`
	CC          []string                 `json:"cc"`
	BCC         []string                 `json:"bcc"`
	ReplyTo     []string                 `json:"reply_to"`
	Subject     string                   `json:"subject"`
	MessageID   string                   `json:"message_id"`
	Attachments []WebhookInboxAttachment `json:"attachments"`
	Read        bool                     `json:"read"`
	ReceivedAt  time.Time                `json:"received_at"`
}

type WebhookInboxAttachment struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
}

type WebhookInboxLabel struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// webhookDomainData contains the shared data for domain lifecycle events.
type webhookDomainData struct {
	ID           string                    `json:"id"`
	Name         string                    `json:"name"`
	Status       string                    `json:"status"`
	CreatedAt    time.Time                 `json:"created_at"`
	Region       string                    `json:"region"`
	Capabilities WebhookDomainCapabilities `json:"capabilities"`
	Records      []WebhookDomainRecord     `json:"records"`
}

// WebhookPayloadDomainCreated is the data payload for domain.created.
type WebhookPayloadDomainCreated struct{ webhookDomainData }

func (WebhookPayloadDomainCreated) IsWebhookPayloadData() {}

// WebhookPayloadDomainDeleted is the data payload for domain.deleted.
type WebhookPayloadDomainDeleted struct{ webhookDomainData }

func (WebhookPayloadDomainDeleted) IsWebhookPayloadData() {}

// WebhookPayloadDomainUpdated is the data payload for domain.updated.
type WebhookPayloadDomainUpdated struct{ webhookDomainData }

func (WebhookPayloadDomainUpdated) IsWebhookPayloadData() {}

// WebhookPayloadEmailBounced is the data payload for email.bounced.
type WebhookPayloadEmailBounced struct {
	BroadcastID string         `json:"broadcast_id"`
	CreatedAt   time.Time      `json:"created_at"`
	EmailID     string         `json:"email_id"`
	MessageID   string         `json:"message_id"`
	From        string         `json:"from"`
	To          []string       `json:"to"`
	Subject     string         `json:"subject"`
	TemplateID  string         `json:"template_id"`
	Bounce      Bounce         `json:"bounce"`
	Tags        map[string]any `json:"tags"`
}

func (WebhookPayloadEmailBounced) IsWebhookPayloadData() {}

// WebhookPayloadEmailClicked is the data payload for email.clicked.
type WebhookPayloadEmailClicked struct {
	BroadcastID string            `json:"broadcast_id"`
	CreatedAt   time.Time         `json:"created_at"`
	EmailID     string            `json:"email_id"`
	MessageID   string            `json:"message_id"`
	From        string            `json:"from"`
	To          []string          `json:"to"`
	Click       WebhookEmailClick `json:"click"`
	Subject     string            `json:"subject"`
	TemplateID  string            `json:"template_id"`
	Tags        map[string]any    `json:"tags"`
}

func (WebhookPayloadEmailClicked) IsWebhookPayloadData() {}

// webhookEmailDeliveryData contains the shared data for email delivery events
// that do not have an event-specific nested object.
type webhookEmailDeliveryData struct {
	BroadcastID string         `json:"broadcast_id"`
	CreatedAt   time.Time      `json:"created_at"`
	EmailID     string         `json:"email_id"`
	MessageID   string         `json:"message_id"`
	From        string         `json:"from"`
	To          []string       `json:"to"`
	Subject     string         `json:"subject"`
	TemplateID  string         `json:"template_id"`
	Tags        map[string]any `json:"tags"`
}

// WebhookPayloadEmailComplained is the data payload for email.complained.
type WebhookPayloadEmailComplained struct{ webhookEmailDeliveryData }

func (WebhookPayloadEmailComplained) IsWebhookPayloadData() {}

// WebhookPayloadEmailDelivered is the data payload for email.delivered.
type WebhookPayloadEmailDelivered struct{ webhookEmailDeliveryData }

func (WebhookPayloadEmailDelivered) IsWebhookPayloadData() {}

// WebhookPayloadEmailDeliveryDelayed is the data payload for email.delivery_delayed.
type WebhookPayloadEmailDeliveryDelayed struct{ webhookEmailDeliveryData }

func (WebhookPayloadEmailDeliveryDelayed) IsWebhookPayloadData() {}

// WebhookPayloadEmailFailed is the data payload for email.failed.
type WebhookPayloadEmailFailed struct {
	BroadcastID string              `json:"broadcast_id"`
	CreatedAt   time.Time           `json:"created_at"`
	EmailID     string              `json:"email_id"`
	From        string              `json:"from"`
	To          []string            `json:"to"`
	Subject     string              `json:"subject"`
	TemplateID  string              `json:"template_id"`
	Failed      WebhookEmailFailure `json:"failed"`
	Tags        map[string]any      `json:"tags"`
}

func (WebhookPayloadEmailFailed) IsWebhookPayloadData() {}

// WebhookPayloadEmailOpened is the data payload for email.opened.
type WebhookPayloadEmailOpened struct {
	BroadcastID string         `json:"broadcast_id"`
	CreatedAt   time.Time      `json:"created_at"`
	EmailID     string         `json:"email_id"`
	MessageID   string         `json:"message_id"`
	From        string         `json:"from"`
	To          []string       `json:"to"`
	Subject     string         `json:"subject"`
	TemplateID  string         `json:"template_id"`
	Tags        map[string]any `json:"tags"`
}

func (WebhookPayloadEmailOpened) IsWebhookPayloadData() {}

// WebhookPayloadEmailReceived is the data payload for email.received.
type WebhookPayloadEmailReceived struct {
	EmailID     string                   `json:"email_id"`
	CreatedAt   time.Time                `json:"created_at"`
	From        string                   `json:"from"`
	To          []string                 `json:"to"`
	BCC         []string                 `json:"bcc"`
	CC          []string                 `json:"cc"`
	ReceivedFor []string                 `json:"received_for"`
	MessageID   string                   `json:"message_id"`
	Subject     string                   `json:"subject"`
	Attachments []WebhookEmailAttachment `json:"attachments"`
}

func (WebhookPayloadEmailReceived) IsWebhookPayloadData() {}

// WebhookPayloadEmailScheduled is the data payload for email.scheduled.
type WebhookPayloadEmailScheduled struct {
	BroadcastID string         `json:"broadcast_id"`
	CreatedAt   time.Time      `json:"created_at"`
	EmailID     string         `json:"email_id"`
	From        string         `json:"from"`
	To          []string       `json:"to"`
	Subject     string         `json:"subject"`
	TemplateID  string         `json:"template_id"`
	Tags        map[string]any `json:"tags"`
}

func (WebhookPayloadEmailScheduled) IsWebhookPayloadData() {}

// WebhookPayloadEmailSent is the data payload for email.sent.
type WebhookPayloadEmailSent struct {
	BroadcastID string         `json:"broadcast_id"`
	CreatedAt   time.Time      `json:"created_at"`
	EmailID     string         `json:"email_id"`
	MessageID   string         `json:"message_id"`
	From        string         `json:"from"`
	To          []string       `json:"to"`
	Subject     string         `json:"subject"`
	TemplateID  string         `json:"template_id"`
	Tags        map[string]any `json:"tags"`
}

func (WebhookPayloadEmailSent) IsWebhookPayloadData() {}

// WebhookPayloadEmailSuppressed is the data payload for email.suppressed.
type WebhookPayloadEmailSuppressed struct {
	BroadcastID string                  `json:"broadcast_id"`
	CreatedAt   time.Time               `json:"created_at"`
	EmailID     string                  `json:"email_id"`
	From        string                  `json:"from"`
	To          []string                `json:"to"`
	Subject     string                  `json:"subject"`
	TemplateID  string                  `json:"template_id"`
	Suppressed  WebhookEmailSuppression `json:"suppressed"`
	Tags        map[string]any          `json:"tags"`
}

func (WebhookPayloadEmailSuppressed) IsWebhookPayloadData() {}

// webhookInboxLifecycleData contains the shared data for inbox lifecycle events.
type webhookInboxLifecycleData struct {
	Source  string       `json:"source"`
	InboxID string       `json:"inbox_id"`
	Inbox   WebhookInbox `json:"inbox"`
}

// WebhookPayloadInboxCreated is the data payload for inbox.created.
type WebhookPayloadInboxCreated struct{ webhookInboxLifecycleData }

func (WebhookPayloadInboxCreated) IsWebhookPayloadData() {}

// WebhookPayloadInboxDeleted is the data payload for inbox.deleted.
type WebhookPayloadInboxDeleted struct {
	Source  string `json:"source"`
	InboxID string `json:"inbox_id"`
}

func (WebhookPayloadInboxDeleted) IsWebhookPayloadData() {}

// WebhookPayloadInboxDraftCreated is the data payload for inbox.draft.created.
type webhookInboxDraftLifecycleData struct {
	Source   string             `json:"source"`
	InboxID  string             `json:"inbox_id"`
	ThreadID string             `json:"thread_id"`
	DraftID  string             `json:"draft_id"`
	Thread   WebhookInboxThread `json:"thread"`
	Draft    WebhookInboxDraft  `json:"draft"`
}

// WebhookPayloadInboxDraftCreated is the data payload for inbox.draft.created.
type WebhookPayloadInboxDraftCreated struct{ webhookInboxDraftLifecycleData }

func (WebhookPayloadInboxDraftCreated) IsWebhookPayloadData() {}

// WebhookPayloadInboxDraftDeleted is the data payload for inbox.draft.deleted.
type WebhookPayloadInboxDraftDeleted struct {
	Source   string             `json:"source"`
	InboxID  string             `json:"inbox_id"`
	ThreadID string             `json:"thread_id"`
	DraftID  string             `json:"draft_id"`
	Thread   WebhookInboxThread `json:"thread"`
}

func (WebhookPayloadInboxDraftDeleted) IsWebhookPayloadData() {}

// WebhookPayloadInboxDraftSent is the data payload for inbox.draft.sent.
type WebhookPayloadInboxDraftSent struct {
	Source   string             `json:"source"`
	InboxID  string             `json:"inbox_id"`
	ThreadID string             `json:"thread_id"`
	EmailID  string             `json:"email_id"`
	DraftID  string             `json:"draft_id"`
	Thread   WebhookInboxThread `json:"thread"`
	Draft    WebhookInboxDraft  `json:"draft"`
}

func (WebhookPayloadInboxDraftSent) IsWebhookPayloadData() {}

// WebhookPayloadInboxDraftUpdated is the data payload for inbox.draft.updated.
type WebhookPayloadInboxDraftUpdated struct{ webhookInboxDraftLifecycleData }

func (WebhookPayloadInboxDraftUpdated) IsWebhookPayloadData() {}

// webhookInboxEmailData contains the shared data for inbox email/thread events.
type webhookInboxEmailData struct {
	Source   string             `json:"source"`
	InboxID  string             `json:"inbox_id"`
	ThreadID string             `json:"thread_id"`
	EmailID  string             `json:"email_id"`
	Thread   WebhookInboxThread `json:"thread"`
	Email    WebhookInboxEmail  `json:"email"`
}

// WebhookPayloadInboxEmailReceived is the data payload for inbox.email.received.
type WebhookPayloadInboxEmailReceived struct{ webhookInboxEmailData }

func (WebhookPayloadInboxEmailReceived) IsWebhookPayloadData() {}

// WebhookPayloadInboxEmailSent is the data payload for inbox.email.sent.
type WebhookPayloadInboxEmailSent struct{ webhookInboxEmailData }

func (WebhookPayloadInboxEmailSent) IsWebhookPayloadData() {}

// WebhookPayloadInboxThreadAssigned is the data payload for inbox.thread.assigned.
type WebhookPayloadInboxThreadAssigned struct {
	AssigneeEmail         string             `json:"assignee_email"`
	AssignedByEmail       string             `json:"assigned_by_email"`
	PreviousAssigneeEmail *string            `json:"previous_assignee_email"`
	Source                string             `json:"source"`
	InboxID               string             `json:"inbox_id"`
	ThreadID              string             `json:"thread_id"`
	Thread                WebhookInboxThread `json:"thread"`
}

func (WebhookPayloadInboxThreadAssigned) IsWebhookPayloadData() {}

// WebhookPayloadInboxThreadCreated is the data payload for inbox.thread.created.
type WebhookPayloadInboxThreadCreated struct{ webhookInboxEmailData }

func (WebhookPayloadInboxThreadCreated) IsWebhookPayloadData() {}

// WebhookPayloadInboxThreadFolderUpdated is the data payload for inbox.thread.folder.updated.
type WebhookPayloadInboxThreadFolderUpdated struct {
	From     string             `json:"from"`
	To       string             `json:"to"`
	Source   string             `json:"source"`
	InboxID  string             `json:"inbox_id"`
	ThreadID string             `json:"thread_id"`
	Thread   WebhookInboxThread `json:"thread"`
}

func (WebhookPayloadInboxThreadFolderUpdated) IsWebhookPayloadData() {}

// WebhookPayloadInboxThreadLabelsUpdated is the data payload for inbox.thread.labels.updated.
type WebhookPayloadInboxThreadLabelsUpdated struct {
	Added    []map[string]any   `json:"added"`
	Removed  []string           `json:"removed"`
	Source   string             `json:"source"`
	InboxID  string             `json:"inbox_id"`
	ThreadID string             `json:"thread_id"`
	Thread   WebhookInboxThread `json:"thread"`
}

func (WebhookPayloadInboxThreadLabelsUpdated) IsWebhookPayloadData() {}

// WebhookPayloadInboxThreadUnassigned is the data payload for inbox.thread.unassigned.
type WebhookPayloadInboxThreadUnassigned struct {
	PreviousAssigneeEmail string             `json:"previous_assignee_email"`
	Source                string             `json:"source"`
	InboxID               string             `json:"inbox_id"`
	ThreadID              string             `json:"thread_id"`
	Thread                WebhookInboxThread `json:"thread"`
}

func (WebhookPayloadInboxThreadUnassigned) IsWebhookPayloadData() {}

// WebhookPayloadInboxUpdated is the data payload for inbox.updated.
type WebhookPayloadInboxUpdated struct{ webhookInboxLifecycleData }

func (WebhookPayloadInboxUpdated) IsWebhookPayloadData() {}

// WebhookPayloadSuppressionAdded is the data payload for suppression.added.
type WebhookPayloadSuppressionAdded struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Origin    string    `json:"origin"`
	SourceID  string    `json:"source_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (WebhookPayloadSuppressionAdded) IsWebhookPayloadData() {}

// WebhookPayloadSuppressionRemoved is the data payload for suppression.removed.
type WebhookPayloadSuppressionRemoved struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Origin    string    `json:"origin"`
	SourceID  *string   `json:"source_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (WebhookPayloadSuppressionRemoved) IsWebhookPayloadData() {}

// WebhookTopicData contains the shared data for topic lifecycle events.
type webhookTopicData struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	DefaultSubscription string    `json:"default_subscription"`
	Deleted             bool      `json:"deleted"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// WebhookPayloadTopicCreated is the data payload for topic.created.
type WebhookPayloadTopicCreated struct{ webhookTopicData }

func (WebhookPayloadTopicCreated) IsWebhookPayloadData() {}

// WebhookPayloadTopicDeleted is the data payload for topic.deleted.
type WebhookPayloadTopicDeleted struct{ webhookTopicData }

func (WebhookPayloadTopicDeleted) IsWebhookPayloadData() {}

// WebhookPayloadTopicUpdated is the data payload for topic.updated.
type WebhookPayloadTopicUpdated struct{ webhookTopicData }

func (WebhookPayloadTopicUpdated) IsWebhookPayloadData() {}

// WebhookContactData contains the fields shared by contact lifecycle events.
// FirstName and LastName are pointers because Resend can return null.
type webhookContactData struct {
	ID           string    `json:"id"`
	AudienceID   string    `json:"audience_id"`
	SegmentIDs   []string  `json:"segment_ids"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	FirstName    *string   `json:"first_name"`
	LastName     *string   `json:"last_name"`
	Unsubscribed bool      `json:"unsubscribed"`
}

// WebhookPayloadContactCreated is the data payload for contacts.created.
type WebhookPayloadContactCreated struct{ webhookContactData }

func (WebhookPayloadContactCreated) IsWebhookPayloadData() {}

// WebhookPayloadContactUpdated is the data payload for contacts.updated.
type WebhookPayloadContactUpdated struct{ webhookContactData }

func (WebhookPayloadContactUpdated) IsWebhookPayloadData() {}

// WebhookPayloadContactDeleted is the data payload for contacts.deleted.
type WebhookPayloadContactDeleted struct{ webhookContactData }

func (WebhookPayloadContactDeleted) IsWebhookPayloadData() {}

// WebhookContactTopicSubscription represents a contact's subscription to a topic.
type WebhookContactTopicSubscription struct {
	ID           string `json:"id"`
	Subscription string `json:"subscription"` // opt_in or opt_out
}

// WebhookPayloadContactTopicsUpdated is the data payload for contacts.topics.updated.
type WebhookPayloadContactTopicsUpdated struct {
	Email  string                            `json:"email"`
	Topics []WebhookContactTopicSubscription `json:"topics"`
}

func (WebhookPayloadContactTopicsUpdated) IsWebhookPayloadData() {}

// WebhookPayloadDataTypes maps Resend event names to their concrete Go data types.
var WebhookPayloadDataTypes = map[string]reflect.Type{
	"contact.created":             reflect.TypeOf(WebhookPayloadContactCreated{}),
	"contact.updated":             reflect.TypeOf(WebhookPayloadContactUpdated{}),
	"contact.deleted":             reflect.TypeOf(WebhookPayloadContactDeleted{}),
	"contact.topics.updated":      reflect.TypeOf(WebhookPayloadContactTopicsUpdated{}),
	"domain.created":              reflect.TypeOf(WebhookPayloadDomainCreated{}),
	"domain.deleted":              reflect.TypeOf(WebhookPayloadDomainDeleted{}),
	"domain.updated":              reflect.TypeOf(WebhookPayloadDomainUpdated{}),
	"email.bounced":               reflect.TypeOf(WebhookPayloadEmailBounced{}),
	"email.clicked":               reflect.TypeOf(WebhookPayloadEmailClicked{}),
	"email.complained":            reflect.TypeOf(WebhookPayloadEmailComplained{}),
	"email.delivered":             reflect.TypeOf(WebhookPayloadEmailDelivered{}),
	"email.delivery_delayed":      reflect.TypeOf(WebhookPayloadEmailDeliveryDelayed{}),
	"email.failed":                reflect.TypeOf(WebhookPayloadEmailFailed{}),
	"email.opened":                reflect.TypeOf(WebhookPayloadEmailOpened{}),
	"email.received":              reflect.TypeOf(WebhookPayloadEmailReceived{}),
	"email.scheduled":             reflect.TypeOf(WebhookPayloadEmailScheduled{}),
	"email.sent":                  reflect.TypeOf(WebhookPayloadEmailSent{}),
	"email.suppressed":            reflect.TypeOf(WebhookPayloadEmailSuppressed{}),
	"inbox.created":               reflect.TypeOf(WebhookPayloadInboxCreated{}),
	"inbox.deleted":               reflect.TypeOf(WebhookPayloadInboxDeleted{}),
	"inbox.draft.created":         reflect.TypeOf(WebhookPayloadInboxDraftCreated{}),
	"inbox.draft.deleted":         reflect.TypeOf(WebhookPayloadInboxDraftDeleted{}),
	"inbox.draft.sent":            reflect.TypeOf(WebhookPayloadInboxDraftSent{}),
	"inbox.draft.updated":         reflect.TypeOf(WebhookPayloadInboxDraftUpdated{}),
	"inbox.email.received":        reflect.TypeOf(WebhookPayloadInboxEmailReceived{}),
	"inbox.email.sent":            reflect.TypeOf(WebhookPayloadInboxEmailSent{}),
	"inbox.thread.assigned":       reflect.TypeOf(WebhookPayloadInboxThreadAssigned{}),
	"inbox.thread.created":        reflect.TypeOf(WebhookPayloadInboxThreadCreated{}),
	"inbox.thread.folder.updated": reflect.TypeOf(WebhookPayloadInboxThreadFolderUpdated{}),
	"inbox.thread.labels.updated": reflect.TypeOf(WebhookPayloadInboxThreadLabelsUpdated{}),
	"inbox.thread.unassigned":     reflect.TypeOf(WebhookPayloadInboxThreadUnassigned{}),
	"inbox.updated":               reflect.TypeOf(WebhookPayloadInboxUpdated{}),
	"suppression.added":           reflect.TypeOf(WebhookPayloadSuppressionAdded{}),
	"suppression.removed":         reflect.TypeOf(WebhookPayloadSuppressionRemoved{}),
	"topic.created":               reflect.TypeOf(WebhookPayloadTopicCreated{}),
	"topic.deleted":               reflect.TypeOf(WebhookPayloadTopicDeleted{}),
	"topic.updated":               reflect.TypeOf(WebhookPayloadTopicUpdated{}),
}
