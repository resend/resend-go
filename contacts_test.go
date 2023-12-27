package resend

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateContact(t *testing.T) {
	setup()
	defer teardown()

	audienceId := "709d076c-2bb1-4be6-94ed-3f8f32622db6"

	mux.HandleFunc("/audiences/"+audienceId+"/contacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		var ret interface{}
		ret = `
		{
			"object": "contact",
			"id": "479e3145-dd38-476b-932c-529ceb705947"
		}`

		fmt.Fprint(w, ret)
	})

	req := &CreateContactRequest{
		Email:      "email@example.com",
		AudienceId: audienceId,
	}
	resp, err := client.Contacts.Create(req)
	if err != nil {
		t.Errorf("Contacts.Create returned error: %v", err)
	}
	assert.Equal(t, resp.Object, "contact")
	assert.Equal(t, resp.Id, "479e3145-dd38-476b-932c-529ceb705947")
}

func TestListContacts(t *testing.T) {
	setup()
	defer teardown()

	audienceId := "709d076c-2bb1-4be6-94ed-3f8f32622db6"

	mux.HandleFunc("/audiences/"+audienceId+"/contacts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
			"data": [
				{
					"id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
					"email": "steve.wozniak@gmail.com",
					"first_name": "Steve",
					"last_name": "Wozniak",
					"created_at": "2023-10-06T23:47:56.678Z",
					"unsubscribed": false
				}
			]
		}`

		fmt.Fprint(w, ret)
	})

	contacts, err := client.Contacts.List(audienceId)
	if err != nil {
		t.Errorf("Contacts.List returned error: %v", err)
	}

	assert.Equal(t, len(contacts.Data), 1)
	assert.Equal(t, contacts.Data[0].Id, "e169aa45-1ecf-4183-9955-b1499d5701d3")
	assert.Equal(t, contacts.Data[0].FirstName, "Steve")
	assert.Equal(t, contacts.Data[0].LastName, "Wozniak")
	assert.Equal(t, contacts.Data[0].CreatedAt, "2023-10-06T23:47:56.678Z")
	assert.Equal(t, contacts.Data[0].Unsubscribed, false)
}

func TestRemoveContact(t *testing.T) {
	setup()
	defer teardown()

	audienceId := "709d076c-2bb1-4be6-94ed-3f8f32622db6"
	contactId := "e169aa45-1ecf-4183-9955-b1499d5701d3"

	mux.HandleFunc("/audiences/"+audienceId+"/contacts/"+contactId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusOK)

		var ret interface{}
		ret = `
		{
			"object": "contact",
			"id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
			"deleted": true
		}`

		fmt.Fprint(w, ret)
	})

	deleted, err := client.Contacts.Remove(audienceId, contactId)
	if err != nil {
		t.Errorf("Domains.Remove returned error: %v", err)
	}
	assert.True(t, deleted.Deleted)
}

func TestGetContact(t *testing.T) {
	setup()
	defer teardown()

	audienceId := "709d076c-2bb1-4be6-94ed-3f8f32622db6"
	contactId := "e169aa45-1ecf-4183-9955-b1499d5701d3"

	mux.HandleFunc("/audiences/"+audienceId+"/contacts/"+contactId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "contact",
			"id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
			"email": "steve.wozniak@gmail.com",
			"first_name": "Steve",
			"last_name": "Wozniak",
			"created_at": "2023-10-06T23:47:56.678Z",
			"unsubscribed": false
		}`

		fmt.Fprint(w, ret)
	})

	contact, err := client.Contacts.Get(audienceId, contactId)
	if err != nil {
		t.Errorf("Contacts.Get returned error: %v", err)
	}

	assert.Equal(t, contact.Id, contactId)
	assert.Equal(t, contact.Object, "contact")
	assert.Equal(t, contact.FirstName, "Steve")
	assert.Equal(t, contact.LastName, "Wozniak")
	assert.Equal(t, contact.CreatedAt, "2023-10-06T23:47:56.678Z")
	assert.Equal(t, contact.Unsubscribed, false)
}

func TestUpdateContact(t *testing.T) {
	setup()
	defer teardown()

	audienceId := "709d076c-2bb1-4be6-94ed-3f8f32622db6"
	id := "109d077c-2bb1-4be6-84ed-ff8f32612dc2"

	mux.HandleFunc("/audiences/"+audienceId+"/contacts/"+id, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		var ret interface{}
		ret = `
		{
			"data": {
				"id": "479e3145-dd38-476b-932c-529ceb705947"
			},
			"error": null
		}`

		fmt.Fprint(w, ret)
	})

	req := &UpdateContactRequest{
		AudienceId: audienceId,
		Id:         id,
		FirstName:  "Updated First Name",
	}
	resp, err := client.Contacts.Update(req)
	if err != nil {
		t.Errorf("Contacts.Update returned error: %v", err)
	}
	assert.NotNil(t, resp.Data)
	assert.Equal(t, resp.Data.Id, "479e3145-dd38-476b-932c-529ceb705947")
	assert.Equal(t, resp.Error, struct{}{})
}
