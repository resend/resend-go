package resend

import (
	"context"
	"net/http"
)

// ReceivedEmail provides the structure for the response from the Receiving.Get call.
type ReceivedEmail struct {
	Id          string               `json:"id"`
	Object      string               `json:"object"` // Always "inbound" for received emails
	To          []string             `json:"to"`
	From        string               `json:"from"`
	CreatedAt   string               `json:"created_at"`
	Subject     string               `json:"subject"`
	Html        string               `json:"html"`
	Text        string               `json:"text"`
	Bcc         []string             `json:"bcc"`
	Cc          []string             `json:"cc"`
	ReplyTo     []string             `json:"reply_to"`
	Headers     map[string]string    `json:"headers"`
	Attachments []ReceivedAttachment `json:"attachments"`
}

// ListReceivedEmail provides the structure for items in the Receiving.List call.
// It omits html, text, and headers fields compared to ReceivedEmail.
type ListReceivedEmail struct {
	Id          string               `json:"id"`
	To          []string             `json:"to"`
	From        string               `json:"from"`
	CreatedAt   string               `json:"created_at"`
	Subject     string               `json:"subject"`
	Bcc         []string             `json:"bcc"`
	Cc          []string             `json:"cc"`
	ReplyTo     []string             `json:"reply_to"`
	Attachments []ReceivedAttachment `json:"attachments"`
}

// ListReceivedEmailsResponse is the response from the Receiving.List call.
type ListReceivedEmailsResponse struct {
	Object  string              `json:"object"`
	HasMore bool                `json:"has_more"`
	Data    []ListReceivedEmail `json:"data"`
}

// ReceivedAttachment represents an attachment in a received email (used in list responses without download URLs)
type ReceivedAttachment struct {
	Id                 string `json:"id"`
	Filename           string `json:"filename"`
	ContentType        string `json:"content_type"`
	ContentDisposition string `json:"content_disposition"`
	ContentId          string `json:"content_id"`
}

// ReceivingSvc handles operations for received/inbound emails
type ReceivingSvc interface {
	GetWithContext(ctx context.Context, emailId string) (*ReceivedEmail, error)
	Get(emailId string) (*ReceivedEmail, error)
	ListWithOptions(ctx context.Context, options *ListOptions) (ListReceivedEmailsResponse, error)
	ListWithContext(ctx context.Context) (ListReceivedEmailsResponse, error)
	List() (ListReceivedEmailsResponse, error)
	GetAttachmentWithContext(ctx context.Context, emailId string, attachmentId string) (*EmailAttachment, error)
	GetAttachment(emailId string, attachmentId string) (*EmailAttachment, error)
	ListAttachmentsWithOptions(ctx context.Context, emailId string, options *ListOptions) (ListEmailAttachmentsResponse, error)
	ListAttachmentsWithContext(ctx context.Context, emailId string) (ListEmailAttachmentsResponse, error)
	ListAttachments(emailId string) (ListEmailAttachmentsResponse, error)
}

// ReceivingSvcImpl is the implementation of the ReceivingSvc interface
type ReceivingSvcImpl struct {
	client *Client
}

// GetWithContext retrieves a received email with the given emailId
// https://resend.com/docs/api-reference/emails/retrieve-received-email
func (s *ReceivingSvcImpl) GetWithContext(ctx context.Context, emailId string) (*ReceivedEmail, error) {
	path := "emails/receiving/" + emailId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, ErrFailedToCreateReceivingGetRequest
	}

	emailResponse := new(ReceivedEmail)

	// Send Request
	_, err = s.client.Perform(req, emailResponse)

	if err != nil {
		return nil, err
	}

	return emailResponse, nil
}

// Get retrieves a received email with the given emailId
// https://resend.com/docs/api-reference/emails/retrieve-received-email
func (s *ReceivingSvcImpl) Get(emailId string) (*ReceivedEmail, error) {
	return s.GetWithContext(context.Background(), emailId)
}

// ListWithOptions retrieves a list of received emails with pagination options
// https://resend.com/docs/api-reference/emails/list-received-emails
func (s *ReceivingSvcImpl) ListWithOptions(ctx context.Context, options *ListOptions) (ListReceivedEmailsResponse, error) {
	path := "emails/receiving" + buildPaginationQuery(options)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListReceivedEmailsResponse{}, ErrFailedToCreateReceivingListRequest
	}

	listEmailsResponse := new(ListReceivedEmailsResponse)

	_, err = s.client.Perform(req, listEmailsResponse)

	if err != nil {
		return ListReceivedEmailsResponse{}, err
	}

	return *listEmailsResponse, nil
}

// ListWithContext retrieves a list of received emails
// https://resend.com/docs/api-reference/emails/list-received-emails
func (s *ReceivingSvcImpl) ListWithContext(ctx context.Context) (ListReceivedEmailsResponse, error) {
	return s.ListWithOptions(ctx, nil)
}

// List retrieves a list of received emails
// https://resend.com/docs/api-reference/emails/list-received-emails
func (s *ReceivingSvcImpl) List() (ListReceivedEmailsResponse, error) {
	return s.ListWithContext(context.Background())
}

// GetAttachmentWithContext retrieves a single attachment from a received email with the given emailId and attachmentId
// https://resend.com/docs/api-reference/attachments/retrieve-received-email-attachment
func (s *ReceivingSvcImpl) GetAttachmentWithContext(ctx context.Context, emailId string, attachmentId string) (*EmailAttachment, error) {
	path := "emails/receiving/" + emailId + "/attachments/" + attachmentId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, ErrFailedToCreateReceivingGetAttachmentRequest
	}

	attachment := new(EmailAttachment)

	// Send Request
	_, err = s.client.Perform(req, attachment)

	if err != nil {
		return nil, err
	}

	return attachment, nil
}

// GetAttachment retrieves a single attachment from a received email with the given emailId and attachmentId
// https://resend.com/docs/api-reference/attachments/retrieve-received-email-attachment
func (s *ReceivingSvcImpl) GetAttachment(emailId string, attachmentId string) (*EmailAttachment, error) {
	return s.GetAttachmentWithContext(context.Background(), emailId, attachmentId)
}

// ListAttachmentsWithOptions retrieves a list of attachments for a received email with pagination options
// https://resend.com/docs/api-reference/attachments/list-received-email-attachments
func (s *ReceivingSvcImpl) ListAttachmentsWithOptions(ctx context.Context, emailId string, options *ListOptions) (ListEmailAttachmentsResponse, error) {
	path := "emails/receiving/" + emailId + "/attachments" + buildPaginationQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListEmailAttachmentsResponse{}, ErrFailedToCreateReceivingListAttachmentsRequest
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

// ListAttachmentsWithContext retrieves a list of attachments for a received email
// https://resend.com/docs/api-reference/attachments/list-received-email-attachments
func (s *ReceivingSvcImpl) ListAttachmentsWithContext(ctx context.Context, emailId string) (ListEmailAttachmentsResponse, error) {
	return s.ListAttachmentsWithOptions(ctx, emailId, nil)
}

// ListAttachments retrieves a list of attachments for a received email
// https://resend.com/docs/api-reference/attachments/list-received-email-attachments
func (s *ReceivingSvcImpl) ListAttachments(emailId string) (ListEmailAttachmentsResponse, error) {
	return s.ListAttachmentsWithContext(context.Background(), emailId)
}
