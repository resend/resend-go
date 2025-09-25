package resend

import (
	"context"
	"encoding/json"
	"net/http"
)

type SendEmailOptions struct {
	IdempotencyKey string `json:"idempotency_key,omitempty"`
}

func (o SendEmailOptions) GetIdempotencyKey() string {
	return o.IdempotencyKey
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
}

// CancelScheduledEmailResponse is the response from the Cancel call.
type CancelScheduledEmailResponse struct {
	Id     string `json:"id"`
	Object string `json:"object"`
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
	CancelWithContext(ctx context.Context, emailID string) (*CancelScheduledEmailResponse, error)
	Cancel(emailID string) (*CancelScheduledEmailResponse, error)
	UpdateWithContext(ctx context.Context, params *UpdateEmailRequest) (*UpdateEmailResponse, error)
	Update(params *UpdateEmailRequest) (*UpdateEmailResponse, error)
	SendWithOptions(ctx context.Context, params *SendEmailRequest, options *SendEmailOptions) (*SendEmailResponse, error)
	SendWithContext(ctx context.Context, params *SendEmailRequest) (*SendEmailResponse, error)
	Send(params *SendEmailRequest) (*SendEmailResponse, error)
	GetWithContext(ctx context.Context, emailID string) (*Email, error)
	Get(emailID string) (*Email, error)

	// Both List and ListWithOptions do the same thing, but since these List methods
	// were introduced after some time, we kept both for overall consistency with
	// the rest of the packages.
	ListWithOptions(ctx context.Context, options *ListOptions) (ListEmailsResponse, error)
	ListWithContext(ctx context.Context) (ListEmailsResponse, error)
	List() (ListEmailsResponse, error)
}

type EmailsSvcImpl struct {
	client *Client
}

// Cancel cancels an email by ID
// https://resend.com/docs/api-reference/emails/cancel-email
func (s *EmailsSvcImpl) Cancel(emailID string) (*CancelScheduledEmailResponse, error) {
	return s.CancelWithContext(context.Background(), emailID)
}

// CancelWithContext cancels an email by ID
// https://resend.com/docs/api-reference/emails/cancel-email
func (s *EmailsSvcImpl) CancelWithContext(ctx context.Context, emailID string) (*CancelScheduledEmailResponse, error) {
	path := "emails/" + emailID + "/cancel"

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

// GetWithContext retrieves an email with the given emailID
// https://resend.com/docs/api-reference/emails/retrieve-email
func (s *EmailsSvcImpl) GetWithContext(ctx context.Context, emailID string) (*Email, error) {
	path := "emails/" + emailID

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

// Get retrieves an email with the given emailID
// https://resend.com/docs/api-reference/emails/retrieve-email
func (s *EmailsSvcImpl) Get(emailID string) (*Email, error) {
	return s.GetWithContext(context.Background(), emailID)
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
