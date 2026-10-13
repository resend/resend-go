package resend

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

// InboxThreadEmailDirection tells whether the inbox received or sent a thread email.
type InboxThreadEmailDirection = string

const (
	InboxThreadEmailDirectionInbound  InboxThreadEmailDirection = "inbound"
	InboxThreadEmailDirectionOutbound InboxThreadEmailDirection = "outbound"
)

type InboxThreadEmailAttachment struct {
	Id       string  `json:"id"`
	Filename *string `json:"filename"`
	// Size is the size of the attachment in bytes.
	Size *int `json:"size"`
}

// InboxThreadEmail is a message in an inbox thread.
type InboxThreadEmail struct {
	Id string `json:"id"`
	// Direction is inbound when the inbox received the message and outbound when it sent it. A
	// message the inbox sent to itself is inbound, so a reply to it goes to its sender.
	Direction InboxThreadEmailDirection `json:"direction"`
	// From is the sender address. It includes the display name when the sender has one, in the
	// format "Name <sender@domain.com>".
	From string `json:"from"`
	// To holds the recipients. For a received message, these are the addresses it was delivered
	// to, not the full `To` header.
	To []string `json:"to"`
	Cc []string `json:"cc"`
	// Bcc is always empty for a received message.
	Bcc       []string `json:"bcc"`
	ReplyTo   []string `json:"reply_to"`
	Subject   *string  `json:"subject"`
	MessageId *string  `json:"message_id"`
	Html      *string  `json:"html"`
	Text      *string  `json:"text"`
	// Attachments holds the attachments of the message.
	Attachments []InboxThreadEmailAttachment `json:"attachments"`
	Read        bool                         `json:"read"`
	// ReceivedAt is the timestamp when the inbox received the message, or when the inbox sent it.
	ReceivedAt string `json:"received_at"`
}

type ListInboxThreadEmailsResponse struct {
	Object  string `json:"object"`
	HasMore bool   `json:"has_more"`
	// Data holds the messages of the thread on this page, oldest first.
	Data []InboxThreadEmail `json:"data"`
}

// ReplyInboxThreadEmailRequest contains params for replying to a thread email. At least one of
// Html or Text must contain more than whitespace.
type ReplyInboxThreadEmailRequest struct {
	// Cc holds CC recipients. Without ReplyAll, the reply copies only these addresses. With
	// ReplyAll, these addresses are merged with the copied recipients, without duplicates. Cc and
	// Bcc combined must be fewer than 50 addresses. After the `to` recipients are added, the reply
	// cannot exceed 50 recipients.
	Cc []string `json:"cc,omitempty"`
	// Bcc holds BCC recipients. BCC recipients of the original message are never copied.
	Bcc  []string `json:"bcc,omitempty"`
	Html string   `json:"html,omitempty"`
	Text string   `json:"text,omitempty"`
	// Subject of the reply. It cannot contain a line break. When omitted, the thread subject is
	// used, prefixed with `Re:` if it does not already start with it.
	Subject string `json:"subject,omitempty"`
	// ReplyAll replies to everyone on the original message. For a message the inbox received, the
	// reply goes to its Reply-To addresses, or to its sender when it has none, and copies the other
	// addresses of its `To` header and its `cc` recipients in `cc`. For a message the inbox sent,
	// the reply goes to that message's Reply-To addresses when they are not the inbox address,
	// otherwise to its `to` recipients, and copies its `cc` recipients. The inbox address is always
	// excluded, and duplicates are removed.
	ReplyAll bool `json:"reply_all,omitempty"`
	// Attachments are files to attach, as on Emails.Send. Each one needs Content or Path. Up to 10
	// files. The Content of all files can total at most 25 MB. A Path file counts toward the 10
	// files only, because it is fetched when the email is sent.
	Attachments []*Attachment `json:"attachments,omitempty"`
	// Headers are custom headers to add to the email. The In-Reply-To, References and Reply-To
	// headers are set by Resend and cannot be set here. Sending one of them, in any letter case,
	// returns a 422.
	Headers map[string]string `json:"headers,omitempty"`
}

// ReplyInboxThreadEmailOptions is the additional options struct for the Inboxes.Threads.Emails.Reply call.
type ReplyInboxThreadEmailOptions struct {
	// IdempotencyKey is a unique key for the request, so the reply is sent only once. 1 to 256
	// characters. A retry with the same key within 24 hours returns the first response, with the
	// `idempotent-replayed: true` header. Reusing the key with a different body, or while the first
	// request is still running, returns a 409.
	// https://resend.com/docs/dashboard/emails/idempotency-keys
	IdempotencyKey string `json:"idempotency_key,omitempty"`
}

// GetIdempotencyKey returns the idempotency key for the Inboxes.Threads.Emails.Reply request.
func (o ReplyInboxThreadEmailOptions) GetIdempotencyKey() string {
	return o.IdempotencyKey
}

type ReplyInboxThreadEmailResponse struct {
	// Id is the ID of the sent email. Use it with the emails methods.
	Id string `json:"id"`
}

// ForwardInboxThreadEmailRequest contains params for forwarding a thread email.
type ForwardInboxThreadEmailRequest struct {
	// To holds the recipients. It must name at least one address. To, Cc and Bcc combined cannot
	// exceed 50 recipients.
	To []string `json:"to"`
	// Cc holds CC recipients. They count toward the 50-recipient limit.
	Cc []string `json:"cc,omitempty"`
	// Bcc holds BCC recipients. They count toward the 50-recipient limit.
	Bcc []string `json:"bcc,omitempty"`
	// Html is an HTML note sent above the quoted original.
	Html string `json:"html,omitempty"`
	// Text is a plain-text note sent above the quoted original.
	Text string `json:"text,omitempty"`
	// Subject of the forward. It cannot contain a line break. When omitted, the thread subject is
	// used, prefixed with `Fwd:` if it does not already start with it.
	Subject string `json:"subject,omitempty"`
	// Attachments are files to attach, as on Emails.Send. Each one needs Content or Path. They
	// share one limit with the attachments of the original message, which the forward also sends.
	// Together they can have up to 10 files and at most 25 MB of Content. A Path file counts
	// toward the 10 files only, because it is fetched when the email is sent.
	Attachments []*Attachment `json:"attachments,omitempty"`
	// Headers are custom headers to add to the email. The In-Reply-To, References and Reply-To
	// headers are set by Resend and cannot be set here. Sending one of them, in any letter case,
	// returns a 422.
	Headers map[string]string `json:"headers,omitempty"`
}

// ForwardInboxThreadEmailOptions is the additional options struct for the Inboxes.Threads.Emails.Forward call.
type ForwardInboxThreadEmailOptions struct {
	// IdempotencyKey is a unique key for the request, so the forward is sent only once. 1 to 256
	// characters. A retry with the same key within 24 hours returns the first response, with the
	// `idempotent-replayed: true` header. Reusing the key with a different body, or while the first
	// request is still running, returns a 409.
	// https://resend.com/docs/dashboard/emails/idempotency-keys
	IdempotencyKey string `json:"idempotency_key,omitempty"`
}

// GetIdempotencyKey returns the idempotency key for the Inboxes.Threads.Emails.Forward request.
func (o ForwardInboxThreadEmailOptions) GetIdempotencyKey() string {
	return o.IdempotencyKey
}

type ForwardInboxThreadEmailResponse struct {
	// Id is the ID of the sent email. Use it with the emails methods.
	Id string `json:"id"`
}

type InboxThreadEmailsSvc interface {
	List(inboxId, threadId string, options *ListOptions) (ListInboxThreadEmailsResponse, error)
	ListWithContext(ctx context.Context, inboxId, threadId string, options *ListOptions) (ListInboxThreadEmailsResponse, error)
	Get(inboxId, threadId, emailId string) (InboxThreadEmail, error)
	GetWithContext(ctx context.Context, inboxId, threadId, emailId string) (InboxThreadEmail, error)
	Reply(inboxId, threadId, emailId string, params *ReplyInboxThreadEmailRequest) (ReplyInboxThreadEmailResponse, error)
	ReplyWithContext(ctx context.Context, inboxId, threadId, emailId string, params *ReplyInboxThreadEmailRequest) (ReplyInboxThreadEmailResponse, error)
	ReplyWithOptions(ctx context.Context, inboxId, threadId, emailId string, params *ReplyInboxThreadEmailRequest, options *ReplyInboxThreadEmailOptions) (ReplyInboxThreadEmailResponse, error)
	Forward(inboxId, threadId, emailId string, params *ForwardInboxThreadEmailRequest) (ForwardInboxThreadEmailResponse, error)
	ForwardWithContext(ctx context.Context, inboxId, threadId, emailId string, params *ForwardInboxThreadEmailRequest) (ForwardInboxThreadEmailResponse, error)
	ForwardWithOptions(ctx context.Context, inboxId, threadId, emailId string, params *ForwardInboxThreadEmailRequest, options *ForwardInboxThreadEmailOptions) (ForwardInboxThreadEmailResponse, error)
}

type InboxThreadEmailsSvcImpl struct {
	client *Client
}

func validateInboxThreadEmailIds(inboxId, threadId, emailId string) error {
	if inboxId == "" {
		return errors.New("[ERROR]: Inbox ID is required")
	}
	if threadId == "" {
		return errors.New("[ERROR]: Thread ID is required")
	}
	if emailId == "" {
		return errors.New("[ERROR]: Email ID is required")
	}
	return nil
}

// List returns the messages of a thread, oldest first. Pass the ID of the last message as After
// to get the next, newer page. Pass the ID of the first message as Before to get the previous,
// older page. You cannot combine After and Before. A cursor that is not a message of the thread
// returns a 404.
// https://resend.com/docs/api-reference/inboxes/list-thread-emails
func (s *InboxThreadEmailsSvcImpl) List(inboxId, threadId string, options *ListOptions) (ListInboxThreadEmailsResponse, error) {
	return s.ListWithContext(context.Background(), inboxId, threadId, options)
}

// ListWithContext returns the messages of a thread with context.
// https://resend.com/docs/api-reference/inboxes/list-thread-emails
func (s *InboxThreadEmailsSvcImpl) ListWithContext(ctx context.Context, inboxId, threadId string, options *ListOptions) (ListInboxThreadEmailsResponse, error) {
	if inboxId == "" {
		return ListInboxThreadEmailsResponse{}, errors.New("[ERROR]: Inbox ID is required")
	}
	if threadId == "" {
		return ListInboxThreadEmailsResponse{}, errors.New("[ERROR]: Thread ID is required")
	}

	path := inboxPath(inboxId, "threads", threadId, "emails") + buildPaginationQuery(options)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListInboxThreadEmailsResponse{}, ErrFailedToCreateInboxesThreadEmailsListRequest
	}

	resp := new(ListInboxThreadEmailsResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return ListInboxThreadEmailsResponse{}, err
	}
	return *resp, nil
}

// Get retrieves a single message of a thread. The emailId param is the ID returned in Data[].Id by
// List.
// https://resend.com/docs/api-reference/inboxes/get-thread-email
func (s *InboxThreadEmailsSvcImpl) Get(inboxId, threadId, emailId string) (InboxThreadEmail, error) {
	return s.GetWithContext(context.Background(), inboxId, threadId, emailId)
}

// GetWithContext retrieves a single message of a thread with context.
// https://resend.com/docs/api-reference/inboxes/get-thread-email
func (s *InboxThreadEmailsSvcImpl) GetWithContext(ctx context.Context, inboxId, threadId, emailId string) (InboxThreadEmail, error) {
	if err := validateInboxThreadEmailIds(inboxId, threadId, emailId); err != nil {
		return InboxThreadEmail{}, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, inboxPath(inboxId, "threads", threadId, "emails", emailId), nil)
	if err != nil {
		return InboxThreadEmail{}, ErrFailedToCreateInboxesThreadEmailsGetRequest
	}

	resp := new(InboxThreadEmail)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return InboxThreadEmail{}, err
	}
	return *resp, nil
}

// Reply sends a reply from the inbox address right away, with no stored draft. The `to`
// recipients come from the message you reply to. For a message the inbox received, the reply goes
// to its Reply-To addresses, or to its sender when it has no Reply-To. For a message the inbox
// sent, the reply goes to that message's `to` recipients. Set ReplyAll to also copy the other
// recipients of the message. The reply threads under the message through the In-Reply-To and
// References headers.
// https://resend.com/docs/api-reference/inboxes/reply-thread-email
func (s *InboxThreadEmailsSvcImpl) Reply(inboxId, threadId, emailId string, params *ReplyInboxThreadEmailRequest) (ReplyInboxThreadEmailResponse, error) {
	return s.ReplyWithContext(context.Background(), inboxId, threadId, emailId, params)
}

// ReplyWithContext sends a reply to a thread email with context.
// https://resend.com/docs/api-reference/inboxes/reply-thread-email
func (s *InboxThreadEmailsSvcImpl) ReplyWithContext(ctx context.Context, inboxId, threadId, emailId string, params *ReplyInboxThreadEmailRequest) (ReplyInboxThreadEmailResponse, error) {
	return s.ReplyWithOptions(ctx, inboxId, threadId, emailId, params, nil)
}

// ReplyWithOptions sends a reply to a thread email with context and additional options, such as
// an idempotency key.
// https://resend.com/docs/api-reference/inboxes/reply-thread-email
func (s *InboxThreadEmailsSvcImpl) ReplyWithOptions(ctx context.Context, inboxId, threadId, emailId string, params *ReplyInboxThreadEmailRequest, options *ReplyInboxThreadEmailOptions) (ReplyInboxThreadEmailResponse, error) {
	if err := validateInboxThreadEmailIds(inboxId, threadId, emailId); err != nil {
		return ReplyInboxThreadEmailResponse{}, err
	}
	if params == nil || (strings.TrimSpace(params.Html) == "" && strings.TrimSpace(params.Text) == "") {
		return ReplyInboxThreadEmailResponse{}, errors.New("[ERROR]: At least one of Html or Text is required")
	}

	path := inboxPath(inboxId, "threads", threadId, "emails", emailId, "reply")

	var opts Options
	if options != nil {
		opts = options
	}

	req, err := s.client.NewRequestWithOptions(ctx, http.MethodPost, path, params, opts)
	if err != nil {
		return ReplyInboxThreadEmailResponse{}, ErrFailedToCreateInboxesThreadEmailsReplyRequest
	}

	resp := new(ReplyInboxThreadEmailResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return ReplyInboxThreadEmailResponse{}, err
	}
	return *resp, nil
}

// Forward forwards a message from the inbox address right away, with no stored draft. The original
// message is quoted under a forwarded-message banner, below the optional Html or Text note. The
// attachments of the original message are sent with the forward, together with any Attachments
// you add. The request fails when they exceed 10 files or 25 MB of content in total. An original
// attachment that was sent from a URL (Path) cannot be forwarded, and the request returns a 422.
// https://resend.com/docs/api-reference/inboxes/forward-thread-email
func (s *InboxThreadEmailsSvcImpl) Forward(inboxId, threadId, emailId string, params *ForwardInboxThreadEmailRequest) (ForwardInboxThreadEmailResponse, error) {
	return s.ForwardWithContext(context.Background(), inboxId, threadId, emailId, params)
}

// ForwardWithContext forwards a thread email with context.
// https://resend.com/docs/api-reference/inboxes/forward-thread-email
func (s *InboxThreadEmailsSvcImpl) ForwardWithContext(ctx context.Context, inboxId, threadId, emailId string, params *ForwardInboxThreadEmailRequest) (ForwardInboxThreadEmailResponse, error) {
	return s.ForwardWithOptions(ctx, inboxId, threadId, emailId, params, nil)
}

// ForwardWithOptions forwards a thread email with context and additional options, such as an
// idempotency key.
// https://resend.com/docs/api-reference/inboxes/forward-thread-email
func (s *InboxThreadEmailsSvcImpl) ForwardWithOptions(ctx context.Context, inboxId, threadId, emailId string, params *ForwardInboxThreadEmailRequest, options *ForwardInboxThreadEmailOptions) (ForwardInboxThreadEmailResponse, error) {
	if err := validateInboxThreadEmailIds(inboxId, threadId, emailId); err != nil {
		return ForwardInboxThreadEmailResponse{}, err
	}
	if params == nil || len(params.To) == 0 {
		return ForwardInboxThreadEmailResponse{}, errors.New("[ERROR]: To is required")
	}

	path := inboxPath(inboxId, "threads", threadId, "emails", emailId, "forward")

	var opts Options
	if options != nil {
		opts = options
	}

	req, err := s.client.NewRequestWithOptions(ctx, http.MethodPost, path, params, opts)
	if err != nil {
		return ForwardInboxThreadEmailResponse{}, ErrFailedToCreateInboxesThreadEmailsForwardRequest
	}

	resp := new(ForwardInboxThreadEmailResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return ForwardInboxThreadEmailResponse{}, err
	}
	return *resp, nil
}
