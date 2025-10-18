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

// ReceivingSvc handles operations for received/inbound emails
type ReceivingSvc interface {
	GetWithContext(ctx context.Context, emailID string) (*ReceivedEmail, error)
	Get(emailID string) (*ReceivedEmail, error)
	ListWithOptions(ctx context.Context, options *ListOptions) (ListReceivedEmailsResponse, error)
	ListWithContext(ctx context.Context) (ListReceivedEmailsResponse, error)
	List() (ListReceivedEmailsResponse, error)
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

	// Build response recipient obj
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

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListReceivedEmailsResponse{}, ErrFailedToCreateReceivingListRequest
	}

	// Build response recipient obj
	listEmailsResponse := new(ListReceivedEmailsResponse)

	// Send Request
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
