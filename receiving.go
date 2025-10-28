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

// ReceivedAttachment represents an attachment in a received email
type ReceivedAttachment struct {
	Id                 string `json:"id"`
	Filename           string `json:"filename"`
	ContentType        string `json:"content_type"`
	ContentDisposition string `json:"content_disposition"`
	ContentId          string `json:"content_id"`
}

// ReceivedEmailAttachment represents the full attachment details including download URL
type ReceivedEmailAttachment struct {
	Id                 string `json:"id"`
	Filename           string `json:"filename"`
	ContentType        string `json:"content_type"`
	ContentDisposition string `json:"content_disposition"`
	ContentId          string `json:"content_id"`
	DownloadUrl        string `json:"download_url"`
	ExpiresAt          string `json:"expires_at"`
}

// receivedEmailAttachmentResponse wraps the API response for a single attachment
type receivedEmailAttachmentResponse struct {
	Object string                   `json:"object"`
	Data   ReceivedEmailAttachment `json:"data"`
}

// ListReceivedEmailAttachmentsResponse is the response from the Receiving.ListAttachments call.
type ListReceivedEmailAttachmentsResponse struct {
	Object  string                    `json:"object"`
	HasMore bool                      `json:"has_more"`
	Data    []ReceivedEmailAttachment `json:"data"`
}

// ReceivingSvc handles operations for received/inbound emails
type ReceivingSvc interface {
	GetWithContext(ctx context.Context, emailID string) (*ReceivedEmail, error)
	Get(emailID string) (*ReceivedEmail, error)
	ListWithOptions(ctx context.Context, options *ListOptions) (ListReceivedEmailsResponse, error)
	ListWithContext(ctx context.Context) (ListReceivedEmailsResponse, error)
	List() (ListReceivedEmailsResponse, error)
	GetAttachmentWithContext(ctx context.Context, emailID string, attachmentID string) (*ReceivedEmailAttachment, error)
	GetAttachment(emailID string, attachmentID string) (*ReceivedEmailAttachment, error)
	ListAttachmentsWithOptions(ctx context.Context, emailID string, options *ListOptions) (ListReceivedEmailAttachmentsResponse, error)
	ListAttachmentsWithContext(ctx context.Context, emailID string) (ListReceivedEmailAttachmentsResponse, error)
	ListAttachments(emailID string) (ListReceivedEmailAttachmentsResponse, error)
}

// ReceivingSvcImpl is the implementation of the ReceivingSvc interface
type ReceivingSvcImpl struct {
	client *Client
}

// GetWithContext retrieves a received email with the given emailID
// https://resend.com/docs/api-reference/emails/retrieve-received-email
func (s *ReceivingSvcImpl) GetWithContext(ctx context.Context, emailID string) (*ReceivedEmail, error) {
	path := "emails/receiving/" + emailID

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

// Get retrieves a received email with the given emailID
// https://resend.com/docs/api-reference/emails/retrieve-received-email
func (s *ReceivingSvcImpl) Get(emailID string) (*ReceivedEmail, error) {
	return s.GetWithContext(context.Background(), emailID)
}

// ListWithOptions retrieves a list of received emails with pagination options
// https://resend.com/docs/api-reference/emails/retrieve-received-email
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
// https://resend.com/docs/api-reference/emails/retrieve-received-email
func (s *ReceivingSvcImpl) ListWithContext(ctx context.Context) (ListReceivedEmailsResponse, error) {
	return s.ListWithOptions(ctx, nil)
}

// List retrieves a list of received emails
// https://resend.com/docs/api-reference/emails/retrieve-received-email
func (s *ReceivingSvcImpl) List() (ListReceivedEmailsResponse, error) {
	return s.ListWithContext(context.Background())
}

// GetAttachmentWithContext retrieves a single attachment from a received email with the given emailID and attachmentID
// https://resend.com/docs/api-reference/attachments/retrieve-received-email-attachment
func (s *ReceivingSvcImpl) GetAttachmentWithContext(ctx context.Context, emailID string, attachmentID string) (*ReceivedEmailAttachment, error) {
	path := "emails/receiving/" + emailID + "/attachments/" + attachmentID

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, ErrFailedToCreateReceivingGetAttachmentRequest
	}

	// Build response wrapper obj
	attachmentResponse := new(receivedEmailAttachmentResponse)

	// Send Request
	_, err = s.client.Perform(req, attachmentResponse)

	if err != nil {
		return nil, err
	}

	return &attachmentResponse.Data, nil
}

// GetAttachment retrieves a single attachment from a received email with the given emailID and attachmentID
// https://resend.com/docs/api-reference/attachments/retrieve-received-email-attachment
func (s *ReceivingSvcImpl) GetAttachment(emailID string, attachmentID string) (*ReceivedEmailAttachment, error) {
	return s.GetAttachmentWithContext(context.Background(), emailID, attachmentID)
}

// ListAttachmentsWithOptions retrieves a list of attachments for a received email with pagination options
// https://resend.com/docs/api-reference/attachments/list-received-email-attachments
func (s *ReceivingSvcImpl) ListAttachmentsWithOptions(ctx context.Context, emailID string, options *ListOptions) (ListReceivedEmailAttachmentsResponse, error) {
	path := "emails/receiving/" + emailID + "/attachments" + buildPaginationQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListReceivedEmailAttachmentsResponse{}, ErrFailedToCreateReceivingListAttachmentsRequest
	}

	// Build response recipient obj
	listAttachmentsResponse := new(ListReceivedEmailAttachmentsResponse)

	// Send Request
	_, err = s.client.Perform(req, listAttachmentsResponse)

	if err != nil {
		return ListReceivedEmailAttachmentsResponse{}, err
	}

	return *listAttachmentsResponse, nil
}

// ListAttachmentsWithContext retrieves a list of attachments for a received email
// https://resend.com/docs/api-reference/attachments/list-received-email-attachments
func (s *ReceivingSvcImpl) ListAttachmentsWithContext(ctx context.Context, emailID string) (ListReceivedEmailAttachmentsResponse, error) {
	return s.ListAttachmentsWithOptions(ctx, emailID, nil)
}

// ListAttachments retrieves a list of attachments for a received email
// https://resend.com/docs/api-reference/attachments/list-received-email-attachments
func (s *ReceivingSvcImpl) ListAttachments(emailID string) (ListReceivedEmailAttachmentsResponse, error) {
	return s.ListAttachmentsWithContext(context.Background(), emailID)
}
