package resend

import (
	"context"
	"net/http"
)

// DefaultSubscription represents the default subscription preference for new contacts
type DefaultSubscription string

const (
	DefaultSubscriptionOptIn  DefaultSubscription = "opt_in"
	DefaultSubscriptionOptOut DefaultSubscription = "opt_out"
)

// CreateTopicRequest is the request payload for creating a topic
type CreateTopicRequest struct {
	Name                string              `json:"name"`
	DefaultSubscription DefaultSubscription `json:"default_subscription"`
	Description         string              `json:"description,omitempty"`
}

// CreateTopicResponse is the response from creating a topic
type CreateTopicResponse struct {
	Id string `json:"id"`
}

// Topic represents a full topic object
type Topic struct {
	Id                  string              `json:"id"`
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	DefaultSubscription DefaultSubscription `json:"default_subscription"`
	CreatedAt           string              `json:"created_at"`
}

// UpdateTopicRequest is the request payload for updating a topic
// Note: default_subscription cannot be changed after creation
type UpdateTopicRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// UpdateTopicResponse is the response from updating a topic
type UpdateTopicResponse struct {
	Id string `json:"id"`
}

// RemoveTopicResponse is the response from removing a topic
type RemoveTopicResponse struct {
	Object  string `json:"object"`
	Id      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// ListTopicsResponse is the response from listing topics
type ListTopicsResponse struct {
	Object  string   `json:"object"`
	HasMore bool     `json:"has_more"`
	Data    []*Topic `json:"data"`
}

// TopicsSvc handles operations for topics
type TopicsSvc interface {
	CreateWithContext(ctx context.Context, params *CreateTopicRequest) (*CreateTopicResponse, error)
	Create(params *CreateTopicRequest) (*CreateTopicResponse, error)
	GetWithContext(ctx context.Context, topicId string) (*Topic, error)
	Get(topicId string) (*Topic, error)
	ListWithContext(ctx context.Context, options *ListOptions) (*ListTopicsResponse, error)
	List(options *ListOptions) (*ListTopicsResponse, error)
	UpdateWithContext(ctx context.Context, topicId string, params *UpdateTopicRequest) (*UpdateTopicResponse, error)
	Update(topicId string, params *UpdateTopicRequest) (*UpdateTopicResponse, error)
	RemoveWithContext(ctx context.Context, topicId string) (*RemoveTopicResponse, error)
	Remove(topicId string) (*RemoveTopicResponse, error)
}

// TopicsSvcImpl is the implementation of the TopicsSvc interface
type TopicsSvcImpl struct {
	client *Client
}

// CreateWithContext creates a new topic with the given parameters
// https://resend.com/docs/api-reference/topics/create-topic
func (s *TopicsSvcImpl) CreateWithContext(ctx context.Context, params *CreateTopicRequest) (*CreateTopicResponse, error) {
	path := "topics"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, params)
	if err != nil {
		return nil, ErrFailedToCreateTopicCreateRequest
	}

	// Build response recipient obj
	topicResponse := new(CreateTopicResponse)

	// Send Request
	_, err = s.client.Perform(req, topicResponse)

	if err != nil {
		return nil, err
	}

	return topicResponse, nil
}

// Create creates a new topic with the given parameters
// https://resend.com/docs/api-reference/topics/create-topic
func (s *TopicsSvcImpl) Create(params *CreateTopicRequest) (*CreateTopicResponse, error) {
	return s.CreateWithContext(context.Background(), params)
}

// GetWithContext retrieves a topic by ID
// https://resend.com/docs/api-reference/topics/get-topic
func (s *TopicsSvcImpl) GetWithContext(ctx context.Context, topicId string) (*Topic, error) {
	path := "topics/" + topicId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, ErrFailedToCreateTopicGetRequest
	}

	// Build response recipient obj
	topicResponse := new(Topic)

	// Send Request
	_, err = s.client.Perform(req, topicResponse)

	if err != nil {
		return nil, err
	}

	return topicResponse, nil
}

// Get retrieves a topic by ID
// https://resend.com/docs/api-reference/topics/get-topic
func (s *TopicsSvcImpl) Get(topicId string) (*Topic, error) {
	return s.GetWithContext(context.Background(), topicId)
}

// ListWithContext retrieves a list of topics with pagination options
// https://resend.com/docs/api-reference/topics/list-topics
func (s *TopicsSvcImpl) ListWithContext(ctx context.Context, options *ListOptions) (*ListTopicsResponse, error) {
	path := "topics" + buildPaginationQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, ErrFailedToCreateTopicListRequest
	}

	// Build response recipient obj
	topicResponse := new(ListTopicsResponse)

	// Send Request
	_, err = s.client.Perform(req, topicResponse)

	if err != nil {
		return nil, err
	}

	return topicResponse, nil
}

// List retrieves a list of topics with pagination options
// https://resend.com/docs/api-reference/topics/list-topics
func (s *TopicsSvcImpl) List(options *ListOptions) (*ListTopicsResponse, error) {
	return s.ListWithContext(context.Background(), options)
}

// UpdateWithContext updates a topic by ID
// https://resend.com/docs/api-reference/topics/update-topic
func (s *TopicsSvcImpl) UpdateWithContext(ctx context.Context, topicId string, params *UpdateTopicRequest) (*UpdateTopicResponse, error) {
	path := "topics/" + topicId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, params)
	if err != nil {
		return nil, ErrFailedToCreateTopicUpdateRequest
	}

	// Build response recipient obj
	topicResponse := new(UpdateTopicResponse)

	// Send Request
	_, err = s.client.Perform(req, topicResponse)

	if err != nil {
		return nil, err
	}

	return topicResponse, nil
}

// Update updates a topic by ID
// https://resend.com/docs/api-reference/topics/update-topic
func (s *TopicsSvcImpl) Update(topicId string, params *UpdateTopicRequest) (*UpdateTopicResponse, error) {
	return s.UpdateWithContext(context.Background(), topicId, params)
}

// RemoveWithContext removes a topic by ID
// https://resend.com/docs/api-reference/topics/delete-topic
func (s *TopicsSvcImpl) RemoveWithContext(ctx context.Context, topicId string) (*RemoveTopicResponse, error) {
	path := "topics/" + topicId

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, ErrFailedToCreateTopicRemoveRequest
	}

	// Build response recipient obj
	topicResponse := new(RemoveTopicResponse)

	// Send Request
	_, err = s.client.Perform(req, topicResponse)

	if err != nil {
		return nil, err
	}

	return topicResponse, nil
}

// Remove removes a topic by ID
// https://resend.com/docs/api-reference/topics/delete-topic
func (s *TopicsSvcImpl) Remove(topicId string) (*RemoveTopicResponse, error) {
	return s.RemoveWithContext(context.Background(), topicId)
}
