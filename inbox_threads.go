package resend

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// InboxThreadLabel is a label as attached to a thread.
type InboxThreadLabel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	// Color is a `#RRGGBB` hex color in uppercase.
	Color string `json:"color"`
}

// InboxThread is a thread as it appears in the list of an inbox's threads.
type InboxThread struct {
	Id      string  `json:"id"`
	Subject *string `json:"subject"`
	// From is the sender email address. It can include the display name.
	From *string `json:"from"`
	// To holds the recipients of the most recent email. For a received email, these are the
	// addresses it was delivered to, not the full `To` header.
	To []string `json:"to"`
	Cc []string `json:"cc"`
	// Bcc is always empty for a received email.
	Bcc    []string           `json:"bcc"`
	Labels []InboxThreadLabel `json:"labels"`
	// MessageCount is the number of emails in the thread.
	MessageCount int `json:"message_count"`
	// HasAttachment tells whether the most recent email in the thread has an attachment.
	HasAttachment bool `json:"has_attachment"`
	// HasDraft tells whether the thread has an unsent draft.
	HasDraft bool `json:"has_draft"`
	// Read is true only when every email in the thread is read.
	Read bool `json:"read"`
	// ReceivedAt is the timestamp when the thread was last active.
	ReceivedAt string `json:"received_at"`
	// Folder is the folder the thread is in. Every thread is in exactly one folder.
	Folder InboxFolder `json:"folder"`
}

// InboxThreadSearchHighlights holds where the search query matched the email in MatchedEmailId,
// as text fragments with each match wrapped in `**`. Only keys with a match are set.
type InboxThreadSearchHighlights struct {
	Subject     []string `json:"subject,omitempty"`
	Body        []string `json:"body,omitempty"`
	From        []string `json:"from,omitempty"`
	To          []string `json:"to,omitempty"`
	Cc          []string `json:"cc,omitempty"`
	Bcc         []string `json:"bcc,omitempty"`
	Attachments []string `json:"attachments,omitempty"`
}

// InboxThreadSearchResult is a thread as returned by Search.
type InboxThreadSearchResult struct {
	InboxThread
	// MatchedEmailId identifies the email in the thread that matched the search. When the query
	// has words or phrases, it is the most relevant match. Otherwise, it is the newest. The get,
	// reply and forward thread email methods accept it as the email ID, but it is not equal to the
	// ID those methods return. nil when the search index returns no matched email for the thread.
	MatchedEmailId *string `json:"matched_email_id"`
	// Highlights is empty when there is no query.
	Highlights InboxThreadSearchHighlights `json:"highlights"`
}

// InboxThreadSummary is the summary of a single thread, as returned by Get and Update.
type InboxThreadSummary struct {
	Object  string             `json:"object"`
	Id      string             `json:"id"`
	Subject *string            `json:"subject"`
	Folder  InboxFolder        `json:"folder"`
	Labels  []InboxThreadLabel `json:"labels"`
	// Read is true only when every email in the thread is read.
	Read bool `json:"read"`
}

// ListInboxThreadsOptions contains the filters and pagination for listing threads.
type ListInboxThreadsOptions struct {
	// Folders to list. A thread is returned when any of its emails is in one of the folders.
	// Defaults to inbox, or to inbox, archive and sent when Labels is set. Spam and trash are only
	// included when you name them.
	Folders []InboxFolder
	// Labels are label IDs. Returns threads with any of these labels. Up to 50 IDs. An ID that
	// does not exist in the inbox returns a 404.
	Labels []string
	// Read set to true returns threads where every email is read. False returns threads with at
	// least one unread email.
	Read  *bool
	Limit *int
	// After and Before take a thread ID. Only one of them can be passed.
	After  *string
	Before *string
}

// SearchInboxThreadsOptions contains the search terms, filters and pagination for searching
// threads. Different parameters combine, so a thread must match all of them. Folders, Labels and
// Read apply to the thread. Query, From, To, Cc, Bcc, HasAttachment, StartDate and EndDate must all
// match the same email in the thread.
type SearchInboxThreadsOptions struct {
	// Query is free text to search for in the subject, body, sender, recipients, and attachment
	// names. Every word must match. Wrap words in double quotes to match an exact phrase. Start a
	// word or quoted phrase with `-` to exclude it. Max 256 characters and 32 words.
	Query string
	// From holds senders. Each value matches part of the sender's address or name, without case.
	// Returns threads that match any of them. Up to 20 values.
	From []string
	// To holds recipients, matched the same way as From. Up to 20 values.
	To []string
	// Cc holds CC recipients, matched the same way as From. Up to 20 values.
	Cc []string
	// Bcc holds BCC recipients, matched the same way as From. Up to 20 values.
	Bcc []string
	// HasAttachment set to true matches emails with an attachment. False matches emails without
	// one.
	HasAttachment *bool
	// StartDate matches emails sent on or after this date. Accepts a date like `2026-09-01` or an
	// ISO 8601 timestamp, in UTC unless the timestamp has an offset.
	StartDate string
	// EndDate matches emails sent on or before this date. Same formats as StartDate. A date covers
	// that whole day. A timestamp is used as given.
	EndDate string
	Folders []InboxFolder
	Labels  []string
	Read    *bool
	Limit   *int
	After   *string
	Before  *string
}

type ListInboxThreadsResponse struct {
	Object  string `json:"object"`
	HasMore bool   `json:"has_more"`
	// Data holds the matching threads, most recently active first.
	Data []InboxThread `json:"data"`
}

type SearchInboxThreadsResponse struct {
	Object  string `json:"object"`
	HasMore bool   `json:"has_more"`
	// Data holds the matching threads, most recently active first.
	Data []InboxThreadSearchResult `json:"data"`
}

// UpdateInboxThreadRequest contains params for updating a thread. At least one field is required.
// A label ID cannot be both added and removed in the same request.
type UpdateInboxThreadRequest struct {
	// Read marks every email in the thread read or unread.
	Read *bool `json:"read,omitempty"`
	// Folder moves the thread to this folder. `sent` is assigned automatically and cannot be set.
	// Moving a thread out of trash restores it and clears its deletion deadline.
	Folder InboxMoveThreadFolder `json:"folder,omitempty"`
	// Deprecated: Use AddLabelIds instead. LabelId applies this label to the thread, the same as
	// one entry in AddLabelIds.
	LabelId string `json:"label_id,omitempty"`
	// AddLabelIds are label IDs to apply to the thread. Between 1 and 10 IDs. A label ID that does
	// not exist in the inbox returns a 404.
	AddLabelIds []string `json:"add_label_ids,omitempty"`
	// RemoveLabelIds are label IDs to remove from the thread. Between 1 and 10 IDs. A label ID
	// that does not exist in the inbox returns a 404.
	RemoveLabelIds []string `json:"remove_label_ids,omitempty"`
}

type RemoveInboxThreadResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
	// Deleted is true when this call moved the thread to trash. False when the thread was already
	// there.
	Deleted bool `json:"deleted"`
}

type InboxThreadsSvc interface {
	List(inboxId string, options *ListInboxThreadsOptions) (ListInboxThreadsResponse, error)
	ListWithContext(ctx context.Context, inboxId string, options *ListInboxThreadsOptions) (ListInboxThreadsResponse, error)
	Search(inboxId string, options *SearchInboxThreadsOptions) (SearchInboxThreadsResponse, error)
	SearchWithContext(ctx context.Context, inboxId string, options *SearchInboxThreadsOptions) (SearchInboxThreadsResponse, error)
	Get(inboxId, threadId string) (InboxThreadSummary, error)
	GetWithContext(ctx context.Context, inboxId, threadId string) (InboxThreadSummary, error)
	Update(inboxId, threadId string, params *UpdateInboxThreadRequest) (InboxThreadSummary, error)
	UpdateWithContext(ctx context.Context, inboxId, threadId string, params *UpdateInboxThreadRequest) (InboxThreadSummary, error)
	Remove(inboxId, threadId string) (RemoveInboxThreadResponse, error)
	RemoveWithContext(ctx context.Context, inboxId, threadId string) (RemoveInboxThreadResponse, error)
}

type InboxThreadsSvcImpl struct {
	client *Client
	Emails InboxThreadEmailsSvc
}

func setInboxCursorQuery(query url.Values, limit *int, after, before *string) {
	if limit != nil {
		query.Set("limit", fmt.Sprintf("%d", *limit))
	}
	if after != nil {
		query.Set("after", *after)
	}
	if before != nil {
		query.Set("before", *before)
	}
}

func setInboxListQuery(query url.Values, key string, values []string) {
	if len(values) > 0 {
		query.Set(key, strings.Join(values, ","))
	}
}

func setInboxBoolQuery(query url.Values, key string, value *bool) {
	if value != nil {
		query.Set(key, strconv.FormatBool(*value))
	}
}

func withQuery(path string, query url.Values) string {
	if len(query) > 0 {
		return path + "?" + query.Encode()
	}
	return path
}

// List returns the threads in an inbox, most recently active first. A thread is a conversation
// made up of one or more emails. To find threads by text, sender, recipient, attachment, or date,
// use Search.
// https://resend.com/docs/api-reference/inboxes/list-threads
func (s *InboxThreadsSvcImpl) List(inboxId string, options *ListInboxThreadsOptions) (ListInboxThreadsResponse, error) {
	return s.ListWithContext(context.Background(), inboxId, options)
}

// ListWithContext returns the threads in an inbox with context.
// https://resend.com/docs/api-reference/inboxes/list-threads
func (s *InboxThreadsSvcImpl) ListWithContext(ctx context.Context, inboxId string, options *ListInboxThreadsOptions) (ListInboxThreadsResponse, error) {
	if inboxId == "" {
		return ListInboxThreadsResponse{}, errors.New("[ERROR]: Inbox ID is required")
	}
	if options == nil {
		options = &ListInboxThreadsOptions{}
	}

	query := make(url.Values)
	setInboxListQuery(query, "folders", options.Folders)
	setInboxListQuery(query, "labels", options.Labels)
	setInboxBoolQuery(query, "read", options.Read)
	setInboxCursorQuery(query, options.Limit, options.After, options.Before)

	req, err := s.client.NewRequest(ctx, http.MethodGet, withQuery(inboxPath(inboxId, "threads"), query), nil)
	if err != nil {
		return ListInboxThreadsResponse{}, ErrFailedToCreateInboxesThreadsListRequest
	}

	resp := new(ListInboxThreadsResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return ListInboxThreadsResponse{}, err
	}
	return *resp, nil
}

// Search finds threads by the text of their emails, who sent or received them, whether they have
// attachments, and their dates. Results are ordered by latest activity, newest first, and paginate
// like List. Search can be up to about a minute behind new emails. When more than 10,000 threads
// match, search covers the 10,000 with the newest matching email.
// https://resend.com/docs/api-reference/inboxes/search-threads
func (s *InboxThreadsSvcImpl) Search(inboxId string, options *SearchInboxThreadsOptions) (SearchInboxThreadsResponse, error) {
	return s.SearchWithContext(context.Background(), inboxId, options)
}

// SearchWithContext finds threads in an inbox with context.
// https://resend.com/docs/api-reference/inboxes/search-threads
func (s *InboxThreadsSvcImpl) SearchWithContext(ctx context.Context, inboxId string, options *SearchInboxThreadsOptions) (SearchInboxThreadsResponse, error) {
	if inboxId == "" {
		return SearchInboxThreadsResponse{}, errors.New("[ERROR]: Inbox ID is required")
	}
	if options == nil {
		options = &SearchInboxThreadsOptions{}
	}

	query := make(url.Values)
	if options.Query != "" {
		query.Set("query", options.Query)
	}
	setInboxListQuery(query, "from", options.From)
	setInboxListQuery(query, "to", options.To)
	setInboxListQuery(query, "cc", options.Cc)
	setInboxListQuery(query, "bcc", options.Bcc)
	setInboxBoolQuery(query, "has_attachment", options.HasAttachment)
	if options.StartDate != "" {
		query.Set("start_date", options.StartDate)
	}
	if options.EndDate != "" {
		query.Set("end_date", options.EndDate)
	}
	setInboxListQuery(query, "folders", options.Folders)
	setInboxListQuery(query, "labels", options.Labels)
	setInboxBoolQuery(query, "read", options.Read)
	setInboxCursorQuery(query, options.Limit, options.After, options.Before)

	req, err := s.client.NewRequest(ctx, http.MethodGet, withQuery(inboxPath(inboxId, "threads", "search"), query), nil)
	if err != nil {
		return SearchInboxThreadsResponse{}, ErrFailedToCreateInboxesThreadsSearchRequest
	}

	resp := new(SearchInboxThreadsResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return SearchInboxThreadsResponse{}, err
	}
	return *resp, nil
}

// Get retrieves the thread summary. The emails in the thread are listed with Threads.Emails.List.
// https://resend.com/docs/api-reference/inboxes/get-thread
func (s *InboxThreadsSvcImpl) Get(inboxId, threadId string) (InboxThreadSummary, error) {
	return s.GetWithContext(context.Background(), inboxId, threadId)
}

// GetWithContext retrieves the thread summary with context.
// https://resend.com/docs/api-reference/inboxes/get-thread
func (s *InboxThreadsSvcImpl) GetWithContext(ctx context.Context, inboxId, threadId string) (InboxThreadSummary, error) {
	if inboxId == "" {
		return InboxThreadSummary{}, errors.New("[ERROR]: Inbox ID is required")
	}
	if threadId == "" {
		return InboxThreadSummary{}, errors.New("[ERROR]: Thread ID is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, inboxPath(inboxId, "threads", threadId), nil)
	if err != nil {
		return InboxThreadSummary{}, ErrFailedToCreateInboxesThreadsGetRequest
	}

	resp := new(InboxThreadSummary)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return InboxThreadSummary{}, err
	}
	return *resp, nil
}

// Update marks a thread read or unread, moves it to a folder, or adds and removes labels. At
// least one field is required. A label ID cannot be both added and removed in the same request.
// Moving a thread out of trash restores it and clears its deletion deadline.
// https://resend.com/docs/api-reference/inboxes/update-thread
func (s *InboxThreadsSvcImpl) Update(inboxId, threadId string, params *UpdateInboxThreadRequest) (InboxThreadSummary, error) {
	return s.UpdateWithContext(context.Background(), inboxId, threadId, params)
}

// UpdateWithContext updates a thread with context.
// https://resend.com/docs/api-reference/inboxes/update-thread
func (s *InboxThreadsSvcImpl) UpdateWithContext(ctx context.Context, inboxId, threadId string, params *UpdateInboxThreadRequest) (InboxThreadSummary, error) {
	if inboxId == "" {
		return InboxThreadSummary{}, errors.New("[ERROR]: Inbox ID is required")
	}
	if threadId == "" {
		return InboxThreadSummary{}, errors.New("[ERROR]: Thread ID is required")
	}
	if params == nil || (params.Read == nil && params.Folder == "" && params.LabelId == "" &&
		len(params.AddLabelIds) == 0 && len(params.RemoveLabelIds) == 0) {
		return InboxThreadSummary{}, errors.New("[ERROR]: At least one of Read, Folder, LabelId, AddLabelIds or RemoveLabelIds is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPatch, inboxPath(inboxId, "threads", threadId), params)
	if err != nil {
		return InboxThreadSummary{}, ErrFailedToCreateInboxesThreadsUpdateRequest
	}

	resp := new(InboxThreadSummary)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return InboxThreadSummary{}, err
	}
	return *resp, nil
}

// Remove moves the thread to trash and schedules it for deletion in 30 days. This is reversible.
// To restore the thread, update its Folder to inbox, archive, or spam, which also clears the
// deletion deadline.
// https://resend.com/docs/api-reference/inboxes/delete-thread
func (s *InboxThreadsSvcImpl) Remove(inboxId, threadId string) (RemoveInboxThreadResponse, error) {
	return s.RemoveWithContext(context.Background(), inboxId, threadId)
}

// RemoveWithContext moves the thread to trash with context.
// https://resend.com/docs/api-reference/inboxes/delete-thread
func (s *InboxThreadsSvcImpl) RemoveWithContext(ctx context.Context, inboxId, threadId string) (RemoveInboxThreadResponse, error) {
	if inboxId == "" {
		return RemoveInboxThreadResponse{}, errors.New("[ERROR]: Inbox ID is required")
	}
	if threadId == "" {
		return RemoveInboxThreadResponse{}, errors.New("[ERROR]: Thread ID is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, inboxPath(inboxId, "threads", threadId), nil)
	if err != nil {
		return RemoveInboxThreadResponse{}, ErrFailedToCreateInboxesThreadsRemoveRequest
	}

	resp := new(RemoveInboxThreadResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return RemoveInboxThreadResponse{}, err
	}
	return *resp, nil
}
