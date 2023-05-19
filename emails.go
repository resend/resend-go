package resend

import (
	"errors"
	"net/http"
)

// https://resend.com/docs/api-reference/emails/send-email
type SendEmailRequest struct {
	From        string       `json:"from"`
	To          []string     `json:"to"`
	Subject     string       `json:"subject"`
	Bcc         []string     `json:"bcc"`
	Cc          []string     `json:"cc"`
	ReplyTo     string       `json:"reply_to"`
	Html        string       `json:"html"`
	Text        string       `json:"text"`
	Tags        []Tag        `json:"tags"`
	Attachments []Attachment `json:"attachments"`
}

type SendEmailResponse struct {
	Id string `json:"id"`
}

type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Attachment struct {

	// Content must be a string representation of a byte array
	Content string `json:"content"`

	// Filename that will appear in the email.
	// Make sure you pick the correct extension otherwise preview
	// make not work as expected
	Filename string `json:"filename"`
}

type EmailsSvc interface {
	Send(*SendEmailRequest) (SendEmailResponse, error)
	Get() error
}

type EmailsSvcImpl struct {
	client *Client
}

// Send sends an email with the givem params
func (s *EmailsSvcImpl) Send(params *SendEmailRequest) (SendEmailResponse, error) {
	path := "emails"

	// Preapre request
	req, err := s.client.NewRequest(http.MethodPost, path, params)
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

func (s *EmailsSvcImpl) Get() error {
	path := "emails/"
	println(path)
	println("fetching email")
	return nil
}
