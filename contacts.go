package resend

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type ContactsSvc interface {
	Create(params *CreateContactRequest) (CreateContactResponse, error)
	CreateWithContext(ctx context.Context, params *CreateContactRequest) (CreateContactResponse, error)
	Get(options *GetContactOptions) (Contact, error)
	GetWithContext(ctx context.Context, options *GetContactOptions) (Contact, error)
	List(options *ListContactsOptions) (ListContactsResponse, error)
	ListWithContext(ctx context.Context, options *ListContactsOptions) (ListContactsResponse, error)
	Update(params *UpdateContactRequest) (UpdateContactResponse, error)
	UpdateWithContext(ctx context.Context, params *UpdateContactRequest) (UpdateContactResponse, error)
	Remove(options *RemoveContactOptions) (RemoveContactResponse, error)
	RemoveWithContext(ctx context.Context, options *RemoveContactOptions) (RemoveContactResponse, error)
}

type ContactsSvcImpl struct {
	client     *Client
	Topics     ContactTopicsSvc
	Segments   ContactSegmentsSvc
	Properties ContactPropertiesSvc
}

// GetContactOptions contains parameters for retrieving a contact
type GetContactOptions struct {
	AudienceId string // Optional - omit for global contacts
	Id         string // Required - can be contact ID or email address
}

// ListContactsOptions contains parameters for listing contacts
type ListContactsOptions struct {
	AudienceId string  // Optional - omit for global contacts
	Limit      *int    // Optional - number of results to return
	After      *string // Optional - cursor for pagination
	Before     *string // Optional - cursor for pagination
}

// RemoveContactOptions contains parameters for removing a contact
type RemoveContactOptions struct {
	AudienceId string // Optional - omit for global contacts
	Id         string // Required - can be contact ID or email address
}

type CreateContactRequest struct {
	Email      string `json:"email"`
	AudienceId string `json:"audience_id,omitempty"` // Deprecated: Optional, use Segments API for contact organization
	Unsubscribed bool `json:"unsubscribed,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	// Properties are custom key-value pairs for global contacts (when audience_id is omitted).
	// NOTE: Currently, the Resend API only accepts string values for properties.
	// Non-string values (numbers, booleans, etc.) will be rejected by the API with a validation error.
	// Example: Properties: map[string]interface{}{"tier": "premium", "age": "30", "active": "true"}
	Properties map[string]interface{} `json:"properties,omitempty"`
}

type UpdateContactRequest struct {
	Id           string `json:"id"`
	Email        string `json:"email,omitempty"`
	AudienceId   string `json:"audience_id,omitempty"` // Deprecated: Optional, use Segments API for contact organization
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Unsubscribed bool   `json:"unsubscribed,omitempty"`
	// Properties are custom key-value pairs for global contacts (when audience_id is omitted).
	// NOTE: Currently, the Resend API only accepts string values for properties.
	// Non-string values (numbers, booleans, etc.) will be rejected by the API with a validation error.
	// Example: Properties: map[string]interface{}{"tier": "premium", "age": "30", "active": "true"}
	Properties      map[string]interface{} `json:"properties,omitempty"`
	unsubscribedSet bool                   `json:"-"`
}

// Temporary setter for the `unsubscribed` field. This is here
// as a backwards compatible way to set the unsubscribed field as false, since the
// default zero value for a bool is false with omitempty, and setting as false
// would omit the field from the JSON representation.
// Proper fix for this is coming in v3.
func (r *UpdateContactRequest) SetUnsubscribed(val bool) {
	r.Unsubscribed = val
	r.unsubscribedSet = true
}

type UpdateContactResponse struct {
	Data  Contact  `json:"data"`
	Error struct{} `json:"error"` // Fix this
}

type CreateContactResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

type RemoveContactResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

type ListContactsResponse struct {
	Object  string    `json:"object"`
	Data    []Contact `json:"data"`
	HasMore bool      `json:"has_more"`
}

type Contact struct {
	Id           string                 `json:"id"`
	Email        string                 `json:"email"`
	Object       string                 `json:"object"`
	FirstName    string                 `json:"first_name"`
	LastName     string                 `json:"last_name"`
	CreatedAt    string                 `json:"created_at"`
	Unsubscribed bool                   `json:"unsubscribed"`
	Properties   map[string]interface{} `json:"properties,omitempty"` // Custom properties for global contacts (currently API only returns string values)
}

// Create creates a new Contact based on the given params
// Supports both global contacts (without audience_id) and audience-specific contacts.
// Global contacts support custom properties while audience-specific contacts do not.
// https://resend.com/docs/api-reference/contacts/create-contact
func (s *ContactsSvcImpl) Create(params *CreateContactRequest) (CreateContactResponse, error) {
	return s.CreateWithContext(context.Background(), params)
}

// CreateWithContext creates a new Contact based on the given params with context
// Supports both global contacts (without audience_id) and audience-specific contacts.
// Global contacts support custom properties while audience-specific contacts do not.
// https://resend.com/docs/api-reference/contacts/create-contact
func (s *ContactsSvcImpl) CreateWithContext(ctx context.Context, params *CreateContactRequest) (CreateContactResponse, error) {
	var path string
	if params.AudienceId != "" {
		// Audience-specific contact (legacy path, no properties support)
		path = "audiences/" + params.AudienceId + "/contacts"
	} else {
		// Global contact (supports properties)
		path = "contacts"
	}

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return CreateContactResponse{}, errors.New("[ERROR]: Failed to create Contacts.Create request")
	}

	// Build response recipient obj
	contactsResp := new(CreateContactResponse)

	// Send Request
	_, err = s.client.Perform(req, contactsResp)

	if err != nil {
		return CreateContactResponse{}, err
	}

	return *contactsResp, nil
}

// List returns the list of all contacts
// If options.AudienceId is empty, lists global contacts. Otherwise lists audience-specific contacts.
// https://resend.com/docs/api-reference/contacts/list-contacts
func (s *ContactsSvcImpl) List(options *ListContactsOptions) (ListContactsResponse, error) {
	return s.ListWithContext(context.Background(), options)
}

// ListWithContext returns the list of all contacts with context
// If options.AudienceId is empty, lists global contacts. Otherwise lists audience-specific contacts.
// https://resend.com/docs/api-reference/contacts/list-contacts
func (s *ContactsSvcImpl) ListWithContext(ctx context.Context, options *ListContactsOptions) (ListContactsResponse, error) {
	if options == nil {
		options = &ListContactsOptions{}
	}

	var path string
	if options.AudienceId != "" {
		// Audience-specific contacts (legacy)
		path = "audiences/" + options.AudienceId + "/contacts"
	} else {
		// Global contacts
		path = "contacts"
	}

	// Build pagination query
	listOpts := &ListOptions{
		Limit:  options.Limit,
		After:  options.After,
		Before: options.Before,
	}
	path += buildPaginationQuery(listOpts)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListContactsResponse{}, errors.New("[ERROR]: Failed to create Contacts.List request")
	}

	contacts := new(ListContactsResponse)

	// Send Request
	_, err = s.client.Perform(req, contacts)

	if err != nil {
		return ListContactsResponse{}, err
	}

	return *contacts, nil
}

// Remove removes a contact
// If options.AudienceId is empty, removes a global contact. Otherwise removes an audience-specific contact.
// The options.Id field can be either a contact ID or email address.
// https://resend.com/docs/api-reference/contacts/delete-contact
func (s *ContactsSvcImpl) Remove(options *RemoveContactOptions) (RemoveContactResponse, error) {
	return s.RemoveWithContext(context.Background(), options)
}

// RemoveWithContext removes a contact with context
// If options.AudienceId is empty, removes a global contact. Otherwise removes an audience-specific contact.
// The options.Id field can be either a contact ID or email address.
// https://resend.com/docs/api-reference/contacts/delete-contact
func (s *ContactsSvcImpl) RemoveWithContext(ctx context.Context, options *RemoveContactOptions) (RemoveContactResponse, error) {
	if options == nil || options.Id == "" {
		return RemoveContactResponse{}, errors.New("[ERROR]: Id is required")
	}

	var path string
	if options.AudienceId != "" {
		// Audience-specific contact (legacy)
		path = "audiences/" + options.AudienceId + "/contacts/" + options.Id
	} else {
		// Global contact
		path = "contacts/" + options.Id
	}

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return RemoveContactResponse{}, errors.New("[ERROR]: Failed to create Contact.Remove request")
	}

	resp := new(RemoveContactResponse)

	// Send Request
	_, err = s.client.Perform(req, resp)

	if err != nil {
		return RemoveContactResponse{}, err
	}

	return *resp, nil
}

// Get retrieves a single contact.
// This method can be used to retrieve a contact by either its ID or email address.
// If options.AudienceId is empty, retrieves a global contact. Otherwise retrieves an audience-specific contact.
// https://resend.com/docs/api-reference/contacts/get-contact
func (s *ContactsSvcImpl) Get(options *GetContactOptions) (Contact, error) {
	return s.GetWithContext(context.Background(), options)
}

// GetWithContext retrieves a single contact with context.
// This method can be used to retrieve a contact by either its ID or email address.
// If options.AudienceId is empty, retrieves a global contact. Otherwise retrieves an audience-specific contact.
// https://resend.com/docs/api-reference/contacts/get-contact
func (s *ContactsSvcImpl) GetWithContext(ctx context.Context, options *GetContactOptions) (Contact, error) {
	if options == nil || options.Id == "" {
		return Contact{}, errors.New("[ERROR]: Id is required")
	}

	var path string
	if options.AudienceId != "" {
		// Audience-specific contact (legacy)
		path = "audiences/" + options.AudienceId + "/contacts/" + options.Id
	} else {
		// Global contact
		path = "contacts/" + options.Id
	}

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return Contact{}, errors.New("[ERROR]: Failed to create Contact.Get request")
	}

	contact := new(Contact)

	// Send Request
	_, err = s.client.Perform(req, contact)

	if err != nil {
		return Contact{}, err
	}

	return *contact, nil
}

// Update updates an existing Contact based on the given params
// Supports both global contacts (without audience_id) and audience-specific contacts.
// https://resend.com/docs/api-reference/contacts/update-contact
func (s *ContactsSvcImpl) Update(params *UpdateContactRequest) (UpdateContactResponse, error) {
	return s.UpdateWithContext(context.Background(), params)
}

// UpdateWithContext updates an existing Contact based on the given params with context
// Supports both global contacts (without audience_id) and audience-specific contacts.
// https://resend.com/docs/api-reference/contacts/update-contact
func (s *ContactsSvcImpl) UpdateWithContext(ctx context.Context, params *UpdateContactRequest) (UpdateContactResponse, error) {
	if params.Id == "" && params.Email == "" {
		return UpdateContactResponse{}, &MissingRequiredFieldsError{message: "[ERROR]: Missing `id` or `email` field."}
	}

	var val string
	if params.Id != "" {
		val = params.Id
	} else {
		val = params.Email
	}

	var path string
	if params.AudienceId != "" {
		// Audience-specific contact (legacy path)
		path = "audiences/" + params.AudienceId + "/contacts/" + val
	} else {
		// Global contact
		path = "contacts/" + val
	}

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, params)
	if err != nil {
		return UpdateContactResponse{}, errors.New("[ERROR]: Failed to create Contacts.Update request")
	}

	// Build response recipient obj
	contactsResp := new(UpdateContactResponse)

	// Send Request
	_, err = s.client.Perform(req, contactsResp)

	if err != nil {
		return UpdateContactResponse{}, err
	}

	return *contactsResp, nil
}

// Patches the JSON representation of the UpdateContactRequest
// in order to properly omit the `unsubscribed` field if it is not set
// This is marked to be fixed properly in v3, since the actual fix would
// require a breaking change
func (r UpdateContactRequest) MarshalJSON() ([]byte, error) {
	type Alias UpdateContactRequest
	aux := make(map[string]interface{})

	aux["id"] = r.Id
	if r.Email != "" {
		aux["email"] = r.Email
	}
	// Only include audience_id if provided (supports global contacts)
	if r.AudienceId != "" {
		aux["audience_id"] = r.AudienceId
	}
	if r.FirstName != "" {
		aux["first_name"] = r.FirstName
	}
	if r.LastName != "" {
		aux["last_name"] = r.LastName
	}
	if r.unsubscribedSet {
		aux["unsubscribed"] = r.Unsubscribed
	}
	// Include properties if provided (for global contacts)
	if r.Properties != nil && len(r.Properties) > 0 {
		aux["properties"] = r.Properties
	}

	return json.Marshal(aux)
}

