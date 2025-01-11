package resend

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBroadcast(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/broadcasts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		var ret interface{}
		ret = `
		{
			"id": "49a3999c-0ce1-4ea6-ab68-afcd6dc2e794"
		}`

		fmt.Fprint(w, ret)
	})

	req := &CreateBroadcastRequest{
		Name:       "New Broadcast",
		AudienceId: "709d076c-2bb1-4be6-94ed-3f8f32622db6",
		From:       "hi@example.com",
		Subject:    "Hello, world!",
	}
	resp, err := client.Broadcasts.Create(req)
	if err != nil {
		t.Errorf("Broadcasts.Create returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "49a3999c-0ce1-4ea6-ab68-afcd6dc2e794")
}

func TestCreateBroadcastValidations(t *testing.T) {
	setup()
	defer teardown()

	req1 := &CreateBroadcastRequest{
		Name:       "New Broadcast",
		AudienceId: "709d076c-2bb1-4be6-94ed-3f8f32622db6",
		From:       "",
	}
	_, err := client.Broadcasts.Create(req1)
	assert.NotNil(t, err)
	if err != nil {
		assert.Equal(t, err.Error(), "[ERROR]: From cannot be empty")
	}

	req2 := &CreateBroadcastRequest{
		Name: "New Broadcast",
		From: "hi@example.com",
	}
	_, err = client.Broadcasts.Create(req2)
	assert.NotNil(t, err)
	if err != nil {
		assert.Equal(t, err.Error(), "[ERROR]: AudienceId cannot be empty")
	}

	req3 := &CreateBroadcastRequest{
		Name:       "New Broadcast",
		From:       "hi@example.com",
		AudienceId: "709d076c-2bb1-4be6-94ed-3f8f32622db6",
		Subject:    "",
	}
	_, err = client.Broadcasts.Create(req3)
	assert.NotNil(t, err)
	if err != nil {
		assert.Equal(t, err.Error(), "[ERROR]: Subject cannot be empty")
	}
}

func TestGetBroadcast(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/broadcasts/559ac32e-9ef5-46fb-82a1-b76b840c0f7b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "broadcast",
			"id": "559ac32e-9ef5-46fb-82a1-b76b840c0f7b",
			"name": "Announcements",
			"audience_id": "78261eea-8f8b-4381-83c6-79fa7120f1cf",
			"from": "Acme <onboarding@resend.dev>",
			"subject": "hello world",
			"reply_to": null,
			"preview_text": "Check out our latest announcements",
			"status": "draft",
			"created_at": "2024-12-01T19:32:22.980Z",
			"scheduled_at": null,
			"sent_at": null
		}`

		fmt.Fprint(w, ret)
	})

	b, err := client.Broadcasts.Get("559ac32e-9ef5-46fb-82a1-b76b840c0f7b")
	if err != nil {
		t.Errorf("Broadcast.Get returned error: %v", err)
	}

	assert.Equal(t, b.Id, "559ac32e-9ef5-46fb-82a1-b76b840c0f7b")
	assert.Equal(t, b.Object, "broadcast")
	assert.Equal(t, b.Name, "Announcements")
	assert.Equal(t, b.AudienceId, "78261eea-8f8b-4381-83c6-79fa7120f1cf")
	assert.Equal(t, b.From, "Acme <onboarding@resend.dev>")
	assert.Equal(t, b.Subject, "hello world")
	assert.Equal(t, b.PreviewText, "Check out our latest announcements")
	assert.Equal(t, b.Status, "draft")
	assert.Equal(t, b.CreatedAt, "2024-12-01T19:32:22.980Z")
}

func TestGetBroadcastValidations(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Broadcasts.Get("")
	assert.NotNil(t, err)
	if err != nil {
		assert.Equal(t, err.Error(), "[ERROR]: broadcastId cannot be empty")
	}
}

func TestSendBroadcast(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/broadcasts/559ac32e-9ef5-46fb-82a1-b76b840c0f7b/send", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "49a3999c-0ce1-4ea6-ab68-afcd6dc2e794"
		}`

		fmt.Fprint(w, ret)
	})

	req := &SendBroadcastRequest{
		BroadcastId: "559ac32e-9ef5-46fb-82a1-b76b840c0f7b",
	}

	b, err := client.Broadcasts.Send(req)
	if err != nil {
		t.Errorf("Broadcast.Send returned error: %v", err)
	}

	assert.Equal(t, b.Id, "49a3999c-0ce1-4ea6-ab68-afcd6dc2e794")
}

func TestSendBroadcastValidations(t *testing.T) {
	setup()
	defer teardown()

	req1 := &SendBroadcastRequest{
		BroadcastId: "",
	}

	_, err := client.Broadcasts.Send(req1)
	assert.NotNil(t, err)
	if err != nil {
		assert.Equal(t, err.Error(), "[ERROR]: BroadcastId cannot be empty")
	}
}

func TestRemoveBroadcast(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/broadcasts/b6d24b8e-af0b-4c3c-be0c-359bbd97381e", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusOK)

		var ret interface{}
		ret = `
		{
			"object": "broadcast",
			"id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e",
			"deleted": true
		}`

		fmt.Fprint(w, ret)
	})

	deleted, err := client.Broadcasts.Remove("b6d24b8e-af0b-4c3c-be0c-359bbd97381e")
	if err != nil {
		t.Errorf("Broadcasts.Remove returned error: %v", err)
	}
	assert.True(t, deleted.Deleted)
	assert.Equal(t, deleted.Id, "b6d24b8e-af0b-4c3c-be0c-359bbd97381e")
	assert.Equal(t, deleted.Object, "broadcast")
}

func TestListBroadcasts(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/broadcasts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
				"data": [
					{
						"id": "49a3999c-0ce1-4ea6-ab68-afcd6dc2e794",
						"audience_id": "78261eea-8f8b-4381-83c6-79fa7120f1cf",
						"status": "draft",
						"created_at": "2024-11-01T15:13:31.723Z",
						"scheduled_at": null,
						"sent_at": null
					},
					{
						"id": "559ac32e-9ef5-46fb-82a1-b76b840c0f7b",
						"audience_id": "78261eea-8f8b-4381-83c6-79fa7120f1cf",
						"status": "sent",
						"created_at": "2024-12-01T19:32:22.980Z",
						"scheduled_at": "2024-12-02T19:32:22.980Z",
						"sent_at": "2024-12-02T19:32:22.980Z"
					}
				]
		}`

		fmt.Fprint(w, ret)
	})

	broadcasts, err := client.Broadcasts.List()
	if err != nil {
		t.Errorf("Broadcasts.List returned error: %v", err)
	}

	assert.Equal(t, len(broadcasts.Data), 2)
	assert.Equal(t, broadcasts.Object, "list")

}
