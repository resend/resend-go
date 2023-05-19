package resend

import (
	"encoding/json"
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
		t.Errorf("ApiKeys.Created returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "1923781293")
	assert.Equal(t, resp.Token, "99999199219192")
}
