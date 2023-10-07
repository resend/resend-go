package resend

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateApiKey(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api-keys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "dacf4072-4119-4d88-932f-6202748ac7c8",
			"token": "re_c1tpEyD8_NKFusih9vKVQknRAQfmFcWCv"
		}`
		fmt.Fprintf(w, ret)
	})

	req := &CreateApiKeyRequest{
		Name: "new api key",
	}
	resp, err := client.ApiKeys.Create(testCtx, req)
	if err != nil {
		t.Errorf("ApiKeys.Create returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "dacf4072-4119-4d88-932f-6202748ac7c8")
	assert.Equal(t, resp.Token, "re_c1tpEyD8_NKFusih9vKVQknRAQfmFcWCv")
}

func TestListApiKeys(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api-keys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"data": [
				{
				  "id": "91f3200a-df72-4654-b0cd-f202395f5354",
				  "name": "Production",
				  "created_at": "2023-04-08T00:11:13.110779+00:00"
				}
			]
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.ApiKeys.List(testCtx)
	if err != nil {
		t.Errorf("ApiKeys.List returned error: %v", err)
	}
	assert.Equal(t, len(resp.Data), 1)
	assert.Equal(t, resp.Data[0].Name, "Production")
}

func TestRemoveApiKey(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api-keys/keyid", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Length", "0")
		fmt.Fprint(w, nil)
	})

	deleted, err := client.ApiKeys.Remove(testCtx, "keyid")
	if err != nil {
		t.Errorf("ApiKeys.Remove returned error: %v", err)
	}
	assert.True(t, deleted)
}
