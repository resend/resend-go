package resend

import (
	"context"
	"errors"
	"net/http"
)

type ContactPropertiesSvc interface {
	CreateWithContext(ctx context.Context, params *CreateContactPropertyRequest) (CreateContactPropertyResponse, error)
	Create(params *CreateContactPropertyRequest) (CreateContactPropertyResponse, error)
	ListWithOptions(ctx context.Context, options *ListOptions) (ListContactPropertiesResponse, error)
	ListWithContext(ctx context.Context) (ListContactPropertiesResponse, error)
	List() (ListContactPropertiesResponse, error)
	GetWithContext(ctx context.Context, id string) (ContactProperty, error)
	Get(id string) (ContactProperty, error)
	UpdateWithContext(ctx context.Context, params *UpdateContactPropertyRequest) (UpdateContactPropertyResponse, error)
	Update(params *UpdateContactPropertyRequest) (UpdateContactPropertyResponse, error)
	RemoveWithContext(ctx context.Context, id string) (RemoveContactPropertyResponse, error)
	Remove(id string) (RemoveContactPropertyResponse, error)
}

type ContactPropertiesSvcImpl struct {
	client *Client
}

type ContactProperty struct {
	Id            string      `json:"id"`
	Key           string      `json:"key"`
	Object        string      `json:"object"`
	CreatedAt     string      `json:"created_at"`
	Type          string      `json:"type"`
	FallbackValue interface{} `json:"fallback_value"`
}

type CreateContactPropertyRequest struct {
	Key           string      `json:"key"`
	Type          string      `json:"type"`
	FallbackValue interface{} `json:"fallback_value"`
}

type CreateContactPropertyResponse struct {
	Id     string `json:"id"`
	Object string `json:"object"`
}

type UpdateContactPropertyRequest struct {
	Id            string      `json:"-"`
	FallbackValue interface{} `json:"fallback_value"`
}

type UpdateContactPropertyResponse struct {
	Id     string `json:"id"`
	Object string `json:"object"`
}

type RemoveContactPropertyResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

type ListContactPropertiesResponse struct {
	Object  string            `json:"object"`
	Data    []ContactProperty `json:"data"`
	HasMore bool              `json:"has_more"`
}

// CreateWithContext creates a new contact property based on the given params
// https://resend.com/docs/api-reference/contact-properties/create-contact-property
func (s *ContactPropertiesSvcImpl) CreateWithContext(ctx context.Context, params *CreateContactPropertyRequest) (CreateContactPropertyResponse, error) {
	if params.Key == "" {
		return CreateContactPropertyResponse{}, errors.New("[ERROR]: Key is missing")
	}

	if params.Type == "" {
		return CreateContactPropertyResponse{}, errors.New("[ERROR]: Type is missing")
	}

	path := "contact-properties"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return CreateContactPropertyResponse{}, errors.New("[ERROR]: Failed to create ContactProperties.Create request")
	}

	// Build response object
	propertyResp := new(CreateContactPropertyResponse)

	// Send Request
	_, err = s.client.Perform(req, propertyResp)

	if err != nil {
		return CreateContactPropertyResponse{}, err
	}

	return *propertyResp, nil
}

// Create creates a new contact property based on the given params
// https://resend.com/docs/api-reference/contact-properties/create-contact-property
func (s *ContactPropertiesSvcImpl) Create(params *CreateContactPropertyRequest) (CreateContactPropertyResponse, error) {
	return s.CreateWithContext(context.Background(), params)
}

// ListWithOptions returns the list of all contact properties with pagination options
// https://resend.com/docs/api-reference/contact-properties/list-contact-properties
func (s *ContactPropertiesSvcImpl) ListWithOptions(ctx context.Context, options *ListOptions) (ListContactPropertiesResponse, error) {
	path := "contact-properties" + buildPaginationQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListContactPropertiesResponse{}, errors.New("[ERROR]: Failed to create ContactProperties.List request")
	}

	properties := new(ListContactPropertiesResponse)

	// Send Request
	_, err = s.client.Perform(req, properties)

	if err != nil {
		return ListContactPropertiesResponse{}, err
	}

	return *properties, nil
}

// ListWithContext returns the list of all contact properties
// https://resend.com/docs/api-reference/contact-properties/list-contact-properties
func (s *ContactPropertiesSvcImpl) ListWithContext(ctx context.Context) (ListContactPropertiesResponse, error) {
	return s.ListWithOptions(ctx, nil)
}

// List returns the list of all contact properties
// https://resend.com/docs/api-reference/contact-properties/list-contact-properties
func (s *ContactPropertiesSvcImpl) List() (ListContactPropertiesResponse, error) {
	return s.ListWithContext(context.Background())
}

// GetWithContext retrieves a single contact property by ID
// https://resend.com/docs/api-reference/contact-properties/get-contact-property
func (s *ContactPropertiesSvcImpl) GetWithContext(ctx context.Context, id string) (ContactProperty, error) {
	if id == "" {
		return ContactProperty{}, errors.New("[ERROR]: ID is missing")
	}

	path := "contact-properties/" + id

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ContactProperty{}, errors.New("[ERROR]: Failed to create ContactProperties.Get request")
	}

	property := new(ContactProperty)

	// Send Request
	_, err = s.client.Perform(req, property)

	if err != nil {
		return ContactProperty{}, err
	}

	return *property, nil
}

// Get retrieves a single contact property by ID
// https://resend.com/docs/api-reference/contact-properties/get-contact-property
func (s *ContactPropertiesSvcImpl) Get(id string) (ContactProperty, error) {
	return s.GetWithContext(context.Background(), id)
}

// UpdateWithContext updates an existing contact property based on the given params
// https://resend.com/docs/api-reference/contact-properties/update-contact-property
func (s *ContactPropertiesSvcImpl) UpdateWithContext(ctx context.Context, params *UpdateContactPropertyRequest) (UpdateContactPropertyResponse, error) {
	if params.Id == "" {
		return UpdateContactPropertyResponse{}, errors.New("[ERROR]: ID is missing")
	}

	path := "contact-properties/" + params.Id

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, params)
	if err != nil {
		return UpdateContactPropertyResponse{}, errors.New("[ERROR]: Failed to create ContactProperties.Update request")
	}

	// Build response object
	propertyResp := new(UpdateContactPropertyResponse)

	// Send Request
	_, err = s.client.Perform(req, propertyResp)

	if err != nil {
		return UpdateContactPropertyResponse{}, err
	}

	return *propertyResp, nil
}

// Update updates an existing contact property based on the given params
// https://resend.com/docs/api-reference/contact-properties/update-contact-property
func (s *ContactPropertiesSvcImpl) Update(params *UpdateContactPropertyRequest) (UpdateContactPropertyResponse, error) {
	return s.UpdateWithContext(context.Background(), params)
}

// RemoveWithContext removes a contact property by ID
// https://resend.com/docs/api-reference/contact-properties/delete-contact-property
func (s *ContactPropertiesSvcImpl) RemoveWithContext(ctx context.Context, id string) (RemoveContactPropertyResponse, error) {
	if id == "" {
		return RemoveContactPropertyResponse{}, errors.New("[ERROR]: ID is missing")
	}

	path := "contact-properties/" + id

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return RemoveContactPropertyResponse{}, errors.New("[ERROR]: Failed to create ContactProperties.Remove request")
	}

	resp := new(RemoveContactPropertyResponse)

	// Send Request
	_, err = s.client.Perform(req, resp)

	if err != nil {
		return RemoveContactPropertyResponse{}, err
	}

	return *resp, nil
}

// Remove removes a contact property by ID
// https://resend.com/docs/api-reference/contact-properties/delete-contact-property
func (s *ContactPropertiesSvcImpl) Remove(id string) (RemoveContactPropertyResponse, error) {
	return s.RemoveWithContext(context.Background(), id)
}
