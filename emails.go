package resend

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// SendEmailRequest is the request object for the SendEmail call.
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
}

// SendEmailResponse is the response from the SendEmail call.
type SendEmailResponse struct {
	Id string `json:"id"`
}

// Email provides the structure for the response from the GetEmail call.
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
}

// MarshalJSON overrides the regular JSON Marshaller to ensure that the
// attachment content is provided in the way Resend expects.
func (a *Attachment) MarshalJSON() ([]byte, error) {
	na := struct {
		Content  []int  `json:"content,omitempty"`
		Filename string `json:"filename,omitempty"`
		Path     string `json:"path,omitempty"`
	}{
		Filename: a.Filename,
		Path:     a.Path,
		Content:  BytesToIntArray(a.Content),
	}
	return json.Marshal(na)
}

type EmailsSvc interface {
	Send(ctx context.Context, params *SendEmailRequest) (*SendEmailResponse, error)
	Get(ctx context.Context, emailID string) (*Email, error)
}

type EmailsSvcImpl struct {
	client *Client
}

// Send sends an email with the given params
func (s *EmailsSvcImpl) Send(ctx context.Context, params *SendEmailRequest) (*SendEmailResponse, error) {
	path := "emails"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return nil, errors.New("[ERROR]: Failed to create SendEmail request")
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

// Get retrieves an email with the given emailID
// https://resend.com/docs/api-reference/emails/retrieve-email
func (s *EmailsSvcImpl) Get(ctx context.Context, emailID string) (*Email, error) {
	path := "emails/" + emailID

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, errors.New("[ERROR]: Failed to create GetEmail request")
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
