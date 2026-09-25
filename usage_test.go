package resend

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUsage(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "usage",
			"emails": {
				"daily": { "used": 258, "limit": null, "sent": 57, "received": 201, "resets_at": "2026-07-17T00:00:00.000Z" },
				"monthly": { "used": 5422, "limit": 10000, "sent": 1000, "received": 4442, "resets_at": "2026-08-01T00:00:00.000Z" }
			},
			"contacts": { "used": 85000, "limit": 150000 },
			"segments": { "used": 2, "limit": 3 },
			"broadcasts": { "used": 100, "limit": null },
			"ai_credits": { "used": 0, "limit": 500, "next_increase_at": "2026-07-18T09:00:00.000Z" },
			"automation_runs": { "used": 0, "limit": 1000, "resets_at": "2026-08-01T00:00:00.000Z" },
			"domains": { "used": 1, "limit": 1000 },
			"rate_limit": { "limit": 10, "duration": "1000ms" }
		}`

		fmt.Fprint(w, ret)
	})

	usage, err := client.Usage.Get()
	if err != nil {
		t.Errorf("Usage.Get returned error: %v", err)
	}

	assert.Equal(t, "usage", usage.Object)
	assert.Equal(t, 258, usage.Emails.Daily.Used)
	assert.Nil(t, usage.Emails.Daily.Limit)
	assert.Equal(t, 57, usage.Emails.Daily.Sent)
	assert.Equal(t, 201, usage.Emails.Daily.Received)
	assert.Equal(t, "2026-07-17T00:00:00.000Z", usage.Emails.Daily.ResetsAt)

	assert.Equal(t, 5422, usage.Emails.Monthly.Used)
	assert.NotNil(t, usage.Emails.Monthly.Limit)
	assert.Equal(t, 10000, *usage.Emails.Monthly.Limit)

	assert.Equal(t, 85000, usage.Contacts.Used)
	assert.Equal(t, 150000, usage.Contacts.Limit)

	assert.Equal(t, 2, usage.Segments.Used)
	assert.NotNil(t, usage.Segments.Limit)
	assert.Equal(t, 3, *usage.Segments.Limit)

	assert.Equal(t, 100, usage.Broadcasts.Used)
	assert.Nil(t, usage.Broadcasts.Limit)

	assert.Equal(t, 0, usage.AiCredits.Used)
	assert.NotNil(t, usage.AiCredits.Limit)
	assert.Equal(t, 500, *usage.AiCredits.Limit)
	assert.NotNil(t, usage.AiCredits.NextIncreaseAt)
	assert.Equal(t, "2026-07-18T09:00:00.000Z", *usage.AiCredits.NextIncreaseAt)

	assert.Equal(t, 0, usage.AutomationRuns.Used)
	assert.Equal(t, 1000, usage.AutomationRuns.Limit)
	assert.Equal(t, "2026-08-01T00:00:00.000Z", usage.AutomationRuns.ResetsAt)

	assert.Equal(t, 1, usage.Domains.Used)
	assert.NotNil(t, usage.Domains.Limit)
	assert.Equal(t, 1000, *usage.Domains.Limit)

	assert.Equal(t, 10, usage.RateLimit.Limit)
	assert.Equal(t, "1000ms", usage.RateLimit.Duration)
}

func TestGetUsageWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/usage", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "usage",
			"emails": {
				"daily": { "used": 0, "limit": null, "sent": 0, "received": 0, "resets_at": "2026-07-17T00:00:00.000Z" },
				"monthly": { "used": 0, "limit": 10000, "sent": 0, "received": 0, "resets_at": "2026-08-01T00:00:00.000Z" }
			},
			"contacts": { "used": 0, "limit": 150000 },
			"segments": { "used": 0, "limit": null },
			"broadcasts": { "used": 0, "limit": null },
			"ai_credits": { "used": 0, "limit": null, "next_increase_at": null },
			"automation_runs": { "used": 0, "limit": 1000, "resets_at": "2026-08-01T00:00:00.000Z" },
			"domains": { "used": 0, "limit": null },
			"rate_limit": { "limit": 10, "duration": "1000ms" }
		}`

		fmt.Fprint(w, ret)
	})

	ctx := context.Background()
	usage, err := client.Usage.GetWithContext(ctx)
	if err != nil {
		t.Errorf("Usage.GetWithContext returned error: %v", err)
	}

	assert.Equal(t, "usage", usage.Object)
	assert.Nil(t, usage.Segments.Limit)
	assert.Nil(t, usage.AiCredits.Limit)
	assert.Nil(t, usage.AiCredits.NextIncreaseAt)
	assert.Nil(t, usage.Domains.Limit)
}
