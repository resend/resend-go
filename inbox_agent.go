package resend

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// InboxAgentAction is an action the inbox agent can take.
type InboxAgentAction = string

const (
	InboxAgentActionDraftReply    InboxAgentAction = "draft_reply"
	InboxAgentActionForwardThread InboxAgentAction = "forward_thread"
	InboxAgentActionAddLabels     InboxAgentAction = "add_labels"
	InboxAgentActionAssignThread  InboxAgentAction = "assign_thread"
	InboxAgentActionMarkAsSpam    InboxAgentAction = "mark_as_spam"
	InboxAgentActionArchiveThread InboxAgentAction = "archive_thread"
	InboxAgentActionDeleteThread  InboxAgentAction = "delete_thread"
)

// InboxAgent holds the agent settings of an inbox.
type InboxAgent struct {
	Object string `json:"object"`
	// Instructions the agent follows when it handles the threads of the inbox. nil when no
	// instructions are set.
	Instructions *string `json:"instructions"`
	// Tone the agent writes in. One of the presets `concise_professional`, `concise_casual`,
	// `expanded_professional`, or `expanded_casual`, or custom text in the owner's own words, such
	// as `friendly and concise`. Older dashboard values such as `Concise · Casual` are custom text.
	// nil when no tone is set.
	Tone *string `json:"tone"`
	// EnabledActions are the actions the agent can take. Empty when no actions are enabled.
	EnabledActions []InboxAgentAction `json:"enabled_actions"`
}

// UpdateInboxAgentRequest contains params for updating the agent settings of an inbox. Provide at
// least one field. Omitted fields keep their current value.
type UpdateInboxAgentRequest struct {
	// Instructions the agent follows when it handles the threads of the inbox. Max 4000
	// characters. The API trims it, and the trimmed value must not be empty. Use
	// ClearInstructions to clear the instructions.
	Instructions *string
	// Tone the agent writes in. Use one of the presets `concise_professional`, `concise_casual`,
	// `expanded_professional`, or `expanded_casual`, or describe the tone in your own words, such
	// as `friendly and concise`. Custom text is max 64 characters. The API trims it, and the
	// trimmed value must not be empty. Use ClearTone to clear the tone.
	Tone *string
	// EnabledActions replaces the whole set of actions the agent can take. Duplicates are
	// removed. A nil slice is omitted. An empty, non-nil slice disables all actions.
	EnabledActions []InboxAgentAction

	instructionsCleared bool
	toneCleared         bool
}

// ClearInstructions sends `instructions: null`, which clears the instructions.
func (r *UpdateInboxAgentRequest) ClearInstructions() {
	r.Instructions = nil
	r.instructionsCleared = true
}

// ClearTone sends `tone: null`, which clears the tone.
func (r *UpdateInboxAgentRequest) ClearTone() {
	r.Tone = nil
	r.toneCleared = true
}

func (r UpdateInboxAgentRequest) isEmpty() bool {
	return r.Instructions == nil && !r.instructionsCleared &&
		r.Tone == nil && !r.toneCleared &&
		r.EnabledActions == nil
}

func (r UpdateInboxAgentRequest) MarshalJSON() ([]byte, error) {
	aux := make(map[string]any)
	setNullableString(aux, "instructions", r.Instructions, r.instructionsCleared)
	setNullableString(aux, "tone", r.Tone, r.toneCleared)
	if r.EnabledActions != nil {
		aux["enabled_actions"] = r.EnabledActions
	}
	return json.Marshal(aux)
}

type UpdateInboxAgentResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

type InboxAgentSvc interface {
	Get(inboxId string) (InboxAgent, error)
	GetWithContext(ctx context.Context, inboxId string) (InboxAgent, error)
	Update(inboxId string, params *UpdateInboxAgentRequest) (UpdateInboxAgentResponse, error)
	UpdateWithContext(ctx context.Context, inboxId string, params *UpdateInboxAgentRequest) (UpdateInboxAgentResponse, error)
}

type InboxAgentSvcImpl struct {
	client *Client
}

// Get retrieves the agent settings of the inbox. An inbox without a configured agent returns
// empty settings: Instructions and Tone are nil and EnabledActions is empty.
// https://resend.com/docs/api-reference/inboxes/get-agent
func (s *InboxAgentSvcImpl) Get(inboxId string) (InboxAgent, error) {
	return s.GetWithContext(context.Background(), inboxId)
}

// GetWithContext retrieves the agent settings of the inbox with context.
// https://resend.com/docs/api-reference/inboxes/get-agent
func (s *InboxAgentSvcImpl) GetWithContext(ctx context.Context, inboxId string) (InboxAgent, error) {
	if inboxId == "" {
		return InboxAgent{}, errors.New("[ERROR]: Inbox ID is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, inboxPath(inboxId, "agent"), nil)
	if err != nil {
		return InboxAgent{}, ErrFailedToCreateInboxesAgentGetRequest
	}

	resp := new(InboxAgent)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return InboxAgent{}, err
	}
	return *resp, nil
}

// Update updates the agent settings of the inbox. The agent is created when the inbox has none.
// Omitted fields keep their current value.
// https://resend.com/docs/api-reference/inboxes/update-agent
func (s *InboxAgentSvcImpl) Update(inboxId string, params *UpdateInboxAgentRequest) (UpdateInboxAgentResponse, error) {
	return s.UpdateWithContext(context.Background(), inboxId, params)
}

// UpdateWithContext updates the agent settings of the inbox with context.
// https://resend.com/docs/api-reference/inboxes/update-agent
func (s *InboxAgentSvcImpl) UpdateWithContext(ctx context.Context, inboxId string, params *UpdateInboxAgentRequest) (UpdateInboxAgentResponse, error) {
	if inboxId == "" {
		return UpdateInboxAgentResponse{}, errors.New("[ERROR]: Inbox ID is required")
	}
	if params == nil || params.isEmpty() {
		return UpdateInboxAgentResponse{}, errors.New("[ERROR]: At least one of Instructions, Tone or EnabledActions is required")
	}

	req, err := s.client.NewRequest(ctx, http.MethodPatch, inboxPath(inboxId, "agent"), params)
	if err != nil {
		return UpdateInboxAgentResponse{}, ErrFailedToCreateInboxesAgentUpdateRequest
	}

	resp := new(UpdateInboxAgentResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return UpdateInboxAgentResponse{}, err
	}
	return *resp, nil
}

func setNullableString(aux map[string]any, key string, value *string, cleared bool) {
	if value != nil {
		aux[key] = *value
	} else if cleared {
		aux[key] = nil
	}
}
