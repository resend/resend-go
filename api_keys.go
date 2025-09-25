package resend

import (
	"context"
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
	Object  string   `json:"object"`
	Data    []ApiKey `json:"data"`
	HasMore bool     `json:"has_more"`
}

type ApiKey struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

type ApiKeysSvc interface {
	CreateWithContext(ctx context.Context, params *CreateApiKeyRequest) (CreateApiKeyResponse, error)
	Create(params *CreateApiKeyRequest) (CreateApiKeyResponse, error)
	ListWithOptions(ctx context.Context, options *ListOptions) (ListApiKeysResponse, error)
	ListWithContext(ctx context.Context) (ListApiKeysResponse, error)
	List() (ListApiKeysResponse, error)
	RemoveWithContext(ctx context.Context, apiKeyId string) (bool, error)
	Remove(apiKeyId string) (bool, error)
}

type ApiKeysSvcImpl struct {
	client *Client
}

// CreateWithContext creates a new API Key based on the given params
// https://resend.com/docs/api-reference/api-keys/create-api-key
func (s *ApiKeysSvcImpl) CreateWithContext(ctx context.Context, params *CreateApiKeyRequest) (CreateApiKeyResponse, error) {
	path := "api-keys"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return CreateApiKeyResponse{}, ErrFailedToCreateApiKeysCreateRequest
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

// Create creates a new API Key based on the given params
func (s *ApiKeysSvcImpl) Create(params *CreateApiKeyRequest) (CreateApiKeyResponse, error) {
	return s.CreateWithContext(context.Background(), params)
}

// ListWithOptions list all API Keys in the project with pagination options
// https://resend.com/docs/api-reference/api-keys/list-api-keys
func (s *ApiKeysSvcImpl) ListWithOptions(ctx context.Context, options *ListOptions) (ListApiKeysResponse, error) {
	path := "api-keys" + buildPaginationQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListApiKeysResponse{}, ErrFailedToCreateApiKeysListRequest
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

// ListWithContext list all API Keys in the project
// https://resend.com/docs/api-reference/api-keys/list-api-keys
func (s *ApiKeysSvcImpl) ListWithContext(ctx context.Context) (ListApiKeysResponse, error) {
	return s.ListWithOptions(ctx, nil)
}

// List all API Keys in the project
func (s *ApiKeysSvcImpl) List() (ListApiKeysResponse, error) {
	return s.ListWithContext(context.Background())
}

// RemoveWithContext deletes a given api key by id
// https://resend.com/docs/api-reference/api-keys/delete-api-key
func (s *ApiKeysSvcImpl) RemoveWithContext(ctx context.Context, apiKeyId string) (bool, error) {
	path := "api-keys/" + apiKeyId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return false, ErrFailedToCreateApiKeysRemoveRequest
	}

	// Send Request
	_, err = s.client.Perform(req, nil)

	if err != nil {
		return false, err
	}

	return true, nil
}

// Remove deletes a given api key by id
func (s *ApiKeysSvcImpl) Remove(apiKeyId string) (bool, error) {
	return s.RemoveWithContext(context.Background(), apiKeyId)
}
