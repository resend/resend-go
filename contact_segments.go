package resend

import (
	"context"
	"errors"
	"net/http"
)

type ContactSegmentsSvc interface {
	AddWithContext(ctx context.Context, params *AddContactSegmentRequest) (AddContactSegmentResponse, error)
	Add(params *AddContactSegmentRequest) (AddContactSegmentResponse, error)
	RemoveWithContext(ctx context.Context, params *RemoveContactSegmentRequest) (RemoveContactSegmentResponse, error)
	Remove(params *RemoveContactSegmentRequest) (RemoveContactSegmentResponse, error)
	ListWithOptions(ctx context.Context, params *ListContactSegmentsRequest, options *ListOptions) (ListContactSegmentsResponse, error)
	ListWithContext(ctx context.Context, params *ListContactSegmentsRequest) (ListContactSegmentsResponse, error)
	List(params *ListContactSegmentsRequest) (ListContactSegmentsResponse, error)
}

type ContactSegmentsSvcImpl struct {
	client *Client
}

type AddContactSegmentRequest struct {
	SegmentId string `json:"segment_id"`
	ContactId string `json:"contact_id,omitempty"`
	Email     string `json:"email,omitempty"`
}

type AddContactSegmentResponse struct {
	Id     string `json:"id"`
	Object string `json:"object"`
}

type RemoveContactSegmentRequest struct {
	SegmentId string `json:"segment_id"`
	ContactId string `json:"contact_id,omitempty"`
	Email     string `json:"email,omitempty"`
}

type RemoveContactSegmentResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

type ListContactSegmentsRequest struct {
	ContactId string `json:"contact_id,omitempty"`
	Email     string `json:"email,omitempty"`
}

type ListContactSegmentsResponse struct {
	Object  string    `json:"object"`
	Data    []Segment `json:"data"`
	HasMore bool      `json:"has_more"`
}

// AddWithContext adds a contact to a segment
// https://resend.com/docs/api-reference/contacts/add-contact-to-segment
func (s *ContactSegmentsSvcImpl) AddWithContext(ctx context.Context, params *AddContactSegmentRequest) (AddContactSegmentResponse, error) {
	if params.SegmentId == "" {
		return AddContactSegmentResponse{}, errors.New("[ERROR]: SegmentId is required")
	}

	if params.ContactId == "" && params.Email == "" {
		return AddContactSegmentResponse{}, errors.New("[ERROR]: Either ContactId or Email must be provided")
	}

	// Determine the identifier to use in the URL
	identifier := params.ContactId
	if identifier == "" {
		identifier = params.Email
	}

	path := "contacts/" + identifier + "/segments/" + params.SegmentId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return AddContactSegmentResponse{}, errors.New("[ERROR]: Failed to create ContactSegments.Add request")
	}

	// Build response recipient obj
	resp := new(AddContactSegmentResponse)

	// Send Request
	_, err = s.client.Perform(req, resp)

	if err != nil {
		return AddContactSegmentResponse{}, err
	}

	return *resp, nil
}

// Add adds a contact to a segment
// https://resend.com/docs/api-reference/contacts/add-contact-to-segment
func (s *ContactSegmentsSvcImpl) Add(params *AddContactSegmentRequest) (AddContactSegmentResponse, error) {
	return s.AddWithContext(context.Background(), params)
}

// RemoveWithContext removes a contact from a segment
// https://resend.com/docs/api-reference/contacts/remove-contact-from-segment
func (s *ContactSegmentsSvcImpl) RemoveWithContext(ctx context.Context, params *RemoveContactSegmentRequest) (RemoveContactSegmentResponse, error) {
	if params.SegmentId == "" {
		return RemoveContactSegmentResponse{}, errors.New("[ERROR]: SegmentId is required")
	}

	if params.ContactId == "" && params.Email == "" {
		return RemoveContactSegmentResponse{}, errors.New("[ERROR]: Either ContactId or Email must be provided")
	}

	// Determine the identifier to use in the URL
	identifier := params.ContactId
	if identifier == "" {
		identifier = params.Email
	}

	path := "contacts/" + identifier + "/segments/" + params.SegmentId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return RemoveContactSegmentResponse{}, errors.New("[ERROR]: Failed to create ContactSegments.Remove request")
	}

	resp := new(RemoveContactSegmentResponse)

	// Send Request
	_, err = s.client.Perform(req, resp)

	if err != nil {
		return RemoveContactSegmentResponse{}, err
	}

	return *resp, nil
}

// Remove removes a contact from a segment
// https://resend.com/docs/api-reference/contacts/remove-contact-from-segment
func (s *ContactSegmentsSvcImpl) Remove(params *RemoveContactSegmentRequest) (RemoveContactSegmentResponse, error) {
	return s.RemoveWithContext(context.Background(), params)
}

// ListWithOptions returns the list of all segments for a contact with pagination options
// https://resend.com/docs/api-reference/contacts/list-contact-segments
func (s *ContactSegmentsSvcImpl) ListWithOptions(ctx context.Context, params *ListContactSegmentsRequest, options *ListOptions) (ListContactSegmentsResponse, error) {
	if params.ContactId == "" && params.Email == "" {
		return ListContactSegmentsResponse{}, errors.New("[ERROR]: Either ContactId or Email must be provided")
	}

	// Determine the identifier to use in the URL
	identifier := params.ContactId
	if identifier == "" {
		identifier = params.Email
	}

	path := "contacts/" + identifier + "/segments" + buildPaginationQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListContactSegmentsResponse{}, errors.New("[ERROR]: Failed to create ContactSegments.List request")
	}

	segments := new(ListContactSegmentsResponse)

	// Send Request
	_, err = s.client.Perform(req, segments)

	if err != nil {
		return ListContactSegmentsResponse{}, err
	}

	return *segments, nil
}

// ListWithContext returns the list of all segments for a contact
// https://resend.com/docs/api-reference/contacts/list-contact-segments
func (s *ContactSegmentsSvcImpl) ListWithContext(ctx context.Context, params *ListContactSegmentsRequest) (ListContactSegmentsResponse, error) {
	return s.ListWithOptions(ctx, params, nil)
}

// List returns the list of all segments for a contact
// https://resend.com/docs/api-reference/contacts/list-contact-segments
func (s *ContactSegmentsSvcImpl) List(params *ListContactSegmentsRequest) (ListContactSegmentsResponse, error) {
	return s.ListWithContext(context.Background(), params)
}
