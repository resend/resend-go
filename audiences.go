package resend

import (
	"context"
)

// Deprecated: Use SegmentsSvc instead. Audiences have been renamed to Segments.
// The audiences API is maintained for backward compatibility and calls the segments API internally.
type AudiencesSvc interface {
	CreateWithContext(ctx context.Context, params *CreateAudienceRequest) (CreateAudienceResponse, error)
	Create(params *CreateAudienceRequest) (CreateAudienceResponse, error)
	ListWithOptions(ctx context.Context, options *ListOptions) (ListAudiencesResponse, error)
	ListWithContext(ctx context.Context) (ListAudiencesResponse, error)
	List() (ListAudiencesResponse, error)
	GetWithContext(ctx context.Context, audienceId string) (Audience, error)
	Get(audienceId string) (Audience, error)
	RemoveWithContext(ctx context.Context, audienceId string) (RemoveAudienceResponse, error)
	Remove(audienceId string) (RemoveAudienceResponse, error)
}

// AudiencesSvcImpl wraps SegmentsSvcImpl to provide backward compatibility
// Deprecated: Use SegmentsSvcImpl instead
type AudiencesSvcImpl struct {
	segments *SegmentsSvcImpl
}

// Type aliases for backward compatibility
// Deprecated: Use CreateSegmentRequest instead
type CreateAudienceRequest = CreateSegmentRequest

// Deprecated: Use CreateSegmentResponse instead
type CreateAudienceResponse = CreateSegmentResponse

// Deprecated: Use RemoveSegmentResponse instead
type RemoveAudienceResponse = RemoveSegmentResponse

// Deprecated: Use ListSegmentsResponse instead
type ListAudiencesResponse = ListSegmentsResponse

// Deprecated: Use Segment instead
type Audience = Segment

// CreateWithContext creates a new Audience entry based on the given params
// Deprecated: Use Segments.CreateWithContext instead
// https://resend.com/docs/api-reference/segments/create-segment
func (s *AudiencesSvcImpl) CreateWithContext(ctx context.Context, params *CreateAudienceRequest) (CreateAudienceResponse, error) {
	return s.segments.CreateWithContext(ctx, params)
}

// Create creates a new Audience entry based on the given params
// Deprecated: Use Segments.Create instead
// https://resend.com/docs/api-reference/segments/create-segment
func (s *AudiencesSvcImpl) Create(params *CreateAudienceRequest) (CreateAudienceResponse, error) {
	return s.segments.Create(params)
}

// ListWithOptions returns the list of all audiences with pagination options
// Deprecated: Use Segments.ListWithOptions instead
// https://resend.com/docs/api-reference/segments/list-segments
func (s *AudiencesSvcImpl) ListWithOptions(ctx context.Context, options *ListOptions) (ListAudiencesResponse, error) {
	return s.segments.ListWithOptions(ctx, options)
}

// ListWithContext returns the list of all audiences
// Deprecated: Use Segments.ListWithContext instead
// https://resend.com/docs/api-reference/segments/list-segments
func (s *AudiencesSvcImpl) ListWithContext(ctx context.Context) (ListAudiencesResponse, error) {
	return s.segments.ListWithContext(ctx)
}

// List returns the list of all audiences
// Deprecated: Use Segments.List instead
// https://resend.com/docs/api-reference/segments/list-segments
func (s *AudiencesSvcImpl) List() (ListAudiencesResponse, error) {
	return s.segments.List()
}

// RemoveWithContext removes a given audience by id
// Deprecated: Use Segments.RemoveWithContext instead
// https://resend.com/docs/api-reference/segments/delete-segment
func (s *AudiencesSvcImpl) RemoveWithContext(ctx context.Context, audienceId string) (RemoveAudienceResponse, error) {
	return s.segments.RemoveWithContext(ctx, audienceId)
}

// Remove removes a given audience entry by id
// Deprecated: Use Segments.Remove instead
// https://resend.com/docs/api-reference/segments/delete-segment
func (s *AudiencesSvcImpl) Remove(audienceId string) (RemoveAudienceResponse, error) {
	return s.segments.Remove(audienceId)
}

// GetWithContext Retrieve a single audience.
// Deprecated: Use Segments.GetWithContext instead
// https://resend.com/docs/api-reference/segments/get-segment
func (s *AudiencesSvcImpl) GetWithContext(ctx context.Context, audienceId string) (Audience, error) {
	return s.segments.GetWithContext(ctx, audienceId)
}

// Get Retrieve a single audience.
// Deprecated: Use Segments.Get instead
// https://resend.com/docs/api-reference/segments/get-segment
func (s *AudiencesSvcImpl) Get(audienceId string) (Audience, error) {
	return s.segments.Get(audienceId)
}
