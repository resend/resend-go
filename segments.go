package resend

import (
	"context"
	"errors"
	"net/http"
)

type SegmentsSvc interface {
	CreateWithContext(ctx context.Context, params *CreateSegmentRequest) (CreateSegmentResponse, error)
	Create(params *CreateSegmentRequest) (CreateSegmentResponse, error)
	ListWithOptions(ctx context.Context, options *ListOptions) (ListSegmentsResponse, error)
	ListWithContext(ctx context.Context) (ListSegmentsResponse, error)
	List() (ListSegmentsResponse, error)
	GetWithContext(ctx context.Context, segmentId string) (Segment, error)
	Get(segmentId string) (Segment, error)
	RemoveWithContext(ctx context.Context, segmentId string) (RemoveSegmentResponse, error)
	Remove(segmentId string) (RemoveSegmentResponse, error)
}

type SegmentsSvcImpl struct {
	client *Client
}

type CreateSegmentRequest struct {
	Name string `json:"name"`
}

type CreateSegmentResponse struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Object string `json:"object"`
}

type RemoveSegmentResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

type ListSegmentsResponse struct {
	Object  string    `json:"object"`
	Data    []Segment `json:"data"`
	HasMore bool      `json:"has_more"`
}

type Segment struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Object    string `json:"object"`
	CreatedAt string `json:"created_at"`
}

// CreateWithContext creates a new Segment entry based on the given params
// https://resend.com/docs/api-reference/segments/create-segment
func (s *SegmentsSvcImpl) CreateWithContext(ctx context.Context, params *CreateSegmentRequest) (CreateSegmentResponse, error) {
	path := "segments"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return CreateSegmentResponse{}, errors.New("[ERROR]: Failed to create Segments.Create request")
	}

	// Build response recipient obj
	segmentsResp := new(CreateSegmentResponse)

	// Send Request
	_, err = s.client.Perform(req, segmentsResp)

	if err != nil {
		return CreateSegmentResponse{}, err
	}

	return *segmentsResp, nil
}

// Create creates a new Segment entry based on the given params
// https://resend.com/docs/api-reference/segments/create-segment
func (s *SegmentsSvcImpl) Create(params *CreateSegmentRequest) (CreateSegmentResponse, error) {
	return s.CreateWithContext(context.Background(), params)
}

// ListWithOptions returns the list of all segments with pagination options
// https://resend.com/docs/api-reference/segments/list-segments
func (s *SegmentsSvcImpl) ListWithOptions(ctx context.Context, options *ListOptions) (ListSegmentsResponse, error) {
	path := "segments" + buildPaginationQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListSegmentsResponse{}, errors.New("[ERROR]: Failed to create Segments.List request")
	}

	segments := new(ListSegmentsResponse)

	// Send Request
	_, err = s.client.Perform(req, segments)

	if err != nil {
		return ListSegmentsResponse{}, err
	}

	return *segments, nil
}

// ListWithContext returns the list of all segments
// https://resend.com/docs/api-reference/segments/list-segments
func (s *SegmentsSvcImpl) ListWithContext(ctx context.Context) (ListSegmentsResponse, error) {
	return s.ListWithOptions(ctx, nil)
}

// List returns the list of all segments
// https://resend.com/docs/api-reference/segments/list-segments
func (s *SegmentsSvcImpl) List() (ListSegmentsResponse, error) {
	return s.ListWithContext(context.Background())
}

// RemoveWithContext removes a given segment by id
// https://resend.com/docs/api-reference/segments/delete-segment
func (s *SegmentsSvcImpl) RemoveWithContext(ctx context.Context, segmentId string) (RemoveSegmentResponse, error) {
	path := "segments/" + segmentId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return RemoveSegmentResponse{}, errors.New("[ERROR]: Failed to create Segment.Remove request")
	}

	resp := new(RemoveSegmentResponse)

	// Send Request
	_, err = s.client.Perform(req, resp)

	if err != nil {
		return RemoveSegmentResponse{}, err
	}

	return *resp, nil
}

// Remove removes a given segment entry by id
// https://resend.com/docs/api-reference/segments/delete-segment
func (s *SegmentsSvcImpl) Remove(segmentId string) (RemoveSegmentResponse, error) {
	return s.RemoveWithContext(context.Background(), segmentId)
}

// GetWithContext Retrieve a single segment.
// https://resend.com/docs/api-reference/segments/get-segment
func (s *SegmentsSvcImpl) GetWithContext(ctx context.Context, segmentId string) (Segment, error) {
	path := "segments/" + segmentId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return Segment{}, errors.New("[ERROR]: Failed to create Segment.Get request")
	}

	segment := new(Segment)

	// Send Request
	_, err = s.client.Perform(req, segment)

	if err != nil {
		return Segment{}, err
	}

	return *segment, nil
}

// Get Retrieve a single segment.
// https://resend.com/docs/api-reference/segments/get-segment
func (s *SegmentsSvcImpl) Get(segmentId string) (Segment, error) {
	return s.GetWithContext(context.Background(), segmentId)
}
