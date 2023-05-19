package resend

import (
	"errors"
	"net/http"
)

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
	List() (ListApiKeysResponse, error)
	Remove(apiKeyId string) (bool, error)
}

type ApiKeysSvcImpl struct {
	client *Client
}

// Create creates a new API Key based on the given params
// https://resend.com/docs/api-reference/api-keys/create-api-key
func (s *ApiKeysSvcImpl) Create(params *CreateApiKeyRequest) (CreateApiKeyResponse, error) {
	path := "api-keys"

	// Prepare request
	req, err := s.client.NewRequest(http.MethodPost, path, params)
	if err != nil {
		return CreateApiKeyResponse{}, errors.New("[ERROR]: Failed to create ApiKeys.Create request")
	}

	// Build response recipient obj
	apiKeysResp := new(CreateApiKeyResponse)

	// Send Request
	_, err = s.client.Perform(req, apiKeysResp)

	if err != nil {
		return CreateApiKeyResponse{}, err
	}

	return *apiKeysResp, nil
}

// List list all API Keys in the project
// https://resend.com/docs/api-reference/api-keys/list-api-keys
func (s *ApiKeysSvcImpl) List() (ListApiKeysResponse, error) {
	path := "api-keys"

	// Prepare request
	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return ListApiKeysResponse{}, errors.New("[ERROR]: Failed to create ApiKeys.List request")
	}

	// Build response recipient obj
	apiKeysResp := new(ListApiKeysResponse)

	// Send Request
	_, err = s.client.Perform(req, apiKeysResp)

	if err != nil {
		return ListApiKeysResponse{}, err
	}

	return *apiKeysResp, nil
}

// Remove deletes a given api key by id
// https://resend.com/docs/api-reference/api-keys/delete-api-key
func (s *ApiKeysSvcImpl) Remove(apiKeyId string) (bool, error) {
	path := "api-keys/" + apiKeyId

	// Prepare request
	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return false, errors.New("[ERROR]: Failed to create ApiKeys.List request")
	}

	// Send Request
	_, err = s.client.Perform(req, nil)

	if err != nil {
		return false, err
	}

	return true, nil
}
