package resend

import (
	"context"
	"net/http"
)

type Log struct {
	Object         string  `json:"object,omitempty"`
	Id             string  `json:"id"`
	CreatedAt      string  `json:"created_at"`
	Endpoint       string  `json:"endpoint"`
	Method         string  `json:"method"`
	ResponseStatus int     `json:"response_status"`
	UserAgent      *string `json:"user_agent"`
	RequestBody    any     `json:"request_body,omitempty"`
	ResponseBody   any     `json:"response_body,omitempty"`
}

type ListLogsResponse struct {
	Object  string `json:"object"`
	Data    []Log  `json:"data"`
	HasMore bool   `json:"has_more"`
}

type LogsSvc interface {
	GetWithContext(ctx context.Context, logId string) (Log, error)
	Get(logId string) (Log, error)
	ListWithOptions(ctx context.Context, options *ListOptions) (ListLogsResponse, error)
	ListWithContext(ctx context.Context) (ListLogsResponse, error)
	List() (ListLogsResponse, error)
}

type LogsSvcImpl struct {
	client *Client
}

// GetWithContext retrieves a single log entry by ID
// https://resend.com/docs/api-reference/logs/retrieve-log
func (s *LogsSvcImpl) GetWithContext(ctx context.Context, logId string) (Log, error) {
	path := "logs/" + logId

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return Log{}, ErrFailedToCreateLogsGetRequest
	}

	logResp := new(Log)

	_, err = s.client.Perform(req, logResp)
	if err != nil {
		return Log{}, err
	}

	return *logResp, nil
}

// Get retrieves a single log entry by ID
func (s *LogsSvcImpl) Get(logId string) (Log, error) {
	return s.GetWithContext(context.Background(), logId)
}

// ListWithOptions lists all logs with pagination options
// https://resend.com/docs/api-reference/logs/list-logs
func (s *LogsSvcImpl) ListWithOptions(ctx context.Context, options *ListOptions) (ListLogsResponse, error) {
	path := "logs" + buildPaginationQuery(options)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return ListLogsResponse{}, ErrFailedToCreateLogsListRequest
	}

	logsResp := new(ListLogsResponse)

	_, err = s.client.Perform(req, logsResp)
	if err != nil {
		return ListLogsResponse{}, err
	}

	return *logsResp, nil
}

// ListWithContext lists all logs
// https://resend.com/docs/api-reference/logs/list-logs
func (s *LogsSvcImpl) ListWithContext(ctx context.Context) (ListLogsResponse, error) {
	return s.ListWithOptions(ctx, nil)
}

// List lists all logs
func (s *LogsSvcImpl) List() (ListLogsResponse, error) {
	return s.ListWithContext(context.Background())
}
