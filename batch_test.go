package resend

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBatchSendEmail(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/batch", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &BatchEmailResponse{
			Data: []SendEmailResponse{
				{Id: "1"},
				{Id: "2"},
			},
		}
		err := json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	req := []*SendEmailRequest{
		{
			To: []string{"d@e.com"},
		},
		{
			To: []string{"d@e.com"},
		},
	}
	resp, err := client.Batch.Send(req)
	if err != nil {
		t.Errorf("Emails.Send returned error: %v", err)
	}
	assert.Equal(t, resp.Data[0].Id, "1")
	assert.Equal(t, resp.Data[1].Id, "2")
}
