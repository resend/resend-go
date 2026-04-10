package resend

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type AutomationStatus = string

const (
	AutomationStatusEnabled  AutomationStatus = "enabled"
	AutomationStatusDisabled AutomationStatus = "disabled"
)

type AutomationRunStatus = string

const (
	AutomationRunStatusRunning   AutomationRunStatus = "running"
	AutomationRunStatusCompleted AutomationRunStatus = "completed"
	AutomationRunStatusFailed    AutomationRunStatus = "failed"
	AutomationRunStatusCancelled AutomationRunStatus = "cancelled"
)

type AutomationStepType = string

const (
	AutomationStepTypeTrigger       AutomationStepType = "trigger"
	AutomationStepTypeSendEmail     AutomationStepType = "send_email"
	AutomationStepTypeDelay         AutomationStepType = "delay"
	AutomationStepTypeWaitForEvent  AutomationStepType = "wait_for_event"
	AutomationStepTypeCondition     AutomationStepType = "condition"
	AutomationStepTypeContactUpdate AutomationStepType = "contact_update"
	AutomationStepTypeContactDelete AutomationStepType = "contact_delete"
	AutomationStepTypeAddToSegment  AutomationStepType = "add_to_segment"
)

type AutomationConnectionType = string

const (
	AutomationConnectionTypeDefault         AutomationConnectionType = "default"
	AutomationConnectionTypeConditionMet    AutomationConnectionType = "condition_met"
	AutomationConnectionTypeConditionNotMet AutomationConnectionType = "condition_not_met"
	AutomationConnectionTypeTimeout         AutomationConnectionType = "timeout"
	AutomationConnectionTypeEventReceived   AutomationConnectionType = "event_received"
)

type AutomationStep struct {
	Key    string             `json:"key"`
	Type   AutomationStepType `json:"type"`
	Config map[string]any     `json:"config"`
}

type AutomationStepResponse struct {
	Key    string             `json:"key"`
	Type   AutomationStepType `json:"type"`
	Config map[string]any     `json:"config"`
}

type AutomationConnection struct {
	From string                   `json:"from"`
	To   string                   `json:"to"`
	Type AutomationConnectionType `json:"type,omitempty"`
}

type CreateAutomationRequest struct {
	Name        string           `json:"name"`
	Status      AutomationStatus `json:"status,omitempty"`
	Steps       []AutomationStep `json:"steps"`
	Connections []AutomationConnection `json:"connections"`
}

type CreateAutomationResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

type UpdateAutomationRequest struct {
	Name        string           `json:"name,omitempty"`
	Status      AutomationStatus `json:"status,omitempty"`
	Steps       []AutomationStep `json:"steps,omitempty"`
	Connections []AutomationConnection `json:"connections,omitempty"`
}

type UpdateAutomationResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

type DeleteAutomationResponse struct {
	Object  string `json:"object"`
	Id      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

type StopAutomationResponse struct {
	Object string `json:"object"`
	Id     string `json:"id"`
	Status string `json:"status"`
}

type AutomationListItem struct {
	Id        string           `json:"id"`
	Name      string           `json:"name"`
	Status    AutomationStatus `json:"status"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
}

type ListAutomationsResponse struct {
	Object  string               `json:"object"`
	HasMore bool                 `json:"has_more"`
	Data    []AutomationListItem `json:"data"`
}

type Automation struct {
	Object      string                   `json:"object"`
	Id          string                   `json:"id"`
	Name        string                   `json:"name"`
	Status      AutomationStatus         `json:"status"`
	CreatedAt   string                   `json:"created_at"`
	UpdatedAt   string                   `json:"updated_at"`
	Steps       []AutomationStepResponse `json:"steps"`
	Connections []AutomationConnection         `json:"connections"`
}

type AutomationRunListItem struct {
	Id          string              `json:"id"`
	Status      AutomationRunStatus `json:"status"`
	StartedAt   *string             `json:"started_at"`
	CompletedAt *string             `json:"completed_at"`
	CreatedAt   string              `json:"created_at"`
}

type ListAutomationRunsResponse struct {
	Object  string                  `json:"object"`
	HasMore bool                    `json:"has_more"`
	Data    []AutomationRunListItem `json:"data"`
}

type AutomationRunStep struct {
	Key         string             `json:"key"`
	Type        AutomationStepType `json:"type"`
	Status      string             `json:"status"`
	StartedAt   *string            `json:"started_at"`
	CompletedAt *string            `json:"completed_at"`
	Output      any                `json:"output"`
	Error       any                `json:"error"`
	CreatedAt   string             `json:"created_at"`
}

type AutomationRun struct {
	Object      string              `json:"object"`
	Id          string              `json:"id"`
	Status      AutomationRunStatus `json:"status"`
	StartedAt   *string             `json:"started_at"`
	CompletedAt *string             `json:"completed_at"`
	CreatedAt   string              `json:"created_at"`
	Steps       []AutomationRunStep `json:"steps"`
}

// ListAutomationsOptions contains pagination and filter parameters for listing automations
type ListAutomationsOptions struct {
	Status *AutomationStatus
	Limit  *int
	After  *string
	Before *string
}

// ListAutomationRunsOptions contains pagination and filter parameters for listing automation runs
type ListAutomationRunsOptions struct {
	Status []AutomationRunStatus
	Limit  *int
	After  *string
	Before *string
}

type AutomationsSvc interface {
	CreateWithContext(ctx context.Context, params *CreateAutomationRequest) (CreateAutomationResponse, error)
	Create(params *CreateAutomationRequest) (CreateAutomationResponse, error)
	GetWithContext(ctx context.Context, automationId string) (Automation, error)
	Get(automationId string) (Automation, error)
	ListWithContext(ctx context.Context) (ListAutomationsResponse, error)
	List() (ListAutomationsResponse, error)
	ListWithOptions(ctx context.Context, options *ListAutomationsOptions) (ListAutomationsResponse, error)
	UpdateWithContext(ctx context.Context, automationId string, params *UpdateAutomationRequest) (UpdateAutomationResponse, error)
	Update(automationId string, params *UpdateAutomationRequest) (UpdateAutomationResponse, error)
	RemoveWithContext(ctx context.Context, automationId string) (DeleteAutomationResponse, error)
	Remove(automationId string) (DeleteAutomationResponse, error)
	StopWithContext(ctx context.Context, automationId string) (StopAutomationResponse, error)
	Stop(automationId string) (StopAutomationResponse, error)
	ListRunsWithContext(ctx context.Context, automationId string, options *ListAutomationRunsOptions) (ListAutomationRunsResponse, error)
	ListRuns(automationId string) (ListAutomationRunsResponse, error)
	GetRunWithContext(ctx context.Context, automationId string, runId string) (AutomationRun, error)
	GetRun(automationId string, runId string) (AutomationRun, error)
}

type AutomationsSvcImpl struct {
	client *Client
}

func buildAutomationsQuery(options *ListAutomationsOptions) string {
	if options == nil {
		return ""
	}
	query := make(url.Values)
	if options.Status != nil {
		query.Set("status", *options.Status)
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
	if len(query) > 0 {
		return "?" + query.Encode()
	}
	return ""
}

func buildAutomationRunsQuery(options *ListAutomationRunsOptions) string {
	if options == nil {
		return ""
	}
	query := make(url.Values)
	if len(options.Status) > 0 {
		statuses := make([]string, len(options.Status))
		for i, s := range options.Status {
			statuses[i] = s
		}
		query.Set("status", strings.Join(statuses, ","))
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
	if len(query) > 0 {
		return "?" + query.Encode()
	}
	return ""
}

// CreateWithContext creates a new Automation
// https://resend.com/docs/api-reference/automations/create-automation
func (s *AutomationsSvcImpl) CreateWithContext(ctx context.Context, params *CreateAutomationRequest) (CreateAutomationResponse, error) {
	req, err := s.client.NewRequest(ctx, http.MethodPost, "automations", params)
	if err != nil {
		return CreateAutomationResponse{}, ErrFailedToCreateAutomationCreateRequest
	}

	resp := new(CreateAutomationResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return CreateAutomationResponse{}, err
	}

	return *resp, nil
}

// Create creates a new Automation
func (s *AutomationsSvcImpl) Create(params *CreateAutomationRequest) (CreateAutomationResponse, error) {
	return s.CreateWithContext(context.Background(), params)
}

// GetWithContext retrieves a single Automation by ID
// https://resend.com/docs/api-reference/automations/get-automation
func (s *AutomationsSvcImpl) GetWithContext(ctx context.Context, automationId string) (Automation, error) {
	path := "automations/" + automationId

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return Automation{}, ErrFailedToCreateAutomationGetRequest
	}

	resp := new(Automation)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return Automation{}, err
	}

	return *resp, nil
}

// Get retrieves a single Automation by ID
func (s *AutomationsSvcImpl) Get(automationId string) (Automation, error) {
	return s.GetWithContext(context.Background(), automationId)
}

// ListWithOptions retrieves a list of Automations with pagination and filter options
// https://resend.com/docs/api-reference/automations/list-automations
func (s *AutomationsSvcImpl) ListWithOptions(ctx context.Context, options *ListAutomationsOptions) (ListAutomationsResponse, error) {
	path := "automations" + buildAutomationsQuery(options)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListAutomationsResponse{}, ErrFailedToCreateAutomationListRequest
	}

	resp := new(ListAutomationsResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return ListAutomationsResponse{}, err
	}

	return *resp, nil
}

// ListWithContext retrieves a list of Automations
// https://resend.com/docs/api-reference/automations/list-automations
func (s *AutomationsSvcImpl) ListWithContext(ctx context.Context) (ListAutomationsResponse, error) {
	return s.ListWithOptions(ctx, nil)
}

// List retrieves a list of Automations
func (s *AutomationsSvcImpl) List() (ListAutomationsResponse, error) {
	return s.ListWithContext(context.Background())
}

// UpdateWithContext updates an existing Automation
// https://resend.com/docs/api-reference/automations/update-automation
func (s *AutomationsSvcImpl) UpdateWithContext(ctx context.Context, automationId string, params *UpdateAutomationRequest) (UpdateAutomationResponse, error) {
	path := "automations/" + automationId

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, params)
	if err != nil {
		return UpdateAutomationResponse{}, ErrFailedToCreateAutomationUpdateRequest
	}

	resp := new(UpdateAutomationResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return UpdateAutomationResponse{}, err
	}

	return *resp, nil
}

// Update updates an existing Automation
func (s *AutomationsSvcImpl) Update(automationId string, params *UpdateAutomationRequest) (UpdateAutomationResponse, error) {
	return s.UpdateWithContext(context.Background(), automationId, params)
}

// RemoveWithContext deletes an Automation by ID
// https://resend.com/docs/api-reference/automations/delete-automation
func (s *AutomationsSvcImpl) RemoveWithContext(ctx context.Context, automationId string) (DeleteAutomationResponse, error) {
	path := "automations/" + automationId

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return DeleteAutomationResponse{}, ErrFailedToCreateAutomationRemoveRequest
	}

	resp := new(DeleteAutomationResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return DeleteAutomationResponse{}, err
	}

	return *resp, nil
}

// Remove deletes an Automation by ID
func (s *AutomationsSvcImpl) Remove(automationId string) (DeleteAutomationResponse, error) {
	return s.RemoveWithContext(context.Background(), automationId)
}

// StopWithContext stops a running Automation by ID
// https://resend.com/docs/api-reference/automations/stop-automation
func (s *AutomationsSvcImpl) StopWithContext(ctx context.Context, automationId string) (StopAutomationResponse, error) {
	path := "automations/" + automationId + "/stop"

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return StopAutomationResponse{}, ErrFailedToCreateAutomationStopRequest
	}

	resp := new(StopAutomationResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return StopAutomationResponse{}, err
	}

	return *resp, nil
}

// Stop stops a running Automation by ID
func (s *AutomationsSvcImpl) Stop(automationId string) (StopAutomationResponse, error) {
	return s.StopWithContext(context.Background(), automationId)
}

// ListRunsWithContext retrieves a list of runs for an Automation
// https://resend.com/docs/api-reference/automations/list-automation-runs
func (s *AutomationsSvcImpl) ListRunsWithContext(ctx context.Context, automationId string, options *ListAutomationRunsOptions) (ListAutomationRunsResponse, error) {
	path := "automations/" + automationId + "/runs" + buildAutomationRunsQuery(options)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListAutomationRunsResponse{}, ErrFailedToCreateAutomationListRunsRequest
	}

	resp := new(ListAutomationRunsResponse)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return ListAutomationRunsResponse{}, err
	}

	return *resp, nil
}

// ListRuns retrieves a list of runs for an Automation
func (s *AutomationsSvcImpl) ListRuns(automationId string) (ListAutomationRunsResponse, error) {
	return s.ListRunsWithContext(context.Background(), automationId, nil)
}

// GetRunWithContext retrieves a single run for an Automation
// https://resend.com/docs/api-reference/automations/get-automation-run
func (s *AutomationsSvcImpl) GetRunWithContext(ctx context.Context, automationId string, runId string) (AutomationRun, error) {
	path := "automations/" + automationId + "/runs/" + runId

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return AutomationRun{}, ErrFailedToCreateAutomationGetRunRequest
	}

	resp := new(AutomationRun)
	_, err = s.client.Perform(req, resp)
	if err != nil {
		return AutomationRun{}, err
	}

	return *resp, nil
}

// GetRun retrieves a single run for an Automation
func (s *AutomationsSvcImpl) GetRun(automationId string, runId string) (AutomationRun, error) {
	return s.GetRunWithContext(context.Background(), automationId, runId)
}
