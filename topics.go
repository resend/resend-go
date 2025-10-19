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

// TopicsSvc handles operations for topics
type TopicsSvc interface {
	CreateWithContext(ctx context.Context, params *CreateTopicRequest) (*CreateTopicResponse, error)
	Create(params *CreateTopicRequest) (*CreateTopicResponse, error)
	GetWithContext(ctx context.Context, topicId string) (*Topic, error)
	Get(topicId string) (*Topic, error)
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
