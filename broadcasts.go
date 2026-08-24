package resend

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type SendBroadcastRequest struct {
	BroadcastId string `json:"broadcast_id"`

	//Schedule email to be sent later. The date should be in language natural (e.g.: in 1 min)
	// or ISO 8601 format (e.g: 2024-08-05T11:52:01.858Z).
	ScheduledAt string `json:"scheduled_at"`
}

type CreateBroadcastRequest struct {
	SegmentId  string   `json:"segment_id,omitempty"`
	AudienceId string   `json:"audience_id,omitempty"` // Deprecated: Use SegmentId instead
	From       string   `json:"from,omitempty"`
	Subject    string   `json:"subject,omitempty"`
	ReplyTo    []string `json:"reply_to,omitempty"`
	Html       string   `json:"html,omitempty"`
	Text       string   `json:"text,omitempty"`
	Name       string   `json:"name,omitempty"`

	// Send the broadcast immediately upon creation instead of creating a draft.
	Send bool `json:"send,omitempty"`

	// Schedule email to be sent later. The date should be in natural language (e.g.: in 1 min)
	// or ISO 8601 format (e.g: 2024-08-05T11:52:01.858Z).
	// Only valid when Send is true.
	ScheduledAt string `json:"scheduled_at,omitempty"`
}

type UpdateBroadcastRequest struct {
	BroadcastId string   `json:"broadcast_id,omitempty"`
	SegmentId   string   `json:"segment_id,omitempty"`
	AudienceId  string   `json:"audience_id,omitempty"` // Deprecated: Use SegmentId instead
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

type CancelBroadcastResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
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

// BroadcastRecipientEventType is the recipient event type to filter by when listing a
// broadcast's recipients.
type BroadcastRecipientEventType = string

const (
	BroadcastRecipientEventTypeSent         BroadcastRecipientEventType = "sent"
	BroadcastRecipientEventTypeDelivered    BroadcastRecipientEventType = "delivered"
	BroadcastRecipientEventTypeOpened       BroadcastRecipientEventType = "opened"
	BroadcastRecipientEventTypeClicked      BroadcastRecipientEventType = "clicked"
	BroadcastRecipientEventTypeBounced      BroadcastRecipientEventType = "bounced"
	BroadcastRecipientEventTypeComplained   BroadcastRecipientEventType = "complained"
	BroadcastRecipientEventTypeUnsubscribed BroadcastRecipientEventType = "unsubscribed"
	BroadcastRecipientEventTypeSuppressed   BroadcastRecipientEventType = "suppressed"
)

// BroadcastRecipientBounceType is the bounce type to filter by. Only meaningful when
// ListBroadcastRecipientsOptions.Type is BroadcastRecipientEventTypeBounced.
type BroadcastRecipientBounceType = string

const (
	BroadcastRecipientBounceTypePermanent    BroadcastRecipientBounceType = "permanent"
	BroadcastRecipientBounceTypeTransient    BroadcastRecipientBounceType = "transient"
	BroadcastRecipientBounceTypeUndetermined BroadcastRecipientBounceType = "undetermined"
)

// BroadcastRecipientClickedLink is a link clicked by a recipient. Only present when Type is
// BroadcastRecipientEventTypeClicked.
type BroadcastRecipientClickedLink struct {
	Url    string `json:"url"`
	Clicks int    `json:"clicks"`
}

// BroadcastRecipient is a single recipient row for the requested event type. Count, BounceType
// and ClickedLinks are only populated depending on ListBroadcastRecipientsOptions.Type.
type BroadcastRecipient struct {
	// Id is an opaque cursor identifying this row, used for pagination. It is not a real entity id.
	Id        string  `json:"id"`
	ContactId *string `json:"contact_id"`
	Email     string  `json:"email"`

	// Count is the number of times this recipient triggered the event. Only present when Type is
	// BroadcastRecipientEventTypeOpened or BroadcastRecipientEventTypeClicked.
	Count int `json:"count,omitempty"`

	// BounceType is only present when Type is BroadcastRecipientEventTypeBounced.
	BounceType BroadcastRecipientBounceType `json:"bounce_type,omitempty"`

	// ClickedLinks is only present when Type is BroadcastRecipientEventTypeClicked.
	ClickedLinks []BroadcastRecipientClickedLink `json:"clicked_links,omitempty"`
}

type ListBroadcastRecipientsResponse struct {
	Object  string               `json:"object"`
	HasMore bool                 `json:"has_more"`
	Data    []BroadcastRecipient `json:"data"`
}

// ListBroadcastRecipientsOptions contains the filter and pagination parameters for listing a
// broadcast's recipients.
type ListBroadcastRecipientsOptions struct {
	// Type is the recipient event type to filter by. Required.
	Type BroadcastRecipientEventType

	// Email filters recipients whose email address contains this substring.
	Email string

	// BounceType filters bounced recipients by bounce type. Only meaningful when Type is
	// BroadcastRecipientEventTypeBounced.
	BounceType BroadcastRecipientBounceType

	Limit  *int
	After  *string
	Before *string
}

type ListBroadcastClickedLinksResponse struct {
	Object  string                 `json:"object"`
	Data    []BroadcastClickedLink `json:"data"`
	HasMore bool                   `json:"has_more"`
}

type BroadcastClickedLink struct {
	// Id is an opaque cursor for this row, used only for pagination. It does
	// not identify any entity in Resend.
	Id           string `json:"id"`
	Url          string `json:"url"`
	Clicks       int    `json:"clicks"`
	UniqueClicks int    `json:"unique_clicks"`
}

type Broadcast struct {
	Object      string   `json:"object"`
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	SegmentId   string   `json:"segment_id"`
	AudienceId  string   `json:"audience_id"` // Deprecated: Use SegmentId instead
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

	ClickedLinksWithOptions(ctx context.Context, broadcastId string, options *ListOptions) (ListBroadcastClickedLinksResponse, error)
	ClickedLinksWithContext(ctx context.Context, broadcastId string) (ListBroadcastClickedLinksResponse, error)
	ClickedLinks(broadcastId string) (ListBroadcastClickedLinksResponse, error)

	GetWithContext(ctx context.Context, broadcastId string) (Broadcast, error)
	Get(broadcastId string) (Broadcast, error)

	SendWithContext(ctx context.Context, params *SendBroadcastRequest) (SendBroadcastResponse, error)
	Send(params *SendBroadcastRequest) (SendBroadcastResponse, error)

	CancelWithContext(ctx context.Context, broadcastId string) (CancelBroadcastResponse, error)
	Cancel(broadcastId string) (CancelBroadcastResponse, error)

	RemoveWithContext(ctx context.Context, broadcastId string) (RemoveBroadcastResponse, error)
	Remove(broadcastId string) (RemoveBroadcastResponse, error)

	RecipientsWithContext(ctx context.Context, broadcastId string, options *ListBroadcastRecipientsOptions) (ListBroadcastRecipientsResponse, error)
	Recipients(broadcastId string, options *ListBroadcastRecipientsOptions) (ListBroadcastRecipientsResponse, error)
}

type BroadcastsSvcImpl struct {
	client *Client
}

// CreateWithContext creates a new Broadcast based on the given params
// https://resend.com/docs/api-reference/broadcasts/create-broadcast
func (s *BroadcastsSvcImpl) CreateWithContext(ctx context.Context, params *CreateBroadcastRequest) (CreateBroadcastResponse, error) {
	path := "/broadcasts"

	if params.SegmentId == "" && params.AudienceId == "" {
		return CreateBroadcastResponse{}, errors.New("[ERROR]: Either SegmentId or AudienceId must be provided")
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

func (s *BroadcastsSvcImpl) CancelWithContext(ctx context.Context, broadcastId string) (CancelBroadcastResponse, error) {
	if broadcastId == "" {
		return CancelBroadcastResponse{}, errors.New("[ERROR]: broadcastId cannot be empty")
	}

	path := "broadcasts/" + broadcastId + "/cancel"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return CancelBroadcastResponse{}, ErrFailedToCreateBroadcastCancelRequest
	}

	resp := new(CancelBroadcastResponse)

	// Send Request
	_, err = s.client.Perform(req, resp)

	if err != nil {
		return CancelBroadcastResponse{}, err
	}

	return *resp, nil
}

func (s *BroadcastsSvcImpl) Cancel(broadcastId string) (CancelBroadcastResponse, error) {
	return s.CancelWithContext(context.Background(), broadcastId)
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

// RecipientsWithContext returns a broadcast's recipients for a given event type, such as who
// opened, clicked, or bounced.
// https://resend.com/docs/api-reference/broadcasts/list-broadcast-recipients
func (s *BroadcastsSvcImpl) RecipientsWithContext(ctx context.Context, broadcastId string, options *ListBroadcastRecipientsOptions) (ListBroadcastRecipientsResponse, error) {
	if broadcastId == "" {
		return ListBroadcastRecipientsResponse{}, errors.New("[ERROR]: broadcastId cannot be empty")
	}

	if options == nil || options.Type == "" {
		return ListBroadcastRecipientsResponse{}, errors.New("[ERROR]: Type cannot be empty")
	}

	query := make(url.Values)
	query.Set("type", options.Type)
	if options.Email != "" {
		query.Set("email", options.Email)
	}
	if options.BounceType != "" {
		query.Set("bounce_type", options.BounceType)
	}
	if options.Limit != nil {
		query.Set("limit", fmt.Sprintf("%d", *options.Limit))
	}
	if options.After != nil {
		query.Set("after", *options.After)
	}
	if options.Before != nil {
		query.Set("before", *options.Before)
	}

	path := "broadcasts/" + broadcastId + "/recipients?" + query.Encode()

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListBroadcastRecipientsResponse{}, ErrFailedToCreateBroadcastRecipientsRequest
	}

	recipients := new(ListBroadcastRecipientsResponse)

	// Send Request
	_, err = s.client.Perform(req, recipients)

	if err != nil {
		return ListBroadcastRecipientsResponse{}, err
	}

	return *recipients, nil
}

// Recipients returns a broadcast's recipients for a given event type, such as who opened,
// clicked, or bounced.
func (s *BroadcastsSvcImpl) Recipients(broadcastId string, options *ListBroadcastRecipientsOptions) (ListBroadcastRecipientsResponse, error) {
	return s.RecipientsWithContext(context.Background(), broadcastId, options)
}

// ClickedLinksWithOptions returns the clicked links for a broadcast with pagination options
// https://resend.com/docs/api-reference/broadcasts/list-broadcast-clicked-links
func (s *BroadcastsSvcImpl) ClickedLinksWithOptions(ctx context.Context, broadcastId string, options *ListOptions) (ListBroadcastClickedLinksResponse, error) {
	if broadcastId == "" {
		return ListBroadcastClickedLinksResponse{}, errors.New("[ERROR]: broadcastId cannot be empty")
	}

	path := "broadcasts/" + broadcastId + "/clicked-links" + buildPaginationQuery(options)

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListBroadcastClickedLinksResponse{}, errors.New("[ERROR]: Failed to create Broadcasts.ClickedLinks request")
	}

	clickedLinks := new(ListBroadcastClickedLinksResponse)

	// Send Request
	_, err = s.client.Perform(req, clickedLinks)

	if err != nil {
		return ListBroadcastClickedLinksResponse{}, err
	}

	return *clickedLinks, nil
}

// ClickedLinksWithContext returns the clicked links for a broadcast
// https://resend.com/docs/api-reference/broadcasts/list-broadcast-clicked-links
func (s *BroadcastsSvcImpl) ClickedLinksWithContext(ctx context.Context, broadcastId string) (ListBroadcastClickedLinksResponse, error) {
	return s.ClickedLinksWithOptions(ctx, broadcastId, nil)
}

// ClickedLinks returns the clicked links for a broadcast
func (s *BroadcastsSvcImpl) ClickedLinks(broadcastId string) (ListBroadcastClickedLinksResponse, error) {
	return s.ClickedLinksWithContext(context.Background(), broadcastId)
}
