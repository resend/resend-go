package resend

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// InboxFolder is the folder an inbox thread is in.
type InboxFolder = string

const (
	InboxFolderInbox   InboxFolder = "inbox"
	InboxFolderArchive InboxFolder = "archive"
	InboxFolderSpam    InboxFolder = "spam"
	// InboxFolderSent is assigned automatically and cannot be set on a thread update.
	InboxFolderSent  InboxFolder = "sent"
	InboxFolderTrash InboxFolder = "trash"
)

// Inbox is a single inbox as returned by Get.
type Inbox struct {
	Object string `json:"object"`
	Id     string `json:"id"`
	// Name is the internal name for the inbox. Recipients do not see it. Falls back to the inbox
	// email address.
	Name         string `json:"name"`
	EmailAddress string `json:"email_address"`
	DomainId     string `json:"domain_id"`
	// ReceivingAddress is the address to forward mail to when forwarding is enabled. nil otherwise.
	ReceivingAddress *string `json:"receiving_address"`
	// FromName is the name recipients see when mail is sent from this inbox. nil when not set.
	FromName *string `json:"from_name"`
	// Unread is the number of unread threads in the Inbox folder.
	Unread int `json:"unread"`
	// Drafts is the number of unsent drafts.
	Drafts int `json:"drafts"`
	// LastReceived is the timestamp of the newest received email in the Inbox folder. nil when the
	// folder has no received email.
	LastReceived *string `json:"last_received"`
	CreatedAt    string  `json:"created_at"`
}

// InboxListItem is an inbox as returned by List. It carries fewer fields than Inbox.
type InboxListItem struct {
	Id           string  `json:"id"`
	Name         *string `json:"name"`
	EmailAddress string  `json:"email_address"`
	FromName     *string `json:"from_name"`
	Unread       int     `json:"unread"`
	LastReceived *string `json:"last_received"`
}

// CreateInboxRequest contains params for creating an inbox.
type CreateInboxRequest struct {
	// EmailAddress is the inbox address on one of your domains. The API trims and lowercases it.
	EmailAddress string `json:"email_address"`
	// Name is the internal name for the inbox. Recipients do not see it. Max 64 characters.
	// Defaults to EmailAddress.
	Name string `json:"name,omitempty"`
	// FromName is the name recipients see when mail is sent from this inbox. A plain name, not a
	// `Name <email>` address. Max 78 characters.
	FromName string `json:"from_name,omitempty"`
	// Forwarding makes Resend provision a receiving address, so you can receive mail without adding
	// an MX record. The address is returned in ReceivingAddress.
	Forwarding bool `json:"forwarding,omitempty"`
}

type CreateInboxResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

type ListInboxesResponse struct {
	Object  string          `json:"object"`
	HasMore bool            `json:"has_more"`
	Data    []InboxListItem `json:"data"`
}

// UpdateInboxRequest contains params for updating an inbox. Provide at least one of Name or
// FromName.
type UpdateInboxRequest struct {
	// Name is the new internal name for the inbox. Between 1 and 64 characters after trimming.
	Name string `json:"name,omitempty"`
	// FromName is the name recipients see when mail is sent from this inbox. Max 78 characters. A
	// pointer to an empty string removes it.
	FromName *string `json:"from_name,omitempty"`
}

type UpdateInboxResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

type RemoveInboxResponse struct {
	Object  string `json:"object"`
	Id      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

type InboxesSvc interface {
	Create(params *CreateInboxRequest) (CreateInboxResponse, error)
	CreateWithContext(ctx context.Context, params *CreateInboxRequest) (CreateInboxResponse, error)
	List(options *ListOptions) (ListInboxesResponse, error)
	ListWithContext(ctx context.Context, options *ListOptions) (ListInboxesResponse, error)
	Get(idOrEmail string) (Inbox, error)
	GetWithContext(ctx context.Context, idOrEmail string) (Inbox, error)
	Update(inboxId string, params *UpdateInboxRequest) (UpdateInboxResponse, error)
	UpdateWithContext(ctx context.Context, inboxId string, params *UpdateInboxRequest) (UpdateInboxResponse, error)
	Remove(inboxId string) (RemoveInboxResponse, error)
	RemoveWithContext(ctx context.Context, inboxId string) (RemoveInboxResponse, error)
}

type InboxesSvcImpl struct {
	client  *Client
	Agent   InboxAgentSvc
	Labels  InboxLabelsSvc
	Threads *InboxThreadsSvcImpl
	Drafts  InboxDraftsSvc
}

func inboxPath(inboxId string, segments ...string) string {
	path := "inboxes/" + url.PathEscape(inboxId)
	for _, segment := range segments {
		path += "/" + url.PathEscape(segment)
	}
	return path
}

// Create creates an inbox that sends and receives email at an address on one of your domains.
// A team that has reached the inbox limit of its plan gets a 403 validation_error.
// https://resend.com/docs/api-reference/inboxes/create-inbox
func (s *InboxesSvcImpl) Create(params *CreateInboxRequest) (CreateInboxResponse, error) {
	return s.CreateWithContext(context.Background(), params)
}

// CreateWithContext creates an inbox with context.
// https://resend.com/docs/api-reference/inboxes/create-inbox
func (s *InboxesSvcImpl) CreateWithContext(ctx context.Context, params *CreateInboxRequest) (CreateInboxResponse, error) {
	if params == nil || params.EmailAddress == "" {
		return CreateInboxResponse{}, errors.New("[ERROR]: EmailAddress is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, "inboxes", params)
	if err != nil {
		return CreateInboxResponse{}, ErrFailedToCreateInboxesCreateRequest
	}

	resp := new(CreateInboxResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return CreateInboxResponse{}, err
	}
	return *resp, nil
}

// List returns the inboxes of the team, newest first. The list items carry fewer fields than a
// single inbox. Retrieve an inbox to get all its fields.
// https://resend.com/docs/api-reference/inboxes/list-inboxes
func (s *InboxesSvcImpl) List(options *ListOptions) (ListInboxesResponse, error) {
	return s.ListWithContext(context.Background(), options)
}

// ListWithContext returns the inboxes of the team with context.
// https://resend.com/docs/api-reference/inboxes/list-inboxes
func (s *InboxesSvcImpl) ListWithContext(ctx context.Context, options *ListOptions) (ListInboxesResponse, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, "inboxes"+buildPaginationQuery(options), nil)
	if err != nil {
		return ListInboxesResponse{}, ErrFailedToCreateInboxesListRequest
	}

	resp := new(ListInboxesResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return ListInboxesResponse{}, err
	}
	return *resp, nil
}

// Get retrieves a single inbox. The idOrEmail param can be the inbox ID or the inbox email
// address. The email match is case-insensitive.
// https://resend.com/docs/api-reference/inboxes/get-inbox
func (s *InboxesSvcImpl) Get(idOrEmail string) (Inbox, error) {
	return s.GetWithContext(context.Background(), idOrEmail)
}

// GetWithContext retrieves a single inbox with context.
// https://resend.com/docs/api-reference/inboxes/get-inbox
func (s *InboxesSvcImpl) GetWithContext(ctx context.Context, idOrEmail string) (Inbox, error) {
	if idOrEmail == "" {
		return Inbox{}, errors.New("[ERROR]: Inbox ID or email is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, inboxPath(idOrEmail), nil)
	if err != nil {
		return Inbox{}, ErrFailedToCreateInboxesGetRequest
	}

	resp := new(Inbox)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return Inbox{}, err
	}
	return *resp, nil
}

// Update updates the internal name of the inbox or the name recipients see when mail is sent from
// it. At least one of Name or FromName is required. The inboxId param must be an ID, not an email.
// https://resend.com/docs/api-reference/inboxes/update-inbox
func (s *InboxesSvcImpl) Update(inboxId string, params *UpdateInboxRequest) (UpdateInboxResponse, error) {
	return s.UpdateWithContext(context.Background(), inboxId, params)
}

// UpdateWithContext updates an inbox with context.
// https://resend.com/docs/api-reference/inboxes/update-inbox
func (s *InboxesSvcImpl) UpdateWithContext(ctx context.Context, inboxId string, params *UpdateInboxRequest) (UpdateInboxResponse, error) {
	if inboxId == "" {
		return UpdateInboxResponse{}, errors.New("[ERROR]: Inbox ID is required")
	}
	if params == nil || (params.Name == "" && params.FromName == nil) {
		return UpdateInboxResponse{}, errors.New("[ERROR]: At least one of Name or FromName is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPatch, inboxPath(inboxId), params)
	if err != nil {
		return UpdateInboxResponse{}, ErrFailedToCreateInboxesUpdateRequest
	}

	resp := new(UpdateInboxResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return UpdateInboxResponse{}, err
	}
	return *resp, nil
}

// Remove removes the inbox. An inbox that does not exist or was already removed returns a 404
// not_found.
// https://resend.com/docs/api-reference/inboxes/delete-inbox
func (s *InboxesSvcImpl) Remove(inboxId string) (RemoveInboxResponse, error) {
	return s.RemoveWithContext(context.Background(), inboxId)
}

// RemoveWithContext removes the inbox with context.
// https://resend.com/docs/api-reference/inboxes/delete-inbox
func (s *InboxesSvcImpl) RemoveWithContext(ctx context.Context, inboxId string) (RemoveInboxResponse, error) {
	if inboxId == "" {
		return RemoveInboxResponse{}, errors.New("[ERROR]: Inbox ID is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, inboxPath(inboxId), nil)
	if err != nil {
		return RemoveInboxResponse{}, ErrFailedToCreateInboxesRemoveRequest
	}

	resp := new(RemoveInboxResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return RemoveInboxResponse{}, err
	}
	return *resp, nil
}
