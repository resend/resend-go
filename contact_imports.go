package resend

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
)

type ContactImportStatus string

const (
	ContactImportStatusQueued     ContactImportStatus = "queued"
	ContactImportStatusInProgress ContactImportStatus = "in_progress"
	ContactImportStatusCompleted  ContactImportStatus = "completed"
	ContactImportStatusFailed     ContactImportStatus = "failed"
)

type ContactImportCounts struct {
	Total   int `json:"total"`
	Created int `json:"created"`
	Updated int `json:"updated"`
	Skipped int `json:"skipped"`
	Failed  int `json:"failed"`
}

type ContactImport struct {
	Object    string               `json:"object"`
	Id        string               `json:"id"`
	Status    ContactImportStatus  `json:"status"`
	CreatedAt string               `json:"created_at"`
	Counts    *ContactImportCounts `json:"counts,omitempty"`
}

// ContactImportSegment represents a segment reference for a contact import.
type ContactImportSegment struct {
	Id string `json:"id"`
}

// CreateContactImportRequest contains params for creating a contact import.
// File is required; all other fields are optional.
type CreateContactImportRequest struct {
	// CSV file content (required). Maximum size is 50MB.
	File []byte
	// Filename to use in the multipart upload. Defaults to "import.csv".
	Filename string
	// ColumnMap maps contact fields to CSV column names.
	// Will be JSON-encoded before sending.
	// Example: map[string]any{"email": "Email", "first_name": "First Name"}
	ColumnMap map[string]any
	// OnConflict strategy: "upsert" or "skip" (default "skip").
	OnConflict string
	// Segments is a list of segment objects to add imported contacts to.
	// Example: []ContactImportSegment{{Id: "60a2ac5e-0774-456e-817d-ebf40f6dba31"}}
	Segments []ContactImportSegment
	// Topics is a list of topic subscriptions to apply to imported contacts.
	// Will be JSON-encoded as [{"id": "...", "subscription": "opt_in|opt_out"}] before sending.
	Topics []TopicSubscriptionUpdate
}

type CreateContactImportResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

type ListContactImportsOptions struct {
	// Status filters imports by status: queued, in_progress, completed, failed.
	Status string
	Limit  *int
	After  *string
	Before *string
}

type ListContactImportsResponse struct {
	Object  string          `json:"object"`
	HasMore bool            `json:"has_more"`
	Data    []ContactImport `json:"data"`
}

type ContactImportsSvc interface {
	Create(params *CreateContactImportRequest) (CreateContactImportResponse, error)
	CreateWithContext(ctx context.Context, params *CreateContactImportRequest) (CreateContactImportResponse, error)
	Get(id string) (ContactImport, error)
	GetWithContext(ctx context.Context, id string) (ContactImport, error)
	List(options *ListContactImportsOptions) (ListContactImportsResponse, error)
	ListWithContext(ctx context.Context, options *ListContactImportsOptions) (ListContactImportsResponse, error)
}

type ContactImportsSvcImpl struct {
	client *Client
}

// Create creates a new contact import from a CSV file.
// https://resend.com/docs/api-reference/contacts/create-contact-import
func (s *ContactImportsSvcImpl) Create(params *CreateContactImportRequest) (CreateContactImportResponse, error) {
	return s.CreateWithContext(context.Background(), params)
}

// CreateWithContext creates a new contact import from a CSV file with context.
// https://resend.com/docs/api-reference/contacts/create-contact-import
func (s *ContactImportsSvcImpl) CreateWithContext(ctx context.Context, params *CreateContactImportRequest) (CreateContactImportResponse, error) {
	if params == nil || len(params.File) == 0 {
		return CreateContactImportResponse{}, errors.New("[ERROR]: File is required")
	}

	filename := params.Filename
	if filename == "" {
		filename = "import.csv"
	}

	fields := make(map[string]string)
	if params.OnConflict != "" {
		fields["on_conflict"] = params.OnConflict
	}
	if len(params.ColumnMap) > 0 {
		b, err := json.Marshal(params.ColumnMap)
		if err != nil {
			return CreateContactImportResponse{}, fmt.Errorf("[ERROR]: Failed to encode column_map: %w", err)
		}
		fields["column_map"] = string(b)
	}
	if len(params.Segments) > 0 {
		b, err := json.Marshal(params.Segments)
		if err != nil {
			return CreateContactImportResponse{}, fmt.Errorf("[ERROR]: Failed to encode segments: %w", err)
		}
		fields["segments"] = string(b)
	}
	if len(params.Topics) > 0 {
		b, err := json.Marshal(params.Topics)
		if err != nil {
			return CreateContactImportResponse{}, fmt.Errorf("[ERROR]: Failed to encode topics: %w", err)
		}
		fields["topics"] = string(b)
	}

	req, err := s.client.NewMultipartRequest(ctx, "contacts/imports", params.File, filename, fields)
	if err != nil {
		return CreateContactImportResponse{}, errors.New("[ERROR]: Failed to create ContactImports.Create request")
	}

	resp := new(CreateContactImportResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return CreateContactImportResponse{}, err
	}
	return *resp, nil
}

// Get retrieves a single contact import by ID.
// https://resend.com/docs/api-reference/contacts/get-contact-import
func (s *ContactImportsSvcImpl) Get(id string) (ContactImport, error) {
	return s.GetWithContext(context.Background(), id)
}

// GetWithContext retrieves a single contact import by ID with context.
// https://resend.com/docs/api-reference/contacts/get-contact-import
func (s *ContactImportsSvcImpl) GetWithContext(ctx context.Context, id string) (ContactImport, error) {
	if id == "" {
		return ContactImport{}, errors.New("[ERROR]: Id is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, "contacts/imports/"+id, nil)
	if err != nil {
		return ContactImport{}, errors.New("[ERROR]: Failed to create ContactImports.Get request")
	}

	resp := new(ContactImport)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return ContactImport{}, err
	}
	return *resp, nil
}

// List retrieves a list of contact imports.
// https://resend.com/docs/api-reference/contacts/list-contact-imports
func (s *ContactImportsSvcImpl) List(options *ListContactImportsOptions) (ListContactImportsResponse, error) {
	return s.ListWithContext(context.Background(), options)
}

// ListWithContext retrieves a list of contact imports with context.
// https://resend.com/docs/api-reference/contacts/list-contact-imports
func (s *ContactImportsSvcImpl) ListWithContext(ctx context.Context, options *ListContactImportsOptions) (ListContactImportsResponse, error) {
	if options == nil {
		options = &ListContactImportsOptions{}
	}

	query := make(url.Values)
	if options.Status != "" {
		query.Set("status", options.Status)
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

	path := "contacts/imports"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListContactImportsResponse{}, errors.New("[ERROR]: Failed to create ContactImports.List request")
	}

	resp := new(ListContactImportsResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return ListContactImportsResponse{}, err
	}
	return *resp, nil
}

// NewMultipartRequest builds an HTTP multipart/form-data request for file uploads.
func (c *Client) NewMultipartRequest(ctx context.Context, path string, fileBytes []byte, filename string, fields map[string]string) (*http.Request, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, bytes.NewReader(fileBytes)); err != nil {
		return nil, err
	}

	for k, v := range fields {
		if err = w.WriteField(k, v); err != nil {
			return nil, err
		}
	}
	w.Close()

	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Accept", contentType)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}

	return req, nil
}
