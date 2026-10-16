package resend

import (
	"context"
	"errors"
	"net/http"
)

// AddInboxThreadLabelsRequest contains params for adding labels to a thread.
type AddInboxThreadLabelsRequest struct {
	// Labels are the IDs of the labels to apply to the thread. Between 1 and 10 IDs, each a label
	// of the inbox. Labels already on the thread stay on it.
	Labels []string `json:"labels"`
}

// SetInboxThreadLabelsRequest contains params for replacing the labels of a thread.
type SetInboxThreadLabelsRequest struct {
	// Labels are the IDs of the labels the thread carries after the change. Up to 100 IDs, each a
	// label of the inbox. A nil or empty slice takes every label off the thread.
	Labels []string `json:"labels"`
}

// InboxThreadLabelsResponse lists every label on the thread after the change.
type InboxThreadLabelsResponse struct {
	Object string             `json:"object"`
	Data   []InboxThreadLabel `json:"data"`
}

type InboxThreadLabelsSvc interface {
	Add(inboxId, threadId string, params *AddInboxThreadLabelsRequest) (InboxThreadLabelsResponse, error)
	AddWithContext(ctx context.Context, inboxId, threadId string, params *AddInboxThreadLabelsRequest) (InboxThreadLabelsResponse, error)
	Set(inboxId, threadId string, params *SetInboxThreadLabelsRequest) (InboxThreadLabelsResponse, error)
	SetWithContext(ctx context.Context, inboxId, threadId string, params *SetInboxThreadLabelsRequest) (InboxThreadLabelsResponse, error)
	Remove(inboxId, threadId, labelId string) (InboxThreadLabelsResponse, error)
	RemoveWithContext(ctx context.Context, inboxId, threadId, labelId string) (InboxThreadLabelsResponse, error)
}

type InboxThreadLabelsSvcImpl struct {
	client *Client
}

func validateInboxThreadIds(inboxId, threadId string) error {
	if inboxId == "" {
		return errors.New("[ERROR]: Inbox ID is required")
	}
	if threadId == "" {
		return errors.New("[ERROR]: Thread ID is required")
	}
	return nil
}

// Add applies one or more labels to the thread. Labels already on the thread stay on it. A label
// ID that is not a label of the inbox returns a 404, and nothing changes. To replace the labels of
// the thread, use Set.
// https://resend.com/docs/api-reference/inboxes/add-thread-labels
func (s *InboxThreadLabelsSvcImpl) Add(inboxId, threadId string, params *AddInboxThreadLabelsRequest) (InboxThreadLabelsResponse, error) {
	return s.AddWithContext(context.Background(), inboxId, threadId, params)
}

// AddWithContext applies one or more labels to the thread with context.
// https://resend.com/docs/api-reference/inboxes/add-thread-labels
func (s *InboxThreadLabelsSvcImpl) AddWithContext(ctx context.Context, inboxId, threadId string, params *AddInboxThreadLabelsRequest) (InboxThreadLabelsResponse, error) {
	if err := validateInboxThreadIds(inboxId, threadId); err != nil {
		return InboxThreadLabelsResponse{}, err
	}
	if params == nil || len(params.Labels) == 0 {
		return InboxThreadLabelsResponse{}, errors.New("[ERROR]: At least one label ID is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, inboxPath(inboxId, "threads", threadId, "labels"), params)
	if err != nil {
		return InboxThreadLabelsResponse{}, ErrFailedToCreateInboxesThreadLabelsAddRequest
	}

	resp := new(InboxThreadLabelsResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return InboxThreadLabelsResponse{}, err
	}
	return *resp, nil
}

// Set replaces every label on the thread with the given labels, in one atomic change. The labels
// themselves are kept. A nil params or an empty Labels slice takes every label off the thread. A
// label ID that is not a label of the inbox returns a 404, and nothing changes.
// https://resend.com/docs/api-reference/inboxes/set-thread-labels
func (s *InboxThreadLabelsSvcImpl) Set(inboxId, threadId string, params *SetInboxThreadLabelsRequest) (InboxThreadLabelsResponse, error) {
	return s.SetWithContext(context.Background(), inboxId, threadId, params)
}

// SetWithContext replaces every label on the thread with context.
// https://resend.com/docs/api-reference/inboxes/set-thread-labels
func (s *InboxThreadLabelsSvcImpl) SetWithContext(ctx context.Context, inboxId, threadId string, params *SetInboxThreadLabelsRequest) (InboxThreadLabelsResponse, error) {
	if err := validateInboxThreadIds(inboxId, threadId); err != nil {
		return InboxThreadLabelsResponse{}, err
	}
	// The API requires the labels field, and a nil slice would marshal to null.
	body := SetInboxThreadLabelsRequest{Labels: []string{}}
	if params != nil && params.Labels != nil {
		body.Labels = params.Labels
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, inboxPath(inboxId, "threads", threadId, "labels"), body)
	if err != nil {
		return InboxThreadLabelsResponse{}, ErrFailedToCreateInboxesThreadLabelsSetRequest
	}

	resp := new(InboxThreadLabelsResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return InboxThreadLabelsResponse{}, err
	}
	return *resp, nil
}

// Remove takes one label off the thread. The label itself is kept, and stays on every other thread
// that carries it. A label that is not on the thread returns the current labels. To delete the
// label, use Inboxes.Labels.Remove.
// https://resend.com/docs/api-reference/inboxes/remove-thread-label
func (s *InboxThreadLabelsSvcImpl) Remove(inboxId, threadId, labelId string) (InboxThreadLabelsResponse, error) {
	return s.RemoveWithContext(context.Background(), inboxId, threadId, labelId)
}

// RemoveWithContext takes one label off the thread with context.
// https://resend.com/docs/api-reference/inboxes/remove-thread-label
func (s *InboxThreadLabelsSvcImpl) RemoveWithContext(ctx context.Context, inboxId, threadId, labelId string) (InboxThreadLabelsResponse, error) {
	if err := validateInboxThreadIds(inboxId, threadId); err != nil {
		return InboxThreadLabelsResponse{}, err
	}
	if labelId == "" {
		return InboxThreadLabelsResponse{}, errors.New("[ERROR]: Label ID is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, inboxPath(inboxId, "threads", threadId, "labels", labelId), nil)
	if err != nil {
		return InboxThreadLabelsResponse{}, ErrFailedToCreateInboxesThreadLabelsRemoveRequest
	}

	resp := new(InboxThreadLabelsResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return InboxThreadLabelsResponse{}, err
	}
	return *resp, nil
}
