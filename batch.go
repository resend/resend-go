package resend

import (
	"errors"
	"net/http"
)

// BatchSendEmailResponse is the response from the BatchSendEmail call.
type BatchEmailResponse struct {
	Data []SendEmailResponse `json:"data"`
}

type BatchSvc interface {
	Send([]*SendEmailRequest) (*BatchEmailResponse, error)
}

type BatchSvcImpl struct {
	client *Client
}

func (s *BatchSvcImpl) Send(params []*SendEmailRequest) (*BatchEmailResponse, error) {
	path := "emails/batch"

	// Prepare request
	req, err := s.client.NewRequest(http.MethodPost, path, params)
	if err != nil {
		return nil, errors.New("[ERROR]: Failed to create SendEmail request")
	}

	// Build response recipient obj
	batchSendEmailResponse := new(BatchEmailResponse)

	// Send Request
	_, err = s.client.Perform(req, batchSendEmailResponse)

	if err != nil {
		return nil, err
	}

	return batchSendEmailResponse, nil
}
