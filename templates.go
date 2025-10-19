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

// DuplicateTemplateResponse is the response from duplicating a template
type DuplicateTemplateResponse struct {
	Id     string `json:"id"`
	Object string `json:"object"`
}

// RemoveTemplateResponse is the response from removing a template
type RemoveTemplateResponse struct {
	Object  string `json:"object"`
	Id      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// TemplateListItem represents a template in a list response
type TemplateListItem struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Status      string  `json:"status"`
	PublishedAt *string `json:"published_at"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Alias       string  `json:"alias"`
}

// ListTemplatesResponse is the response from listing templates
type ListTemplatesResponse struct {
	Object  string              `json:"object"`
	Data    []*TemplateListItem `json:"data"`
	HasMore bool                `json:"has_more"`
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
	ListWithContext(ctx context.Context, options *ListOptions) (*ListTemplatesResponse, error)
	List(options *ListOptions) (*ListTemplatesResponse, error)
	UpdateWithContext(ctx context.Context, identifier string, params *UpdateTemplateRequest) (*UpdateTemplateResponse, error)
	Update(identifier string, params *UpdateTemplateRequest) (*UpdateTemplateResponse, error)
	PublishWithContext(ctx context.Context, identifier string) (*PublishTemplateResponse, error)
	Publish(identifier string) (*PublishTemplateResponse, error)
	DuplicateWithContext(ctx context.Context, identifier string) (*DuplicateTemplateResponse, error)
	Duplicate(identifier string) (*DuplicateTemplateResponse, error)
	RemoveWithContext(ctx context.Context, identifier string) (*RemoveTemplateResponse, error)
	Remove(identifier string) (*RemoveTemplateResponse, error)
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

// ListWithContext retrieves a list of templates with pagination options
// https://resend.com/docs/api-reference/templates/list-templates
func (s *TemplatesSvcImpl) ListWithContext(ctx context.Context, options *ListOptions) (*ListTemplatesResponse, error) {
	path := "templates" + buildPaginationQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, ErrFailedToCreateTemplateListRequest
	}

	// Build response recipient obj
	templateResponse := new(ListTemplatesResponse)

	// Send Request
	_, err = s.client.Perform(req, templateResponse)

	if err != nil {
		return nil, err
	}

	return templateResponse, nil
}

// List retrieves a list of templates with pagination options
// https://resend.com/docs/api-reference/templates/list-templates
func (s *TemplatesSvcImpl) List(options *ListOptions) (*ListTemplatesResponse, error) {
	return s.ListWithContext(context.Background(), options)
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

// DuplicateWithContext duplicates a template by ID or alias
// https://resend.com/docs/api-reference/templates/duplicate-template
func (s *TemplatesSvcImpl) DuplicateWithContext(ctx context.Context, identifier string) (*DuplicateTemplateResponse, error) {
	path := "templates/" + identifier + "/duplicate"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, ErrFailedToCreateTemplateDuplicateRequest
	}

	// Build response recipient obj
	templateResponse := new(DuplicateTemplateResponse)

	// Send Request
	_, err = s.client.Perform(req, templateResponse)

	if err != nil {
		return nil, err
	}

	return templateResponse, nil
}

// Duplicate duplicates a template by ID or alias
// https://resend.com/docs/api-reference/templates/duplicate-template
func (s *TemplatesSvcImpl) Duplicate(identifier string) (*DuplicateTemplateResponse, error) {
	return s.DuplicateWithContext(context.Background(), identifier)
}

// RemoveWithContext removes a template by ID or alias
// https://resend.com/docs/api-reference/templates/delete-template
func (s *TemplatesSvcImpl) RemoveWithContext(ctx context.Context, identifier string) (*RemoveTemplateResponse, error) {
	path := "templates/" + identifier

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, ErrFailedToCreateTemplateRemoveRequest
	}

	// Build response recipient obj
	templateResponse := new(RemoveTemplateResponse)

	// Send Request
	_, err = s.client.Perform(req, templateResponse)

	if err != nil {
		return nil, err
	}

	return templateResponse, nil
}

// Remove removes a template by ID or alias
// https://resend.com/docs/api-reference/templates/delete-template
func (s *TemplatesSvcImpl) Remove(identifier string) (*RemoveTemplateResponse, error) {
	return s.RemoveWithContext(context.Background(), identifier)
}
