package resend

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateContactImport(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/contacts/imports", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"object":"contact_import","id":"479e3145-dd38-476b-932c-529ceb705947"}`)
	})

	req := &CreateContactImportRequest{
		File: []byte("email,first_name\nsteve@example.com,Steve"),
	}
	resp, err := client.Contacts.Imports.Create(req)
	if err != nil {
		t.Errorf("ContactImports.Create returned error: %v", err)
	}
	assert.Equal(t, "contact_import", resp.Object)
	assert.Equal(t, "479e3145-dd38-476b-932c-529ceb705947", resp.Id)
}

func TestCreateContactImportWithOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/contacts/imports", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		r.ParseMultipartForm(10 << 20)
		assert.Equal(t, "upsert", r.FormValue("on_conflict"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"object":"contact_import","id":"479e3145-dd38-476b-932c-529ceb705947"}`)
	})

	req := &CreateContactImportRequest{
		File:       []byte("email\nsteve@example.com"),
		OnConflict: "upsert",
		ColumnMap:  map[string]any{"email": "Email"},
		Segments:   []string{"seg-123"},
	}
	resp, err := client.Contacts.Imports.Create(req)
	if err != nil {
		t.Errorf("ContactImports.Create returned error: %v", err)
	}
	assert.Equal(t, "479e3145-dd38-476b-932c-529ceb705947", resp.Id)
}

func TestCreateContactImportMissingFile(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Contacts.Imports.Create(&CreateContactImportRequest{})
	assert.Error(t, err)
}

func TestGetContactImport(t *testing.T) {
	setup()
	defer teardown()

	importId := "479e3145-dd38-476b-932c-529ceb705947"

	mux.HandleFunc("/contacts/imports/"+importId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"object": "contact_import",
			"id": "%s",
			"status": "completed",
			"created_at": "2023-10-06T23:47:56.678Z",
			"counts": {"total": 100, "created": 80, "updated": 10, "skipped": 5, "failed": 5}
		}`, importId)
	})

	resp, err := client.Contacts.Imports.Get(importId)
	if err != nil {
		t.Errorf("ContactImports.Get returned error: %v", err)
	}
	assert.Equal(t, "contact_import", resp.Object)
	assert.Equal(t, importId, resp.Id)
	assert.Equal(t, ContactImportStatusCompleted, resp.Status)
	assert.NotNil(t, resp.Counts)
	assert.Equal(t, 100, resp.Counts.Total)
	assert.Equal(t, 80, resp.Counts.Created)
}

func TestGetContactImportMissingId(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Contacts.Imports.Get("")
	assert.Error(t, err)
}

func TestListContactImports(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/contacts/imports", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "list",
			"has_more": false,
			"data": [
				{
					"object": "contact_import",
					"id": "479e3145-dd38-476b-932c-529ceb705947",
					"status": "completed",
					"created_at": "2023-10-06T23:47:56.678Z"
				}
			]
		}`)
	})

	resp, err := client.Contacts.Imports.List(nil)
	if err != nil {
		t.Errorf("ContactImports.List returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.False(t, resp.HasMore)
	assert.Len(t, resp.Data, 1)
	assert.Equal(t, "479e3145-dd38-476b-932c-529ceb705947", resp.Data[0].Id)
	assert.Equal(t, ContactImportStatusCompleted, resp.Data[0].Status)
}

func TestListContactImportsWithStatusFilter(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/contacts/imports", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, "completed", r.URL.Query().Get("status"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"object":"list","has_more":false,"data":[]}`)
	})

	resp, err := client.Contacts.Imports.List(&ListContactImportsOptions{Status: "completed"})
	if err != nil {
		t.Errorf("ContactImports.List returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
}
