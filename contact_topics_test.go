package resend

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListContactTopics(t *testing.T) {
	setup()
	defer teardown()

	contactId := "e169aa45-1ecf-4183-9955-b1499d5701d3"

	mux.HandleFunc("/contacts/"+contactId+"/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
			"has_more": false,
			"data": [
				{
					"id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
					"name": "Product Updates",
					"description": "New features, and latest announcements.",
					"subscription": "opt_in"
				},
				{
					"id": "07d84122-7224-4881-9c31-1c048e204602",
					"name": "Newsletter",
					"description": "Weekly newsletter with tips and tricks.",
					"subscription": "opt_out"
				}
			]
		}`

		fmt.Fprint(w, ret)
	})

	topics, err := client.Contacts.Topics.List(contactId)
	if err != nil {
		t.Errorf("ContactTopics.List returned error: %v", err)
	}

	assert.Equal(t, "list", topics.Object)
	assert.False(t, topics.HasMore)
	assert.Equal(t, 2, len(topics.Data))

	// Check first topic
	assert.Equal(t, "b6d24b8e-af0b-4c3c-be0c-359bbd97381e", topics.Data[0].Id)
	assert.Equal(t, "Product Updates", topics.Data[0].Name)
	assert.Equal(t, "New features, and latest announcements.", topics.Data[0].Description)
	assert.Equal(t, "opt_in", topics.Data[0].Subscription)

	// Check second topic
	assert.Equal(t, "07d84122-7224-4881-9c31-1c048e204602", topics.Data[1].Id)
	assert.Equal(t, "Newsletter", topics.Data[1].Name)
	assert.Equal(t, "opt_out", topics.Data[1].Subscription)
}

func TestListContactTopicsByEmail(t *testing.T) {
	setup()
	defer teardown()

	contactEmail := "steve.wozniak@gmail.com"

	mux.HandleFunc("/contacts/"+contactEmail+"/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
			"has_more": false,
			"data": [
				{
					"id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
					"name": "Product Updates",
					"description": "New features, and latest announcements.",
					"subscription": "opt_in"
				}
			]
		}`

		fmt.Fprint(w, ret)
	})

	topics, err := client.Contacts.Topics.List(contactEmail)
	if err != nil {
		t.Errorf("ContactTopics.List returned error: %v", err)
	}

	assert.Equal(t, "list", topics.Object)
	assert.False(t, topics.HasMore)
	assert.Equal(t, 1, len(topics.Data))
	assert.Equal(t, "b6d24b8e-af0b-4c3c-be0c-359bbd97381e", topics.Data[0].Id)
}

func TestListContactTopicsIdMissing(t *testing.T) {
	setup()
	defer teardown()

	topics, err := client.Contacts.Topics.List("")

	assert.Error(t, err)
	assert.Equal(t, ListContactTopicsResponse{}, topics)
	assert.Contains(t, err.Error(), "[ERROR]: Contact ID or email is missing")
}

func TestUpdateContactTopicsById(t *testing.T) {
	setup()
	defer teardown()

	contactId := "e169aa45-1ecf-4183-9955-b1499d5701d3"

	mux.HandleFunc("/contacts/"+contactId+"/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "e169aa45-1ecf-4183-9955-b1499d5701d3"
		}`

		fmt.Fprint(w, ret)
	})

	params := &UpdateContactTopicsRequest{
		Id: contactId,
		Topics: []TopicSubscriptionUpdate{
			{
				Id:           "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
				Subscription: "opt_out",
			},
			{
				Id:           "07d84122-7224-4881-9c31-1c048e204602",
				Subscription: "opt_in",
			},
		},
	}

	resp, err := client.Contacts.Topics.Update(params)
	if err != nil {
		t.Errorf("ContactTopics.Update returned error: %v", err)
	}

	assert.Equal(t, contactId, resp.Id)
}

func TestUpdateContactTopicsByEmail(t *testing.T) {
	setup()
	defer teardown()

	contactEmail := "steve.wozniak@gmail.com"

	mux.HandleFunc("/contacts/"+contactEmail+"/topics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "e169aa45-1ecf-4183-9955-b1499d5701d3"
		}`

		fmt.Fprint(w, ret)
	})

	params := &UpdateContactTopicsRequest{
		Email: contactEmail,
		Topics: []TopicSubscriptionUpdate{
			{
				Id:           "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
				Subscription: "opt_in",
			},
		},
	}

	resp, err := client.Contacts.Topics.Update(params)
	if err != nil {
		t.Errorf("ContactTopics.Update returned error: %v", err)
	}

	assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", resp.Id)
}

func TestUpdateContactTopicsIdMissing(t *testing.T) {
	setup()
	defer teardown()

	params := &UpdateContactTopicsRequest{
		Topics: []TopicSubscriptionUpdate{
			{
				Id:           "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
				Subscription: "opt_in",
			},
		},
	}

	resp, err := client.Contacts.Topics.Update(params)

	assert.Error(t, err)
	assert.Equal(t, UpdateContactTopicsResponse{}, resp)
	assert.Contains(t, err.Error(), "[ERROR]: Contact ID or email is missing")
}

func TestUpdateContactTopicsEmptyArray(t *testing.T) {
	setup()
	defer teardown()

	params := &UpdateContactTopicsRequest{
		Id:     "e169aa45-1ecf-4183-9955-b1499d5701d3",
		Topics: []TopicSubscriptionUpdate{},
	}

	resp, err := client.Contacts.Topics.Update(params)

	assert.Error(t, err)
	assert.Equal(t, UpdateContactTopicsResponse{}, resp)
	assert.Contains(t, err.Error(), "[ERROR]: Topics array is empty")
}
