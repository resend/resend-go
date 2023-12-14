package resend

import (
	"context"
	"errors"
	"net/http"
)

type ContactsSvc interface {
	CreateWithContext(ctx context.Context, audienceId string, params *CreateContactRequest) (CreateContactResponse, error)
	Create(audienceId string, params *CreateContactRequest) (CreateContactResponse, error)
	ListWithContext(ctx context.Context, audienceId string) (ListContactsResponse, error)
	List(audienceId string) (ListContactsResponse, error)
	GetWithContext(ctx context.Context, audienceId, contactId string) (Contact, error)
	Get(audienceId, contactId string) (Contact, error)
	RemoveWithContext(ctx context.Context, audienceId, contactId string) (RemoveContactResponse, error)
	Remove(audienceId, contactId string) (RemoveContactResponse, error)
}

type ContactsSvcImpl struct {
	client *Client
}

type CreateContactRequest struct {
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Unsubscribed bool   `json:"unsubscribed"`
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
func (s *ContactsSvcImpl) CreateWithContext(ctx context.Context, audienceId string, params *CreateContactRequest) (CreateContactResponse, error) {
	path := "audiences/" + audienceId + "/contacts"

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
func (s *ContactsSvcImpl) Create(audienceId string, params *CreateContactRequest) (CreateContactResponse, error) {
	return s.CreateWithContext(context.Background(), audienceId, params)
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

// RemoveWithContext removes a given contact by id
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

// Remove removes a given contact entry by id
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
