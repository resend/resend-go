package resend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const forgeDomainResponse = `{
	"object": "domain",
	"id": "d91cd9bd-1176-453e-8fc1-35364d380206",
	"name": "example.com",
	"status": "pending",
	"created_at": "2026-11-11 10:00:00.000000+00",
	"region": "us-east-1",
	"capabilities": {
		"sending": "enabled",
		"receiving": "disabled"
	},
	"records": [
		{
			"record": "SPF",
			"name": "send",
			"value": "v=spf1 include:spf.resend.com ~all",
			"type": "TXT",
			"ttl": "Auto",
			"status": "pending"
		},
		{
			"record": "DKIM",
			"name": "resend._domainkey",
			"value": "resend._domainkey.example.com.dkim.resend.com.",
			"type": "CNAME",
			"ttl": "Auto",
			"status": "pending"
		}
	]
}`

func TestEnableDomainForge(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/d91cd9bd-1176-453e-8fc1-35364d380206/forge", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		assert.Equal(t, "send", body["return_path"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, forgeDomainResponse)
	})

	resp, err := client.DomainForge.Enable("d91cd9bd-1176-453e-8fc1-35364d380206", &EnableDomainForgeRequest{
		ReturnPath: "send",
	})
	if err != nil {
		t.Errorf("DomainForge.Enable returned error: %v", err)
	}

	assert.Equal(t, "domain", resp.Object)
	assert.Equal(t, "d91cd9bd-1176-453e-8fc1-35364d380206", resp.Id)
	assert.Equal(t, "example.com", resp.Name)
	assert.Equal(t, "pending", resp.Status)
	assert.Len(t, resp.Records, 2)
	assert.Equal(t, "send", resp.Records[0].Name)
	assert.NotNil(t, resp.Capabilities)
	assert.Equal(t, DomainCapabilityStatusEnabled, resp.Capabilities.Sending)
}

func TestEnableDomainForgeWithoutReturnPath(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/d91cd9bd-1176-453e-8fc1-35364d380206/forge", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		assert.Empty(t, body)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, forgeDomainResponse)
	})

	resp, err := client.DomainForge.Enable("d91cd9bd-1176-453e-8fc1-35364d380206", nil)
	if err != nil {
		t.Errorf("DomainForge.Enable returned error: %v", err)
	}

	assert.Equal(t, "d91cd9bd-1176-453e-8fc1-35364d380206", resp.Id)
	assert.Equal(t, "pending", resp.Status)
}

func TestUpdateDomainForge(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/d91cd9bd-1176-453e-8fc1-35364d380206/forge", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		assert.Equal(t, "outbound", body["return_path"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, forgeDomainResponse)
	})

	resp, err := client.DomainForge.Update("d91cd9bd-1176-453e-8fc1-35364d380206", &UpdateDomainForgeRequest{
		ReturnPath: "outbound",
	})
	if err != nil {
		t.Errorf("DomainForge.Update returned error: %v", err)
	}

	assert.Equal(t, "d91cd9bd-1176-453e-8fc1-35364d380206", resp.Id)
	assert.Equal(t, "pending", resp.Status)
}

func TestUpdateDomainForgeRequiresReturnPath(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.DomainForge.Update("d91cd9bd-1176-453e-8fc1-35364d380206", nil)
	assert.EqualError(t, err, "[ERROR]: ReturnPath cannot be empty")

	_, err = client.DomainForge.Update("d91cd9bd-1176-453e-8fc1-35364d380206", &UpdateDomainForgeRequest{})
	assert.EqualError(t, err, "[ERROR]: ReturnPath cannot be empty")
}

func TestVerifyDomainForge(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/d91cd9bd-1176-453e-8fc1-35364d380206/forge/verify", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"object": "domain",
			"id": "d91cd9bd-1176-453e-8fc1-35364d380206"
		}`)
	})

	resp, err := client.DomainForge.Verify("d91cd9bd-1176-453e-8fc1-35364d380206")
	if err != nil {
		t.Errorf("DomainForge.Verify returned error: %v", err)
	}

	assert.Equal(t, "domain", resp.Object)
	assert.Equal(t, "d91cd9bd-1176-453e-8fc1-35364d380206", resp.Id)
}
