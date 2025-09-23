package resend

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type ContactsSvc interface {
	CreateWithContext(ctx context.Context, params *CreateContactRequest) (CreateContactResponse, error)
	Create(params *CreateContactRequest) (CreateContactResponse, error)
	ListWithOptions(ctx context.Context, audienceId string, options *ListOptions) (ListContactsResponse, error)
	ListWithContext(ctx context.Context, audienceId string) (ListContactsResponse, error)
	List(audienceId string) (ListContactsResponse, error)
	GetWithContext(ctx context.Context, audienceId, id string) (Contact, error)
	Get(audienceId, id string) (Contact, error)
	RemoveWithContext(ctx context.Context, audienceId, id string) (RemoveContactResponse, error)
	Remove(audienceId, id string) (RemoveContactResponse, error)
	UpdateWithContext(ctx context.Context, params *UpdateContactRequest) (UpdateContactResponse, error)
	Update(params *UpdateContactRequest) (UpdateContactResponse, error)
}

type ContactsSvcImpl struct {
	client *Client
}

type CreateContactRequest struct {
	Email        string `json:"email"`
	AudienceId   string `json:"audience_id"`
	Unsubscribed bool   `json:"unsubscribed,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
}

type UpdateContactRequest struct {
	Id           string `json:"id"`
	Email        string `json:"email,omitempty"`
	AudienceId   string `json:"audience_id"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Unsubscribed bool   `json:"unsubscribed,omitempty"`

	unsubscribedSet bool `json:"-"`
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
	Id           string `json:"id"`
	Email        string `json:"email"`
	Object       string `json:"object"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	CreatedAt    string `json:"created_at"`
	Unsubscribed bool   `json:"unsubscribed"`
}

// CreateWithContext creates a new Contact based on the given params
// https://resend.com/docs/api-reference/contacts/create-contact
func (s *ContactsSvcImpl) CreateWithContext(ctx context.Context, params *CreateContactRequest) (CreateContactResponse, error) {
	if params.AudienceId == "" {
		return CreateContactResponse{}, errors.New("[ERROR]: AudienceId is missing")
	}

	path := "audiences/" + params.AudienceId + "/contacts"

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

// Create creates a new Contact entry based on the given params
// https://resend.com/docs/api-reference/contacts/create-contact
func (s *ContactsSvcImpl) Create(params *CreateContactRequest) (CreateContactResponse, error) {
	return s.CreateWithContext(context.Background(), params)
}

// ListWithOptions returns the list of all contacts in an audience with pagination options
// https://resend.com/docs/api-reference/contacts/list-contacts
func (s *ContactsSvcImpl) ListWithOptions(ctx context.Context, audienceId string, options *ListOptions) (ListContactsResponse, error) {
	path := "audiences/" + audienceId + "/contacts" + buildPaginationQuery(options)

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

// ListWithContext returns the list of all contacts in an audience
// https://resend.com/docs/api-reference/contacts/list-contacts
func (s *ContactsSvcImpl) ListWithContext(ctx context.Context, audienceId string) (ListContactsResponse, error) {
	return s.ListWithOptions(ctx, audienceId, nil)
}

// List returns the list of all contacts in an audience
// https://resend.com/docs/api-reference/contacts/list-contacts
func (s *ContactsSvcImpl) List(audienceId string) (ListContactsResponse, error) {
	return s.ListWithContext(context.Background(), audienceId)
}

// RemoveWithContext same as Remove but with context
// https://resend.com/docs/api-reference/contacts/delete-contact
func (s *ContactsSvcImpl) RemoveWithContext(ctx context.Context, audienceId, id string) (RemoveContactResponse, error) {
	path := "audiences/" + audienceId + "/contacts/" + id

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

// Remove removes a given contact entry by id or email
//
// @param [id] - can be either a contact id or email
//
// https://resend.com/docs/api-reference/contacts/delete-contact
func (s *ContactsSvcImpl) Remove(audienceId, id string) (RemoveContactResponse, error) {
	return s.RemoveWithContext(context.Background(), audienceId, id)
}

// GetWithContext Retrieve a single contact.
// This method can be used to retrieve a contact by either its ID or email address.
//
// @param [id] - can be either a contact id or email
//
// https://resend.com/docs/api-reference/contacts/get-contact
func (s *ContactsSvcImpl) GetWithContext(ctx context.Context, audienceId, id string) (Contact, error) {
	path := "audiences/" + audienceId + "/contacts/" + id

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

// Get Retrieve a single contact.
// This method can be used to retrieve a contact by either its ID or email address.
//
// @param [id] - can be either a contact id or email
// https://resend.com/docs/api-reference/contacts/get-contact
func (s *ContactsSvcImpl) Get(audienceId, id string) (Contact, error) {
	return s.GetWithContext(context.Background(), audienceId, id)
}

// UpdateWithContext updates an existing Contact based on the given params
// https://resend.com/docs/api-reference/contacts/update-contact
func (s *ContactsSvcImpl) UpdateWithContext(ctx context.Context, params *UpdateContactRequest) (UpdateContactResponse, error) {
	if params.AudienceId == "" {
		return UpdateContactResponse{}, errors.New("[ERROR]: AudienceId is missing")
	}

	if params.Id == "" && params.Email == "" {
		return UpdateContactResponse{}, &MissingRequiredFieldsError{message: "[ERROR]: Missing `id` or `email` field."}
	}

	var val string
	if params.Id != "" {
		val = params.Id
	} else {
		val = params.Email
	}

	path := "audiences/" + params.AudienceId + "/contacts/" + val

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
	aux["audience_id"] = r.AudienceId
	if r.FirstName != "" {
		aux["first_name"] = r.FirstName
	}
	if r.LastName != "" {
		aux["last_name"] = r.LastName
	}
	if r.unsubscribedSet {
		aux["unsubscribed"] = r.Unsubscribed
	}

	return json.Marshal(aux)
}

// UpdateWithContext updates an existing Contact based on the given params
// https://resend.com/docs/api-reference/contacts/update-contact
func (s *ContactsSvcImpl) Update(params *UpdateContactRequest) (UpdateContactResponse, error) {
	return s.UpdateWithContext(context.Background(), params)
}
