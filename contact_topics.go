package resend

import (
	"context"
	"errors"
	"net/http"
)

// ContactTopic represents a topic subscription for a contact
type ContactTopic struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Subscription string `json:"subscription"`
}

// ListContactTopicsResponse is the response from listing contact topics
type ListContactTopicsResponse struct {
	Object  string         `json:"object"`
	HasMore bool           `json:"has_more"`
	Data    []ContactTopic `json:"data"`
}

// TopicSubscriptionUpdate represents a single topic subscription update
type TopicSubscriptionUpdate struct {
	Id           string `json:"id"`
	Subscription string `json:"subscription"`
}

// UpdateContactTopicsRequest is the request for updating contact topics
type UpdateContactTopicsRequest struct {
	Id     string                    `json:"-"`
	Email  string                    `json:"-"`
	Topics []TopicSubscriptionUpdate `json:"topics"`
}

// UpdateContactTopicsResponse is the response from updating contact topics
type UpdateContactTopicsResponse struct {
	Id string `json:"id"`
}

// ContactTopicsSvc handles operations for contact topics
type ContactTopicsSvc interface {
	ListWithOptions(ctx context.Context, id string, options *ListOptions) (ListContactTopicsResponse, error)
	ListWithContext(ctx context.Context, id string) (ListContactTopicsResponse, error)
	List(id string) (ListContactTopicsResponse, error)
	UpdateWithContext(ctx context.Context, params *UpdateContactTopicsRequest) (UpdateContactTopicsResponse, error)
	Update(params *UpdateContactTopicsRequest) (UpdateContactTopicsResponse, error)
}

// ContactTopicsSvcImpl is the implementation of ContactTopicsSvc
type ContactTopicsSvcImpl struct {
	client *Client
}

// ListWithOptions retrieves a list of topics subscriptions for a contact with pagination options.
// The id parameter can be either a contact ID or email address.
// https://resend.com/docs/api-reference/contacts/get-contact-topics
func (s *ContactTopicsSvcImpl) ListWithOptions(ctx context.Context, id string, options *ListOptions) (ListContactTopicsResponse, error) {
	if id == "" {
		return ListContactTopicsResponse{}, errors.New("[ERROR]: Contact ID or email is missing")
	}

	path := "contacts/" + id + "/topics" + buildPaginationQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListContactTopicsResponse{}, errors.New("[ERROR]: Failed to create ContactTopics.List request")
	}

	topics := new(ListContactTopicsResponse)

	// Send Request
	_, err = s.client.Perform(req, topics)

	if err != nil {
		return ListContactTopicsResponse{}, err
	}

	return *topics, nil
}

// ListWithContext retrieves a list of topics subscriptions for a contact.
// The id parameter can be either a contact ID or email address.
// https://resend.com/docs/api-reference/contacts/get-contact-topics
func (s *ContactTopicsSvcImpl) ListWithContext(ctx context.Context, id string) (ListContactTopicsResponse, error) {
	return s.ListWithOptions(ctx, id, nil)
}

// List retrieves a list of topics subscriptions for a contact.
// The id parameter can be either a contact ID or email address.
// https://resend.com/docs/api-reference/contacts/get-contact-topics
func (s *ContactTopicsSvcImpl) List(id string) (ListContactTopicsResponse, error) {
	return s.ListWithContext(context.Background(), id)
}

// UpdateWithContext updates topic subscriptions for a contact.
// Either Id or Email must be provided in the params.
// https://resend.com/docs/api-reference/contacts/update-contact-topics
func (s *ContactTopicsSvcImpl) UpdateWithContext(ctx context.Context, params *UpdateContactTopicsRequest) (UpdateContactTopicsResponse, error) {
	if params.Id == "" && params.Email == "" {
		return UpdateContactTopicsResponse{}, errors.New("[ERROR]: Contact ID or email is missing")
	}

	if len(params.Topics) == 0 {
		return UpdateContactTopicsResponse{}, errors.New("[ERROR]: Topics array is empty")
	}

	var identifier string
	if params.Id != "" {
		identifier = params.Id
	} else {
		identifier = params.Email
	}

	path := "contacts/" + identifier + "/topics"

	// Prepare request - send only the topics array as body
	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, params.Topics)
	if err != nil {
		return UpdateContactTopicsResponse{}, errors.New("[ERROR]: Failed to create ContactTopics.Update request")
	}

	resp := new(UpdateContactTopicsResponse)

	// Send Request
	_, err = s.client.Perform(req, resp)

	if err != nil {
		return UpdateContactTopicsResponse{}, err
	}

	return *resp, nil
}

// Update updates topic subscriptions for a contact.
// Either Id or Email must be provided in the params.
// https://resend.com/docs/api-reference/contacts/update-contact-topics
func (s *ContactTopicsSvcImpl) Update(params *UpdateContactTopicsRequest) (UpdateContactTopicsResponse, error) {
	return s.UpdateWithContext(context.Background(), params)
}
