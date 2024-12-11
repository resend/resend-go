package resend

import (
	"context"
	"net/http"
)

type CreateBroadcastRequest struct {
	AudienceId string   `json:"audience_id"`
	From       string   `json:"from"`
	Subject    string   `json:"subject"`
	ReplyTo    []string `json:"reply_to"`
	Html       string   `json:"html"`
	Text       string   `json:"text"`
	Name       string   `json:"name"`
}

type CreateBroadcastResponse struct {
	Id string `json:"id"`
}

type BroadcastsSvc interface {
	CreateWithContext(ctx context.Context, params *CreateBroadcastRequest) (CreateBroadcastResponse, error)
	Create(params *CreateBroadcastRequest) (CreateBroadcastResponse, error)
}

type BroadcastsSvcImpl struct {
	client *Client
}

// CreateWithContext creates a new Broadcast based on the given params
// https://resend.com/docs/api-reference/broadcasts/create-broadcast
func (s *BroadcastsSvcImpl) CreateWithContext(ctx context.Context, params *CreateBroadcastRequest) (CreateBroadcastResponse, error) {
	path := "/broadcasts"

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
