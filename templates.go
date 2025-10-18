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

// TemplatesSvc handles operations for templates
type TemplatesSvc interface {
	CreateWithContext(ctx context.Context, params *CreateTemplateRequest) (*CreateTemplateResponse, error)
	Create(params *CreateTemplateRequest) (*CreateTemplateResponse, error)
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
