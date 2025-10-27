package resend

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// Email events
	EventEmailSent            = "email.sent"
	EventEmailDelivered       = "email.delivered"
	EventEmailDeliveryDelayed = "email.delivery_delayed"
	EventEmailComplained      = "email.complained"
	EventEmailBounced         = "email.bounced"
	EventEmailOpened          = "email.opened"
	EventEmailClicked         = "email.clicked"
	EventEmailReceived        = "email.received"
	EventEmailFailed          = "email.failed"

	// Contact events
	EventContactCreated = "contact.created"
	EventContactUpdated = "contact.updated"
	EventContactDeleted = "contact.deleted"

	// Domain events
	EventDomainCreated = "domain.created"
	EventDomainUpdated = "domain.updated"
	EventDomainDeleted = "domain.deleted"
)

// Default tolerance for timestamp validation (5 minutes)
const DefaultWebhookToleranceSeconds = 300

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

// WebhookHeaders represents the webhook verification headers
type WebhookHeaders struct {
	Id        string // svix-id header
	Timestamp string // svix-timestamp header
	Signature string // svix-signature header
}

// VerifyWebhookOptions represents the parameters for webhook verification
type VerifyWebhookOptions struct {
	Payload       string         // Raw webhook payload body
	Headers       WebhookHeaders // Webhook headers from the request
	WebhookSecret string         // Signing secret (from webhook creation response)
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
	Verify(options *VerifyWebhookOptions) error
}

// WebhooksSvcImpl implements the WebhooksSvc interface
type WebhooksSvcImpl struct {
	client *Client
}

// CreateWithContext creates a new webhook with the given context
// https://resend.com/docs/api-reference/webhooks/create-webhook
func (s *WebhooksSvcImpl) CreateWithContext(ctx context.Context, params *CreateWebhookRequest) (*CreateWebhookResponse, error) {
	path := "webhooks"

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return nil, err
	}

	webhookResp := new(CreateWebhookResponse)

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

// Verify validates a webhook payload using HMAC-SHA256 signature verification
// This implements manual verification without external dependencies
// https://docs.svix.com/receiving/verifying-payloads/how-manual
func (s *WebhooksSvcImpl) Verify(options *VerifyWebhookOptions) error {
	if options == nil {
		return errors.New("options cannot be nil")
	}

	if options.Payload == "" {
		return errors.New("payload cannot be empty")
	}

	if options.WebhookSecret == "" {
		return errors.New("webhook secret cannot be empty")
	}

	if options.Headers.Id == "" {
		return errors.New("svix-id header is required")
	}

	if options.Headers.Timestamp == "" {
		return errors.New("svix-timestamp header is required")
	}

	if options.Headers.Signature == "" {
		return errors.New("svix-signature header is required")
	}

	// Step 1: Validate timestamp to prevent replay attacks
	timestamp, err := strconv.ParseInt(options.Headers.Timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid timestamp format: %w", err)
	}

	now := time.Now().Unix()
	diff := now - timestamp
	if diff > DefaultWebhookToleranceSeconds || diff < -DefaultWebhookToleranceSeconds {
		return fmt.Errorf("timestamp outside tolerance window: difference of %d seconds", diff)
	}

	// Step 2: Construct signed content: {id}.{timestamp}.{payload}
	signedContent := fmt.Sprintf("%s.%s.%s", options.Headers.Id, options.Headers.Timestamp, options.Payload)

	// Step 3: Decode the signing secret (strip whsec_ prefix and base64 decode)
	secret := strings.TrimPrefix(options.WebhookSecret, "whsec_")
	decodedSecret, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return fmt.Errorf("failed to decode webhook secret: %w", err)
	}

	// Step 4: Calculate expected signature using HMAC-SHA256
	expectedSignature := generateSignature(decodedSecret, []byte(signedContent))

	// Step 5: Compare signatures using constant-time comparison
	// The signature header contains space-separated signatures with version prefixes (e.g., "v1,sig1 v1,sig2")
	signatures := strings.Split(options.Headers.Signature, " ")
	for _, sig := range signatures {
		// Strip version prefix (e.g., "v1,")
		parts := strings.SplitN(sig, ",", 2)
		if len(parts) != 2 {
			continue
		}

		receivedSignature := parts[1]
		if subtle.ConstantTimeCompare([]byte(expectedSignature), []byte(receivedSignature)) == 1 {
			return nil // Signature matches
		}
	}

	return errors.New("no matching signature found")
}

// generateSignature creates an HMAC-SHA256 signature and returns it as base64
func generateSignature(secret, content []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write(content)
	signature := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature)
}
