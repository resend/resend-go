package resend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateContactSubscriptionIntent(t *testing.T) {
	targets := []struct {
		name       string
		id         string
		email      string
		audienceID string
		path       string
	}{
		{name: "global ID", id: "contact-id", path: "/contacts/contact-id"},
		{name: "global email", email: "contact@example.com", path: "/contacts/contact@example.com"},
		{name: "legacy ID", id: "contact-id", audienceID: "audience-id", path: "/audiences/audience-id/contacts/contact-id"},
		{name: "legacy email", email: "contact@example.com", audienceID: "audience-id", path: "/audiences/audience-id/contacts/contact@example.com"},
	}
	states := []struct {
		name        string
		value       bool
		useSetter   bool
		wantPresent bool
	}{
		{name: "direct true", value: true, wantPresent: true},
		{name: "omitted"},
		{name: "setter false", useSetter: true, wantPresent: true},
		{name: "setter true", value: true, useSetter: true, wantPresent: true},
	}

	for _, target := range targets {
		for _, state := range states {
			t.Run(target.name+"/"+state.name, func(t *testing.T) {
				received := make(chan map[string]any, 1)
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodPatch, r.Method)
					assert.Equal(t, target.path, r.URL.Path)
					assert.Empty(t, r.URL.RawQuery)
					var body map[string]any
					if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
						t.Errorf("decode update: %v", err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					received <- body
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprint(w, `{"object":"contact","id":"contact-id"}`)
				}))
				defer server.Close()

				client := NewCustomClient(server.Client(), "test-key")
				var err error
				client.BaseURL, err = url.Parse(server.URL + "/")
				require.NoError(t, err)
				params := &UpdateContactRequest{
					Id:           target.id,
					Email:        target.email,
					AudienceId:   target.audienceID,
					Unsubscribed: state.value,
					FirstName:    "Updated",
					LastName:     "Contact",
					Properties:   map[string]any{"tier": "pro", "clear": nil},
				}
				if state.useSetter {
					params.SetUnsubscribed(state.value)
				}

				_, err = client.Contacts.Update(params)
				require.NoError(t, err)
				body := <-received
				value, present := body["unsubscribed"]
				assert.Equal(t, state.wantPresent, present, "subscription intent must reach the HTTP request")
				if state.wantPresent {
					assert.Equal(t, state.value, value, "subscription intent must be sent as a JSON boolean")
				}
				assert.Equal(t, target.id, body["id"])
				assert.Equal(t, "Updated", body["first_name"])
				assert.Equal(t, "Contact", body["last_name"])
				assert.Equal(t, params.Properties, body["properties"])
				if target.email != "" {
					assert.Equal(t, target.email, body["email"])
				} else {
					assert.NotContains(t, body, "email")
				}
				if target.audienceID != "" {
					assert.Equal(t, target.audienceID, body["audience_id"])
				} else {
					assert.NotContains(t, body, "audience_id")
				}
			})
		}
	}
}
