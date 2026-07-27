package resend

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type SuppressionOrigin = string

const (
	SuppressionOriginBounce    SuppressionOrigin = "bounce"
	SuppressionOriginComplaint SuppressionOrigin = "complaint"
	SuppressionOriginManual    SuppressionOrigin = "manual"
)

// SuppressionListEntry is a suppressed email address as returned by List. It has no Object field
// because the list endpoint does not send one; only Get does.
type SuppressionListEntry struct {
	Id     string            `json:"id"`
	Email  string            `json:"email"`
	Origin SuppressionOrigin `json:"origin"`
	// SourceId identifies the event that caused the suppression, such as the email that
	// bounced or complained. It is null for manual suppressions.
	SourceId  *string `json:"source_id"`
	CreatedAt string  `json:"created_at"`
}

// Suppression is a suppressed email address as returned by Get.
type Suppression struct {
	Object string            `json:"object"`
	Id     string            `json:"id"`
	Email  string            `json:"email"`
	Origin SuppressionOrigin `json:"origin"`
	// SourceId identifies the event that caused the suppression, such as the email that
	// bounced or complained. It is null for manual suppressions.
	SourceId  *string `json:"source_id"`
	CreatedAt string  `json:"created_at"`
}

// AddSuppressionRequest contains params for suppressing an email address.
type AddSuppressionRequest struct {
	Email string `json:"email"`
}

type AddSuppressionResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

// ListSuppressionsOptions contains parameters for listing suppressions.
type ListSuppressionsOptions struct {
	// Origin filters suppressions by origin: bounce, complaint, manual.
	Origin SuppressionOrigin
	Limit  *int
	After  *string
	Before *string
}

type ListSuppressionsResponse struct {
	Object  string                 `json:"object"`
	HasMore bool                   `json:"has_more"`
	Data    []SuppressionListEntry `json:"data"`
}

type RemoveSuppressionResponse struct {
	Object  string `json:"object"`
	Id      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// BatchAddSuppressionsRequest contains params for suppressing up to 100 email addresses at once.
type BatchAddSuppressionsRequest struct {
	Emails []string `json:"emails"`
}

type BatchAddSuppressionsResponse struct {
	Data []AddSuppressionResponse `json:"data"`
}

// BatchRemoveSuppressionsRequest contains params for removing up to 100 suppressions at once.
// Provide either Emails or Ids, but not both. Ids must be suppression IDs in UUID form. The unset
// field is omitted from the request body; the API rejects it being sent as null.
type BatchRemoveSuppressionsRequest struct {
	Emails []string `json:"emails,omitempty"`
	Ids    []string `json:"ids,omitempty"`
}

type BatchRemoveSuppressionsResponse struct {
	Data []RemoveSuppressionResponse `json:"data"`
}

type SuppressionsSvc interface {
	Add(params *AddSuppressionRequest) (AddSuppressionResponse, error)
	AddWithContext(ctx context.Context, params *AddSuppressionRequest) (AddSuppressionResponse, error)
	List(options *ListSuppressionsOptions) (ListSuppressionsResponse, error)
	ListWithContext(ctx context.Context, options *ListSuppressionsOptions) (ListSuppressionsResponse, error)
	Get(idOrEmail string) (Suppression, error)
	GetWithContext(ctx context.Context, idOrEmail string) (Suppression, error)
	Remove(idOrEmail string) (RemoveSuppressionResponse, error)
	RemoveWithContext(ctx context.Context, idOrEmail string) (RemoveSuppressionResponse, error)
}

type SuppressionsSvcImpl struct {
	client *Client
	Batch  SuppressionsBatchSvc
}

// Add suppresses an email address. The API lowercases and trims the address, and the call is an
// upsert: re-adding an already-suppressed address succeeds and returns the existing ID.
// https://resend.com/docs/api-reference/suppressions/add-suppression
func (s *SuppressionsSvcImpl) Add(params *AddSuppressionRequest) (AddSuppressionResponse, error) {
	return s.AddWithContext(context.Background(), params)
}

// AddWithContext suppresses an email address with context. The API lowercases and trims the
// address, and the call is an upsert: re-adding an already-suppressed address succeeds and returns
// the existing ID.
// https://resend.com/docs/api-reference/suppressions/add-suppression
func (s *SuppressionsSvcImpl) AddWithContext(ctx context.Context, params *AddSuppressionRequest) (AddSuppressionResponse, error) {
	if params == nil || params.Email == "" {
		return AddSuppressionResponse{}, errors.New("[ERROR]: Email is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, "suppressions", params)
	if err != nil {
		return AddSuppressionResponse{}, errors.New("[ERROR]: Failed to create Suppressions.Add request")
	}

	resp := new(AddSuppressionResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return AddSuppressionResponse{}, err
	}
	return *resp, nil
}

// List retrieves a list of suppressions.
// https://resend.com/docs/api-reference/suppressions/list-suppressions
func (s *SuppressionsSvcImpl) List(options *ListSuppressionsOptions) (ListSuppressionsResponse, error) {
	return s.ListWithContext(context.Background(), options)
}

// ListWithContext retrieves a list of suppressions with context.
// https://resend.com/docs/api-reference/suppressions/list-suppressions
func (s *SuppressionsSvcImpl) ListWithContext(ctx context.Context, options *ListSuppressionsOptions) (ListSuppressionsResponse, error) {
	if options == nil {
		options = &ListSuppressionsOptions{}
	}

	query := make(url.Values)
	if options.Origin != "" {
		query.Set("origin", options.Origin)
	}
	if options.Limit != nil {
		query.Set("limit", fmt.Sprintf("%d", *options.Limit))
	}
	if options.After != nil {
		query.Set("after", *options.After)
	}
	if options.Before != nil {
		query.Set("before", *options.Before)
	}

	path := "suppressions"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListSuppressionsResponse{}, errors.New("[ERROR]: Failed to create Suppressions.List request")
	}

	resp := new(ListSuppressionsResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return ListSuppressionsResponse{}, err
	}
	return *resp, nil
}

// Get retrieves a single suppression. The idOrEmail param can be either a suppression ID or an
// email address. An unknown identifier returns a 404 "Suppression not found" error.
// https://resend.com/docs/api-reference/suppressions/get-suppression
func (s *SuppressionsSvcImpl) Get(idOrEmail string) (Suppression, error) {
	return s.GetWithContext(context.Background(), idOrEmail)
}

// GetWithContext retrieves a single suppression with context. The idOrEmail param can be either a
// suppression ID or an email address. An unknown identifier returns a 404 "Suppression not found"
// error.
// https://resend.com/docs/api-reference/suppressions/get-suppression
func (s *SuppressionsSvcImpl) GetWithContext(ctx context.Context, idOrEmail string) (Suppression, error) {
	if idOrEmail == "" {
		return Suppression{}, errors.New("[ERROR]: Id or email is required")
	}

	path := "suppressions/" + url.PathEscape(idOrEmail)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return Suppression{}, errors.New("[ERROR]: Failed to create Suppressions.Get request")
	}

	resp := new(Suppression)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return Suppression{}, err
	}
	return *resp, nil
}

// Remove removes a single suppression. The idOrEmail param can be either a suppression ID or an
// email address. An unknown identifier returns a 404 "Suppression not found" error.
// https://resend.com/docs/api-reference/suppressions/remove-suppression
func (s *SuppressionsSvcImpl) Remove(idOrEmail string) (RemoveSuppressionResponse, error) {
	return s.RemoveWithContext(context.Background(), idOrEmail)
}

// RemoveWithContext removes a single suppression with context. The idOrEmail param can be either a
// suppression ID or an email address. An unknown identifier returns a 404 "Suppression not found"
// error.
// https://resend.com/docs/api-reference/suppressions/remove-suppression
func (s *SuppressionsSvcImpl) RemoveWithContext(ctx context.Context, idOrEmail string) (RemoveSuppressionResponse, error) {
	if idOrEmail == "" {
		return RemoveSuppressionResponse{}, errors.New("[ERROR]: Id or email is required")
	}

	path := "suppressions/" + url.PathEscape(idOrEmail)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return RemoveSuppressionResponse{}, errors.New("[ERROR]: Failed to create Suppressions.Remove request")
	}

	resp := new(RemoveSuppressionResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return RemoveSuppressionResponse{}, err
	}
	return *resp, nil
}

type SuppressionsBatchSvc interface {
	Add(params *BatchAddSuppressionsRequest) (BatchAddSuppressionsResponse, error)
	AddWithContext(ctx context.Context, params *BatchAddSuppressionsRequest) (BatchAddSuppressionsResponse, error)
	Remove(params *BatchRemoveSuppressionsRequest) (BatchRemoveSuppressionsResponse, error)
	RemoveWithContext(ctx context.Context, params *BatchRemoveSuppressionsRequest) (BatchRemoveSuppressionsResponse, error)
}

type SuppressionsBatchSvcImpl struct {
	client *Client
}

// Add suppresses up to 100 email addresses at once. The API lowercases, trims and dedupes the
// addresses, and upserts them, so already-suppressed addresses succeed and return their existing
// IDs. The returned data can therefore be shorter than Emails.
// https://resend.com/docs/api-reference/suppressions/add-suppressions
func (s *SuppressionsBatchSvcImpl) Add(params *BatchAddSuppressionsRequest) (BatchAddSuppressionsResponse, error) {
	return s.AddWithContext(context.Background(), params)
}

// AddWithContext suppresses up to 100 email addresses at once with context. The API lowercases,
// trims and dedupes the addresses, and upserts them, so already-suppressed addresses succeed and
// return their existing IDs. The returned data can therefore be shorter than Emails.
// https://resend.com/docs/api-reference/suppressions/add-suppressions
func (s *SuppressionsBatchSvcImpl) AddWithContext(ctx context.Context, params *BatchAddSuppressionsRequest) (BatchAddSuppressionsResponse, error) {
	if params == nil || len(params.Emails) == 0 {
		return BatchAddSuppressionsResponse{}, errors.New("[ERROR]: Emails is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, "suppressions/batch/add", params)
	if err != nil {
		return BatchAddSuppressionsResponse{}, errors.New("[ERROR]: Failed to create Suppressions.Batch.Add request")
	}

	resp := new(BatchAddSuppressionsResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return BatchAddSuppressionsResponse{}, err
	}
	return *resp, nil
}

// Remove removes up to 100 suppressions at once, by email address or by ID. Unlike the single
// Remove, identifiers that are not suppressed are not an error: the returned data only covers the
// entries that were actually deleted.
// https://resend.com/docs/api-reference/suppressions/remove-suppressions
func (s *SuppressionsBatchSvcImpl) Remove(params *BatchRemoveSuppressionsRequest) (BatchRemoveSuppressionsResponse, error) {
	return s.RemoveWithContext(context.Background(), params)
}

// RemoveWithContext removes up to 100 suppressions at once with context, by email address or by ID.
// Unlike the single RemoveWithContext, identifiers that are not suppressed are not an error: the
// returned data only covers the entries that were actually deleted.
// https://resend.com/docs/api-reference/suppressions/remove-suppressions
func (s *SuppressionsBatchSvcImpl) RemoveWithContext(ctx context.Context, params *BatchRemoveSuppressionsRequest) (BatchRemoveSuppressionsResponse, error) {
	if params == nil || (len(params.Emails) == 0 && len(params.Ids) == 0) {
		return BatchRemoveSuppressionsResponse{}, errors.New("[ERROR]: Either Emails or Ids is required")
	}
	if len(params.Emails) > 0 && len(params.Ids) > 0 {
		return BatchRemoveSuppressionsResponse{}, errors.New("[ERROR]: Provide either `emails` or `ids`, but not both.")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, "suppressions/batch/remove", params)
	if err != nil {
		return BatchRemoveSuppressionsResponse{}, errors.New("[ERROR]: Failed to create Suppressions.Batch.Remove request")
	}

	resp := new(BatchRemoveSuppressionsResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return BatchRemoveSuppressionsResponse{}, err
	}
	return *resp, nil
}
