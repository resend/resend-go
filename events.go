package resend

import (
	"context"
	"net/http"
)

type EventSchemaType = string

const (
	EventSchemaTypeString  EventSchemaType = "string"
	EventSchemaTypeNumber  EventSchemaType = "number"
	EventSchemaTypeBoolean EventSchemaType = "boolean"
	EventSchemaTypeDate    EventSchemaType = "date"
)

type CreateEventRequest struct {
	Name   string            `json:"name"`
	Schema map[string]string `json:"schema,omitempty"`
}

type CreateEventResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

type UpdateEventRequest struct {
	Schema map[string]string `json:"schema,omitempty"`
}

type UpdateEventResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

type DeleteEventResponse struct {
	Object  string `json:"object"`
	Id      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

type SendEventRequest struct {
	Event     string         `json:"event"`
	ContactId string         `json:"contact_id,omitempty"`
	Email     string         `json:"email,omitempty"`
	Payload   map[string]any `json:"payload,omitempty"`
}

type SendEventResponse struct {
	Object string `json:"object"`
	Event  string `json:"event"`
}

type EventSummary struct {
	Id        string            `json:"id"`
	Name      string            `json:"name"`
	Schema    map[string]string `json:"schema"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt *string           `json:"updated_at"`
}

type ListEventsResponse struct {
	Object  string         `json:"object"`
	HasMore bool           `json:"has_more"`
	Data    []EventSummary `json:"data"`
}

type Event struct {
	Object    string            `json:"object"`
	Id        string            `json:"id"`
	Name      string            `json:"name"`
	Schema    map[string]string `json:"schema"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt *string           `json:"updated_at"`
}

type EventsSvc interface {
	CreateWithContext(ctx context.Context, params *CreateEventRequest) (CreateEventResponse, error)
	Create(params *CreateEventRequest) (CreateEventResponse, error)
	GetWithContext(ctx context.Context, identifier string) (Event, error)
	Get(identifier string) (Event, error)
	ListWithContext(ctx context.Context) (ListEventsResponse, error)
	List() (ListEventsResponse, error)
	ListWithOptions(ctx context.Context, options *ListOptions) (ListEventsResponse, error)
	UpdateWithContext(ctx context.Context, identifier string, params *UpdateEventRequest) (UpdateEventResponse, error)
	Update(identifier string, params *UpdateEventRequest) (UpdateEventResponse, error)
	RemoveWithContext(ctx context.Context, identifier string) (DeleteEventResponse, error)
	Remove(identifier string) (DeleteEventResponse, error)
	SendWithContext(ctx context.Context, params *SendEventRequest) (SendEventResponse, error)
	Send(params *SendEventRequest) (SendEventResponse, error)
}

type EventsSvcImpl struct {
	client *Client
}

// CreateWithContext creates a new Event definition
// https://resend.com/docs/api-reference/events/create-event
func (s *EventsSvcImpl) CreateWithContext(ctx context.Context, params *CreateEventRequest) (CreateEventResponse, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "events", params)
	if err != nil {
		return CreateEventResponse{}, ErrFailedToCreateEventCreateRequest
	}

	resp := new(CreateEventResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return CreateEventResponse{}, err
	}

	return *resp, nil
}

// Create creates a new Event definition
func (s *EventsSvcImpl) Create(params *CreateEventRequest) (CreateEventResponse, error) {
	return s.CreateWithContext(context.Background(), params)
}

// GetWithContext retrieves an Event by ID or name
// https://resend.com/docs/api-reference/events/get-event
func (s *EventsSvcImpl) GetWithContext(ctx context.Context, identifier string) (Event, error) {
	path := "events/" + identifier

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return Event{}, ErrFailedToCreateEventGetRequest
	}

	resp := new(Event)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return Event{}, err
	}

	return *resp, nil
}

// Get retrieves an Event by ID or name
func (s *EventsSvcImpl) Get(identifier string) (Event, error) {
	return s.GetWithContext(context.Background(), identifier)
}

// ListWithOptions retrieves a list of Events with pagination options
// https://resend.com/docs/api-reference/events/list-events
func (s *EventsSvcImpl) ListWithOptions(ctx context.Context, options *ListOptions) (ListEventsResponse, error) {
	path := "events" + buildPaginationQuery(options)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListEventsResponse{}, ErrFailedToCreateEventListRequest
	}

	resp := new(ListEventsResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return ListEventsResponse{}, err
	}

	return *resp, nil
}

// ListWithContext retrieves a list of Events
// https://resend.com/docs/api-reference/events/list-events
func (s *EventsSvcImpl) ListWithContext(ctx context.Context) (ListEventsResponse, error) {
	return s.ListWithOptions(ctx, nil)
}

// List retrieves a list of Events
func (s *EventsSvcImpl) List() (ListEventsResponse, error) {
	return s.ListWithContext(context.Background())
}

// UpdateWithContext updates an Event's schema by ID or name
// https://resend.com/docs/api-reference/events/update-event
func (s *EventsSvcImpl) UpdateWithContext(ctx context.Context, identifier string, params *UpdateEventRequest) (UpdateEventResponse, error) {
	path := "events/" + identifier

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, params)
	if err != nil {
		return UpdateEventResponse{}, ErrFailedToCreateEventUpdateRequest
	}

	resp := new(UpdateEventResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return UpdateEventResponse{}, err
	}

	return *resp, nil
}

// Update updates an Event's schema by ID or name
func (s *EventsSvcImpl) Update(identifier string, params *UpdateEventRequest) (UpdateEventResponse, error) {
	return s.UpdateWithContext(context.Background(), identifier, params)
}

// RemoveWithContext deletes an Event by ID or name
// https://resend.com/docs/api-reference/events/delete-event
func (s *EventsSvcImpl) RemoveWithContext(ctx context.Context, identifier string) (DeleteEventResponse, error) {
	path := "events/" + identifier

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return DeleteEventResponse{}, ErrFailedToCreateEventRemoveRequest
	}

	resp := new(DeleteEventResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return DeleteEventResponse{}, err
	}

	return *resp, nil
}

// Remove deletes an Event by ID or name
func (s *EventsSvcImpl) Remove(identifier string) (DeleteEventResponse, error) {
	return s.RemoveWithContext(context.Background(), identifier)
}

// SendWithContext sends an Event to trigger automations
// https://resend.com/docs/api-reference/events/send-event
func (s *EventsSvcImpl) SendWithContext(ctx context.Context, params *SendEventRequest) (SendEventResponse, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "events/send", params)
	if err != nil {
		return SendEventResponse{}, ErrFailedToCreateEventSendRequest
	}

	resp := new(SendEventResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return SendEventResponse{}, err
	}

	return *resp, nil
}

// Send sends an Event to trigger automations
func (s *EventsSvcImpl) Send(params *SendEventRequest) (SendEventResponse, error) {
	return s.SendWithContext(context.Background(), params)
}
