package resend

import (
	"errors"
	"net/http"
)

// _convertedSendEmailRequest is the same as SendEmailRequest
// but with Attachments type fixed, for backwards compatibility
// used only internally
type _convertedSendEmailRequest struct {
	From        string            `json:"from"`
	To          []string          `json:"to"`
	Subject     string            `json:"subject"`
	Bcc         []string          `json:"bcc"`
	Cc          []string          `json:"cc"`
	ReplyTo     string            `json:"reply_to"`
	Html        string            `json:"html"`
	Text        string            `json:"text"`
	Tags        []Tag             `json:"tags"`
	Attachments []_attachment     `json:"attachments"`
	Headers     map[string]string `json:"headers"`
}

// https://resend.com/docs/api-reference/emails/send-email
type SendEmailRequest struct {
	From        string            `json:"from"`
	To          []string          `json:"to"`
	Subject     string            `json:"subject"`
	Bcc         []string          `json:"bcc"`
	Cc          []string          `json:"cc"`
	ReplyTo     string            `json:"reply_to"`
	Html        string            `json:"html"`
	Text        string            `json:"text"`
	Tags        []Tag             `json:"tags"`
	Attachments []Attachment      `json:"attachments"`
	Headers     map[string]string `json:"headers"`
}

type SendEmailResponse struct {
	Id string `json:"id"`
}

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

type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// _attachment is the struct for internal use
// Different from Attachment, Content here is an array of strings
// which represent an array of bytes in string form
type _attachment struct {
	Content  []string `json:"content,omitempty"`
	Filename string   `json:"filename"`
	Path     string   `json:"path,omitempty"`
}

// Attachment is the public struct used for adding attachments to emails
type Attachment struct {

	// Content is a string here, but this will be converted back to
	// an array of strings representing an array of bytes
	Content string `json:"content,omitempty"`

	// Filename that will appear in the email.
	// Make sure you pick the correct extension otherwise preview
	// may not work as expected
	Filename string `json:"filename"`

	// Path where the attachment file is hosted
	Path string `json:"path,omitempty"`
}

type EmailsSvc interface {
	Send(*SendEmailRequest) (SendEmailResponse, error)
	Get(emailId string) (Email, error)
}

type EmailsSvcImpl struct {
	client *Client
}

// convertRequest gets a SendEmailRequest and builds an internal sendEmailRequest
func convertRequest(params *SendEmailRequest) (_convertedSendEmailRequest, error) {

	newReq := _convertedSendEmailRequest{}

	if params.To != nil {
		newReq.To = params.To
	}
	if params.From != "" {
		newReq.From = params.From
	}
	if params.Subject != "" {
		newReq.Subject = params.Subject
	}
	if params.Bcc != nil {
		newReq.Bcc = params.Bcc
	}
	if params.Cc != nil {
		newReq.Cc = params.Cc
	}
	if params.ReplyTo != "" {
		newReq.ReplyTo = params.ReplyTo
	}
	if params.Html != "" {
		newReq.Html = params.Html
	}
	if params.Text != "" {
		newReq.Text = params.Text
	}
	if params.Tags != nil {
		newReq.Tags = params.Tags
	}
	if params.Headers != nil {
		newReq.Headers = params.Headers
	}
	// Backwards compatibility Attachment handling
	if params.Attachments != nil {
		newReq.Attachments = PrepareAttachments(params.Attachments)
	}
	return newReq, nil
}

// Send sends an email with the given params
func (s *EmailsSvcImpl) Send(params *SendEmailRequest) (SendEmailResponse, error) {
	path := "emails"

	convertedParams, err := convertRequest(params)
	if err != nil {
		return SendEmailResponse{}, errors.New("[ERROR]: Failed to create SendEmail request")
	}

	// Prepare request
	req, err := s.client.NewRequest(http.MethodPost, path, convertedParams)
	if err != nil {
		return SendEmailResponse{}, errors.New("[ERROR]: Failed to create SendEmail request")
	}

	// Build response recipient obj
	emailResponse := new(SendEmailResponse)

	// Send Request
	_, err = s.client.Perform(req, emailResponse)

	if err != nil {
		return SendEmailResponse{}, err
	}

	return *emailResponse, nil
}

// Get retrives an email with the given emailId
// https://resend.com/docs/api-reference/emails/retrieve-email
func (s *EmailsSvcImpl) Get(emailId string) (Email, error) {
	path := "emails/" + emailId

	// Prepare request
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return Email{}, errors.New("[ERROR]: Failed to create GetEmail request")
	}

	// Build response recipient obj
	emailResponse := new(Email)

	// Send Request
	_, err = s.client.Perform(req, emailResponse)

	if err != nil {
		return Email{}, err
	}

	return *emailResponse, nil
}

// PrepareAttachments converts a Attachment into _attachment
func PrepareAttachments(attachments []Attachment) []_attachment {

	var atts []_attachment

	// Loop through attachments and transform
	for _, a := range attachments {
		attachment := _attachment{}

		if a.Content != "" {
			attachment.Content = ByteArrayToStringArray([]byte(a.Content))
		}

		if a.Filename != "" {
			attachment.Filename = a.Filename
		}

		if a.Path != "" {
			attachment.Path = a.Path
		}

		atts = append(atts, attachment)
	}
	return atts
}
