package resend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDomainClaim(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/claim", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, `{
			"object": "domain_claim",
			"id": "dacf4072-4119-4d88-932f-6c6126d3a9d1",
			"name": "example.com",
			"status": "pending",
			"domain_id": "d91cd9bd-1176-453e-8fc1-35364d380206",
			"region": "us-east-1",
			"record": {
				"type": "TXT",
				"name": "example.com",
				"value": "resend-domain-verification=3f8a1c2d4e5b6a7f8091a2b3c4d5e6f7",
				"ttl": "Auto"
			},
			"blocked_reason": null,
			"failure_reason": null,
			"created_at": "2026-06-16T17:12:02.059593+00:00",
			"expires_at": "2026-06-23T17:12:02.059593+00:00"
		}`)
	})

	req := &CreateDomainClaimRequest{
		Name:   "example.com",
		Region: "us-east-1",
	}
	resp, err := client.Domains.Claims.Create(req)
	if err != nil {
		t.Errorf("Domains.Claims.Create returned error: %v", err)
	}

	assert.Equal(t, "domain_claim", resp.Object)
	assert.Equal(t, "dacf4072-4119-4d88-932f-6c6126d3a9d1", resp.Id)
	assert.Equal(t, "example.com", resp.Name)
	assert.Equal(t, DomainClaimStatusPending, resp.Status)
	assert.Equal(t, "d91cd9bd-1176-453e-8fc1-35364d380206", resp.DomainId)
	assert.Equal(t, "us-east-1", resp.Region)
	assert.NotNil(t, resp.Record)
	assert.Equal(t, "TXT", resp.Record.Type)
	assert.Equal(t, "example.com", resp.Record.Name)
	assert.Equal(t, "resend-domain-verification=3f8a1c2d4e5b6a7f8091a2b3c4d5e6f7", resp.Record.Value)
	assert.Equal(t, "Auto", resp.Record.Ttl)
	assert.Equal(t, "", resp.BlockedReason)
	assert.Equal(t, "", resp.FailureReason)
}

func TestCreateDomainClaimSendsAllOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/claim", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		assert.Equal(t, "example.com", body["name"])
		assert.Equal(t, "us-east-1", body["region"])
		assert.Equal(t, "send", body["custom_return_path"])
		assert.Equal(t, true, body["open_tracking"])
		assert.Equal(t, false, body["click_tracking"])
		assert.Equal(t, "links", body["tracking_subdomain"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"object": "domain_claim",
			"id": "dacf4072-4119-4d88-932f-6c6126d3a9d1",
			"name": "example.com",
			"status": "pending"
		}`)
	})

	req := &CreateDomainClaimRequest{
		Name:              "example.com",
		Region:            "us-east-1",
		CustomReturnPath:  "send",
		OpenTracking:      Bool(true),
		ClickTracking:     Bool(false),
		TrackingSubdomain: "links",
	}
	_, err := client.Domains.Claims.Create(req)
	if err != nil {
		t.Errorf("Domains.Claims.Create returned error: %v", err)
	}
}

func TestGetDomainClaim(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/d91cd9bd-1176-453e-8fc1-35364d380206/claim", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"object": "domain_claim",
			"id": "dacf4072-4119-4d88-932f-6c6126d3a9d1",
			"name": "example.com",
			"status": "blocked",
			"domain_id": "d91cd9bd-1176-453e-8fc1-35364d380206",
			"region": "us-east-1",
			"blocked_reason": "grace_period",
			"failure_reason": null,
			"created_at": "2026-06-16T17:12:02.059593+00:00",
			"expires_at": "2026-06-23T17:12:02.059593+00:00"
		}`)
	})

	resp, err := client.Domains.Claims.Get("d91cd9bd-1176-453e-8fc1-35364d380206")
	if err != nil {
		t.Errorf("Domains.Claims.Get returned error: %v", err)
	}

	assert.Equal(t, "dacf4072-4119-4d88-932f-6c6126d3a9d1", resp.Id)
	assert.Equal(t, DomainClaimStatusBlocked, resp.Status)
	assert.Equal(t, DomainClaimBlockedReasonGracePeriod, resp.BlockedReason)
}

func TestVerifyDomainClaim(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/d91cd9bd-1176-453e-8fc1-35364d380206/claim/verify", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"object": "domain_claim",
			"id": "dacf4072-4119-4d88-932f-6c6126d3a9d1",
			"name": "example.com",
			"status": "pending",
			"domain_id": "d91cd9bd-1176-453e-8fc1-35364d380206",
			"region": "us-east-1"
		}`)
	})

	resp, err := client.Domains.Claims.Verify("d91cd9bd-1176-453e-8fc1-35364d380206")
	if err != nil {
		t.Errorf("Domains.Claims.Verify returned error: %v", err)
	}

	assert.Equal(t, "dacf4072-4119-4d88-932f-6c6126d3a9d1", resp.Id)
	assert.Equal(t, DomainClaimStatusPending, resp.Status)
}
