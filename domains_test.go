package resend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDomain(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		var ret any
		ret = `
		{
			"id": "4dd369bc-aa82-4ff3-97de-514ae3000ee0",
			"name": "example.com",
			"createdAt": "2023-03-28T17:12:02.059593+00:00",
			"status": "not_started",
			"records": [
			  {
				"record": "SPF",
				"name": "bounces",
				"type": "MX",
				"ttl": "Auto",
				"status": "not_started",
				"value": "feedback-smtp.us-east-1.amazonses.com",
				"priority": 10
			  },
			  {
				"record": "SPF",
				"name": "bounces",
				"value": "\"v=spf1 include:amazonses.com ~all\"",
				"type": "TXT",
				"ttl": "Auto",
				"status": "not_started"
			  },
			  {
				"record": "DKIM",
				"name": "nhapbbryle57yxg3fbjytyodgbt2kyyg._domainkey",
				"value": "nhapbbryle57yxg3fbjytyodgbt2kyyg.dkim.amazonses.com.",
				"type": "CNAME",
				"status": "not_started",
				"ttl": "Auto"
			  },
			  {
				"record": "DKIM",
				"name": "xbakwbe5fcscrhzshpap6kbxesf6pfgn._domainkey",
				"value": "xbakwbe5fcscrhzshpap6kbxesf6pfgn.dkim.amazonses.com.",
				"type": "CNAME",
				"status": "not_started",
				"ttl": "Auto"
			  },
			  {
				"record": "DKIM",
				"name": "txrcreso3dqbvcve45tqyosxwaegvhgn._domainkey",
				"value": "txrcreso3dqbvcve45tqyosxwaegvhgn.dkim.amazonses.com.",
				"type": "CNAME",
				"status": "not_started",
				"ttl": "Auto"
			  }
			],
			"region": "us-east-1",
			"dnsProvider": "Unidentified"
		  }`

		fmt.Fprint(w, ret)
	})

	req := &CreateDomainRequest{
		Name:             "example.com",
		Region:           "us-east-1",
		CustomReturnPath: "outbound",
	}
	resp, err := client.Domains.Create(req)
	if err != nil {
		t.Errorf("Domains.Create returned error: %v", err)
	}
	assert.Equal(t, resp.DnsProvider, "Unidentified")
	assert.Equal(t, resp.Id, "4dd369bc-aa82-4ff3-97de-514ae3000ee0")
	assert.Equal(t, resp.Region, "us-east-1")
	assert.Equal(t, resp.Status, "not_started")
	assert.Equal(t, resp.CreatedAt, "2023-03-28T17:12:02.059593+00:00")
	assert.NotNil(t, resp.Records)
	assert.Equal(t, len(resp.Records), 5)

	assert.Equal(t, resp.Records[0].Record, RecordTypeSPF)
	assert.Equal(t, resp.Records[0].Name, "bounces")
	assert.Equal(t, resp.Records[0].Type, "MX")
	assert.Equal(t, resp.Records[0].Ttl, "Auto")
	assert.Equal(t, resp.Records[0].Status, "not_started")
	assert.Equal(t, resp.Records[0].Value, "feedback-smtp.us-east-1.amazonses.com")
	assert.Equal(t, resp.Records[0].Priority, json.Number("10"))

	assert.Equal(t, resp.Records[1].Priority, json.Number(""))
}

func TestCreateDomainWithCapabilities(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		capabilities, ok := body["capabilities"].(map[string]any)
		assert.True(t, ok, "capabilities should be present in payload as an object")
		assert.Equal(t, "enabled", capabilities["sending"])
		assert.Equal(t, "enabled", capabilities["receiving"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, `{
			"id": "4dd369bc-aa82-4ff3-97de-514ae3000ee0",
			"name": "example.com",
			"createdAt": "2023-03-28T17:12:02.059593+00:00",
			"status": "not_started",
			"region": "us-east-1",
			"dnsProvider": "Unidentified"
		}`)
	})

	req := &CreateDomainRequest{
		Name:   "example.com",
		Region: "us-east-1",
		Capabilities: &DomainCapabilities{
			Sending:   DomainCapabilityStatusEnabled,
			Receiving: DomainCapabilityStatusEnabled,
		},
	}
	resp, err := client.Domains.Create(req)
	if err != nil {
		t.Errorf("Domains.Create returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "4dd369bc-aa82-4ff3-97de-514ae3000ee0")
}

func TestCreateDomainWithoutCapabilitiesOmitsKey(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		_, present := body["capabilities"]
		assert.False(t, present, "capabilities should be omitted when not set")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, `{
			"id": "4dd369bc-aa82-4ff3-97de-514ae3000ee0",
			"name": "example.com",
			"createdAt": "2023-03-28T17:12:02.059593+00:00",
			"status": "not_started",
			"region": "us-east-1",
			"dnsProvider": "Unidentified"
		}`)
	})

	req := &CreateDomainRequest{
		Name:   "example.com",
		Region: "us-east-1",
	}
	_, err := client.Domains.Create(req)
	if err != nil {
		t.Errorf("Domains.Create returned error: %v", err)
	}
}

func TestCreateDomainResponseDeserializesCapabilities(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{
			"id": "4dd369bc-aa82-4ff3-97de-514ae3000ee0",
			"name": "example.com",
			"createdAt": "2023-03-28T17:12:02.059593+00:00",
			"status": "not_started",
			"region": "us-east-1",
			"dnsProvider": "Unidentified",
			"capabilities": {
				"sending": "enabled",
				"receiving": "disabled"
			}
		}`)
	})

	req := &CreateDomainRequest{Name: "example.com"}
	resp, err := client.Domains.Create(req)
	if err != nil {
		t.Errorf("Domains.Create returned error: %v", err)
	}
	assert.NotNil(t, resp.Capabilities)
	assert.Equal(t, DomainCapabilityStatusEnabled, resp.Capabilities.Sending)
	assert.Equal(t, DomainCapabilityStatusDisabled, resp.Capabilities.Receiving)
}

func TestVerifyDomain(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/d91cd9bd-1176-453e-8fc1-35364d380206/verify", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Length", "0")
		fmt.Fprint(w, nil)
	})

	verified, err := client.Domains.Verify("d91cd9bd-1176-453e-8fc1-35364d380206")
	if err != nil {
		t.Errorf("Domains.Verify returned error: %v", err)
	}
	assert.True(t, verified)
}

func TestListDomains(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"data": [
			  {
				"id": "d91cd9bd-1176-453e-8fc1-35364d380206",
				"name": "example.com",
				"status": "not_started",
				"created_at": "2023-04-26T20:21:26.347412+00:00",
				"region": "us-east-1",
				"records": [
					{
						"record": "SPF",
						"name": "bounces"
					}
				]
			  }
			]
		}`

		fmt.Fprint(w, ret)
	})

	domains, err := client.Domains.List()
	if err != nil {
		t.Errorf("Domains.List returned error: %v", err)
	}

	assert.Equal(t, len(domains.Data), 1)
	assert.Equal(t, domains.Data[0].Id, "d91cd9bd-1176-453e-8fc1-35364d380206")
	assert.Equal(t, domains.Data[0].Name, "example.com")
	assert.Equal(t, domains.Data[0].Status, "not_started")
	assert.Equal(t, domains.Data[0].CreatedAt, "2023-04-26T20:21:26.347412+00:00")
	assert.Equal(t, domains.Data[0].Region, "us-east-1")
	assert.Equal(t, domains.Data[0].Records[0].Record, RecordTypeSPF)
}

func TestRemoveDomain(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/b6d24b8e-af0b-4c3c-be0c-359bbd97381e", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Length", "0")
		fmt.Fprint(w, nil)
	})

	deleted, err := client.Domains.Remove("b6d24b8e-af0b-4c3c-be0c-359bbd97381e")
	if err != nil {
		t.Errorf("Domains.Remove returned error: %v", err)
	}
	assert.True(t, deleted)
}

func TestGetDomain(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/d91cd9bd-1176-453e-8fc1-35364d380206", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "domain",
			"id": "d91cd9bd-1176-453e-8fc1-35364d380206",
			"name": "example.com",
			"status": "not_started",
			"created_at": "2023-04-26T20:21:26.347412+00:00",
			"region": "us-east-1",
			"records": [
				{
					"record": "SPF",
					"name": "bounces"
				}
			]
		}`

		fmt.Fprint(w, ret)
	})

	domain, err := client.Domains.Get("d91cd9bd-1176-453e-8fc1-35364d380206")
	if err != nil {
		t.Errorf("Domains.Get returned error: %v", err)
	}

	assert.Equal(t, domain.Id, "d91cd9bd-1176-453e-8fc1-35364d380206")
	assert.Equal(t, domain.Object, "domain")
	assert.Equal(t, domain.Name, "example.com")
	assert.Equal(t, domain.Status, "not_started")
	assert.Equal(t, domain.CreatedAt, "2023-04-26T20:21:26.347412+00:00")
	assert.Equal(t, domain.Region, "us-east-1")
	assert.Equal(t, domain.Records[0].Record, RecordTypeSPF)
	assert.Equal(t, domain.Records[0].Name, "bounces")
}

func TestGetDomainDeserializesCapabilities(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/d91cd9bd-1176-453e-8fc1-35364d380206", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "domain",
			"id": "d91cd9bd-1176-453e-8fc1-35364d380206",
			"name": "example.com",
			"status": "not_started",
			"created_at": "2023-04-26T20:21:26.347412+00:00",
			"region": "us-east-1",
			"capabilities": {
				"sending": "enabled",
				"receiving": "enabled"
			}
		}`)
	})

	domain, err := client.Domains.Get("d91cd9bd-1176-453e-8fc1-35364d380206")
	if err != nil {
		t.Errorf("Domains.Get returned error: %v", err)
	}
	assert.NotNil(t, domain.Capabilities)
	assert.Equal(t, DomainCapabilityStatusEnabled, domain.Capabilities.Sending)
	assert.Equal(t, DomainCapabilityStatusEnabled, domain.Capabilities.Receiving)
}

func TestGetDomainPartiallyVerified(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/fd61172c-cafc-40f5-b049-b45947779a29", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "domain",
			"id": "fd61172c-cafc-40f5-b049-b45947779a29",
			"name": "resend.com",
			"status": "partially_verified",
			"created_at": "2023-06-21T06:10:36.144Z",
			"region": "us-east-1",
			"records": [
				{
					"record": "DKIM",
					"name": "resend._domainkey",
					"value": "p=MIG...",
					"type": "TXT",
					"status": "verified",
					"ttl": "Auto"
				},
				{
					"record": "Tracking",
					"name": "track.resend.com",
					"value": "tracking.resend.com",
					"type": "CNAME",
					"ttl": "Auto",
					"status": "pending"
				}
			]
		}`)
	})

	domain, err := client.Domains.Get("fd61172c-cafc-40f5-b049-b45947779a29")
	if err != nil {
		t.Errorf("Domains.Get returned error: %v", err)
	}

	assert.Equal(t, domain.Status, DomainStatusPartiallyVerified)
	assert.Equal(t, domain.Records[0].Status, DomainRecordStatusVerified)
	assert.Equal(t, domain.Records[1].Status, DomainRecordStatusPending)
}

func TestGetDomainPartiallyFailed(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/fd61172c-cafc-40f5-b049-b45947779a30", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "domain",
			"id": "fd61172c-cafc-40f5-b049-b45947779a30",
			"name": "resend.com",
			"status": "partially_failed",
			"created_at": "2023-06-21T06:10:36.144Z",
			"region": "us-east-1",
			"records": [
				{
					"record": "DKIM",
					"name": "resend._domainkey",
					"value": "p=MIG...",
					"type": "TXT",
					"status": "verified",
					"ttl": "Auto"
				},
				{
					"record": "Receiving",
					"name": "resend.com",
					"value": "inbound-mx.resend.com",
					"type": "MX",
					"ttl": "Auto",
					"status": "failed",
					"priority": 10
				}
			]
		}`)
	})

	domain, err := client.Domains.Get("fd61172c-cafc-40f5-b049-b45947779a30")
	if err != nil {
		t.Errorf("Domains.Get returned error: %v", err)
	}

	assert.Equal(t, domain.Status, DomainStatusPartiallyFailed)
	assert.Equal(t, domain.Records[0].Status, DomainRecordStatusVerified)
	assert.Equal(t, domain.Records[1].Status, DomainRecordStatusFailed)
}

func TestCreateDomainWithTrackingSubdomain(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, `{
			"id": "4dd369bc-aa82-4ff3-97de-514ae3000ee0",
			"name": "example.com",
			"createdAt": "2023-03-28T17:12:02.059593+00:00",
			"status": "not_started",
			"region": "us-east-1",
			"dnsProvider": "Unidentified",
			"open_tracking": true,
			"click_tracking": true,
			"tracking_subdomain": "links",
			"records": [
				{
					"record": "SPF",
					"name": "bounces",
					"type": "MX",
					"ttl": "Auto",
					"status": "not_started",
					"value": "feedback-smtp.us-east-1.amazonses.com",
					"priority": 10
				},
				{
					"record": "Tracking",
					"name": "links.example.com",
					"value": "links1.resend-dns.com",
					"type": "CNAME",
					"ttl": "Auto",
					"status": "not_started"
				},
				{
					"record": "TrackingCAA",
					"name": "",
					"value": "0 issue \"amazon.com\"",
					"type": "CAA",
					"ttl": "Auto",
					"status": "not_started"
				}
			]
		}`)
	})

	req := &CreateDomainRequest{
		Name:              "example.com",
		Region:            "us-east-1",
		TrackingSubdomain: "links",
		OpenTracking:      Bool(true),
		ClickTracking:     Bool(true),
	}
	resp, err := client.Domains.Create(req)
	if err != nil {
		t.Errorf("Domains.Create returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "4dd369bc-aa82-4ff3-97de-514ae3000ee0")
	assert.True(t, resp.OpenTracking)
	assert.True(t, resp.ClickTracking)
	assert.Equal(t, resp.TrackingSubdomain, "links")
	assert.Equal(t, len(resp.Records), 3)
	assert.Equal(t, resp.Records[1].Record, RecordTypeTracking)
	assert.Equal(t, resp.Records[1].Name, "links.example.com")
	assert.Equal(t, resp.Records[1].Value, "links1.resend-dns.com")
	assert.Equal(t, resp.Records[1].Type, "CNAME")
	assert.Equal(t, resp.Records[1].Priority, json.Number(""))
	assert.Equal(t, resp.Records[2].Record, RecordTypeTrackingCAA)
	assert.Equal(t, resp.Records[2].Name, "")
	assert.Equal(t, resp.Records[2].Value, "0 issue \"amazon.com\"")
	assert.Equal(t, resp.Records[2].Type, "CAA")
	assert.Equal(t, resp.Records[2].Ttl, "Auto")
}

func TestGetDomainWithTrackingFields(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/d91cd9bd-1176-453e-8fc1-35364d380206", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"object": "domain",
			"id": "d91cd9bd-1176-453e-8fc1-35364d380206",
			"name": "example.com",
			"status": "not_started",
			"created_at": "2023-04-26T20:21:26.347412+00:00",
			"region": "us-east-1",
			"open_tracking": true,
			"click_tracking": true,
			"tracking_subdomain": "links",
			"records": [
				{
					"record": "SPF",
					"name": "bounces"
				},
				{
					"record": "Tracking",
					"name": "links.example.com",
					"value": "links1.resend-dns.com",
					"type": "CNAME",
					"ttl": "Auto",
					"status": "not_started"
				},
				{
					"record": "TrackingCAA",
					"name": "",
					"value": "0 issue \"amazon.com\"",
					"type": "CAA",
					"ttl": "Auto",
					"status": "verified"
				}
			]
		}`)
	})

	domain, err := client.Domains.Get("d91cd9bd-1176-453e-8fc1-35364d380206")
	if err != nil {
		t.Errorf("Domains.Get returned error: %v", err)
	}
	assert.Equal(t, domain.Id, "d91cd9bd-1176-453e-8fc1-35364d380206")
	assert.True(t, domain.OpenTracking)
	assert.True(t, domain.ClickTracking)
	assert.Equal(t, domain.TrackingSubdomain, "links")
	assert.Equal(t, len(domain.Records), 3)
	assert.Equal(t, domain.Records[2].Record, RecordTypeTrackingCAA)
	assert.Equal(t, domain.Records[2].Type, "CAA")
	assert.Equal(t, domain.Records[2].Value, "0 issue \"amazon.com\"")
}

func TestUpdateDomainWithTrackingSubdomain(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/d91cd9bd-1176-453e-8fc1-35364d380206", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"id": "d91cd9bd-1176-453e-8fc1-35364d380206",
			"object": "domain"
		}`)
	})

	params := &UpdateDomainRequest{
		TrackingSubdomain: "links",
	}
	updated, err := client.Domains.Update("d91cd9bd-1176-453e-8fc1-35364d380206", params)
	if err != nil {
		t.Errorf("Domains.Update returned error: %v", err)
	}
	assert.Equal(t, updated.Id, "d91cd9bd-1176-453e-8fc1-35364d380206")
}

func TestUpdateDomainSetTrackingToFalse(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/d91cd9bd-1176-453e-8fc1-35364d380206", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		// Both fields must be present and false in the payload
		openTracking, ok := body["open_tracking"]
		assert.True(t, ok, "open_tracking should be present in payload")
		assert.Equal(t, false, openTracking)

		clickTracking, ok := body["click_tracking"]
		assert.True(t, ok, "click_tracking should be present in payload")
		assert.Equal(t, false, clickTracking)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"id": "d91cd9bd-1176-453e-8fc1-35364d380206", "object": "domain"}`)
	})

	params := &UpdateDomainRequest{}
	params.SetOpenTracking(false)
	params.SetClickTracking(false)

	updated, err := client.Domains.Update("d91cd9bd-1176-453e-8fc1-35364d380206", params)
	if err != nil {
		t.Errorf("Domains.Update returned error: %v", err)
	}
	assert.Equal(t, updated.Id, "d91cd9bd-1176-453e-8fc1-35364d380206")
}

func TestUpdateDomain(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/domains/d91cd9bd-1176-453e-8fc1-35364d380206", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "d91cd9bd-1176-453e-8fc1-35364d380206",
			"object": "domain"
		}`

		fmt.Fprint(w, ret)
	})

	params := &UpdateDomainRequest{
		OpenTracking: true,
		Tls:          Opportunistic,
	}
	updated, err := client.Domains.Update("d91cd9bd-1176-453e-8fc1-35364d380206", params)
	if err != nil {
		t.Errorf("Domains.Update returned error: %v", err)
	}
	assert.True(t, updated.Id == "d91cd9bd-1176-453e-8fc1-35364d380206")
	assert.True(t, updated.Object == "domain")
}
