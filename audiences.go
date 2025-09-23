package resend

import (
	"context"
	"errors"
	"net/http"
)

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

type AudiencesSvcImpl struct {
	client *Client
}

type CreateAudienceRequest struct {
	Name string `json:"name"`
}

type CreateAudienceResponse struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Object string `json:"object"`
}

type RemoveAudienceResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

type ListAudiencesResponse struct {
	Object  string     `json:"object"`
	Data    []Audience `json:"data"`
	HasMore bool       `json:"has_more"`
}

type Audience struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Object    string `json:"object"`
	CreatedAt string `json:"created_at"`
}

// CreateWithContext creates a new Audience entry based on the given params
// https://resend.com/docs/api-reference/audiences/create-audience
func (s *AudiencesSvcImpl) CreateWithContext(ctx context.Context, params *CreateAudienceRequest) (CreateAudienceResponse, error) {
	path := "audiences"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return CreateAudienceResponse{}, errors.New("[ERROR]: Failed to create Audiences.Create request")
	}

	// Build response recipient obj
	audiencesResp := new(CreateAudienceResponse)

	// Send Request
	_, err = s.client.Perform(req, audiencesResp)

	if err != nil {
		return CreateAudienceResponse{}, err
	}

	return *audiencesResp, nil
}

// Create creates a new Audience entry based on the given params
// https://resend.com/docs/api-reference/audiences/create-audience
func (s *AudiencesSvcImpl) Create(params *CreateAudienceRequest) (CreateAudienceResponse, error) {
	return s.CreateWithContext(context.Background(), params)
}

// ListWithOptions returns the list of all audiences with pagination options
// https://resend.com/docs/api-reference/audiences/list-audiences
func (s *AudiencesSvcImpl) ListWithOptions(ctx context.Context, options *ListOptions) (ListAudiencesResponse, error) {
	path := "audiences" + buildPaginationQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListAudiencesResponse{}, errors.New("[ERROR]: Failed to create Audiences.List request")
	}

	audiences := new(ListAudiencesResponse)

	// Send Request
	_, err = s.client.Perform(req, audiences)

	if err != nil {
		return ListAudiencesResponse{}, err
	}

	return *audiences, nil
}

// ListWithContext returns the list of all audiences
// https://resend.com/docs/api-reference/audiences/list-audiences
func (s *AudiencesSvcImpl) ListWithContext(ctx context.Context) (ListAudiencesResponse, error) {
	return s.ListWithOptions(ctx, nil)
}

// List returns the list of all audiences
// https://resend.com/docs/api-reference/audiences/list-audiences
func (s *AudiencesSvcImpl) List() (ListAudiencesResponse, error) {
	return s.ListWithContext(context.Background())
}

// RemoveWithContext removes a given audience by id
// https://resend.com/docs/api-reference/audiences/delete-audience
func (s *AudiencesSvcImpl) RemoveWithContext(ctx context.Context, audienceId string) (RemoveAudienceResponse, error) {
	path := "audiences/" + audienceId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return RemoveAudienceResponse{}, errors.New("[ERROR]: Failed to create Audience.Remove request")
	}

	resp := new(RemoveAudienceResponse)

	// Send Request
	_, err = s.client.Perform(req, resp)

	if err != nil {
		return RemoveAudienceResponse{}, err
	}

	return *resp, nil
}

// Remove removes a given audience entry by id
// https://resend.com/docs/api-reference/audiences/delete-audience
func (s *AudiencesSvcImpl) Remove(audienceId string) (RemoveAudienceResponse, error) {
	return s.RemoveWithContext(context.Background(), audienceId)
}

// GetWithContext Retrieve a single audience.
// https://resend.com/docs/api-reference/audiences/get-audience
func (s *AudiencesSvcImpl) GetWithContext(ctx context.Context, audienceId string) (Audience, error) {
	path := "audiences/" + audienceId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return Audience{}, errors.New("[ERROR]: Failed to create Audience.Get request")
	}

	audience := new(Audience)

	// Send Request
	_, err = s.client.Perform(req, audience)

	if err != nil {
		return Audience{}, err
	}

	return *audience, nil
}

// Get Retrieve a single audience.
// https://resend.com/docs/api-reference/audiences/get-audience
func (s *AudiencesSvcImpl) Get(audienceId string) (Audience, error) {
	return s.GetWithContext(context.Background(), audienceId)
}
