package resend

import (
	"context"
	"errors"
	"net/http"
)

type ContactsSvc interface {
	CreateWithContext(ctx context.Context, params *CreateContactRequest) (CreateContactResponse, error)
	Create(params *CreateContactRequest) (CreateContactResponse, error)
	ListWithContext(ctx context.Context, audienceId string) (ListContactsResponse, error)
	List(audienceId string) (ListContactsResponse, error)
	GetWithContext(ctx context.Context, audienceId, contactId string) (Contact, error)
	Get(audienceId, contactId string) (Contact, error)
	RemoveWithContext(ctx context.Context, audienceId, contactId string) (RemoveContactResponse, error)
	Remove(audienceId, contactId string) (RemoveContactResponse, error)
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
	AudienceId   string `json:"audience_id"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	Unsubscribed bool   `json:"unsubscribed,omitempty"`
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
	Object string    `json:"object"`
	Data   []Contact `json:"data"`
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

// ListWithContext returns the list of all contacts in an audience
// https://resend.com/docs/api-reference/contacts/list-contacts
func (s *ContactsSvcImpl) ListWithContext(ctx context.Context, audienceId string) (ListContactsResponse, error) {
	path := "audiences/" + audienceId + "/contacts"

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

// List returns the list of all contacts in an audience
// https://resend.com/docs/api-reference/contacts/list-contacts
func (s *ContactsSvcImpl) List(audienceId string) (ListContactsResponse, error) {
	return s.ListWithContext(context.Background(), audienceId)
}

// RemoveWithContext same as Remove but with context
// https://resend.com/docs/api-reference/contacts/delete-contact
func (s *ContactsSvcImpl) RemoveWithContext(ctx context.Context, audienceId, contactId string) (RemoveContactResponse, error) {
	path := "audiences/" + audienceId + "/contacts/" + contactId

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
// @param [contactId] - the contact id or contact email
//
// https://resend.com/docs/api-reference/contacts/delete-contact
func (s *ContactsSvcImpl) Remove(audienceId, contactId string) (RemoveContactResponse, error) {
	return s.RemoveWithContext(context.Background(), audienceId, contactId)
}

// GetWithContext Retrieve a single contact.
// https://resend.com/docs/api-reference/contacts/get-contact
func (s *ContactsSvcImpl) GetWithContext(ctx context.Context, audienceId, contactId string) (Contact, error) {
	path := "audiences/" + audienceId + "/contacts/" + contactId

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
// https://resend.com/docs/api-reference/contacts/get-contact
func (s *ContactsSvcImpl) Get(audienceId, contactId string) (Contact, error) {
	return s.GetWithContext(context.Background(), audienceId, contactId)
}

// UpdateWithContext updates an existing Contact based on the given params
// https://resend.com/docs/api-reference/contacts/update-contact
func (s *ContactsSvcImpl) UpdateWithContext(ctx context.Context, params *UpdateContactRequest) (UpdateContactResponse, error) {
	if params.AudienceId == "" {
		return UpdateContactResponse{}, errors.New("[ERROR]: AudienceId is missing")
	}

	if params.Id == "" {
		return UpdateContactResponse{}, errors.New("[ERROR]: Id is missing")
	}

	path := "audiences/" + params.AudienceId + "/contacts/" + params.Id

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

// UpdateWithContext updates an existing Contact based on the given params
// https://resend.com/docs/api-reference/contacts/update-contact
func (s *ContactsSvcImpl) Update(params *UpdateContactRequest) (UpdateContactResponse, error) {
	return s.UpdateWithContext(context.Background(), params)
}
