package resend

// Webhook event payload types for parsing verified webhook requests.
// See https://resend.com/docs/webhooks/event-types

// BaseEmailEventData is the common data payload for outbound email webhook events.
type BaseEmailEventData struct {
	BroadcastId string            `json:"broadcast_id,omitempty"`
	CreatedAt   string            `json:"created_at"`
	EmailId     string            `json:"email_id"`
	MessageId   string            `json:"message_id"`
	From        string            `json:"from"`
	To          []string          `json:"to"`
	Subject     string            `json:"subject"`
	TemplateId  string            `json:"template_id,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
}

// EmailBounce contains bounce details for email.bounced events.
type EmailBounce struct {
	Message string `json:"message"`
	SubType string `json:"subType"`
	Type    string `json:"type"`
}

// EmailClick contains click tracking details for email.clicked events.
type EmailClick struct {
	IpAddress string `json:"ipAddress"`
	Link      string `json:"link"`
	Timestamp string `json:"timestamp"`
	UserAgent string `json:"userAgent"`
}

// EmailFailed contains failure details for email.failed events.
type EmailFailed struct {
	Reason string `json:"reason"`
}

// EmailSuppressed contains suppression details for email.suppressed events.
type EmailSuppressed struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// ReceivedEmailWebhookAttachment represents an attachment in email.received events.
type ReceivedEmailWebhookAttachment struct {
	Id                 string `json:"id"`
	Filename           string `json:"filename"`
	ContentType        string `json:"content_type"`
	ContentDisposition string `json:"content_disposition"`
	ContentId          string `json:"content_id"`
}

// ReceivedEmailEventData is the data payload for email.received events.
type ReceivedEmailEventData struct {
	EmailId     string                           `json:"email_id"`
	CreatedAt   string                           `json:"created_at"`
	From        string                           `json:"from"`
	To          []string                         `json:"to"`
	Bcc         []string                         `json:"bcc"`
	Cc          []string                         `json:"cc"`
	ReceivedFor []string                         `json:"received_for"`
	MessageId   string                           `json:"message_id"`
	Subject     string                           `json:"subject"`
	Attachments []ReceivedEmailWebhookAttachment `json:"attachments"`
}

// ContactEventData is the data payload for contact webhook events.
type ContactEventData struct {
	Id            string   `json:"id"`
	AudienceId    string   `json:"audience_id"`
	SegmentIds    []string `json:"segment_ids"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
	Email         string   `json:"email"`
	FirstName     string   `json:"first_name,omitempty"`
	LastName      string   `json:"last_name,omitempty"`
	Unsubscribed  bool     `json:"unsubscribed"`
}

// DomainRecord represents a DNS record in domain webhook events.
type DomainRecord struct {
	Record   string `json:"record"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Ttl      string `json:"ttl"`
	Status   string `json:"status"`
	Value    string `json:"value"`
	Priority *int   `json:"priority,omitempty"`
}

// DomainEventData is the data payload for domain webhook events.
type DomainEventData struct {
	Id        string         `json:"id"`
	Name      string         `json:"name"`
	Status    string         `json:"status"`
	CreatedAt string         `json:"created_at"`
	Region    string         `json:"region"`
	Records   []DomainRecord `json:"records"`
}

// EmailSentEvent is the payload for email.sent webhook events.
type EmailSentEvent struct {
	Type      string             `json:"type"`
	CreatedAt string             `json:"created_at"`
	Data      BaseEmailEventData `json:"data"`
}

// EmailScheduledEvent is the payload for email.scheduled webhook events.
type EmailScheduledEvent struct {
	Type      string             `json:"type"`
	CreatedAt string             `json:"created_at"`
	Data      BaseEmailEventData `json:"data"`
}

// EmailDeliveredEvent is the payload for email.delivered webhook events.
type EmailDeliveredEvent struct {
	Type      string             `json:"type"`
	CreatedAt string             `json:"created_at"`
	Data      BaseEmailEventData `json:"data"`
}

// EmailDeliveryDelayedEvent is the payload for email.delivery_delayed webhook events.
type EmailDeliveryDelayedEvent struct {
	Type      string             `json:"type"`
	CreatedAt string             `json:"created_at"`
	Data      BaseEmailEventData `json:"data"`
}

// EmailComplainedEvent is the payload for email.complained webhook events.
type EmailComplainedEvent struct {
	Type      string             `json:"type"`
	CreatedAt string             `json:"created_at"`
	Data      BaseEmailEventData `json:"data"`
}

// EmailBouncedEventData extends BaseEmailEventData with bounce details.
type EmailBouncedEventData struct {
	BaseEmailEventData
	Bounce EmailBounce `json:"bounce"`
}

// EmailBouncedEvent is the payload for email.bounced webhook events.
type EmailBouncedEvent struct {
	Type      string                `json:"type"`
	CreatedAt string                `json:"created_at"`
	Data      EmailBouncedEventData `json:"data"`
}

// EmailOpenedEvent is the payload for email.opened webhook events.
type EmailOpenedEvent struct {
	Type      string             `json:"type"`
	CreatedAt string             `json:"created_at"`
	Data      BaseEmailEventData `json:"data"`
}

// EmailClickedEventData extends BaseEmailEventData with click details.
type EmailClickedEventData struct {
	BaseEmailEventData
	Click EmailClick `json:"click"`
}

// EmailClickedEvent is the payload for email.clicked webhook events.
type EmailClickedEvent struct {
	Type      string                `json:"type"`
	CreatedAt string                `json:"created_at"`
	Data      EmailClickedEventData `json:"data"`
}

// EmailReceivedEvent is the payload for email.received webhook events.
type EmailReceivedEvent struct {
	Type      string                 `json:"type"`
	CreatedAt string                 `json:"created_at"`
	Data      ReceivedEmailEventData `json:"data"`
}

// EmailFailedEventData extends BaseEmailEventData with failure details.
type EmailFailedEventData struct {
	BaseEmailEventData
	Failed EmailFailed `json:"failed"`
}

// EmailFailedEvent is the payload for email.failed webhook events.
type EmailFailedEvent struct {
	Type      string               `json:"type"`
	CreatedAt string               `json:"created_at"`
	Data      EmailFailedEventData `json:"data"`
}

// EmailSuppressedEventData extends BaseEmailEventData with suppression details.
type EmailSuppressedEventData struct {
	BaseEmailEventData
	Suppressed EmailSuppressed `json:"suppressed"`
}

// EmailSuppressedEvent is the payload for email.suppressed webhook events.
type EmailSuppressedEvent struct {
	Type      string                   `json:"type"`
	CreatedAt string                   `json:"created_at"`
	Data      EmailSuppressedEventData `json:"data"`
}

// ContactCreatedEvent is the payload for contact.created webhook events.
type ContactCreatedEvent struct {
	Type      string           `json:"type"`
	CreatedAt string           `json:"created_at"`
	Data      ContactEventData `json:"data"`
}

// ContactUpdatedEvent is the payload for contact.updated webhook events.
type ContactUpdatedEvent struct {
	Type      string           `json:"type"`
	CreatedAt string           `json:"created_at"`
	Data      ContactEventData `json:"data"`
}

// ContactDeletedEvent is the payload for contact.deleted webhook events.
type ContactDeletedEvent struct {
	Type      string           `json:"type"`
	CreatedAt string           `json:"created_at"`
	Data      ContactEventData `json:"data"`
}

// DomainCreatedEvent is the payload for domain.created webhook events.
type DomainCreatedEvent struct {
	Type      string          `json:"type"`
	CreatedAt string          `json:"created_at"`
	Data      DomainEventData `json:"data"`
}

// DomainUpdatedEvent is the payload for domain.updated webhook events.
type DomainUpdatedEvent struct {
	Type      string          `json:"type"`
	CreatedAt string          `json:"created_at"`
	Data      DomainEventData `json:"data"`
}

// DomainDeletedEvent is the payload for domain.deleted webhook events.
type DomainDeletedEvent struct {
	Type      string          `json:"type"`
	CreatedAt string          `json:"created_at"`
	Data      DomainEventData `json:"data"`
}
