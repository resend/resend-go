package resend

import (
	"errors"
	"net/http"
)

// https://resend.com/docs/api-reference/api-keys/create-api-key
type CreateApiKeyRequest struct {
	Name       string `json:"name"`
	Permission string `json:"permission,omitempty"` // TODO: update permission to type
	DomainId   string `json:"domain_id,omitempty"`
}

type CreateApiKeyResponse struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}

type ListApiKeysResponse struct {
	Data []ApiKey `json:"data"`
}

type ApiKey struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type ApiKeysSvc interface {
	Create(*CreateApiKeyRequest) (CreateApiKeyResponse, error)
}

type ApiKeysSvcImpl struct {
	client *Client
}

// Create creates a new API Key based on the given params
func (s *ApiKeysSvcImpl) Create(params *CreateApiKeyRequest) (CreateApiKeyResponse, error) {
	path := "api-keys"

	// Prepare request
	req, err := s.client.NewRequest(http.MethodPost, path, params)
	if err != nil {
		return CreateApiKeyResponse{}, errors.New("[ERROR]: Failed to create CreateApiKey request")
	}

	// Build response recipient obj
	emailResponse := new(CreateApiKeyResponse)

	// Send Request
	_, err = s.client.Perform(req, emailResponse)

	if err != nil {
		return CreateApiKeyResponse{}, err
	}

	return *emailResponse, nil
}
