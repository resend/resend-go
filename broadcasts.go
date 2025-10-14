package resend

import (
	"context"
	"errors"
	"net/http"
)

type SendBroadcastRequest struct {
	BroadcastId string `json:"broadcast_id"`

	//Schedule email to be sent later. The date should be in language natural (e.g.: in 1 min)
	// or ISO 8601 format (e.g: 2024-08-05T11:52:01.858Z).
	ScheduledAt string `json:"scheduled_at"`
}

type CreateBroadcastRequest struct {
	AudienceId string   `json:"audience_id,omitempty"`
	From       string   `json:"from,omitempty"`
	Subject    string   `json:"subject,omitempty"`
	ReplyTo    []string `json:"reply_to,omitempty"`
	Html       string   `json:"html,omitempty"`
	Text       string   `json:"text,omitempty"`
	Name       string   `json:"name,omitempty""`
}

type UpdateBroadcastRequest struct {
	BroadcastId string   `json:"broadcast_id,omitempty"`
	AudienceId  string   `json:"audience_id,omitempty"`
	From        string   `json:"from,omitempty"`
	Subject     string   `json:"subject,omitempty"`
	ReplyTo     []string `json:"reply_to,omitempty"`
	Html        string   `json:"html,omitempty"`
	Text        string   `json:"text,omitempty"`
	Name        string   `json:"name,omitempty"`
}

type CreateBroadcastResponse struct {
	Id string `json:"id"`
}

type UpdateBroadcastResponse struct {
	Id string `json:"id"`
}

type SendBroadcastResponse struct {
	Id string `json:"id"`
}

type RemoveBroadcastResponse struct {
	Object  string `json:"object"`
	Id      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

type ListBroadcastsResponse struct {
	Object  string      `json:"object"`
	Data    []Broadcast `json:"data"`
	HasMore bool        `json:"has_more"`
}

type Broadcast struct {
	Object      string   `json:"object"`
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	AudienceId  string   `json:"audience_id"`
	From        string   `json:"from"`
	Subject     string   `json:"subject"`
	ReplyTo     []string `json:"reply_to"`
	PreviewText string   `json:"preview_text"`
	Status      string   `json:"status"`
	CreatedAt   string   `json:"created_at"`
	ScheduledAt string   `json:"scheduled_at"`
	SentAt      string   `json:"sent_at"`
	Html        string   `json:"html"`
	Text        string   `json:"text"`
}

type BroadcastsSvc interface {
	CreateWithContext(ctx context.Context, params *CreateBroadcastRequest) (CreateBroadcastResponse, error)
	Create(params *CreateBroadcastRequest) (CreateBroadcastResponse, error)

	UpdateWithContext(ctx context.Context, params *UpdateBroadcastRequest) (UpdateBroadcastResponse, error)
	Update(params *UpdateBroadcastRequest) (UpdateBroadcastResponse, error)

	ListWithOptions(ctx context.Context, options *ListOptions) (ListBroadcastsResponse, error)
	ListWithContext(ctx context.Context) (ListBroadcastsResponse, error)
	List() (ListBroadcastsResponse, error)

	GetWithContext(ctx context.Context, broadcastId string) (Broadcast, error)
	Get(broadcastId string) (Broadcast, error)

	SendWithContext(ctx context.Context, params *SendBroadcastRequest) (SendBroadcastResponse, error)
	Send(params *SendBroadcastRequest) (SendBroadcastResponse, error)

	RemoveWithContext(ctx context.Context, broadcastId string) (RemoveBroadcastResponse, error)
	Remove(broadcastId string) (RemoveBroadcastResponse, error)
}

type BroadcastsSvcImpl struct {
	client *Client
}

// CreateWithContext creates a new Broadcast based on the given params
// https://resend.com/docs/api-reference/broadcasts/create-broadcast
func (s *BroadcastsSvcImpl) CreateWithContext(ctx context.Context, params *CreateBroadcastRequest) (CreateBroadcastResponse, error) {
	path := "/broadcasts"

	if params.AudienceId == "" {
		return CreateBroadcastResponse{}, errors.New("[ERROR]: AudienceId cannot be empty")
	}

	if params.From == "" {
		return CreateBroadcastResponse{}, errors.New("[ERROR]: From cannot be empty")
	}

	if params.Subject == "" {
		return CreateBroadcastResponse{}, errors.New("[ERROR]: Subject cannot be empty")
	}

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return CreateBroadcastResponse{}, ErrFailedToCreateBroadcastCreateRequest
	}

	// Build response recipient obj
	broadcastResp := new(CreateBroadcastResponse)

	// Send Request
	_, err = s.client.Perform(req, broadcastResp)

	if err != nil {
		return CreateBroadcastResponse{}, err
	}

	return *broadcastResp, nil
}

// Create creates a new Broadcast based on the given params
func (s *BroadcastsSvcImpl) Create(params *CreateBroadcastRequest) (CreateBroadcastResponse, error) {
	return s.CreateWithContext(context.Background(), params)
}

// UpdateWithContext updates a given broadcast entry
// https://resend.com/docs/api-reference/broadcasts/update-broadcast
func (s *BroadcastsSvcImpl) UpdateWithContext(ctx context.Context, params *UpdateBroadcastRequest) (UpdateBroadcastResponse, error) {
	if params.BroadcastId == "" {
		return UpdateBroadcastResponse{}, errors.New("[ERROR]: BroadcastId cannot be empty")
	}

	path := "/broadcasts/" + params.BroadcastId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, params)
	if err != nil {
		return UpdateBroadcastResponse{}, ErrFailedToCreateBroadcastUpdateRequest
	}

	// Build response recipient obj
	broadcastResp := new(UpdateBroadcastResponse)

	// Send Request
	_, err = s.client.Perform(req, broadcastResp)

	if err != nil {
		return UpdateBroadcastResponse{}, err
	}

	return *broadcastResp, nil
}

func (s *BroadcastsSvcImpl) Update(params *UpdateBroadcastRequest) (UpdateBroadcastResponse, error) {
	return s.UpdateWithContext(context.Background(), params)
}

// GetWithContext Retrieve a single broadcast.
// https://resend.com/docs/api-reference/broadcasts/get-broadcast
func (s *BroadcastsSvcImpl) GetWithContext(ctx context.Context, broadcastId string) (Broadcast, error) {

	if broadcastId == "" {
		return Broadcast{}, errors.New("[ERROR]: broadcastId cannot be empty")
	}

	path := "broadcasts/" + broadcastId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return Broadcast{}, errors.New("[ERROR]: Failed to create Broadcast.Get request")
	}

	broadcast := new(Broadcast)

	// Send Request
	_, err = s.client.Perform(req, broadcast)

	if err != nil {
		return Broadcast{}, err
	}

	return *broadcast, nil
}

// Get retrieves a single broadcast.
func (s *BroadcastsSvcImpl) Get(broadcastId string) (Broadcast, error) {
	return s.GetWithContext(context.Background(), broadcastId)
}

// SendWithContext Sends broadcasts to your audience.
// https://resend.com/docs/api-reference/broadcasts/send-broadcast
func (s *BroadcastsSvcImpl) SendWithContext(ctx context.Context, params *SendBroadcastRequest) (SendBroadcastResponse, error) {
	if params.BroadcastId == "" {
		return SendBroadcastResponse{}, errors.New("[ERROR]: BroadcastId cannot be empty")
	}

	path := "/broadcasts/" + params.BroadcastId + "/send"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return SendBroadcastResponse{}, ErrFailedToCreateBroadcastSendRequest
	}

	// Build response recipient obj
	broadcastResp := new(SendBroadcastResponse)

	// Send Request
	_, err = s.client.Perform(req, broadcastResp)

	if err != nil {
		return SendBroadcastResponse{}, err
	}

	return *broadcastResp, nil
}

// Send sends broadcasts to your audience.
func (s *BroadcastsSvcImpl) Send(params *SendBroadcastRequest) (SendBroadcastResponse, error) {
	return s.SendWithContext(context.Background(), params)
}

// RemoveWithContext removes a given broadcast by id
// https://resend.com/docs/api-reference/broadcasts/delete-broadcast
func (s *BroadcastsSvcImpl) RemoveWithContext(ctx context.Context, broadcastId string) (RemoveBroadcastResponse, error) {
	path := "broadcasts/" + broadcastId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return RemoveBroadcastResponse{}, errors.New("[ERROR]: Failed to create Broadcast.Remove request")
	}

	resp := new(RemoveBroadcastResponse)

	// Send Request
	_, err = s.client.Perform(req, resp)

	if err != nil {
		return RemoveBroadcastResponse{}, err
	}

	return *resp, nil
}

// Remove removes a given broadcast entry by id
func (s *BroadcastsSvcImpl) Remove(broadcastId string) (RemoveBroadcastResponse, error) {
	return s.RemoveWithContext(context.Background(), broadcastId)
}

// ListWithOptions returns the list of all broadcasts with pagination options
// https://resend.com/docs/api-reference/broadcasts/list-broadcasts
func (s *BroadcastsSvcImpl) ListWithOptions(ctx context.Context, options *ListOptions) (ListBroadcastsResponse, error) {
	path := "broadcasts" + buildPaginationQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListBroadcastsResponse{}, errors.New("[ERROR]: Failed to create Broadcasts.List request")
	}

	broadcasts := new(ListBroadcastsResponse)

	// Send Request
	_, err = s.client.Perform(req, broadcasts)

	if err != nil {
		return ListBroadcastsResponse{}, err
	}

	return *broadcasts, nil
}

// ListWithContext returns the list of all broadcasts
// https://resend.com/docs/api-reference/broadcasts/list-broadcasts
func (s *BroadcastsSvcImpl) ListWithContext(ctx context.Context) (ListBroadcastsResponse, error) {
	return s.ListWithOptions(ctx, nil)
}

// List returns the list of all broadcasts
func (s *BroadcastsSvcImpl) List() (ListBroadcastsResponse, error) {
	return s.ListWithContext(context.Background())
}
