package resend

import (
	"encoding/json"
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

		ret := &CreateApiKeyResponse{
			Id:    "1923781293",
			Token: "99999199219192",
		}
		err := json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	req := &CreateApiKeyRequest{
		Name: "new api key",
	}
	resp, err := client.ApiKeys.Create(req)
	if err != nil {
		t.Errorf("ApiKeys.Create returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "1923781293")
	assert.Equal(t, resp.Token, "99999199219192")
}

func TestListApiKeys(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api-keys", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &ListApiKeysResponse{
			Data: []ApiKey{
				{
					Name:      "prod",
					Id:        "91f3200a-df72-4654-b0cd-f202395f5354",
					CreatedAt: "2023-04-08T00:11:13.110779+00:00",
				},
				{
					Name:      "stage",
					Id:        "91f3200a-df72-4654-b0cd-f402395f5354",
					CreatedAt: "2023-07-08T00:11:13.110779+00:00",
				},
			},
		}
		err := json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	resp, err := client.ApiKeys.List()
	if err != nil {
		t.Errorf("ApiKeys.List returned error: %v", err)
	}
	assert.Equal(t, len(resp.Data), 2)
	assert.Equal(t, resp.Data[0].Name, "prod")
}

func TestDeleteApiKey(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api-keys/keyid", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Length", "0")
		fmt.Fprint(w, nil)
	})

	deleted, err := client.ApiKeys.Delete("keyid")
	if err != nil {
		t.Errorf("ApiKeys.Delete returned error: %v", err)
	}
	assert.True(t, deleted)
}
