package resend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddSuppression(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		assert.Equal(t, "steve.wozniak@gmail.com", body["email"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, `{
			"object": "suppression",
			"id": "e169aa45-1ecf-4183-9955-b1499d5701d3"
		}`)
	})

	req := &AddSuppressionRequest{Email: "steve.wozniak@gmail.com"}
	resp, err := client.Suppressions.Add(req)
	if err != nil {
		t.Errorf("Suppressions.Add returned error: %v", err)
	}

	assert.Equal(t, "suppression", resp.Object)
	assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", resp.Id)
}

func TestAddSuppressionWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, `{
			"object": "suppression",
			"id": "e169aa45-1ecf-4183-9955-b1499d5701d3"
		}`)
	})

	ctx := context.Background()
	resp, err := client.Suppressions.AddWithContext(ctx, &AddSuppressionRequest{
		Email: "steve.wozniak@gmail.com",
	})
	if err != nil {
		t.Errorf("Suppressions.AddWithContext returned error: %v", err)
	}

	assert.Equal(t, "suppression", resp.Object)
	assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", resp.Id)
}

func TestAddSuppressionMissingEmail(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Suppressions.Add(&AddSuppressionRequest{})
	assert.Error(t, err)
	assert.Equal(t, "[ERROR]: Email is required", err.Error())
}

func TestListSuppressions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"object": "list",
			"has_more": true,
			"data": [
				{
					"id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
					"email": "steve.wozniak@gmail.com",
					"origin": "bounce",
					"source_id": "479e3145-dd38-476b-932c-529ceb705947",
					"created_at": "2023-10-06T23:47:56.678Z"
				},
				{
					"id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
					"email": "steve.jobs@gmail.com",
					"origin": "manual",
					"source_id": null,
					"created_at": "2023-10-07T10:12:31.001Z"
				}
			]
		}`)
	})

	resp, err := client.Suppressions.List(nil)
	if err != nil {
		t.Errorf("Suppressions.List returned error: %v", err)
	}

	assert.Equal(t, "list", resp.Object)
	assert.True(t, resp.HasMore)
	assert.Len(t, resp.Data, 2)

	assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", resp.Data[0].Id)
	assert.Equal(t, "steve.wozniak@gmail.com", resp.Data[0].Email)
	assert.Equal(t, SuppressionOriginBounce, resp.Data[0].Origin)
	assert.NotNil(t, resp.Data[0].SourceId)
	assert.Equal(t, "479e3145-dd38-476b-932c-529ceb705947", *resp.Data[0].SourceId)
	assert.Equal(t, "2023-10-06T23:47:56.678Z", resp.Data[0].CreatedAt)

	assert.Equal(t, SuppressionOriginManual, resp.Data[1].Origin)
	assert.Nil(t, resp.Data[1].SourceId)
}

func TestListSuppressionsWithOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, "complaint", r.URL.Query().Get("origin"))
		assert.Equal(t, "20", r.URL.Query().Get("limit"))
		assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", r.URL.Query().Get("after"))
		assert.Equal(t, "b6d24b8e-af0b-4c3c-be0c-359bbd97381e", r.URL.Query().Get("before"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"object":"list","has_more":false,"data":[]}`)
	})

	limit := 20
	after := "e169aa45-1ecf-4183-9955-b1499d5701d3"
	before := "b6d24b8e-af0b-4c3c-be0c-359bbd97381e"
	resp, err := client.Suppressions.List(&ListSuppressionsOptions{
		Origin: SuppressionOriginComplaint,
		Limit:  &limit,
		After:  &after,
		Before: &before,
	})
	if err != nil {
		t.Errorf("Suppressions.List returned error: %v", err)
	}

	assert.Equal(t, "list", resp.Object)
	assert.False(t, resp.HasMore)
	assert.Empty(t, resp.Data)
}

func TestListSuppressionsWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"object": "list",
			"has_more": false,
			"data": [
				{
					"id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
					"email": "steve.wozniak@gmail.com",
					"origin": "bounce",
					"source_id": "479e3145-dd38-476b-932c-529ceb705947",
					"created_at": "2023-10-06T23:47:56.678Z"
				}
			]
		}`)
	})

	ctx := context.Background()
	resp, err := client.Suppressions.ListWithContext(ctx, nil)
	if err != nil {
		t.Errorf("Suppressions.ListWithContext returned error: %v", err)
	}

	assert.Equal(t, "list", resp.Object)
	assert.False(t, resp.HasMore)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", resp.Data[0].Id)
}

func TestSuppressionListEntryHasNoObjectField(t *testing.T) {
	entry, err := json.Marshal(SuppressionListEntry{Id: "e169aa45-1ecf-4183-9955-b1499d5701d3"})
	if err != nil {
		t.Fatalf("failed to marshal list entry: %v", err)
	}
	assert.NotContains(t, string(entry), "object")

	single, err := json.Marshal(Suppression{Object: "suppression", Id: "e169aa45-1ecf-4183-9955-b1499d5701d3"})
	if err != nil {
		t.Fatalf("failed to marshal suppression: %v", err)
	}
	assert.Contains(t, string(single), `"object":"suppression"`)
}

func TestGetSuppression(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions/e169aa45-1ecf-4183-9955-b1499d5701d3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"object": "suppression",
			"id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
			"email": "steve.wozniak@gmail.com",
			"origin": "complaint",
			"source_id": "479e3145-dd38-476b-932c-529ceb705947",
			"created_at": "2023-10-06T23:47:56.678Z"
		}`)
	})

	resp, err := client.Suppressions.Get("e169aa45-1ecf-4183-9955-b1499d5701d3")
	if err != nil {
		t.Errorf("Suppressions.Get returned error: %v", err)
	}

	assert.Equal(t, "suppression", resp.Object)
	assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", resp.Id)
	assert.Equal(t, "steve.wozniak@gmail.com", resp.Email)
	assert.Equal(t, SuppressionOriginComplaint, resp.Origin)
	assert.NotNil(t, resp.SourceId)
	assert.Equal(t, "479e3145-dd38-476b-932c-529ceb705947", *resp.SourceId)
	assert.Equal(t, "2023-10-06T23:47:56.678Z", resp.CreatedAt)
}

func TestGetSuppressionWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions/e169aa45-1ecf-4183-9955-b1499d5701d3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"object": "suppression",
			"id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
			"email": "steve.wozniak@gmail.com",
			"origin": "complaint",
			"source_id": "479e3145-dd38-476b-932c-529ceb705947",
			"created_at": "2023-10-06T23:47:56.678Z"
		}`)
	})

	ctx := context.Background()
	resp, err := client.Suppressions.GetWithContext(ctx, "e169aa45-1ecf-4183-9955-b1499d5701d3")
	if err != nil {
		t.Errorf("Suppressions.GetWithContext returned error: %v", err)
	}

	assert.Equal(t, "suppression", resp.Object)
	assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", resp.Id)
	assert.Equal(t, "steve.wozniak@gmail.com", resp.Email)
	assert.Equal(t, SuppressionOriginComplaint, resp.Origin)
}

func TestGetSuppressionByEmail(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions/steve.wozniak+news@gmail.com", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, "/suppressions/steve.wozniak+news@gmail.com", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"object": "suppression",
			"id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
			"email": "steve.wozniak+news@gmail.com",
			"origin": "manual",
			"source_id": null,
			"created_at": "2023-10-06T23:47:56.678Z"
		}`)
	})

	resp, err := client.Suppressions.Get("steve.wozniak+news@gmail.com")
	if err != nil {
		t.Errorf("Suppressions.Get returned error: %v", err)
	}

	assert.Equal(t, "steve.wozniak+news@gmail.com", resp.Email)
	assert.Equal(t, SuppressionOriginManual, resp.Origin)
	assert.Nil(t, resp.SourceId)
}

func TestGetSuppressionEscapesReservedPathCharacters(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		assert.Equal(t, "/suppressions/steve%2Fwoz%20iak@gmail.com", r.URL.EscapedPath())
		assert.Equal(t, []string{"", "suppressions", "steve%2Fwoz%20iak@gmail.com"}, strings.Split(r.URL.EscapedPath(), "/"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"object": "suppression",
			"id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
			"email": "steve/woz iak@gmail.com",
			"origin": "manual",
			"source_id": null,
			"created_at": "2023-10-06T23:47:56.678Z"
		}`)
	})

	resp, err := client.Suppressions.Get("steve/woz iak@gmail.com")
	if err != nil {
		t.Errorf("Suppressions.Get returned error: %v", err)
	}

	assert.Equal(t, "steve/woz iak@gmail.com", resp.Email)
}

func TestRemoveSuppressionEscapesReservedPathCharacters(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)

		assert.Equal(t, "/suppressions/steve%2Fwozniak@gmail.com", r.URL.EscapedPath())
		assert.Equal(t, []string{"", "suppressions", "steve%2Fwozniak@gmail.com"}, strings.Split(r.URL.EscapedPath(), "/"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"object": "suppression",
			"id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
			"deleted": true
		}`)
	})

	resp, err := client.Suppressions.Remove("steve/wozniak@gmail.com")
	if err != nil {
		t.Errorf("Suppressions.Remove returned error: %v", err)
	}

	assert.True(t, resp.Deleted)
}

func TestGetSuppressionMissingIdentifier(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Suppressions.Get("")
	assert.Error(t, err)
	assert.Equal(t, "[ERROR]: Id or email is required", err.Error())
}

func TestRemoveSuppression(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions/e169aa45-1ecf-4183-9955-b1499d5701d3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"object": "suppression",
			"id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
			"deleted": true
		}`)
	})

	resp, err := client.Suppressions.Remove("e169aa45-1ecf-4183-9955-b1499d5701d3")
	if err != nil {
		t.Errorf("Suppressions.Remove returned error: %v", err)
	}

	assert.Equal(t, "suppression", resp.Object)
	assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", resp.Id)
	assert.True(t, resp.Deleted)
}

func TestRemoveSuppressionWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions/e169aa45-1ecf-4183-9955-b1499d5701d3", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"object": "suppression",
			"id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
			"deleted": true
		}`)
	})

	ctx := context.Background()
	resp, err := client.Suppressions.RemoveWithContext(ctx, "e169aa45-1ecf-4183-9955-b1499d5701d3")
	if err != nil {
		t.Errorf("Suppressions.RemoveWithContext returned error: %v", err)
	}

	assert.Equal(t, "suppression", resp.Object)
	assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", resp.Id)
	assert.True(t, resp.Deleted)
}

func TestRemoveSuppressionMissingIdentifier(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Suppressions.Remove("")
	assert.Error(t, err)
	assert.Equal(t, "[ERROR]: Id or email is required", err.Error())
}

func TestBatchAddSuppressions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions/batch/add", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		assert.Equal(t, []any{"steve.wozniak@gmail.com", "steve.jobs@gmail.com"}, body["emails"])

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, `{
			"data": [
				{ "object": "suppression", "id": "e169aa45-1ecf-4183-9955-b1499d5701d3" },
				{ "object": "suppression", "id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e" }
			]
		}`)
	})

	req := &BatchAddSuppressionsRequest{
		Emails: []string{"steve.wozniak@gmail.com", "steve.jobs@gmail.com"},
	}
	resp, err := client.Suppressions.Batch.Add(req)
	if err != nil {
		t.Errorf("Suppressions.Batch.Add returned error: %v", err)
	}

	assert.Len(t, resp.Data, 2)
	assert.Equal(t, "suppression", resp.Data[0].Object)
	assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", resp.Data[0].Id)
	assert.Equal(t, "b6d24b8e-af0b-4c3c-be0c-359bbd97381e", resp.Data[1].Id)
}

func TestBatchAddSuppressionsWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions/batch/add", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, `{
			"data": [
				{ "object": "suppression", "id": "e169aa45-1ecf-4183-9955-b1499d5701d3" }
			]
		}`)
	})

	ctx := context.Background()
	resp, err := client.Suppressions.Batch.AddWithContext(ctx, &BatchAddSuppressionsRequest{
		Emails: []string{"steve.wozniak@gmail.com"},
	})
	if err != nil {
		t.Errorf("Suppressions.Batch.AddWithContext returned error: %v", err)
	}

	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "suppression", resp.Data[0].Object)
	assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", resp.Data[0].Id)
}

func TestBatchAddSuppressionsMissingEmails(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Suppressions.Batch.Add(&BatchAddSuppressionsRequest{})
	assert.Error(t, err)
	assert.Equal(t, "[ERROR]: Emails is required", err.Error())
}

func TestBatchRemoveSuppressionsWithEmails(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions/batch/remove", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		assert.Equal(t, []any{"steve.wozniak@gmail.com"}, body["emails"])
		assert.NotContains(t, body, "ids")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"data": [
				{ "object": "suppression", "id": "e169aa45-1ecf-4183-9955-b1499d5701d3", "deleted": true }
			]
		}`)
	})

	req := &BatchRemoveSuppressionsRequest{Emails: []string{"steve.wozniak@gmail.com"}}
	resp, err := client.Suppressions.Batch.Remove(req)
	if err != nil {
		t.Errorf("Suppressions.Batch.Remove returned error: %v", err)
	}

	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "suppression", resp.Data[0].Object)
	assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", resp.Data[0].Id)
	assert.True(t, resp.Data[0].Deleted)
}

func TestBatchRemoveSuppressionsWithIds(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions/batch/remove", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		assert.Equal(t, []any{"e169aa45-1ecf-4183-9955-b1499d5701d3"}, body["ids"])
		assert.NotContains(t, body, "emails")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"data": [
				{ "object": "suppression", "id": "e169aa45-1ecf-4183-9955-b1499d5701d3", "deleted": true }
			]
		}`)
	})

	req := &BatchRemoveSuppressionsRequest{Ids: []string{"e169aa45-1ecf-4183-9955-b1499d5701d3"}}
	resp, err := client.Suppressions.Batch.Remove(req)
	if err != nil {
		t.Errorf("Suppressions.Batch.Remove returned error: %v", err)
	}

	assert.Len(t, resp.Data, 1)
	assert.True(t, resp.Data[0].Deleted)
}

func TestBatchRemoveSuppressionsWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/suppressions/batch/remove", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, `{
			"data": [
				{ "object": "suppression", "id": "e169aa45-1ecf-4183-9955-b1499d5701d3", "deleted": true }
			]
		}`)
	})

	ctx := context.Background()
	resp, err := client.Suppressions.Batch.RemoveWithContext(ctx, &BatchRemoveSuppressionsRequest{
		Emails: []string{"steve.wozniak@gmail.com"},
	})
	if err != nil {
		t.Errorf("Suppressions.Batch.RemoveWithContext returned error: %v", err)
	}

	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", resp.Data[0].Id)
	assert.True(t, resp.Data[0].Deleted)
}

func TestBatchRemoveSuppressionsOmitsUnsetKey(t *testing.T) {
	body, err := json.Marshal(&BatchRemoveSuppressionsRequest{Ids: []string{"e169aa45-1ecf-4183-9955-b1499d5701d3"}})
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	assert.JSONEq(t, `{"ids":["e169aa45-1ecf-4183-9955-b1499d5701d3"]}`, string(body))
	assert.NotContains(t, string(body), "emails")
}

func TestBatchRemoveSuppressionsRequiresOneOf(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Suppressions.Batch.Remove(&BatchRemoveSuppressionsRequest{})
	assert.Error(t, err)
	assert.Equal(t, "[ERROR]: Either Emails or Ids is required", err.Error())

	_, err = client.Suppressions.Batch.Remove(&BatchRemoveSuppressionsRequest{
		Emails: []string{"steve.wozniak@gmail.com"},
		Ids:    []string{"e169aa45-1ecf-4183-9955-b1499d5701d3"},
	})
	assert.Error(t, err)
	assert.Equal(t, "[ERROR]: Provide either `emails` or `ids`, but not both.", err.Error())
}

func TestSuppressionsShouldReturnErrorIfContextIsCancelled(t *testing.T) {
	setup()
	defer teardown()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cases := []struct {
		desc string
		call func() error
	}{
		{
			desc: "Add",
			call: func() error {
				_, err := client.Suppressions.AddWithContext(ctx, &AddSuppressionRequest{Email: "steve.wozniak@gmail.com"})
				return err
			},
		},
		{
			desc: "List",
			call: func() error {
				_, err := client.Suppressions.ListWithContext(ctx, nil)
				return err
			},
		},
		{
			desc: "Get",
			call: func() error {
				_, err := client.Suppressions.GetWithContext(ctx, "e169aa45-1ecf-4183-9955-b1499d5701d3")
				return err
			},
		},
		{
			desc: "Remove",
			call: func() error {
				_, err := client.Suppressions.RemoveWithContext(ctx, "e169aa45-1ecf-4183-9955-b1499d5701d3")
				return err
			},
		},
		{
			desc: "Batch.Add",
			call: func() error {
				_, err := client.Suppressions.Batch.AddWithContext(ctx, &BatchAddSuppressionsRequest{
					Emails: []string{"steve.wozniak@gmail.com"},
				})
				return err
			},
		},
		{
			desc: "Batch.Remove",
			call: func() error {
				_, err := client.Suppressions.Batch.RemoveWithContext(ctx, &BatchRemoveSuppressionsRequest{
					Emails: []string{"steve.wozniak@gmail.com"},
				})
				return err
			},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			err := c.call()
			assert.Error(t, err)
			assert.True(t, errors.Is(err, context.Canceled))
		})
	}
}
