package resend

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

type SendEmailOptions struct {
	IdempotencyKey string `json:"idempotency_key,omitempty"`
}

func (o SendEmailOptions) GetIdempotencyKey() string {
	return o.IdempotencyKey
}

// EmailTemplate represents a template configuration for sending emails.
type EmailTemplate struct {
	// Id is the template ID or alias to use for this email
	Id string `json:"id"`
	// Variables are the key-value pairs to populate the template placeholders
	Variables map[string]any `json:"variables,omitempty"`
}

// SendEmailRequest is the request object for the Send call.
//
// See also https://resend.com/docs/api-reference/emails/send-email
type SendEmailRequest struct {
	From        string            `json:"from"`
	To          []string          `json:"to"`
	Subject     string            `json:"subject"`
	Bcc         []string          `json:"bcc,omitempty"`
	Cc          []string          `json:"cc,omitempty"`
	ReplyTo     string            `json:"reply_to,omitempty"`
	Html        string            `json:"html,omitempty"`
	Text        string            `json:"text,omitempty"`
	Tags        []Tag             `json:"tags,omitempty"`
	Attachments []*Attachment     `json:"attachments,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	ScheduledAt string            `json:"scheduled_at,omitempty"`
	Template    *EmailTemplate    `json:"template,omitempty"`
	TopicId     string            `json:"topic_id,omitempty"`
}

// CancelScheduledEmailResponse is the response from the Cancel call.
type CancelScheduledEmailResponse struct {
	Id     string `json:"id"`
	Object string `json:"object"`
}

// ShareEmailRequest is the request object for the Share call.
type ShareEmailRequest struct {
	ExpiresIn string `json:"expires_in,omitempty"`
}

// ShareEmailResponse is the response from the Share call.
type ShareEmailResponse struct {
	Id     string `json:"id"`
	Object string `json:"object"`
	Url    string `json:"url"`
}

// SendEmailResponse is the response from the Send call.
type SendEmailResponse struct {
	Id string `json:"id"`
}

// UpdateEmailRequest is the request object for the Update call.
type UpdateEmailRequest struct {
	Id          string `json:"id"`
	ScheduledAt string `json:"scheduled_at"`
}

// UpdateEmailResponse is the type that represents the response from the Update call.
type UpdateEmailResponse struct {
	Id     string `json:"id"`
	Object string `json:"object"`
}

// Email provides the structure for the response from the Get call.
type Email struct {
	Id        string   `json:"id"`
	Object    string   `json:"object"`
	MessageId string   `json:"message_id"`
	To        []string `json:"to"`
	From      string   `json:"from"`
	CreatedAt string   `json:"created_at"`
	Subject   string   `json:"subject"`
	Html      string   `json:"html"`
	Text      string   `json:"text"`
	Bcc       []string `json:"bcc"`
	Cc        []string `json:"cc"`
	ReplyTo   []string `json:"reply_to"`
	LastEvent string   `json:"last_event"`
}

// ListEmailsResponse is the response from the List call.
type ListEmailsResponse struct {
	Object  string  `json:"object"`
	HasMore bool    `json:"has_more"`
	Data    []Email `json:"data"`
}

// Tags are used to define custom metadata for emails
type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// EmailAttachment represents an attachment for both sent and received emails.
// When returned from GET /emails/:id/attachments/:attachmentId, it includes the Object field.
// When returned in a list (Data array), the Object field is omitted.
type EmailAttachment struct {
	Object             string `json:"object,omitempty"`
	Id                 string `json:"id"`
	Filename           string `json:"filename"`
	ContentType        string `json:"content_type"`
	ContentDisposition string `json:"content_disposition"`
	ContentId          string `json:"content_id"`
	DownloadUrl        string `json:"download_url"`
	ExpiresAt          string `json:"expires_at"`
}

// ListEmailAttachmentsResponse is the response from the ListAttachments call.
type ListEmailAttachmentsResponse struct {
	Object  string            `json:"object"`
	HasMore bool              `json:"has_more"`
	Data    []EmailAttachment `json:"data"`
}

// MetricName identifies a metric that can be requested from the Metrics call.
type MetricName = string

const (
	MetricReceived            MetricName = "received"
	MetricDelivered           MetricName = "delivered"
	MetricComplained          MetricName = "complained"
	MetricSuppressed          MetricName = "suppressed"
	MetricBounced             MetricName = "bounced"
	MetricBouncedTransient    MetricName = "bounced_transient"
	MetricBouncedPermanent    MetricName = "bounced_permanent"
	MetricBouncedUndetermined MetricName = "bounced_undetermined"
	MetricOpened              MetricName = "opened"
	MetricClicked             MetricName = "clicked"
	MetricUnsubscribed        MetricName = "unsubscribed"
	MetricDeliveryDelayed     MetricName = "delivery_delayed"
	MetricFailed              MetricName = "failed"
	MetricSent                MetricName = "sent"
	MetricUniqueOpened        MetricName = "unique_opened"
	MetricUniqueClicked       MetricName = "unique_clicked"
	MetricDeliveryRate        MetricName = "delivery_rate"
	MetricOpenRate            MetricName = "open_rate"
	MetricClickRate           MetricName = "click_rate"
	MetricBounceRate          MetricName = "bounce_rate"
	MetricComplaintRate       MetricName = "complaint_rate"
	MetricUnsubscribeRate     MetricName = "unsubscribe_rate"
)

// MetricsDimension identifies a dimension to group metrics by in the Metrics call.
type MetricsDimension = string

const (
	MetricsDimensionPeriod    MetricsDimension = "period"
	MetricsDimensionDomain    MetricsDimension = "domain"
	MetricsDimensionEmail     MetricsDimension = "email"
	MetricsDimensionBroadcast MetricsDimension = "broadcast"
)

// MetricsGranularity is the bucket size used for the `period` dimension.
type MetricsGranularity = string

const (
	MetricsGranularityHourly  MetricsGranularity = "hourly"
	MetricsGranularityDaily   MetricsGranularity = "daily"
	MetricsGranularityWeekly  MetricsGranularity = "weekly"
	MetricsGranularityMonthly MetricsGranularity = "monthly"
)

// MetricsOptions contains the optional query parameters for the Metrics call.
//
// See also https://resend.com/docs/api-reference/emails/get-metrics
type MetricsOptions struct {
	// StartDate filters metrics to on/after this ISO 8601 date or datetime.
	// Defaults server-side to 6 days before EndDate.
	StartDate *string

	// EndDate filters metrics to on/before this ISO 8601 date or datetime.
	// Defaults server-side to now.
	EndDate *string

	// Timezone is an IANA timezone identifier, e.g. "America/New_York".
	// Defaults server-side to "UTC".
	Timezone *string

	// Granularity is the bucket size used when the `period` dimension is
	// requested. Defaults server-side to MetricsGranularityDaily.
	Granularity *MetricsGranularity

	// Metrics selects which metrics to compute. Defaults server-side to all
	// metrics.
	Metrics []MetricName

	// Dimensions selects which dimensions to group results by. Defaults
	// server-side to none, in which case the response only contains Totals.
	Dimensions []MetricsDimension

	// DomainId restricts results to these sending domain IDs (max 100).
	DomainId []string

	// EmailId restricts results to these email IDs (max 100). Cannot be
	// combined with MetricsDimensionBroadcast or BroadcastId.
	EmailId []string

	// BroadcastId restricts results to these broadcast IDs (max 100). Cannot
	// be combined with MetricsDimensionEmail or EmailId.
	BroadcastId []string
}

// metricsHasDimension reports whether dimensions includes the given dimension.
func metricsHasDimension(dimensions []MetricsDimension, dimension MetricsDimension) bool {
	for _, d := range dimensions {
		if d == dimension {
			return true
		}
	}
	return false
}

// metricsInvolvesEmailAndBroadcast reports whether options combines the
// `email` dimension/EmailId with the `broadcast` dimension/BroadcastId,
// which the Metrics endpoint rejects.
func metricsInvolvesEmailAndBroadcast(options *MetricsOptions) bool {
	involvesEmail := len(options.EmailId) > 0 || metricsHasDimension(options.Dimensions, MetricsDimensionEmail)
	involvesBroadcast := len(options.BroadcastId) > 0 || metricsHasDimension(options.Dimensions, MetricsDimensionBroadcast)
	return involvesEmail && involvesBroadcast
}

// buildMetricsQuery constructs query parameters for the Metrics call
func buildMetricsQuery(options *MetricsOptions) string {
	if options == nil {
		return ""
	}

	query := make(url.Values)
	if options.StartDate != nil {
		query.Set("start_date", *options.StartDate)
	}
	if options.EndDate != nil {
		query.Set("end_date", *options.EndDate)
	}
	if options.Timezone != nil {
		query.Set("timezone", *options.Timezone)
	}
	if options.Granularity != nil {
		query.Set("granularity", *options.Granularity)
	}
	if len(options.Metrics) > 0 {
		query.Set("metrics", strings.Join(options.Metrics, ","))
	}
	if len(options.Dimensions) > 0 {
		query.Set("dimensions", strings.Join(options.Dimensions, ","))
	}
	if len(options.DomainId) > 0 {
		query.Set("domain_id", strings.Join(options.DomainId, ","))
	}
	if len(options.EmailId) > 0 {
		query.Set("email_id", strings.Join(options.EmailId, ","))
	}
	if len(options.BroadcastId) > 0 {
		query.Set("broadcast_id", strings.Join(options.BroadcastId, ","))
	}

	if len(query) > 0 {
		return "?" + query.Encode()
	}
	return ""
}

// MetricsDataPoint is a single row in MetricsResponse.Data. Which fields are
// populated depends on the requested dimensions and metrics: dimension key
// fields are set only when that dimension was requested, and metric fields
// are set only when that metric was requested.
type MetricsDataPoint struct {
	// Period is set when MetricsDimensionPeriod is requested.
	Period *string `json:"period,omitempty"`

	// DomainId and DomainName are set when MetricsDimensionDomain is requested.
	DomainId   *string `json:"domain_id,omitempty"`
	DomainName *string `json:"domain_name,omitempty"`

	// EmailId is set when MetricsDimensionEmail is requested.
	EmailId *string `json:"email_id,omitempty"`

	// BroadcastId and BroadcastName are set when MetricsDimensionBroadcast is requested.
	BroadcastId   *string `json:"broadcast_id,omitempty"`
	BroadcastName *string `json:"broadcast_name,omitempty"`

	Received            *int64   `json:"received,omitempty"`
	Delivered           *int64   `json:"delivered,omitempty"`
	Complained          *int64   `json:"complained,omitempty"`
	Suppressed          *int64   `json:"suppressed,omitempty"`
	Bounced             *int64   `json:"bounced,omitempty"`
	BouncedTransient    *int64   `json:"bounced_transient,omitempty"`
	BouncedPermanent    *int64   `json:"bounced_permanent,omitempty"`
	BouncedUndetermined *int64   `json:"bounced_undetermined,omitempty"`
	Opened              *int64   `json:"opened,omitempty"`
	Clicked             *int64   `json:"clicked,omitempty"`
	Unsubscribed        *int64   `json:"unsubscribed,omitempty"`
	DeliveryDelayed     *int64   `json:"delivery_delayed,omitempty"`
	Failed              *int64   `json:"failed,omitempty"`
	Sent                *int64   `json:"sent,omitempty"`
	UniqueOpened        *int64   `json:"unique_opened,omitempty"`
	UniqueClicked       *int64   `json:"unique_clicked,omitempty"`
	DeliveryRate        *float64 `json:"delivery_rate,omitempty"`
	OpenRate            *float64 `json:"open_rate,omitempty"`
	ClickRate           *float64 `json:"click_rate,omitempty"`
	BounceRate          *float64 `json:"bounce_rate,omitempty"`
	ComplaintRate       *float64 `json:"complaint_rate,omitempty"`
	UnsubscribeRate     *float64 `json:"unsubscribe_rate,omitempty"`
}

// MetricsResponse is the response from the Metrics call.
//
// See also https://resend.com/docs/api-reference/emails/get-metrics
type MetricsResponse struct {
	Object      string             `json:"object"`
	StartDate   string             `json:"start_date"`
	EndDate     string             `json:"end_date"`
	Metrics     []string           `json:"metrics"`
	Dimensions  []string           `json:"dimensions"`
	Granularity string             `json:"granularity"`
	Totals      map[string]float64 `json:"totals"`
	// Data is omitted from the response when no dimensions were requested.
	Data []MetricsDataPoint `json:"data,omitempty"`
}

// Attachment is the public struct used for adding attachments to emails
type Attachment struct {
	// Content is the binary content of the attachment to use when a Path
	// is not available.
	Content []byte

	// Filename that will appear in the email.
	// Make sure you pick the correct extension otherwise preview
	// may not work as expected
	Filename string

	// Path where the attachment file is hosted instead of providing the
	// content directly.
	Path string

	// Content type for the attachment, if not set will be derived from
	// the filename property
	ContentType string

	// Optional content ID for the attachment, to be used as a reference in the HTML content.
	// If set, this attachment will be sent as an inline attachment and you can reference it
	// in the HTML content using the `cid:` prefix.
	ContentId string

	// Deprecated: Use ContentId instead. Kept for backwards compatibility.
	// Optional content ID for the attachment, to be used as a reference in the HTML content.
	// If set, this attachment will be sent as an inline attachment and you can reference it
	// in the HTML content using the `cid:` prefix.
	InlineContentId string
}

// MarshalJSON overrides the regular JSON Marshaller to ensure that the
// attachment content is provided in the way Resend expects.
func (a *Attachment) MarshalJSON() ([]byte, error) {
	na := struct {
		Content         []int  `json:"content,omitempty"`
		Filename        string `json:"filename,omitempty"`
		Path            string `json:"path,omitempty"`
		ContentType     string `json:"content_type,omitempty"`
		ContentId       string `json:"content_id,omitempty"`
		InlineContentId string `json:"inline_content_id,omitempty"`
	}{
		Filename:        a.Filename,
		Path:            a.Path,
		Content:         BytesToIntArray(a.Content),
		ContentType:     a.ContentType,
		ContentId:       a.ContentId,
		InlineContentId: a.InlineContentId,
	}
	return json.Marshal(na)
}

type EmailsSvc interface {
	CancelWithContext(ctx context.Context, emailId string) (*CancelScheduledEmailResponse, error)
	Cancel(emailId string) (*CancelScheduledEmailResponse, error)
	ShareWithContext(ctx context.Context, emailId string, params *ShareEmailRequest) (*ShareEmailResponse, error)
	Share(emailId string, params *ShareEmailRequest) (*ShareEmailResponse, error)
	UpdateWithContext(ctx context.Context, params *UpdateEmailRequest) (*UpdateEmailResponse, error)
	Update(params *UpdateEmailRequest) (*UpdateEmailResponse, error)
	SendWithOptions(ctx context.Context, params *SendEmailRequest, options *SendEmailOptions) (*SendEmailResponse, error)
	SendWithContext(ctx context.Context, params *SendEmailRequest) (*SendEmailResponse, error)
	Send(params *SendEmailRequest) (*SendEmailResponse, error)
	GetWithContext(ctx context.Context, emailId string) (*Email, error)
	Get(emailId string) (*Email, error)

	// Both List and ListWithOptions do the same thing, but since these List methods
	// were introduced after some time, we kept both for overall consistency with
	// the rest of the packages.
	ListWithOptions(ctx context.Context, options *ListOptions) (ListEmailsResponse, error)
	ListWithContext(ctx context.Context) (ListEmailsResponse, error)
	List() (ListEmailsResponse, error)

	// Attachment methods for sent emails
	GetAttachmentWithContext(ctx context.Context, emailId string, attachmentId string) (*EmailAttachment, error)
	GetAttachment(emailId string, attachmentId string) (*EmailAttachment, error)
	ListAttachmentsWithOptions(ctx context.Context, emailId string, options *ListOptions) (ListEmailAttachmentsResponse, error)
	ListAttachmentsWithContext(ctx context.Context, emailId string) (ListEmailAttachmentsResponse, error)
	ListAttachments(emailId string) (ListEmailAttachmentsResponse, error)

	// Metrics methods
	MetricsWithOptions(ctx context.Context, options *MetricsOptions) (*MetricsResponse, error)
	MetricsWithContext(ctx context.Context) (*MetricsResponse, error)
	Metrics() (*MetricsResponse, error)
}

type EmailsSvcImpl struct {
	client    *Client
	Receiving ReceivingSvc
}

// Cancel cancels an email by ID
// https://resend.com/docs/api-reference/emails/cancel-email
func (s *EmailsSvcImpl) Cancel(emailId string) (*CancelScheduledEmailResponse, error) {
	return s.CancelWithContext(context.Background(), emailId)
}

// CancelWithContext cancels an email by ID
// https://resend.com/docs/api-reference/emails/cancel-email
func (s *EmailsSvcImpl) CancelWithContext(ctx context.Context, emailId string) (*CancelScheduledEmailResponse, error) {
	path := "emails/" + emailId + "/cancel"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, ErrFailedToCreateEmailsSendRequest
	}

	// Build response recipient obj
	resp := new(CancelScheduledEmailResponse)

	// Send Request
	_, err = s.client.Perform(req, resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Share creates a shareable link for a sent or received email by ID
// https://resend.com/docs/api-reference/emails/share-email
func (s *EmailsSvcImpl) Share(emailId string, params *ShareEmailRequest) (*ShareEmailResponse, error) {
	return s.ShareWithContext(context.Background(), emailId, params)
}

// ShareWithContext creates a shareable link for a sent or received email by ID
// https://resend.com/docs/api-reference/emails/share-email
func (s *EmailsSvcImpl) ShareWithContext(ctx context.Context, emailId string, params *ShareEmailRequest) (*ShareEmailResponse, error) {
	path := "emails/" + emailId + "/share"

	var body any
	if params != nil {
		body = params
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, ErrFailedToCreateEmailsShareRequest
	}

	resp := new(ShareEmailResponse)
	_, err = s.client.Perform(req, resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Update updates an email with the given params
// https://resend.com/docs/api-reference/emails/update-email
func (s *EmailsSvcImpl) Update(params *UpdateEmailRequest) (*UpdateEmailResponse, error) {
	return s.UpdateWithContext(context.Background(), params)
}

// UpdateWithContext updates an email with the given params
// https://resend.com/docs/api-reference/emails/update-email
func (s *EmailsSvcImpl) UpdateWithContext(ctx context.Context, params *UpdateEmailRequest) (*UpdateEmailResponse, error) {
	path := "emails/" + params.Id

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, params)
	if err != nil {
		return nil, ErrFailedToCreateUpdateEmailRequest
	}

	// Build response recipient obj
	updateEmailResponse := new(UpdateEmailResponse)

	// Send Request
	_, err = s.client.Perform(req, updateEmailResponse)

	if err != nil {
		return nil, err
	}

	return updateEmailResponse, nil
}

// SendWithOptions sends an email with the given params
// and additional options
// https://resend.com/docs/api-reference/emails/send-email
func (s *EmailsSvcImpl) SendWithOptions(ctx context.Context, params *SendEmailRequest, options *SendEmailOptions) (*SendEmailResponse, error) {
	path := "emails"

	// Prepare request
	req, err := s.client.NewRequestWithOptions(ctx, http.MethodPost, path, params, options)
	if err != nil {
		return nil, ErrFailedToCreateEmailsSendRequest
	}

	// Build response recipient obj
	emailResponse := new(SendEmailResponse)

	// Send Request
	_, err = s.client.Perform(req, emailResponse)

	if err != nil {
		return nil, err
	}

	return emailResponse, nil
}

// SendWithContext sends an email with the given params
// https://resend.com/docs/api-reference/emails/send-email
func (s *EmailsSvcImpl) SendWithContext(ctx context.Context, params *SendEmailRequest) (*SendEmailResponse, error) {
	path := "emails"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return nil, ErrFailedToCreateEmailsSendRequest
	}

	// Build response recipient obj
	emailResponse := new(SendEmailResponse)

	// Send Request
	_, err = s.client.Perform(req, emailResponse)

	if err != nil {
		return nil, err
	}

	return emailResponse, nil
}

// Send sends an email with the given params
// https://resend.com/docs/api-reference/emails/send-email
func (s *EmailsSvcImpl) Send(params *SendEmailRequest) (*SendEmailResponse, error) {
	return s.SendWithContext(context.Background(), params)
}

// GetWithContext retrieves an email with the given emailId
// https://resend.com/docs/api-reference/emails/retrieve-email
func (s *EmailsSvcImpl) GetWithContext(ctx context.Context, emailId string) (*Email, error) {
	path := "emails/" + emailId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, ErrFailedToCreateEmailsGetRequest
	}

	// Build response recipient obj
	emailResponse := new(Email)

	// Send Request
	_, err = s.client.Perform(req, emailResponse)

	if err != nil {
		return nil, err
	}

	return emailResponse, nil
}

// Get retrieves an email with the given emailId
// https://resend.com/docs/api-reference/emails/retrieve-email
func (s *EmailsSvcImpl) Get(emailId string) (*Email, error) {
	return s.GetWithContext(context.Background(), emailId)
}

// ListWithOptions retrieves a list of emails with pagination options
// https://resend.com/docs/api-reference/emails/list-emails
func (s *EmailsSvcImpl) ListWithOptions(ctx context.Context, options *ListOptions) (ListEmailsResponse, error) {
	path := "emails" + buildPaginationQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListEmailsResponse{}, ErrFailedToCreateEmailsListRequest
	}

	// Build response recipient obj
	listEmailsResponse := new(ListEmailsResponse)

	// Send Request
	_, err = s.client.Perform(req, listEmailsResponse)

	if err != nil {
		return ListEmailsResponse{}, err
	}

	return *listEmailsResponse, nil
}

// ListWithContext retrieves a list of emails
// https://resend.com/docs/api-reference/emails/list-emails
func (s *EmailsSvcImpl) ListWithContext(ctx context.Context) (ListEmailsResponse, error) {
	return s.ListWithOptions(ctx, nil)
}

// List retrieves a list of emails
// https://resend.com/docs/api-reference/emails/list-emails
func (s *EmailsSvcImpl) List() (ListEmailsResponse, error) {
	return s.ListWithContext(context.Background())
}

// GetAttachmentWithContext retrieves a single attachment from a sent email with the given emailId and attachmentId
// https://resend.com/docs/api-reference/attachments/retrieve-sent-email-attachment
func (s *EmailsSvcImpl) GetAttachmentWithContext(ctx context.Context, emailId string, attachmentId string) (*EmailAttachment, error) {
	path := "emails/" + emailId + "/attachments/" + attachmentId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, ErrFailedToCreateEmailsGetAttachmentRequest
	}

	attachment := new(EmailAttachment)

	// Send Request
	_, err = s.client.Perform(req, attachment)

	if err != nil {
		return nil, err
	}

	return attachment, nil
}

// GetAttachment retrieves a single attachment from a sent email with the given emailId and attachmentId
// https://resend.com/docs/api-reference/attachments/retrieve-sent-email-attachment
func (s *EmailsSvcImpl) GetAttachment(emailId string, attachmentId string) (*EmailAttachment, error) {
	return s.GetAttachmentWithContext(context.Background(), emailId, attachmentId)
}

// ListAttachmentsWithOptions retrieves a list of attachments for a sent email with pagination options
// https://resend.com/docs/api-reference/attachments/list-sent-email-attachments
func (s *EmailsSvcImpl) ListAttachmentsWithOptions(ctx context.Context, emailId string, options *ListOptions) (ListEmailAttachmentsResponse, error) {
	path := "emails/" + emailId + "/attachments" + buildPaginationQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListEmailAttachmentsResponse{}, ErrFailedToCreateEmailsListAttachmentsRequest
	}

	// Build response recipient obj
	listAttachmentsResponse := new(ListEmailAttachmentsResponse)

	// Send Request
	_, err = s.client.Perform(req, listAttachmentsResponse)

	if err != nil {
		return ListEmailAttachmentsResponse{}, err
	}

	return *listAttachmentsResponse, nil
}

// ListAttachmentsWithContext retrieves a list of attachments for a sent email
// https://resend.com/docs/api-reference/attachments/list-sent-email-attachments
func (s *EmailsSvcImpl) ListAttachmentsWithContext(ctx context.Context, emailId string) (ListEmailAttachmentsResponse, error) {
	return s.ListAttachmentsWithOptions(ctx, emailId, nil)
}

// ListAttachments retrieves a list of attachments for a sent email
// https://resend.com/docs/api-reference/attachments/list-sent-email-attachments
func (s *EmailsSvcImpl) ListAttachments(emailId string) (ListEmailAttachmentsResponse, error) {
	return s.ListAttachmentsWithContext(context.Background(), emailId)
}

// MetricsWithOptions retrieves email metrics with the given options.
// https://resend.com/docs/api-reference/emails/get-metrics
func (s *EmailsSvcImpl) MetricsWithOptions(ctx context.Context, options *MetricsOptions) (*MetricsResponse, error) {
	if options != nil && metricsInvolvesEmailAndBroadcast(options) {
		return nil, errors.New("[ERROR]: The `email` dimension/EmailId cannot be combined with the `broadcast` dimension/BroadcastId.")
	}

	path := "emails/metrics" + buildMetricsQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, ErrFailedToCreateEmailsMetricsRequest
	}

	// Build response recipient obj
	metricsResponse := new(MetricsResponse)

	// Send Request
	_, err = s.client.Perform(req, metricsResponse)
	if err != nil {
		return nil, err
	}

	return metricsResponse, nil
}

// MetricsWithContext retrieves email metrics.
// https://resend.com/docs/api-reference/emails/get-metrics
func (s *EmailsSvcImpl) MetricsWithContext(ctx context.Context) (*MetricsResponse, error) {
	return s.MetricsWithOptions(ctx, nil)
}

// Metrics retrieves email metrics.
// https://resend.com/docs/api-reference/emails/get-metrics
func (s *EmailsSvcImpl) Metrics() (*MetricsResponse, error) {
	return s.MetricsWithContext(context.Background())
}
