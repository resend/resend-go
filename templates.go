package resend

import (
	"context"
	"net/http"
)

// VariableType represents the type of a template variable
type VariableType string

const (
	VariableTypeString  VariableType = "string"
	VariableTypeNumber  VariableType = "number"
	VariableTypeBoolean VariableType = "boolean"
	VariableTypeObject  VariableType = "object"
	VariableTypeList    VariableType = "list"
)

// TemplateVariable represents a variable in a template
// Important:
// - All variables used in the template HTML (e.g., {{{NAME}}}) must be declared in the Variables array
// - Variables of type 'object' and 'list' REQUIRE a FallbackValue, or the API will return an error
// - For 'list' type: FallbackValue must be a non-empty array (e.g., []interface{}{"item"})
// - For 'object' type: FallbackValue must be a valid object (e.g., map[string]interface{}{"key": "value"})
// - Variables of type 'string', 'number', and 'boolean' can have optional FallbackValue
type TemplateVariable struct {
	Key           string       `json:"key"`
	Type          VariableType `json:"type"`
	FallbackValue interface{}  `json:"fallback_value,omitempty"`
}

// CreateTemplateRequest is the request payload for creating a template
// Important: All variables referenced in Html (e.g., {{{NAME}}}) must be
// declared in the Variables array, or the API will return a validation error.
type CreateTemplateRequest struct {
	Name      string              `json:"name"`
	Alias     string              `json:"alias,omitempty"`
	From      string              `json:"from,omitempty"`
	Subject   string              `json:"subject,omitempty"`
	ReplyTo   interface{}         `json:"reply_to,omitempty"` // string or []string
	Html      string              `json:"html"`
	Text      string              `json:"text,omitempty"`
	Variables []*TemplateVariable `json:"variables,omitempty"`
}

// CreateTemplateResponse is the response from creating a template
type CreateTemplateResponse struct {
	Id     string `json:"id"`
	Object string `json:"object"`
}

// UpdateTemplateRequest is the request payload for updating a template
// Important: All variables referenced in Html (e.g., {{{NAME}}}) must be
// declared in the Variables array, or the API will return a validation error.
type UpdateTemplateRequest struct {
	Name      string              `json:"name"`
	Alias     string              `json:"alias,omitempty"`
	From      string              `json:"from,omitempty"`
	Subject   string              `json:"subject,omitempty"`
	ReplyTo   interface{}         `json:"reply_to,omitempty"` // string or []string
	Html      string              `json:"html"`
	Text      string              `json:"text,omitempty"`
	Variables []*TemplateVariable `json:"variables,omitempty"`
}

// UpdateTemplateResponse is the response from updating a template
type UpdateTemplateResponse struct {
	Id     string `json:"id"`
	Object string `json:"object"`
}

// PublishTemplateResponse is the response from publishing a template
type PublishTemplateResponse struct {
	Id     string `json:"id"`
	Object string `json:"object"`
}

// TemplateVariableResponse represents a variable in a template response (with additional fields)
type TemplateVariableResponse struct {
	Id            string       `json:"id"`
	Key           string       `json:"key"`
	Type          VariableType `json:"type"`
	FallbackValue interface{}  `json:"fallback_value"`
	CreatedAt     string       `json:"created_at"`
	UpdatedAt     string       `json:"updated_at"`
}

// Template represents a full template object returned by the Get endpoint
type Template struct {
	Object      string                      `json:"object"`
	Id          string                      `json:"id"`
	Alias       string                      `json:"alias"`
	Name        string                      `json:"name"`
	CreatedAt   string                      `json:"created_at"`
	UpdatedAt   string                      `json:"updated_at"`
	Status      string                      `json:"status"`
	PublishedAt string                      `json:"published_at"`
	From        string                      `json:"from"`
	Subject     string                      `json:"subject"`
	ReplyTo     interface{}                 `json:"reply_to"` // string, []string, or null
	Html        string                      `json:"html"`
	Text        string                      `json:"text"`
	Variables   []*TemplateVariableResponse `json:"variables"`
}

// TemplatesSvc handles operations for templates
type TemplatesSvc interface {
	CreateWithContext(ctx context.Context, params *CreateTemplateRequest) (*CreateTemplateResponse, error)
	Create(params *CreateTemplateRequest) (*CreateTemplateResponse, error)
	GetWithContext(ctx context.Context, identifier string) (*Template, error)
	Get(identifier string) (*Template, error)
	UpdateWithContext(ctx context.Context, identifier string, params *UpdateTemplateRequest) (*UpdateTemplateResponse, error)
	Update(identifier string, params *UpdateTemplateRequest) (*UpdateTemplateResponse, error)
	PublishWithContext(ctx context.Context, identifier string) (*PublishTemplateResponse, error)
	Publish(identifier string) (*PublishTemplateResponse, error)
}

// TemplatesSvcImpl is the implementation of the TemplatesSvc interface
type TemplatesSvcImpl struct {
	client *Client
}

// CreateWithContext creates a new template with the given parameters
// https://resend.com/docs/api-reference/templates/create-template
func (s *TemplatesSvcImpl) CreateWithContext(ctx context.Context, params *CreateTemplateRequest) (*CreateTemplateResponse, error) {
	path := "templates"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return nil, ErrFailedToCreateTemplateCreateRequest
	}

	// Build response recipient obj
	templateResponse := new(CreateTemplateResponse)

	// Send Request
	_, err = s.client.Perform(req, templateResponse)

	if err != nil {
		return nil, err
	}

	return templateResponse, nil
}

// Create creates a new template with the given parameters
// https://resend.com/docs/api-reference/templates/create-template
func (s *TemplatesSvcImpl) Create(params *CreateTemplateRequest) (*CreateTemplateResponse, error) {
	return s.CreateWithContext(context.Background(), params)
}

// GetWithContext retrieves a template by ID or alias
// https://resend.com/docs/api-reference/templates/get-template
func (s *TemplatesSvcImpl) GetWithContext(ctx context.Context, identifier string) (*Template, error) {
	path := "templates/" + identifier

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, ErrFailedToCreateTemplateGetRequest
	}

	// Build response recipient obj
	templateResponse := new(Template)

	// Send Request
	_, err = s.client.Perform(req, templateResponse)

	if err != nil {
		return nil, err
	}

	return templateResponse, nil
}

// Get retrieves a template by ID or alias
// https://resend.com/docs/api-reference/templates/get-template
func (s *TemplatesSvcImpl) Get(identifier string) (*Template, error) {
	return s.GetWithContext(context.Background(), identifier)
}

// UpdateWithContext updates a template by ID or alias
// https://resend.com/docs/api-reference/templates/update-template
func (s *TemplatesSvcImpl) UpdateWithContext(ctx context.Context, identifier string, params *UpdateTemplateRequest) (*UpdateTemplateResponse, error) {
	path := "templates/" + identifier

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, params)
	if err != nil {
		return nil, ErrFailedToCreateTemplateUpdateRequest
	}

	// Build response recipient obj
	templateResponse := new(UpdateTemplateResponse)

	// Send Request
	_, err = s.client.Perform(req, templateResponse)

	if err != nil {
		return nil, err
	}

	return templateResponse, nil
}

// Update updates a template by ID or alias
// https://resend.com/docs/api-reference/templates/update-template
func (s *TemplatesSvcImpl) Update(identifier string, params *UpdateTemplateRequest) (*UpdateTemplateResponse, error) {
	return s.UpdateWithContext(context.Background(), identifier, params)
}

// PublishWithContext publishes a template by ID or alias
// https://resend.com/docs/api-reference/templates/publish-template
func (s *TemplatesSvcImpl) PublishWithContext(ctx context.Context, identifier string) (*PublishTemplateResponse, error) {
	path := "templates/" + identifier + "/publish"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, ErrFailedToCreateTemplatePublishRequest
	}

	// Build response recipient obj
	templateResponse := new(PublishTemplateResponse)

	// Send Request
	_, err = s.client.Perform(req, templateResponse)

	if err != nil {
		return nil, err
	}

	return templateResponse, nil
}

// Publish publishes a template by ID or alias
// https://resend.com/docs/api-reference/templates/publish-template
func (s *TemplatesSvcImpl) Publish(identifier string) (*PublishTemplateResponse, error) {
	return s.PublishWithContext(context.Background(), identifier)
}
