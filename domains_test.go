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

		var ret interface{}
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
		Name: "example.com",
	}
	resp, err := client.Domains.Create(testCtx, req)
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

	assert.Equal(t, resp.Records[0].Record, "SPF")
	assert.Equal(t, resp.Records[0].Name, "bounces")
	assert.Equal(t, resp.Records[0].Type, "MX")
	assert.Equal(t, resp.Records[0].Ttl, "Auto")
	assert.Equal(t, resp.Records[0].Status, "not_started")
	assert.Equal(t, resp.Records[0].Value, "feedback-smtp.us-east-1.amazonses.com")
	assert.Equal(t, resp.Records[0].Priority, json.Number("10"))

	assert.Equal(t, resp.Records[1].Priority, json.Number(""))
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

	verified, err := client.Domains.Verify(testCtx, "d91cd9bd-1176-453e-8fc1-35364d380206")
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
				"region": "us-east-1"
			  }
			]
		}`

		fmt.Fprint(w, ret)
	})

	domains, err := client.Domains.List(testCtx)
	if err != nil {
		t.Errorf("Domains.List returned error: %v", err)
	}

	assert.Equal(t, len(domains.Data), 1)
	assert.Equal(t, domains.Data[0].Id, "d91cd9bd-1176-453e-8fc1-35364d380206")
	assert.Equal(t, domains.Data[0].Name, "example.com")
	assert.Equal(t, domains.Data[0].Status, "not_started")
	assert.Equal(t, domains.Data[0].CreatedAt, "2023-04-26T20:21:26.347412+00:00")
	assert.Equal(t, domains.Data[0].Region, "us-east-1")
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

	deleted, err := client.Domains.Remove(testCtx, "b6d24b8e-af0b-4c3c-be0c-359bbd97381e")
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
			"region": "us-east-1"
		}`

		fmt.Fprint(w, ret)
	})

	domain, err := client.Domains.Get(testCtx, "d91cd9bd-1176-453e-8fc1-35364d380206")
	if err != nil {
		t.Errorf("Domains.Get returned error: %v", err)
	}

	assert.Equal(t, domain.Id, "d91cd9bd-1176-453e-8fc1-35364d380206")
	assert.Equal(t, domain.Object, "domain")
	assert.Equal(t, domain.Name, "example.com")
	assert.Equal(t, domain.Status, "not_started")
	assert.Equal(t, domain.CreatedAt, "2023-04-26T20:21:26.347412+00:00")
	assert.Equal(t, domain.Region, "us-east-1")
}
