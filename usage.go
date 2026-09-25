package resend

import (
	"context"
	"errors"
	"net/http"
)

type UsageSvc interface {
	GetWithContext(ctx context.Context) (Usage, error)
	Get() (Usage, error)
}

type UsageSvcImpl struct {
	client *Client
}

// UsageEmailsPeriod is the usage for emails over a given period (daily or monthly).
type UsageEmailsPeriod struct {
	Used     int    `json:"used"`
	Limit    *int   `json:"limit"`
	Sent     int    `json:"sent"`
	Received int    `json:"received"`
	ResetsAt string `json:"resets_at"`
}

// UsageEmails is the emails usage, broken down by period.
type UsageEmails struct {
	Daily   UsageEmailsPeriod `json:"daily"`
	Monthly UsageEmailsPeriod `json:"monthly"`
}

type UsageContacts struct {
	Used  int `json:"used"`
	Limit int `json:"limit"`
}

type UsageSegments struct {
	Used  int  `json:"used"`
	Limit *int `json:"limit"`
}

type UsageBroadcasts struct {
	Used  int  `json:"used"`
	Limit *int `json:"limit"`
}

type UsageAiCredits struct {
	Used           int     `json:"used"`
	Limit          *int    `json:"limit"`
	NextIncreaseAt *string `json:"next_increase_at"`
}

type UsageAutomationRuns struct {
	Used     int    `json:"used"`
	Limit    int    `json:"limit"`
	ResetsAt string `json:"resets_at"`
}

type UsageDomains struct {
	Used  int  `json:"used"`
	Limit *int `json:"limit"`
}

// UsageRateLimit is the rate limit applied to the caller's API key.
type UsageRateLimit struct {
	Limit    int    `json:"limit"`
	Duration string `json:"duration"`
}

// Usage provides the structure for the response from the Get call.
type Usage struct {
	Object         string              `json:"object"`
	Emails         UsageEmails         `json:"emails"`
	Contacts       UsageContacts       `json:"contacts"`
	Segments       UsageSegments       `json:"segments"`
	Broadcasts     UsageBroadcasts     `json:"broadcasts"`
	AiCredits      UsageAiCredits      `json:"ai_credits"`
	AutomationRuns UsageAutomationRuns `json:"automation_runs"`
	Domains        UsageDomains        `json:"domains"`
	RateLimit      UsageRateLimit      `json:"rate_limit"`
}

// GetWithContext retrieves the caller's account-level usage and quota data.
// https://resend.com/docs/api-reference/usage/get-usage
func (s *UsageSvcImpl) GetWithContext(ctx context.Context) (Usage, error) {
	path := "usage"

	// Prepare request
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return Usage{}, errors.New("[ERROR]: Failed to create Usage.Get request")
	}

	usage := new(Usage)

	// Send Request
	_, err = s.client.Perform(req, usage)

	if err != nil {
		return Usage{}, err
	}

	return *usage, nil
}

// Get retrieves the caller's account-level usage and quota data.
// https://resend.com/docs/api-reference/usage/get-usage
func (s *UsageSvcImpl) Get() (Usage, error) {
	return s.GetWithContext(context.Background())
}
