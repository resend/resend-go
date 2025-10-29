package resend

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateContactProperty(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/contact-properties", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		ret := `
		{
			"object": "contact_property",
			"id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
		}`

		fmt.Fprint(w, ret)
	})

	req := &CreateContactPropertyRequest{
		Key:           "age",
		Type:          "number",
		FallbackValue: 0,
	}
	resp, err := client.ContactProperties.Create(req)
	if err != nil {
		t.Errorf("ContactProperties.Create returned error: %v", err)
	}
	assert.Equal(t, "contact_property", resp.Object)
	assert.Equal(t, "a1b2c3d4-e5f6-7890-abcd-ef1234567890", resp.Id)
}

func TestCreateContactPropertyKeyMissing(t *testing.T) {
	setup()
	defer teardown()

	req := &CreateContactPropertyRequest{
		Type:          "string",
		FallbackValue: "default",
	}
	resp, err := client.ContactProperties.Create(req)

	assert.Error(t, err)
	assert.Equal(t, CreateContactPropertyResponse{}, resp)
	assert.Contains(t, err.Error(), "[ERROR]: Key is missing")
}

func TestCreateContactPropertyTypeMissing(t *testing.T) {
	setup()
	defer teardown()

	req := &CreateContactPropertyRequest{
		Key:           "age",
		FallbackValue: 0,
	}
	resp, err := client.ContactProperties.Create(req)

	assert.Error(t, err)
	assert.Equal(t, CreateContactPropertyResponse{}, resp)
	assert.Contains(t, err.Error(), "[ERROR]: Type is missing")
}

func TestListContactProperties(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/contact-properties", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
			"data": [
				{
					"id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
					"key": "age",
					"object": "contact_property",
					"created_at": "2025-10-22T15:30:00.000Z",
					"type": "number",
					"fallback_value": 0
				},
				{
					"id": "b2c3d4e5-f6a7-8901-bcde-f12345678901",
					"key": "country",
					"object": "contact_property",
					"created_at": "2025-10-22T15:31:00.000Z",
					"type": "string",
					"fallback_value": "US"
				}
			],
			"has_more": false
		}`

		fmt.Fprint(w, ret)
	})

	properties, err := client.ContactProperties.List()
	if err != nil {
		t.Errorf("ContactProperties.List returned error: %v", err)
	}

	assert.Equal(t, "list", properties.Object)
	assert.Equal(t, 2, len(properties.Data))
	assert.False(t, properties.HasMore)

	// Check first property
	assert.Equal(t, "a1b2c3d4-e5f6-7890-abcd-ef1234567890", properties.Data[0].Id)
	assert.Equal(t, "age", properties.Data[0].Key)
	assert.Equal(t, "contact_property", properties.Data[0].Object)
	assert.Equal(t, "number", properties.Data[0].Type)
	assert.Equal(t, float64(0), properties.Data[0].FallbackValue)
	assert.Equal(t, "2025-10-22T15:30:00.000Z", properties.Data[0].CreatedAt)

	// Check second property
	assert.Equal(t, "b2c3d4e5-f6a7-8901-bcde-f12345678901", properties.Data[1].Id)
	assert.Equal(t, "country", properties.Data[1].Key)
	assert.Equal(t, "string", properties.Data[1].Type)
	assert.Equal(t, "US", properties.Data[1].FallbackValue)
}

func TestGetContactProperty(t *testing.T) {
	setup()
	defer teardown()

	propertyId := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"

	mux.HandleFunc("/contact-properties/"+propertyId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			"key": "age",
			"object": "contact_property",
			"created_at": "2025-10-22T15:30:00.000Z",
			"type": "number",
			"fallback_value": 0
		}`

		fmt.Fprint(w, ret)
	})

	property, err := client.ContactProperties.Get(propertyId)
	if err != nil {
		t.Errorf("ContactProperties.Get returned error: %v", err)
	}

	assert.Equal(t, propertyId, property.Id)
	assert.Equal(t, "age", property.Key)
	assert.Equal(t, "contact_property", property.Object)
	assert.Equal(t, "number", property.Type)
	assert.Equal(t, float64(0), property.FallbackValue)
	assert.Equal(t, "2025-10-22T15:30:00.000Z", property.CreatedAt)
}

func TestGetContactPropertyIdMissing(t *testing.T) {
	setup()
	defer teardown()

	property, err := client.ContactProperties.Get("")

	assert.Error(t, err)
	assert.Equal(t, ContactProperty{}, property)
	assert.Contains(t, err.Error(), "[ERROR]: ID is missing")
}

func TestUpdateContactProperty(t *testing.T) {
	setup()
	defer teardown()

	propertyId := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"

	mux.HandleFunc("/contact-properties/"+propertyId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			"object": "contact_property"
		}`

		fmt.Fprint(w, ret)
	})

	req := &UpdateContactPropertyRequest{
		Id:            propertyId,
		FallbackValue: 25,
	}
	resp, err := client.ContactProperties.Update(req)
	if err != nil {
		t.Errorf("ContactProperties.Update returned error: %v", err)
	}

	assert.Equal(t, propertyId, resp.Id)
	assert.Equal(t, "contact_property", resp.Object)
}

func TestUpdateContactPropertyIdMissing(t *testing.T) {
	setup()
	defer teardown()

	req := &UpdateContactPropertyRequest{
		FallbackValue: "new value",
	}
	resp, err := client.ContactProperties.Update(req)

	assert.Error(t, err)
	assert.Equal(t, UpdateContactPropertyResponse{}, resp)
	assert.Contains(t, err.Error(), "[ERROR]: ID is missing")
}

func TestRemoveContactProperty(t *testing.T) {
	setup()
	defer teardown()

	propertyId := "a1b2c3d4-e5f6-7890-abcd-ef1234567890"

	mux.HandleFunc("/contact-properties/"+propertyId, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "contact_property",
			"id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
			"deleted": true
		}`

		fmt.Fprint(w, ret)
	})

	deleted, err := client.ContactProperties.Remove(propertyId)
	if err != nil {
		t.Errorf("ContactProperties.Remove returned error: %v", err)
	}

	assert.Equal(t, propertyId, deleted.Id)
	assert.Equal(t, "contact_property", deleted.Object)
	assert.True(t, deleted.Deleted)
}

func TestRemoveContactPropertyIdMissing(t *testing.T) {
	setup()
	defer teardown()

	deleted, err := client.ContactProperties.Remove("")

	assert.Error(t, err)
	assert.Equal(t, RemoveContactPropertyResponse{}, deleted)
	assert.Contains(t, err.Error(), "[ERROR]: ID is missing")
}
