package resend

import (
	_ "embed"
	"encoding/json"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var webhookSampleRE = regexp.MustCompile("(?ms)^## ([^\\n]+)\\n\\n```json\\n(.*?)\\n```")

//go:embed testdata/webhook_examples.md
var webhookSampleMarkdown string

func loadWebhookSamples() map[string]string {

	samples := make(map[string]string)
	for _, match := range webhookSampleRE.FindAllStringSubmatch(webhookSampleMarkdown, -1) {
		samples[match[1]] = match[2]
	}
	return samples
}

func TestWebhookPayloadSamples(t *testing.T) {
	samples := loadWebhookSamples()
	require.Len(t, samples, len(WebhookPayloadDataTypes))

	for eventType, sample := range samples {
		t.Run(eventType, func(t *testing.T) {
			var payload WebhookPayload
			require.NoError(t, json.Unmarshal([]byte(sample), &payload))
			require.Equal(t, eventType, payload.Type)
			require.NotNil(t, payload.Data)

			expectedType := reflect.New(WebhookPayloadDataTypes[eventType]).Interface()
			assert.IsType(t, expectedType, payload.Data)

			encoded, err := json.Marshal(payload)
			require.NoError(t, err)
			assert.JSONEq(t, canonicalizeWebhookTimes(t, sample), string(encoded))
		})
	}
}

func canonicalizeWebhookTimes(t *testing.T, input string) string {
	t.Helper()

	var value interface{}
	require.NoError(t, json.Unmarshal([]byte(input), &value))
	canonicalizeWebhookTimeFields(value)

	b, err := json.Marshal(value)
	require.NoError(t, err)
	return string(b)
}

func canonicalizeWebhookTimeFields(value interface{}) {
	switch value := value.(type) {
	case map[string]interface{}:
		for key, child := range value {
			if isWebhookTimeField(key) {
				if text, ok := child.(string); ok {
					if parsed, err := time.Parse(time.RFC3339Nano, text); err == nil {
						value[key] = parsed.Format(time.RFC3339Nano)
					}
				}
			}
			canonicalizeWebhookTimeFields(child)
		}
	case []interface{}:
		for _, child := range value {
			canonicalizeWebhookTimeFields(child)
		}
	}
}

func isWebhookTimeField(name string) bool {
	switch name {
	case "created_at", "updated_at", "received_at", "last_received":
		return true
	default:
		return false
	}
}

func TestWebhookPayloadUnknownType(t *testing.T) {
	var payload WebhookPayload
	err := json.Unmarshal([]byte(`{
		"type": "unknown.event",
		"created_at": "2026-11-22T23:41:12.126Z",
		"data": {}
	}`), &payload)

	require.Error(t, err)
	assert.ErrorContains(t, err, "not implemented type")
}

func TestWebhookPayloadRequiresObjectData(t *testing.T) {
	var payload WebhookPayload
	err := json.Unmarshal([]byte(`{
		"type": "email.bounced",
		"created_at": "2026-11-22T23:41:12.126Z",
		"data": null
	}`), &payload)

	require.Error(t, err)
	assert.ErrorContains(t, err, "data must be a JSON object")
}

func TestWebhookPayloadNilDataMarshalsAsNull(t *testing.T) {
	payload := WebhookPayload{Type: "email.bounced"}
	b, err := json.Marshal(payload)
	require.NoError(t, err)
	assert.JSONEq(t, `{"created_at":"0001-01-01T00:00:00Z","data":null,"type":"email.bounced"}`, string(b))
}
