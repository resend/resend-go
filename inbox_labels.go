package resend

import (
	"context"
	"errors"
	"net/http"
)

// InboxLabel is a named, colored tag that belongs to one inbox.
type InboxLabel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	// Color is a `#RRGGBB` hex color in uppercase.
	Color     string `json:"color"`
	CreatedAt string `json:"created_at"`
}

type ListInboxLabelsResponse struct {
	Object string `json:"object"`
	// Data holds the labels on the inbox, oldest first.
	Data []InboxLabel `json:"data"`
}

// CreateInboxLabelRequest contains params for creating a label on an inbox.
type CreateInboxLabelRequest struct {
	// Name of the label. Max 64 characters. The API trims it. It must not match another label on
	// the inbox, compared without case.
	Name string `json:"name"`
	// Color of the label, as a `#RRGGBB` hex color. It is stored in uppercase. When omitted, a
	// random color from the dashboard palette is used.
	Color string `json:"color,omitempty"`
}

type CreateInboxLabelResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

// UpdateInboxLabelRequest contains params for renaming or recoloring a label. Omitted fields keep
// their current value.
type UpdateInboxLabelRequest struct {
	// Name is a new name for the label. Max 64 characters. The API trims it. It must not match
	// another label on the inbox, compared without case.
	Name string `json:"name,omitempty"`
	// Color is a new color for the label, as a `#RRGGBB` hex color. It is stored in uppercase.
	Color string `json:"color,omitempty"`
}

type UpdateInboxLabelResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

type RemoveInboxLabelResponse struct {
	Object  string `json:"object"`
	Id      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

type InboxLabelsSvc interface {
	List(inboxId string) (ListInboxLabelsResponse, error)
	ListWithContext(ctx context.Context, inboxId string) (ListInboxLabelsResponse, error)
	Create(inboxId string, params *CreateInboxLabelRequest) (CreateInboxLabelResponse, error)
	CreateWithContext(ctx context.Context, inboxId string, params *CreateInboxLabelRequest) (CreateInboxLabelResponse, error)
	Update(inboxId, labelId string, params *UpdateInboxLabelRequest) (UpdateInboxLabelResponse, error)
	UpdateWithContext(ctx context.Context, inboxId, labelId string, params *UpdateInboxLabelRequest) (UpdateInboxLabelResponse, error)
	Remove(inboxId, labelId string) (RemoveInboxLabelResponse, error)
	RemoveWithContext(ctx context.Context, inboxId, labelId string) (RemoveInboxLabelResponse, error)
}

type InboxLabelsSvcImpl struct {
	client *Client
}

// List returns every label on the inbox, oldest first. The list is not paginated.
// https://resend.com/docs/api-reference/inboxes/list-labels
func (s *InboxLabelsSvcImpl) List(inboxId string) (ListInboxLabelsResponse, error) {
	return s.ListWithContext(context.Background(), inboxId)
}

// ListWithContext returns every label on the inbox with context.
// https://resend.com/docs/api-reference/inboxes/list-labels
func (s *InboxLabelsSvcImpl) ListWithContext(ctx context.Context, inboxId string) (ListInboxLabelsResponse, error) {
	if inboxId == "" {
		return ListInboxLabelsResponse{}, errors.New("[ERROR]: Inbox ID is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, inboxPath(inboxId, "labels"), nil)
	if err != nil {
		return ListInboxLabelsResponse{}, ErrFailedToCreateInboxesLabelsListRequest
	}

	resp := new(ListInboxLabelsResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return ListInboxLabelsResponse{}, err
	}
	return *resp, nil
}

// Create creates a label on the inbox. An inbox can have at most 100 labels. Label names are
// unique per inbox, compared without case. A duplicate name or a full inbox returns a 422
// validation_error.
// https://resend.com/docs/api-reference/inboxes/create-label
func (s *InboxLabelsSvcImpl) Create(inboxId string, params *CreateInboxLabelRequest) (CreateInboxLabelResponse, error) {
	return s.CreateWithContext(context.Background(), inboxId, params)
}

// CreateWithContext creates a label on the inbox with context.
// https://resend.com/docs/api-reference/inboxes/create-label
func (s *InboxLabelsSvcImpl) CreateWithContext(ctx context.Context, inboxId string, params *CreateInboxLabelRequest) (CreateInboxLabelResponse, error) {
	if inboxId == "" {
		return CreateInboxLabelResponse{}, errors.New("[ERROR]: Inbox ID is required")
	}
	if params == nil || params.Name == "" {
		return CreateInboxLabelResponse{}, errors.New("[ERROR]: Name is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, inboxPath(inboxId, "labels"), params)
	if err != nil {
		return CreateInboxLabelResponse{}, ErrFailedToCreateInboxesLabelsCreateRequest
	}

	resp := new(CreateInboxLabelResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return CreateInboxLabelResponse{}, err
	}
	return *resp, nil
}

// Update renames or recolors a label. Omitted fields keep their current value. A new name must not
// match another label on the inbox, compared without case, or the request returns a 422
// validation_error.
// https://resend.com/docs/api-reference/inboxes/update-label
func (s *InboxLabelsSvcImpl) Update(inboxId, labelId string, params *UpdateInboxLabelRequest) (UpdateInboxLabelResponse, error) {
	return s.UpdateWithContext(context.Background(), inboxId, labelId, params)
}

// UpdateWithContext renames or recolors a label with context.
// https://resend.com/docs/api-reference/inboxes/update-label
func (s *InboxLabelsSvcImpl) UpdateWithContext(ctx context.Context, inboxId, labelId string, params *UpdateInboxLabelRequest) (UpdateInboxLabelResponse, error) {
	if inboxId == "" {
		return UpdateInboxLabelResponse{}, errors.New("[ERROR]: Inbox ID is required")
	}
	if labelId == "" {
		return UpdateInboxLabelResponse{}, errors.New("[ERROR]: Label ID is required")
	}
	if params == nil {
		params = &UpdateInboxLabelRequest{}
	}

	req, err := s.client.NewRequest(ctx, http.MethodPatch, inboxPath(inboxId, "labels", labelId), params)
	if err != nil {
		return UpdateInboxLabelResponse{}, ErrFailedToCreateInboxesLabelsUpdateRequest
	}

	resp := new(UpdateInboxLabelResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return UpdateInboxLabelResponse{}, err
	}
	return *resp, nil
}

// Remove deletes the label and removes it from every thread that carries it.
// https://resend.com/docs/api-reference/inboxes/delete-label
func (s *InboxLabelsSvcImpl) Remove(inboxId, labelId string) (RemoveInboxLabelResponse, error) {
	return s.RemoveWithContext(context.Background(), inboxId, labelId)
}

// RemoveWithContext deletes the label with context.
// https://resend.com/docs/api-reference/inboxes/delete-label
func (s *InboxLabelsSvcImpl) RemoveWithContext(ctx context.Context, inboxId, labelId string) (RemoveInboxLabelResponse, error) {
	if inboxId == "" {
		return RemoveInboxLabelResponse{}, errors.New("[ERROR]: Inbox ID is required")
	}
	if labelId == "" {
		return RemoveInboxLabelResponse{}, errors.New("[ERROR]: Label ID is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, inboxPath(inboxId, "labels", labelId), nil)
	if err != nil {
		return RemoveInboxLabelResponse{}, ErrFailedToCreateInboxesLabelsRemoveRequest
	}

	resp := new(RemoveInboxLabelResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return RemoveInboxLabelResponse{}, err
	}
	return *resp, nil
}
