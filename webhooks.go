package resend

import (
	"context"
	"net/http"
)

// CreateWebhookRequest represents the parameters for creating a webhook
type CreateWebhookRequest struct {
	Endpoint string   `json:"endpoint"`
	Events   []string `json:"events"`
}

// CreateWebhookResponse represents the response from creating a webhook
type CreateWebhookResponse struct {
	Object        string `json:"object"`
	Id            string `json:"id"`
	SigningSecret string `json:"signing_secret"`
}

// Webhook represents a webhook object
type Webhook struct {
	Object        string   `json:"object"`
	Id            string   `json:"id"`
	CreatedAt     string   `json:"created_at,omitempty"`
	Status        string   `json:"status,omitempty"`
	Endpoint      string   `json:"endpoint,omitempty"`
	Events        []string `json:"events,omitempty"`
	SigningSecret string   `json:"signing_secret,omitempty"`
}

// UpdateWebhookRequest represents the parameters for updating a webhook
type UpdateWebhookRequest struct {
	Endpoint *string  `json:"endpoint,omitempty"`
	Events   []string `json:"events,omitempty"`
	Status   *string  `json:"status,omitempty"`
}

// UpdateWebhookResponse represents the response from updating a webhook
type UpdateWebhookResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

// ListWebhooksResponse represents the response from listing webhooks
type ListWebhooksResponse struct {
	Object  string          `json:"object"`
	HasMore bool            `json:"has_more"`
	Data    []WebhookInList `json:"data"`
}

// WebhookInList represents a webhook in the list response
type WebhookInList struct {
	Id        string   `json:"id"`
	CreatedAt string   `json:"created_at"`
	Status    string   `json:"status"`
	Endpoint  string   `json:"endpoint"`
	Events    []string `json:"events"`
}

// DeleteWebhookResponse represents the response from deleting a webhook
type DeleteWebhookResponse struct {
	Object  string `json:"object"`
	Id      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// WebhooksSvc defines the interface for webhook operations
type WebhooksSvc interface {
	CreateWithContext(ctx context.Context, params *CreateWebhookRequest) (*CreateWebhookResponse, error)
	Create(params *CreateWebhookRequest) (*CreateWebhookResponse, error)
	GetWithContext(ctx context.Context, webhookId string) (*Webhook, error)
	Get(webhookId string) (*Webhook, error)
	UpdateWithContext(ctx context.Context, webhookId string, params *UpdateWebhookRequest) (*UpdateWebhookResponse, error)
	Update(webhookId string, params *UpdateWebhookRequest) (*UpdateWebhookResponse, error)
	ListWithOptions(ctx context.Context, options *ListOptions) (*ListWebhooksResponse, error)
	ListWithContext(ctx context.Context) (*ListWebhooksResponse, error)
	List() (*ListWebhooksResponse, error)
	RemoveWithContext(ctx context.Context, webhookId string) (*DeleteWebhookResponse, error)
	Remove(webhookId string) (*DeleteWebhookResponse, error)
}

// WebhooksSvcImpl implements the WebhooksSvc interface
type WebhooksSvcImpl struct {
	client *Client
}

// CreateWithContext creates a new webhook with the given context
// https://resend.com/docs/api-reference/webhooks/create-webhook
func (s *WebhooksSvcImpl) CreateWithContext(ctx context.Context, params *CreateWebhookRequest) (*CreateWebhookResponse, error) {
	path := "webhooks"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return nil, err
	}

	// Build response object
	webhookResp := new(CreateWebhookResponse)

	// Send Request
	_, err = s.client.Perform(req, webhookResp)
	if err != nil {
		return nil, err
	}

	return webhookResp, nil
}

// Create creates a new webhook
func (s *WebhooksSvcImpl) Create(params *CreateWebhookRequest) (*CreateWebhookResponse, error) {
	return s.CreateWithContext(context.Background(), params)
}

// GetWithContext retrieves a webhook by ID with the given context
// https://resend.com/docs/api-reference/webhooks/get-webhook
func (s *WebhooksSvcImpl) GetWithContext(ctx context.Context, webhookId string) (*Webhook, error) {
	path := "webhooks/" + webhookId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// Build response object
	webhookResp := new(Webhook)

	// Send Request
	_, err = s.client.Perform(req, webhookResp)
	if err != nil {
		return nil, err
	}

	return webhookResp, nil
}

// Get retrieves a webhook by ID
func (s *WebhooksSvcImpl) Get(webhookId string) (*Webhook, error) {
	return s.GetWithContext(context.Background(), webhookId)
}

// UpdateWithContext updates a webhook with the given context
// https://resend.com/docs/api-reference/webhooks/update-webhook
func (s *WebhooksSvcImpl) UpdateWithContext(ctx context.Context, webhookId string, params *UpdateWebhookRequest) (*UpdateWebhookResponse, error) {
	path := "webhooks/" + webhookId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, params)
	if err != nil {
		return nil, err
	}

	// Build response object
	webhookResp := new(UpdateWebhookResponse)

	// Send Request
	_, err = s.client.Perform(req, webhookResp)
	if err != nil {
		return nil, err
	}

	return webhookResp, nil
}

// Update updates a webhook
func (s *WebhooksSvcImpl) Update(webhookId string, params *UpdateWebhookRequest) (*UpdateWebhookResponse, error) {
	return s.UpdateWithContext(context.Background(), webhookId, params)
}

// ListWithOptions lists all webhooks with pagination options
// https://resend.com/docs/api-reference/webhooks/list-webhooks
func (s *WebhooksSvcImpl) ListWithOptions(ctx context.Context, options *ListOptions) (*ListWebhooksResponse, error) {
	path := "webhooks" + buildPaginationQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// Build response object
	webhooksResp := new(ListWebhooksResponse)

	// Send Request
	_, err = s.client.Perform(req, webhooksResp)
	if err != nil {
		return nil, err
	}

	return webhooksResp, nil
}

// ListWithContext lists all webhooks with the given context
func (s *WebhooksSvcImpl) ListWithContext(ctx context.Context) (*ListWebhooksResponse, error) {
	return s.ListWithOptions(ctx, nil)
}

// List lists all webhooks
func (s *WebhooksSvcImpl) List() (*ListWebhooksResponse, error) {
	return s.ListWithContext(context.Background())
}

// RemoveWithContext deletes a webhook by ID with the given context
// https://resend.com/docs/api-reference/webhooks/delete-webhook
func (s *WebhooksSvcImpl) RemoveWithContext(ctx context.Context, webhookId string) (*DeleteWebhookResponse, error) {
	path := "webhooks/" + webhookId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	// Build response object
	webhookResp := new(DeleteWebhookResponse)

	// Send Request
	_, err = s.client.Perform(req, webhookResp)
	if err != nil {
		return nil, err
	}

	return webhookResp, nil
}

// Remove deletes a webhook by ID
func (s *WebhooksSvcImpl) Remove(webhookId string) (*DeleteWebhookResponse, error) {
	return s.RemoveWithContext(context.Background(), webhookId)
}
