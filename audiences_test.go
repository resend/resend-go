package resend

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateAudience tests the deprecated Audiences API
// which internally calls the Segments API
func TestCreateAudience(t *testing.T) {
	setup()
	defer teardown()

	// Note: This still hits /segments endpoint because Audiences wraps Segments
	mux.HandleFunc("/segments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		var ret interface{}
		ret = `
		{
			"object": "segment",
			"id": "78261eea-8f8b-4381-83c6-79fa7120f1c",
			"name": "Registered Users"
		}`

		fmt.Fprint(w, ret)
	})

	req := &CreateAudienceRequest{
		Name: "New Audience",
	}
	resp, err := client.Audiences.Create(req)
	if err != nil {
		t.Errorf("Audiences.Create returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "78261eea-8f8b-4381-83c6-79fa7120f1c")
	assert.Equal(t, resp.Object, "segment")
	assert.Equal(t, resp.Name, "Registered Users")
}

func TestListAudiences(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/segments", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
			"data": [
			  {
					"id": "d91cd9bd-1176-453e-8fc1-35364d380206",
					"name": "Registered Users",
					"created_at": "2023-04-26T20:21:26.347412+00:00"
			  }
			]
		}`

		fmt.Fprint(w, ret)
	})

	audiences, err := client.Audiences.List()
	if err != nil {
		t.Errorf("Audiences.List returned error: %v", err)
	}

	assert.Equal(t, len(audiences.Data), 1)
	assert.Equal(t, audiences.Object, "list")
	assert.Equal(t, audiences.Data[0].Id, "d91cd9bd-1176-453e-8fc1-35364d380206")
	assert.Equal(t, audiences.Data[0].Name, "Registered Users")
	assert.Equal(t, audiences.Data[0].CreatedAt, "2023-04-26T20:21:26.347412+00:00")
}

func TestRemoveAudience(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/segments/b6d24b8e-af0b-4c3c-be0c-359bbd97381e", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusOK)

		var ret interface{}
		ret = `
		{
			"object": "segment",
			"id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
			"deleted": true
		}`

		fmt.Fprint(w, ret)
	})

	deleted, err := client.Audiences.Remove("b6d24b8e-af0b-4c3c-be0c-359bbd97381e")
	if err != nil {
		t.Errorf("Audiences.Remove returned error: %v", err)
	}
	assert.True(t, deleted.Deleted)
	assert.Equal(t, deleted.Id, "b6d24b8e-af0b-4c3c-be0c-359bbd97381e")
	assert.Equal(t, deleted.Object, "segment")
}

func TestGetAudience(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/segments/d91cd9bd-1176-453e-8fc1-35364d380206", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "segment",
			"id": "d91cd9bd-1176-453e-8fc1-35364d380206",
			"name": "Registered Users",
			"created_at": "2023-10-06T22:59:55.977Z"
		}`

		fmt.Fprint(w, ret)
	})

	audience, err := client.Audiences.Get("d91cd9bd-1176-453e-8fc1-35364d380206")
	if err != nil {
		t.Errorf("Audience.Get returned error: %v", err)
	}

	assert.Equal(t, audience.Id, "d91cd9bd-1176-453e-8fc1-35364d380206")
	assert.Equal(t, audience.Object, "segment")
	assert.Equal(t, audience.Name, "Registered Users")
	assert.Equal(t, audience.CreatedAt, "2023-10-06T22:59:55.977Z")
}
