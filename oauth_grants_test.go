package resend

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListOAuthGrants(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/oauth/grants", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
			"has_more": false,
			"data": [
				{
					"id": "650e8400-e29b-41d4-a716-446655440001",
					"client_id": "430eed87-632a-4ea6-90db-0aace67ec228",
					"scopes": ["emails:send"],
					"created_at": "2023-06-21T06:10:36.144Z",
					"revoked_at": null,
					"revoked_reason": null,
					"client": {
						"name": "Resend CLI",
						"logo_uri": "https://example.com/logo.png"
					}
				},
				{
					"id": "650e8400-e29b-41d4-a716-446655440002",
					"client_id": "430eed87-632a-4ea6-90db-0aace67ec228",
					"scopes": ["emails:send", "domains:read"],
					"created_at": "2023-06-20T06:10:36.144Z",
					"revoked_at": "2023-06-22T06:10:36.144Z",
					"revoked_reason": "revoked_from_api",
					"client": {
						"name": "Resend CLI",
						"logo_uri": null
					}
				}
			]
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.OAuthGrants.List()
	if err != nil {
		t.Errorf("OAuthGrants.List returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.False(t, resp.HasMore)
	assert.Equal(t, len(resp.Data), 2)

	assert.Equal(t, resp.Data[0].Id, "650e8400-e29b-41d4-a716-446655440001")
	assert.Equal(t, resp.Data[0].ClientId, "430eed87-632a-4ea6-90db-0aace67ec228")
	assert.Equal(t, resp.Data[0].Scopes, []string{"emails:send"})
	assert.Nil(t, resp.Data[0].RevokedAt)
	assert.Nil(t, resp.Data[0].RevokedReason)
	assert.Equal(t, resp.Data[0].Client.Name, "Resend CLI")
	assert.NotNil(t, resp.Data[0].Client.LogoUri)
	assert.Equal(t, *resp.Data[0].Client.LogoUri, "https://example.com/logo.png")

	assert.Equal(t, resp.Data[1].Scopes, []string{"emails:send", "domains:read"})
	assert.NotNil(t, resp.Data[1].RevokedAt)
	assert.Equal(t, *resp.Data[1].RevokedAt, "2023-06-22T06:10:36.144Z")
	assert.NotNil(t, resp.Data[1].RevokedReason)
	assert.Equal(t, *resp.Data[1].RevokedReason, "revoked_from_api")
	assert.Nil(t, resp.Data[1].Client.LogoUri)
}

func TestRevokeOAuthGrant(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/oauth/grants/650e8400-e29b-41d4-a716-446655440001", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "oauth_grant",
			"id": "650e8400-e29b-41d4-a716-446655440001",
			"revoked_at": "2026-04-08T00:11:13.110Z",
			"revoked_reason": "revoked_from_api"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.OAuthGrants.Revoke("650e8400-e29b-41d4-a716-446655440001")
	if err != nil {
		t.Errorf("OAuthGrants.Revoke returned error: %v", err)
	}
	assert.Equal(t, resp.Object, "oauth_grant")
	assert.Equal(t, resp.Id, "650e8400-e29b-41d4-a716-446655440001")
	assert.NotNil(t, resp.RevokedAt)
	assert.Equal(t, *resp.RevokedAt, "2026-04-08T00:11:13.110Z")
	assert.NotNil(t, resp.RevokedReason)
	assert.Equal(t, *resp.RevokedReason, "revoked_from_api")
}

func TestRevokeOAuthGrantNullFields(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/oauth/grants/650e8400-e29b-41d4-a716-446655440001", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "oauth_grant",
			"id": "650e8400-e29b-41d4-a716-446655440001",
			"revoked_at": null,
			"revoked_reason": null
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.OAuthGrants.Revoke("650e8400-e29b-41d4-a716-446655440001")
	if err != nil {
		t.Errorf("OAuthGrants.Revoke returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "650e8400-e29b-41d4-a716-446655440001")
	assert.Nil(t, resp.RevokedAt)
	assert.Nil(t, resp.RevokedReason)
}
