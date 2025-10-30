package resend

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSegment(t *testing.T) {
	setup()
	defer teardown()

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

	req := &CreateSegmentRequest{
		Name: "New Segment",
	}
	resp, err := client.Segments.Create(req)
	if err != nil {
		t.Errorf("Segments.Create returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "78261eea-8f8b-4381-83c6-79fa7120f1c")
	assert.Equal(t, resp.Object, "segment")
	assert.Equal(t, resp.Name, "Registered Users")
}

func TestListSegments(t *testing.T) {
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

	segments, err := client.Segments.List()
	if err != nil {
		t.Errorf("Segments.List returned error: %v", err)
	}

	assert.Equal(t, len(segments.Data), 1)
	assert.Equal(t, segments.Object, "list")
	assert.Equal(t, segments.Data[0].Id, "d91cd9bd-1176-453e-8fc1-35364d380206")
	assert.Equal(t, segments.Data[0].Name, "Registered Users")
	assert.Equal(t, segments.Data[0].CreatedAt, "2023-04-26T20:21:26.347412+00:00")
}

func TestRemoveSegment(t *testing.T) {
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

	deleted, err := client.Segments.Remove("b6d24b8e-af0b-4c3c-be0c-359bbd97381e")
	if err != nil {
		t.Errorf("Segments.Remove returned error: %v", err)
	}
	assert.True(t, deleted.Deleted)
	assert.Equal(t, deleted.Id, "b6d24b8e-af0b-4c3c-be0c-359bbd97381e")
	assert.Equal(t, deleted.Object, "segment")
}

func TestGetSegment(t *testing.T) {
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

	segment, err := client.Segments.Get("d91cd9bd-1176-453e-8fc1-35364d380206")
	if err != nil {
		t.Errorf("Segment.Get returned error: %v", err)
	}

	assert.Equal(t, segment.Id, "d91cd9bd-1176-453e-8fc1-35364d380206")
	assert.Equal(t, segment.Object, "segment")
	assert.Equal(t, segment.Name, "Registered Users")
	assert.Equal(t, segment.CreatedAt, "2023-10-06T22:59:55.977Z")
}
