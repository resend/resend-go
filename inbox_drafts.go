package resend

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// InboxDraftType tells whether a draft replies in an existing thread or starts a new one.
type InboxDraftType = string

const (
	InboxDraftTypeStandalone InboxDraftType = "standalone"
	InboxDraftTypeReply      InboxDraftType = "reply"
)

// InboxDraft is an inbox draft as returned by Get and Update.
type InboxDraft struct {
	Object string `json:"object"`
	Id     string `json:"id"`
	// Type is reply when the draft replies to an email in an existing thread, and standalone when
	// the draft starts a new thread.
	Type InboxDraftType `json:"type"`
	// To holds the recipients. nil or empty when no recipients are set.
	To      []string `json:"to"`
	Cc      []string `json:"cc"`
	Bcc     []string `json:"bcc"`
	Subject *string  `json:"subject"`
	Html    *string  `json:"html"`
	Text    *string  `json:"text"`
	// ThreadId is the thread ID when the draft is a reply. nil otherwise.
	ThreadId *string `json:"thread_id"`
	// ReplyToEmailId is the email ID the draft replies to. nil otherwise.
	ReplyToEmailId *string `json:"reply_to_email_id"`
	// EmailId is the ID of the outbound email queued when the draft was sent. nil until the draft
	// is sent.
	EmailId   *string `json:"email_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// InboxDraftListItem is an inbox draft as returned by List.
type InboxDraftListItem struct {
	Id      string         `json:"id"`
	Type    InboxDraftType `json:"type"`
	To      []string       `json:"to"`
	Cc      []string       `json:"cc"`
	Bcc     []string       `json:"bcc"`
	Subject *string        `json:"subject"`
	// Snippet is a short preview of the body.
	Snippet        *string `json:"snippet"`
	ThreadId       *string `json:"thread_id"`
	ReplyToEmailId *string `json:"reply_to_email_id"`
	UpdatedAt      string  `json:"updated_at"`
}

type ListInboxDraftsResponse struct {
	Object  string `json:"object"`
	HasMore bool   `json:"has_more"`
	// Data holds the unsent drafts on this page, most recently updated first.
	Data []InboxDraftListItem `json:"data"`
}

// CreateInboxDraftRequest contains params for creating a draft. At least one of To, Cc, Bcc,
// Subject, Text or Html must be non-empty. The combined recipients of To, Cc and Bcc cannot exceed
// 50. ThreadId and ReplyToEmailId must be sent together or not at all.
type CreateInboxDraftRequest struct {
	// To holds the recipients. Each value can include a friendly name, for example
	// `Ada Lovelace <ada@example.org>`. A nil slice is omitted, so a reply draft is addressed to
	// the sender of the email it replies to, or to its Reply-To when set. For an email the inbox
	// sent, it defaults to that email's `to` recipients. An empty, non-nil slice leaves it empty.
	To  []string
	Cc  []string
	Bcc []string
	// Subject is max 2000 characters. A string that is only whitespace, or that contains a line
	// break, is rejected.
	Subject string
	Html    string
	Text    string
	// ThreadId is the thread ID to reply in. It must be sent with ReplyToEmailId.
	ThreadId string
	// ReplyToEmailId is the email ID to reply to, as returned in Data[].Id by
	// Threads.Emails.List. It must be sent with ThreadId.
	ReplyToEmailId string
}

func (r CreateInboxDraftRequest) MarshalJSON() ([]byte, error) {
	aux := make(map[string]any)
	if r.To != nil {
		aux["to"] = r.To
	}
	if r.Cc != nil {
		aux["cc"] = r.Cc
	}
	if r.Bcc != nil {
		aux["bcc"] = r.Bcc
	}
	if r.Subject != "" {
		aux["subject"] = r.Subject
	}
	if r.Html != "" {
		aux["html"] = r.Html
	}
	if r.Text != "" {
		aux["text"] = r.Text
	}
	if r.ThreadId != "" {
		aux["thread_id"] = r.ThreadId
	}
	if r.ReplyToEmailId != "" {
		aux["reply_to_email_id"] = r.ReplyToEmailId
	}
	return json.Marshal(aux)
}

type CreateInboxDraftResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

// UpdateInboxDraftRequest contains params for updating a draft. Send at least one field. Omitted
// fields are left unchanged. Setting Text or Html replaces the whole body, so the body part you
// omit becomes null. The merged draft must contain at least one non-empty field, and its combined
// recipients cannot exceed 50.
type UpdateInboxDraftRequest struct {
	// To holds the recipients. Each value can include a friendly name. A nil slice is omitted. An
	// empty, non-nil slice clears it.
	To []string
	// Cc holds the CC recipients. A nil slice is omitted. An empty, non-nil slice clears it.
	Cc []string
	// Bcc holds the BCC recipients. A nil slice is omitted. An empty, non-nil slice clears it.
	Bcc []string
	// Subject is max 2000 characters and must not contain a line break. Use ClearSubject to
	// clear it.
	Subject *string
	// Html is the HTML body. Use ClearHtml to clear it.
	Html *string
	// Text is the plain text body. Use ClearText to clear it.
	Text *string

	subjectCleared bool
	htmlCleared    bool
	textCleared    bool
}

// ClearSubject sends `subject: null`, which clears the subject.
func (r *UpdateInboxDraftRequest) ClearSubject() {
	r.Subject = nil
	r.subjectCleared = true
}

// ClearHtml sends `html: null`, which clears the HTML body.
func (r *UpdateInboxDraftRequest) ClearHtml() {
	r.Html = nil
	r.htmlCleared = true
}

// ClearText sends `text: null`, which clears the plain text body.
func (r *UpdateInboxDraftRequest) ClearText() {
	r.Text = nil
	r.textCleared = true
}

func (r UpdateInboxDraftRequest) isEmpty() bool {
	return r.To == nil && r.Cc == nil && r.Bcc == nil &&
		r.Subject == nil && !r.subjectCleared &&
		r.Html == nil && !r.htmlCleared &&
		r.Text == nil && !r.textCleared
}

func (r UpdateInboxDraftRequest) MarshalJSON() ([]byte, error) {
	aux := make(map[string]any)
	if r.To != nil {
		aux["to"] = r.To
	}
	if r.Cc != nil {
		aux["cc"] = r.Cc
	}
	if r.Bcc != nil {
		aux["bcc"] = r.Bcc
	}
	setNullableString(aux, "subject", r.Subject, r.subjectCleared)
	setNullableString(aux, "html", r.Html, r.htmlCleared)
	setNullableString(aux, "text", r.Text, r.textCleared)
	return json.Marshal(aux)
}

type RemoveInboxDraftResponse struct {
	Object  string `json:"object"`
	Id      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

type SendInboxDraftResponse struct {
	// Id is the ID of the sent email. Use it with the emails methods.
	Id string `json:"id"`
}

type InboxDraftsSvc interface {
	List(inboxId string, options *ListOptions) (ListInboxDraftsResponse, error)
	ListWithContext(ctx context.Context, inboxId string, options *ListOptions) (ListInboxDraftsResponse, error)
	Create(inboxId string, params *CreateInboxDraftRequest) (CreateInboxDraftResponse, error)
	CreateWithContext(ctx context.Context, inboxId string, params *CreateInboxDraftRequest) (CreateInboxDraftResponse, error)
	Get(inboxId, draftId string) (InboxDraft, error)
	GetWithContext(ctx context.Context, inboxId, draftId string) (InboxDraft, error)
	Update(inboxId, draftId string, params *UpdateInboxDraftRequest) (InboxDraft, error)
	UpdateWithContext(ctx context.Context, inboxId, draftId string, params *UpdateInboxDraftRequest) (InboxDraft, error)
	Remove(inboxId, draftId string) (RemoveInboxDraftResponse, error)
	RemoveWithContext(ctx context.Context, inboxId, draftId string) (RemoveInboxDraftResponse, error)
	Send(inboxId, draftId string) (SendInboxDraftResponse, error)
	SendWithContext(ctx context.Context, inboxId, draftId string) (SendInboxDraftResponse, error)
}

type InboxDraftsSvcImpl struct {
	client *Client
}

func validateInboxDraftIds(inboxId, draftId string) error {
	if inboxId == "" {
		return errors.New("[ERROR]: Inbox ID is required")
	}
	if draftId == "" {
		return errors.New("[ERROR]: Draft ID is required")
	}
	return nil
}

// List returns the unsent drafts on the inbox, most recently updated first. The After and Before
// cursors take a draft ID from the current page. A cursor that is not an unsent draft on this
// inbox returns a 404.
// https://resend.com/docs/api-reference/inboxes/list-drafts
func (s *InboxDraftsSvcImpl) List(inboxId string, options *ListOptions) (ListInboxDraftsResponse, error) {
	return s.ListWithContext(context.Background(), inboxId, options)
}

// ListWithContext returns the unsent drafts on the inbox with context.
// https://resend.com/docs/api-reference/inboxes/list-drafts
func (s *InboxDraftsSvcImpl) ListWithContext(ctx context.Context, inboxId string, options *ListOptions) (ListInboxDraftsResponse, error) {
	if inboxId == "" {
		return ListInboxDraftsResponse{}, errors.New("[ERROR]: Inbox ID is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, inboxPath(inboxId, "drafts")+buildPaginationQuery(options), nil)
	if err != nil {
		return ListInboxDraftsResponse{}, ErrFailedToCreateInboxesDraftsListRequest
	}

	resp := new(ListInboxDraftsResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return ListInboxDraftsResponse{}, err
	}
	return *resp, nil
}

// Create creates an unsent draft on the inbox. Omit ThreadId and ReplyToEmailId for a draft that
// starts a new thread. Pass both together for a reply draft in an existing thread. Only one unsent
// reply draft can target an email, so a reply draft request for an email that already has one
// returns the ID of the existing draft.
// https://resend.com/docs/api-reference/inboxes/create-draft
func (s *InboxDraftsSvcImpl) Create(inboxId string, params *CreateInboxDraftRequest) (CreateInboxDraftResponse, error) {
	return s.CreateWithContext(context.Background(), inboxId, params)
}

// CreateWithContext creates an unsent draft on the inbox with context.
// https://resend.com/docs/api-reference/inboxes/create-draft
func (s *InboxDraftsSvcImpl) CreateWithContext(ctx context.Context, inboxId string, params *CreateInboxDraftRequest) (CreateInboxDraftResponse, error) {
	if inboxId == "" {
		return CreateInboxDraftResponse{}, errors.New("[ERROR]: Inbox ID is required")
	}
	if params == nil {
		return CreateInboxDraftResponse{}, errors.New("[ERROR]: Draft params are required")
	}
	if (params.ThreadId == "") != (params.ReplyToEmailId == "") {
		return CreateInboxDraftResponse{}, errors.New("[ERROR]: ThreadId and ReplyToEmailId must be provided together")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, inboxPath(inboxId, "drafts"), params)
	if err != nil {
		return CreateInboxDraftResponse{}, ErrFailedToCreateInboxesDraftsCreateRequest
	}

	resp := new(CreateInboxDraftResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return CreateInboxDraftResponse{}, err
	}
	return *resp, nil
}

// Get retrieves a single draft.
// https://resend.com/docs/api-reference/inboxes/get-draft
func (s *InboxDraftsSvcImpl) Get(inboxId, draftId string) (InboxDraft, error) {
	return s.GetWithContext(context.Background(), inboxId, draftId)
}

// GetWithContext retrieves a single draft with context.
// https://resend.com/docs/api-reference/inboxes/get-draft
func (s *InboxDraftsSvcImpl) GetWithContext(ctx context.Context, inboxId, draftId string) (InboxDraft, error) {
	if err := validateInboxDraftIds(inboxId, draftId); err != nil {
		return InboxDraft{}, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, inboxPath(inboxId, "drafts", draftId), nil)
	if err != nil {
		return InboxDraft{}, ErrFailedToCreateInboxesDraftsGetRequest
	}

	resp := new(InboxDraft)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return InboxDraft{}, err
	}
	return *resp, nil
}

// Update updates the content of an unsent draft. Omitted fields are left unchanged. A request that
// sets Text or Html replaces the whole body. The body part you omit from that request becomes
// null. A draft that was already sent, or that changed during the request, returns a 409.
// https://resend.com/docs/api-reference/inboxes/update-draft
func (s *InboxDraftsSvcImpl) Update(inboxId, draftId string, params *UpdateInboxDraftRequest) (InboxDraft, error) {
	return s.UpdateWithContext(context.Background(), inboxId, draftId, params)
}

// UpdateWithContext updates the content of an unsent draft with context.
// https://resend.com/docs/api-reference/inboxes/update-draft
func (s *InboxDraftsSvcImpl) UpdateWithContext(ctx context.Context, inboxId, draftId string, params *UpdateInboxDraftRequest) (InboxDraft, error) {
	if err := validateInboxDraftIds(inboxId, draftId); err != nil {
		return InboxDraft{}, err
	}
	if params == nil || params.isEmpty() {
		return InboxDraft{}, errors.New("[ERROR]: At least one of To, Cc, Bcc, Subject, Text or Html is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPatch, inboxPath(inboxId, "drafts", draftId), params)
	if err != nil {
		return InboxDraft{}, ErrFailedToCreateInboxesDraftsUpdateRequest
	}

	resp := new(InboxDraft)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return InboxDraft{}, err
	}
	return *resp, nil
}

// Remove discards an unsent draft. A draft that was already sent, or that changed during the
// request, returns a 409.
// https://resend.com/docs/api-reference/inboxes/delete-draft
func (s *InboxDraftsSvcImpl) Remove(inboxId, draftId string) (RemoveInboxDraftResponse, error) {
	return s.RemoveWithContext(context.Background(), inboxId, draftId)
}

// RemoveWithContext discards an unsent draft with context.
// https://resend.com/docs/api-reference/inboxes/delete-draft
func (s *InboxDraftsSvcImpl) RemoveWithContext(ctx context.Context, inboxId, draftId string) (RemoveInboxDraftResponse, error) {
	if err := validateInboxDraftIds(inboxId, draftId); err != nil {
		return RemoveInboxDraftResponse{}, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, inboxPath(inboxId, "drafts", draftId), nil)
	if err != nil {
		return RemoveInboxDraftResponse{}, ErrFailedToCreateInboxesDraftsRemoveRequest
	}

	resp := new(RemoveInboxDraftResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return RemoveInboxDraftResponse{}, err
	}
	return *resp, nil
}

// Send sends the draft from the inbox address. A reply draft is sent in its thread. A draft without
// a thread starts a new thread. The draft must have at least one `to` recipient, or the request
// returns a 422. A draft that was already sent returns a 409.
// https://resend.com/docs/api-reference/inboxes/send-draft
func (s *InboxDraftsSvcImpl) Send(inboxId, draftId string) (SendInboxDraftResponse, error) {
	return s.SendWithContext(context.Background(), inboxId, draftId)
}

// SendWithContext sends the draft with context.
// https://resend.com/docs/api-reference/inboxes/send-draft
func (s *InboxDraftsSvcImpl) SendWithContext(ctx context.Context, inboxId, draftId string) (SendInboxDraftResponse, error) {
	if err := validateInboxDraftIds(inboxId, draftId); err != nil {
		return SendInboxDraftResponse{}, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, inboxPath(inboxId, "drafts", draftId, "send"), nil)
	if err != nil {
		return SendInboxDraftResponse{}, ErrFailedToCreateInboxesDraftsSendRequest
	}

	resp := new(SendInboxDraftResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return SendInboxDraftResponse{}, err
	}
	return *resp, nil
}
