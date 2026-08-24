package resend

import (
	"context"
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

		var ret any
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

func TestCreateAndSendBroadcast(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/broadcasts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		var ret any
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
		Send:       true,
	}
	resp, err := client.Broadcasts.Create(req)
	if err != nil {
		t.Errorf("Broadcasts.Create returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "49a3999c-0ce1-4ea6-ab68-afcd6dc2e794")
}

func TestCreateAndScheduleBroadcast(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/broadcasts", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		var ret any
		ret = `
		{
			"id": "49a3999c-0ce1-4ea6-ab68-afcd6dc2e794"
		}`

		fmt.Fprint(w, ret)
	})

	req := &CreateBroadcastRequest{
		Name:        "New Broadcast",
		AudienceId:  "709d076c-2bb1-4be6-94ed-3f8f32622db6",
		From:        "hi@example.com",
		Subject:     "Hello, world!",
		Send:        true,
		ScheduledAt: "2024-12-01T19:32:22.980Z",
	}
	resp, err := client.Broadcasts.Create(req)
	if err != nil {
		t.Errorf("Broadcasts.Create returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "49a3999c-0ce1-4ea6-ab68-afcd6dc2e794")
}

func TestUpdateBroadcast(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/broadcasts/559ac32e-9ef5-46fb-82a1-b76b840c0f7b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "559ac32e-9ef5-46fb-82a1-b76b840c0f7b"
		}`

		fmt.Fprint(w, ret)
	})

	req := &UpdateBroadcastRequest{
		BroadcastId: "559ac32e-9ef5-46fb-82a1-b76b840c0f7b",
		Name:        "Updated Broadcast",
	}
	resp, err := client.Broadcasts.Update(req)
	if err != nil {
		t.Errorf("Broadcasts.Update returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "559ac32e-9ef5-46fb-82a1-b76b840c0f7b")
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
		assert.Equal(t, err.Error(), "[ERROR]: Either SegmentId or AudienceId must be provided")
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
			"created_at": "2024-12-01 19:32:22.98+00",
			"scheduled_at": null,
			"sent_at": null,
			"html": "<h1>Hello world</h1>",
			"text": "Hello world"
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
	assert.Equal(t, b.CreatedAt, "2024-12-01 19:32:22.98+00")
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

		var ret any
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

func TestCancelBroadcast(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/broadcasts/b6d24b8e-af0b-4c3c-be0c-359bbd97381e/cancel", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.WriteHeader(http.StatusOK)

		var ret any
		ret = `
		{
			"object": "broadcast",
			"id": "b6d24b8e-af0b-4c3c-be0c-359bbd97381e"
		}`

		fmt.Fprint(w, ret)
	})

	canceled, err := client.Broadcasts.Cancel("b6d24b8e-af0b-4c3c-be0c-359bbd97381e")
	if err != nil {
		t.Errorf("Broadcasts.Cancel returned error: %v", err)
	}
	assert.Equal(t, canceled.Id, "b6d24b8e-af0b-4c3c-be0c-359bbd97381e")
	assert.Equal(t, canceled.Object, "broadcast")
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
						"created_at": "2024-11-01 15:13:31.723+00",
						"scheduled_at": null,
						"sent_at": null
					},
					{
						"id": "559ac32e-9ef5-46fb-82a1-b76b840c0f7b",
						"audience_id": "78261eea-8f8b-4381-83c6-79fa7120f1cf",
						"status": "sent",
						"created_at": "2024-12-01 19:32:22.98+00",
						"scheduled_at": "2024-12-02 19:32:22.98+00",
						"sent_at": "2024-12-02 19:32:22.98+00"
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

func TestBroadcastRecipients(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/broadcasts/559ac32e-9ef5-46fb-82a1-b76b840c0f7b/recipients", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, "sent", r.URL.Query().Get("type"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
			"has_more": false,
			"data": [
				{
					"id": "b2Zmc2V0OjA",
					"contact_id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
					"email": "carter@example.com"
				}
			]
		}`

		fmt.Fprint(w, ret)
	})

	resp, err := client.Broadcasts.Recipients("559ac32e-9ef5-46fb-82a1-b76b840c0f7b", &ListBroadcastRecipientsOptions{
		Type: BroadcastRecipientEventTypeSent,
	})
	if err != nil {
		t.Errorf("Broadcasts.Recipients returned error: %v", err)
	}

	assert.Equal(t, "list", resp.Object)
	assert.False(t, resp.HasMore)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "b2Zmc2V0OjA", resp.Data[0].Id)
	assert.Equal(t, "carter@example.com", resp.Data[0].Email)
	assert.NotNil(t, resp.Data[0].ContactId)
	assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", *resp.Data[0].ContactId)
	assert.Equal(t, 0, resp.Data[0].Count)
	assert.Equal(t, "", resp.Data[0].BounceType)
	assert.Nil(t, resp.Data[0].ClickedLinks)
}

func TestBroadcastRecipientsWithFilters(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/broadcasts/559ac32e-9ef5-46fb-82a1-b76b840c0f7b/recipients", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		q := r.URL.Query()
		assert.Equal(t, "opened", q.Get("type"))
		assert.Equal(t, "carter", q.Get("email"))
		assert.Equal(t, "10", q.Get("limit"))
		assert.Equal(t, "b2Zmc2V0OjA", q.Get("after"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
			"has_more": true,
			"data": [
				{
					"id": "b2Zmc2V0OjE",
					"contact_id": null,
					"email": "carter@example.com",
					"count": 3
				}
			]
		}`

		fmt.Fprint(w, ret)
	})

	limit := 10
	after := "b2Zmc2V0OjA"
	resp, err := client.Broadcasts.Recipients("559ac32e-9ef5-46fb-82a1-b76b840c0f7b", &ListBroadcastRecipientsOptions{
		Type:  BroadcastRecipientEventTypeOpened,
		Email: "carter",
		Limit: &limit,
		After: &after,
	})
	if err != nil {
		t.Errorf("Broadcasts.Recipients returned error: %v", err)
	}

	assert.True(t, resp.HasMore)
	assert.Equal(t, 1, len(resp.Data))
	assert.Nil(t, resp.Data[0].ContactId)
	assert.Equal(t, 3, resp.Data[0].Count)
}

func TestBroadcastRecipientsClicked(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/broadcasts/559ac32e-9ef5-46fb-82a1-b76b840c0f7b/recipients", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, "clicked", r.URL.Query().Get("type"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
			"has_more": false,
			"data": [
				{
					"id": "b2Zmc2V0OjA",
					"contact_id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
					"email": "carter@example.com",
					"count": 3,
					"clicked_links": [
						{"url": "https://resend.com/pricing", "clicks": 2}
					]
				}
			]
		}`

		fmt.Fprint(w, ret)
	})

	resp, err := client.Broadcasts.Recipients("559ac32e-9ef5-46fb-82a1-b76b840c0f7b", &ListBroadcastRecipientsOptions{
		Type: BroadcastRecipientEventTypeClicked,
	})
	if err != nil {
		t.Errorf("Broadcasts.Recipients returned error: %v", err)
	}

	assert.Equal(t, 3, resp.Data[0].Count)
	assert.Equal(t, 1, len(resp.Data[0].ClickedLinks))
	assert.Equal(t, "https://resend.com/pricing", resp.Data[0].ClickedLinks[0].Url)
	assert.Equal(t, 2, resp.Data[0].ClickedLinks[0].Clicks)
}

func TestBroadcastRecipientsBounced(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/broadcasts/559ac32e-9ef5-46fb-82a1-b76b840c0f7b/recipients", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		q := r.URL.Query()
		assert.Equal(t, "bounced", q.Get("type"))
		assert.Equal(t, "permanent", q.Get("bounce_type"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
			"has_more": false,
			"data": [
				{
					"id": "b2Zmc2V0OjA",
					"contact_id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
					"email": "carter@example.com",
					"bounce_type": "permanent"
				}
			]
		}`

		fmt.Fprint(w, ret)
	})

	resp, err := client.Broadcasts.Recipients("559ac32e-9ef5-46fb-82a1-b76b840c0f7b", &ListBroadcastRecipientsOptions{
		Type:       BroadcastRecipientEventTypeBounced,
		BounceType: BroadcastRecipientBounceTypePermanent,
	})
	if err != nil {
		t.Errorf("Broadcasts.Recipients returned error: %v", err)
	}

	assert.Equal(t, BroadcastRecipientBounceTypePermanent, resp.Data[0].BounceType)
}

func TestBroadcastRecipientsNotFound(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/broadcasts/00000000-0000-0000-0000-000000000000/recipients", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		ret := `
		{
			"statusCode": 404,
			"name": "not_found",
			"message": "Broadcast not found"
		}`

		fmt.Fprint(w, ret)
	})

	resp, err := client.Broadcasts.Recipients("00000000-0000-0000-0000-000000000000", &ListBroadcastRecipientsOptions{
		Type: BroadcastRecipientEventTypeSent,
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Broadcast not found")
	assert.Equal(t, ListBroadcastRecipientsResponse{}, resp)
}

func TestBroadcastRecipientsValidations(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Broadcasts.Recipients("", &ListBroadcastRecipientsOptions{
		Type: BroadcastRecipientEventTypeSent,
	})
	assert.NotNil(t, err)
	if err != nil {
		assert.Equal(t, err.Error(), "[ERROR]: broadcastId cannot be empty")
	}

	_, err = client.Broadcasts.Recipients("559ac32e-9ef5-46fb-82a1-b76b840c0f7b", nil)
	assert.NotNil(t, err)
	if err != nil {
		assert.Equal(t, err.Error(), "[ERROR]: Type cannot be empty")
	}

	_, err = client.Broadcasts.Recipients("559ac32e-9ef5-46fb-82a1-b76b840c0f7b", &ListBroadcastRecipientsOptions{})
	assert.NotNil(t, err)
	if err != nil {
		assert.Equal(t, err.Error(), "[ERROR]: Type cannot be empty")
	}
}

func TestBroadcastsClickedLinks(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/broadcasts/559ac32e-9ef5-46fb-82a1-b76b840c0f7b/clicked-links", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
			"has_more": false,
			"data": [
				{
					"id": "b2Zmc2V0OjA",
					"url": "https://resend.com/pricing",
					"clicks": 42,
					"unique_clicks": 30
				},
				{
					"id": "b2Zmc2V0OjE",
					"url": "https://resend.com/docs",
					"clicks": 17,
					"unique_clicks": 15
				}
			]
		}`

		fmt.Fprint(w, ret)
	})

	clickedLinks, err := client.Broadcasts.ClickedLinks("559ac32e-9ef5-46fb-82a1-b76b840c0f7b")
	if err != nil {
		t.Errorf("Broadcasts.ClickedLinks returned error: %v", err)
	}

	assert.Equal(t, clickedLinks.Object, "list")
	assert.Equal(t, len(clickedLinks.Data), 2)
	assert.Equal(t, clickedLinks.Data[0].Url, "https://resend.com/pricing")
	assert.Equal(t, clickedLinks.Data[0].Clicks, 42)
	assert.Equal(t, clickedLinks.Data[0].UniqueClicks, 30)
}

func TestBroadcastsClickedLinksWithOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/broadcasts/559ac32e-9ef5-46fb-82a1-b76b840c0f7b/clicked-links", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, r.URL.Query().Get("limit"), "1")
		assert.Equal(t, r.URL.Query().Get("after"), "cursor-value")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
			"has_more": true,
			"data": [
				{
					"id": "b2Zmc2V0OjA",
					"url": "https://resend.com/pricing",
					"clicks": 42,
					"unique_clicks": 30
				}
			]
		}`

		fmt.Fprint(w, ret)
	})

	limit := 1
	after := "cursor-value"
	clickedLinks, err := client.Broadcasts.ClickedLinksWithOptions(context.Background(), "559ac32e-9ef5-46fb-82a1-b76b840c0f7b", &ListOptions{
		Limit: &limit,
		After: &after,
	})
	if err != nil {
		t.Errorf("Broadcasts.ClickedLinksWithOptions returned error: %v", err)
	}

	assert.Equal(t, clickedLinks.HasMore, true)
	assert.Equal(t, len(clickedLinks.Data), 1)
}

func TestBroadcastsClickedLinksValidations(t *testing.T) {
	setup()
	defer teardown()

	_, err := client.Broadcasts.ClickedLinks("")
	assert.NotNil(t, err)
	if err != nil {
		assert.Equal(t, err.Error(), "[ERROR]: broadcastId cannot be empty")
	}
}
