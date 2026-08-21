package resend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client = NewClient("")
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func TestScheduleEmail(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &SendEmailResponse{
			Id: "1923781293",
		}
		err := json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	req := &SendEmailRequest{
		To:          []string{"d@e.com"},
		ScheduledAt: "2024-09-05T11:52:01.858Z",
	}
	resp, err := client.Emails.Send(req)
	if err != nil {
		t.Errorf("Emails.Send returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "1923781293")
}

func TestSendEmail(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &SendEmailResponse{
			Id: "1923781293",
		}
		err := json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	req := &SendEmailRequest{
		To: []string{"d@e.com"},
	}
	resp, err := client.Emails.Send(req)
	if err != nil {
		t.Errorf("Emails.Send returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "1923781293")
}

func TestSendEmailWithAttachment(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		content, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
		exp := `"attachments":[{"content":[104,101,108,108,111],"filename":"hello.txt","content_type":"text/plain"}]`
		if !bytes.Contains(content, []byte(exp)) {
			t.Errorf("request body does not include attachment data")
		}
		w.WriteHeader(http.StatusOK)
		ret := &SendEmailResponse{
			Id: "1923781293",
		}
		if err := json.NewEncoder(w).Encode(&ret); err != nil {
			panic(err)
		}
	})

	req := &SendEmailRequest{
		To: []string{"d@e.com"},
		Attachments: []*Attachment{
			{
				Content:     []byte("hello"),
				Filename:    "hello.txt",
				ContentType: "text/plain",
			},
		},
	}
	resp, err := client.Emails.Send(req)
	if err != nil {
		t.Errorf("Emails.Send returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "1923781293")
}

func TestSendEmailWithInlineAttachmentUsingContentId(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		content, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
		// Check that content_id is sent when ContentId is used
		expContentId := `"content_id":"test-cid"`
		if !bytes.Contains(content, []byte(expContentId)) {
			t.Errorf("request body does not include content_id field, got: %s", string(content))
		}
		w.WriteHeader(http.StatusOK)
		ret := &SendEmailResponse{
			Id: "1923781293",
		}
		if err := json.NewEncoder(w).Encode(&ret); err != nil {
			panic(err)
		}
	})

	req := &SendEmailRequest{
		To: []string{"d@e.com"},
		Attachments: []*Attachment{
			{
				Content:     []byte("hello"),
				Filename:    "hello.txt",
				ContentType: "text/plain",
				ContentId:   "test-cid",
			},
		},
	}
	resp, err := client.Emails.Send(req)
	if err != nil {
		t.Errorf("Emails.Send returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "1923781293")
}

func TestSendEmailWithInlineAttachmentUsingInlineContentId(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		content, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
		// Check that inline_content_id is sent when InlineContentId is used
		expInlineContentId := `"inline_content_id":"legacy-cid"`
		if !bytes.Contains(content, []byte(expInlineContentId)) {
			t.Errorf("request body does not include inline_content_id field, got: %s", string(content))
		}
		w.WriteHeader(http.StatusOK)
		ret := &SendEmailResponse{
			Id: "1923781293",
		}
		if err := json.NewEncoder(w).Encode(&ret); err != nil {
			panic(err)
		}
	})

	req := &SendEmailRequest{
		To: []string{"d@e.com"},
		Attachments: []*Attachment{
			{
				Content:         []byte("hello"),
				Filename:        "hello.txt",
				ContentType:     "text/plain",
				InlineContentId: "legacy-cid",
			},
		},
	}
	resp, err := client.Emails.Send(req)
	if err != nil {
		t.Errorf("Emails.Send returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "1923781293")
}

func TestSendEmailWithBothContentIdAndInlineContentId(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		content, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
		// When both are set, both should be sent to maintain compatibility
		expContentId := `"content_id":"preferred-cid"`
		expInlineContentId := `"inline_content_id":"legacy-cid"`
		if !bytes.Contains(content, []byte(expContentId)) {
			t.Errorf("request body does not include content_id field, got: %s", string(content))
		}
		if !bytes.Contains(content, []byte(expInlineContentId)) {
			t.Errorf("request body does not include inline_content_id field, got: %s", string(content))
		}
		w.WriteHeader(http.StatusOK)
		ret := &SendEmailResponse{
			Id: "1923781293",
		}
		if err := json.NewEncoder(w).Encode(&ret); err != nil {
			panic(err)
		}
	})

	req := &SendEmailRequest{
		To: []string{"d@e.com"},
		Attachments: []*Attachment{
			{
				Content:         []byte("hello"),
				Filename:        "hello.txt",
				ContentType:     "text/plain",
				ContentId:       "preferred-cid",
				InlineContentId: "legacy-cid",
			},
		},
	}
	resp, err := client.Emails.Send(req)
	if err != nil {
		t.Errorf("Emails.Send returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "1923781293")
}

func TestGetEmail(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/49a3999c-0ce1-4ea6-ab68-afcd6dc2e794", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id":"49a3999c-0ce1-4ea6-ab68-afcd6dc2e794",
			"message_id":"<111-222-333@email.example.com>",
			"from":"from@example.com",
			"to":["james@bond.com"],
			"created_at":"2023-04-03 22:13:42.674981+00",
			"subject": "Hello World",
			"html":"html"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Emails.Get("49a3999c-0ce1-4ea6-ab68-afcd6dc2e794")
	if err != nil {
		t.Errorf("Emails.Get returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "49a3999c-0ce1-4ea6-ab68-afcd6dc2e794")
	assert.Equal(t, resp.MessageId, "<111-222-333@email.example.com>")
	assert.Equal(t, resp.From, "from@example.com")
	assert.Equal(t, resp.Html, "html")
	assert.Equal(t, resp.To[0], "james@bond.com")
	assert.Equal(t, resp.Subject, "Hello World")
}

func TestCancelScheduledEmail(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/dacf4072-4119-4d88-932f-6202748ac7c8/cancel", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "dacf4072-4119-4d88-932f-6202748ac7c8",
			"object": "email"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Emails.Cancel("dacf4072-4119-4d88-932f-6202748ac7c8")
	if err != nil {
		t.Errorf("Emails.Cancel returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "dacf4072-4119-4d88-932f-6202748ac7c8")
	assert.Equal(t, resp.Object, "email")
}

func TestShareEmailDefaultExpiresIn(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/dacf4072-4119-4d88-932f-6202748ac7c8/share", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		content, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("failed to read request body: %v", err)
		}
		if len(content) != 0 {
			t.Errorf("expected an empty request body when params is nil, got: %s", string(content))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "dacf4072-4119-4d88-932f-6202748ac7c8",
			"object": "email",
			"url": "https://resend.com/share/dacf4072-4119-4d88-932f-6202748ac7c8"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Emails.Share("dacf4072-4119-4d88-932f-6202748ac7c8", nil)
	if err != nil {
		t.Errorf("Emails.Share returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "dacf4072-4119-4d88-932f-6202748ac7c8")
	assert.Equal(t, resp.Object, "email")
	assert.Equal(t, resp.Url, "https://resend.com/share/dacf4072-4119-4d88-932f-6202748ac7c8")
}

func TestShareEmailWithCustomExpiresIn(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/dacf4072-4119-4d88-932f-6202748ac7c8/share", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		var req ShareEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		assert.Equal(t, "1h 30m", req.ExpiresIn)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "dacf4072-4119-4d88-932f-6202748ac7c8",
			"object": "email",
			"url": "https://resend.com/share/dacf4072-4119-4d88-932f-6202748ac7c8"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Emails.Share("dacf4072-4119-4d88-932f-6202748ac7c8", &ShareEmailRequest{
		ExpiresIn: "1h 30m",
	})
	if err != nil {
		t.Errorf("Emails.Share returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "dacf4072-4119-4d88-932f-6202748ac7c8")
	assert.Equal(t, resp.Object, "email")
	assert.Equal(t, resp.Url, "https://resend.com/share/dacf4072-4119-4d88-932f-6202748ac7c8")
}

func TestShareEmailWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/dacf4072-4119-4d88-932f-6202748ac7c8/share", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "dacf4072-4119-4d88-932f-6202748ac7c8",
			"object": "email",
			"url": "https://resend.com/share/dacf4072-4119-4d88-932f-6202748ac7c8"
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	resp, err := client.Emails.ShareWithContext(ctx, "dacf4072-4119-4d88-932f-6202748ac7c8", nil)
	if err != nil {
		t.Errorf("Emails.ShareWithContext returned error: %v", err)
	}
	assert.Equal(t, resp.Id, "dacf4072-4119-4d88-932f-6202748ac7c8")
}

func TestShareEmailMalformedExpiresInReturnsError(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/dacf4072-4119-4d88-932f-6202748ac7c8/share", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)

		ret := `
		{
			"statusCode": 422,
			"name": "validation_error",
			"message": "expires_in must not exceed 48 hours"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Emails.Share("dacf4072-4119-4d88-932f-6202748ac7c8", &ShareEmailRequest{
		ExpiresIn: "72h",
	})
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expires_in must not exceed 48 hours")
}

func TestShareEmailNotFound(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/00000000-0000-0000-0000-000000000000/share", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		ret := `
		{
			"statusCode": 404,
			"name": "not_found",
			"message": "Email not found"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Emails.Share("00000000-0000-0000-0000-000000000000", nil)
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Email not found")
}

func TestSendEmailWithOptions(t *testing.T) {
	ctx := context.TODO()
	client := NewClient("123")
	params := &SendEmailRequest{
		To: []string{"email@example.com", "email2@example.com"},
	}
	options := &SendEmailOptions{
		IdempotencyKey: "unique-idempotency-key",
	}

	req, err := client.NewRequestWithOptions(ctx, "POST", "/emails/", params, options)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, req.Header["Accept"][0], "application/json")
	assert.Equal(t, req.Header["Content-Type"][0], "application/json")
	assert.Equal(t, req.Method, http.MethodPost)
	assert.Equal(t, req.URL.String(), "https://api.resend.com/emails/")
	assert.Equal(t, req.Header["Authorization"][0], "Bearer 123")
	assert.Equal(t, req.Header["Idempotency-Key"][0], "unique-idempotency-key")
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

func TestListEmails(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &ListEmailsResponse{
			Object:  "list",
			HasMore: true,
			Data: []Email{
				{
					Id:        "1",
					Object:    "email",
					MessageId: "<111-222-333@email.example.com>",
					To:        []string{"recipient@example.com"},
					From:      "sender@example.com",
					CreatedAt: "2024-01-01 00:00:00+00",
					Subject:   "Test Email 1",
					Html:      "<p>Test content</p>",
					Text:      "Test content",
					LastEvent: "delivered",
				},
				{
					Id:        "2",
					Object:    "email",
					MessageId: "<111-222-333@email.example.com>",
					To:        []string{"recipient2@example.com"},
					From:      "sender@example.com",
					CreatedAt: "2024-01-02 00:00:00+00",
					Subject:   "Test Email 2",
					Html:      "<p>Test content 2</p>",
					Text:      "Test content 2",
					LastEvent: "delivered",
				},
			},
		}
		err := json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	resp, err := client.Emails.List()
	if err != nil {
		t.Errorf("Emails.List returned error: %v", err)
	}

	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, true, resp.HasMore)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "1", resp.Data[0].Id)
	assert.Equal(t, "<111-222-333@email.example.com>", resp.Data[0].MessageId)
	assert.Equal(t, "Test Email 1", resp.Data[0].Subject)
	assert.Equal(t, "2", resp.Data[1].Id)
	assert.Equal(t, "Test Email 2", resp.Data[1].Subject)
}

func TestListEmailsWithParameters(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "50", query.Get("limit"))
		assert.Equal(t, "cursor123", query.Get("after"))
		assert.Equal(t, "cursor456", query.Get("before"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &ListEmailsResponse{
			Object:  "list",
			HasMore: false,
			Data: []Email{
				{
					Id:        "3",
					Object:    "email",
					To:        []string{"recipient3@example.com"},
					From:      "sender@example.com",
					CreatedAt: "2024-01-03 00:00:00+00",
					Subject:   "Test Email 3",
					Html:      "<p>Test content 3</p>",
					Text:      "Test content 3",
					LastEvent: "delivered",
				},
			},
		}
		err := json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	limit := 50
	after := "cursor123"
	before := "cursor456"
	options := &ListOptions{
		Limit:  &limit,
		After:  &after,
		Before: &before,
	}
	resp, err := client.Emails.ListWithOptions(context.Background(), options)
	if err != nil {
		t.Errorf("Emails.List returned error: %v", err)
	}

	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, false, resp.HasMore)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "3", resp.Data[0].Id)
	assert.Equal(t, "Test Email 3", resp.Data[0].Subject)
}

func TestListEmailsWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := &ListEmailsResponse{
			Object:  "list",
			HasMore: false,
			Data:    []Email{},
		}
		err := json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	ctx := context.Background()
	resp, err := client.Emails.ListWithContext(ctx)
	if err != nil {
		t.Errorf("Emails.ListWithContext returned error: %v", err)
	}

	assert.Equal(t, "list", resp.Object)
	assert.Equal(t, false, resp.HasMore)
	assert.Equal(t, 0, len(resp.Data))
}

func TestSendEmailWithTemplate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify template fields
		var req SendEmailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Verify template fields
		assert.NotNil(t, req.Template)
		assert.Equal(t, "welcome-template", req.Template.Id)
		assert.Nil(t, req.Template.Variables)

		ret := &SendEmailResponse{
			Id: "template-email-123",
		}
		err = json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	req := &SendEmailRequest{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Welcome!",
		Template: &EmailTemplate{
			Id: "welcome-template",
		},
	}
	resp, err := client.Emails.Send(req)
	if err != nil {
		t.Errorf("Emails.Send returned error: %v", err)
	}
	assert.Equal(t, "template-email-123", resp.Id)
}

func TestSendEmailWithTemplateAndVariables(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify template fields
		var req SendEmailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Verify template fields
		assert.NotNil(t, req.Template)
		assert.Equal(t, "user-welcome", req.Template.Id)
		assert.NotNil(t, req.Template.Variables)
		assert.Equal(t, "John Doe", req.Template.Variables["name"])
		assert.Equal(t, float64(25), req.Template.Variables["age"])

		ret := &SendEmailResponse{
			Id: "template-email-456",
		}
		err = json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	req := &SendEmailRequest{
		From:    "noreply@example.com",
		To:      []string{"john@example.com"},
		Subject: "Welcome to our service",
		Template: &EmailTemplate{
			Id: "user-welcome",
			Variables: map[string]any{
				"name": "John Doe",
				"age":  25,
			},
		},
	}
	resp, err := client.Emails.Send(req)
	if err != nil {
		t.Errorf("Emails.Send returned error: %v", err)
	}
	assert.Equal(t, "template-email-456", resp.Id)
}

func TestSendEmailWithTemplateByAlias(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify template alias
		var req SendEmailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Verify template uses alias
		assert.NotNil(t, req.Template)
		assert.Equal(t, "welcome-v2", req.Template.Id)

		ret := &SendEmailResponse{
			Id: "template-alias-789",
		}
		err = json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	req := &SendEmailRequest{
		From:    "team@example.com",
		To:      []string{"user@example.com"},
		Subject: "Hello!",
		Template: &EmailTemplate{
			Id: "welcome-v2",
		},
	}
	resp, err := client.Emails.Send(req)
	if err != nil {
		t.Errorf("Emails.Send returned error: %v", err)
	}
	assert.Equal(t, "template-alias-789", resp.Id)
}

func TestSendEmailWithTemplateAndOverrides(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify overrides
		var req SendEmailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		// Verify template and overrides
		assert.NotNil(t, req.Template)
		assert.Equal(t, "newsletter-template", req.Template.Id)
		assert.Equal(t, "Custom Newsletter", req.Template.Variables["title"])

		// Verify overridden fields
		assert.Equal(t, "custom-sender@example.com", req.From)
		assert.Equal(t, "Custom Subject Line", req.Subject)
		assert.Equal(t, "reply@example.com", req.ReplyTo)
		assert.Equal(t, []string{"bcc@example.com"}, req.Bcc)
		assert.Equal(t, 1, len(req.Tags))
		assert.Equal(t, "campaign", req.Tags[0].Name)
		assert.Equal(t, "2024-q1", req.Tags[0].Value)

		ret := &SendEmailResponse{
			Id: "template-override-999",
		}
		err = json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	req := &SendEmailRequest{
		From:    "custom-sender@example.com",
		To:      []string{"subscriber@example.com"},
		Subject: "Custom Subject Line",
		Bcc:     []string{"bcc@example.com"},
		ReplyTo: "reply@example.com",
		Tags: []Tag{
			{Name: "campaign", Value: "2024-q1"},
		},
		Template: &EmailTemplate{
			Id: "newsletter-template",
			Variables: map[string]any{
				"title": "Custom Newsletter",
			},
		},
	}
	resp, err := client.Emails.Send(req)
	if err != nil {
		t.Errorf("Emails.Send returned error: %v", err)
	}
	assert.Equal(t, "template-override-999", resp.Id)
}

func TestSendEmailWithTemplateAndContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify template
		var req SendEmailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		assert.NotNil(t, req.Template)
		assert.Equal(t, "context-template", req.Template.Id)

		ret := &SendEmailResponse{
			Id: "context-template-email-111",
		}
		err = json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	ctx := context.Background()
	req := &SendEmailRequest{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Context Test",
		Template: &EmailTemplate{
			Id: "context-template",
			Variables: map[string]any{
				"contextVar": "test",
			},
		},
	}
	resp, err := client.Emails.SendWithContext(ctx, req)
	if err != nil {
		t.Errorf("Emails.SendWithContext returned error: %v", err)
	}
	assert.Equal(t, "context-template-email-111", resp.Id)
}

func TestSendEmailWithTemplateComplexVariables(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify complex variables
		var req SendEmailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		assert.NotNil(t, req.Template)
		assert.Equal(t, "complex-template", req.Template.Id)
		assert.NotNil(t, req.Template.Variables)
		assert.Equal(t, float64(42), req.Template.Variables["count"])

		ret := &SendEmailResponse{
			Id: "complex-vars-email-222",
		}
		err = json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	req := &SendEmailRequest{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Complex Variables Test",
		Template: &EmailTemplate{
			Id: "complex-template",
			Variables: map[string]any{
				"count": 42,
			},
		},
	}
	resp, err := client.Emails.Send(req)
	if err != nil {
		t.Errorf("Emails.Send returned error: %v", err)
	}
	assert.Equal(t, "complex-vars-email-222", resp.Id)
}

func TestSendEmailWithTopicId(t *testing.T) {
	setup()
	defer teardown()

	topicId := "1234567890abcdef"

	mux.HandleFunc("/emails", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Verify that topic_id is in the request body
		bodyBytes, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Contains(t, string(bodyBytes), `"topic_id"`+":"+`"`+topicId+`"`)

		ret := &SendEmailResponse{
			Id: "topic-email-001",
		}
		err = json.NewEncoder(w).Encode(&ret)
		if err != nil {
			panic(err)
		}
	})

	req := &SendEmailRequest{
		From:    "sender@example.com",
		To:      []string{"recipient@example.com"},
		Subject: "Email with Topic",
		Html:    "<p>This email is sent to a topic subscriber</p>",
		TopicId: topicId,
	}
	resp, err := client.Emails.Send(req)
	if err != nil {
		t.Errorf("Emails.Send returned error: %v", err)
	}
	assert.Equal(t, "topic-email-001", resp.Id)
}

func TestEmailsMetrics(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, "", r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "metrics",
			"start_date": "2026-07-01T00:00:00.000Z",
			"end_date": "2026-07-08T00:00:00.000Z",
			"metrics": ["delivered", "opened"],
			"dimensions": [],
			"granularity": "daily",
			"totals": {"delivered": 100, "opened": 40}
		}`)
	})

	resp, err := client.Emails.Metrics()
	if err != nil {
		t.Fatalf("Emails.Metrics returned error: %v", err)
	}
	assert.Equal(t, "metrics", resp.Object)
	assert.Equal(t, "2026-07-01T00:00:00.000Z", resp.StartDate)
	assert.Equal(t, "2026-07-08T00:00:00.000Z", resp.EndDate)
	assert.Equal(t, []string{"delivered", "opened"}, resp.Metrics)
	assert.Equal(t, "daily", resp.Granularity)
	assert.Equal(t, float64(100), resp.Totals["delivered"])
	assert.Equal(t, float64(40), resp.Totals["opened"])
	assert.Nil(t, resp.Data)
}

func TestEmailsMetricsWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "metrics",
			"start_date": "2026-07-01T00:00:00.000Z",
			"end_date": "2026-07-08T00:00:00.000Z",
			"metrics": ["delivered"],
			"dimensions": [],
			"granularity": "daily",
			"totals": {"delivered": 5}
		}`)
	})

	resp, err := client.Emails.MetricsWithContext(context.Background())
	if err != nil {
		t.Fatalf("Emails.MetricsWithContext returned error: %v", err)
	}
	assert.Equal(t, float64(5), resp.Totals["delivered"])
}

func TestEmailsMetricsWithPeriodDimension(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, "period", r.URL.Query().Get("dimensions"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "metrics",
			"start_date": "2026-07-01T00:00:00.000Z",
			"end_date": "2026-07-08T00:00:00.000Z",
			"metrics": ["delivered"],
			"dimensions": ["period"],
			"granularity": "daily",
			"totals": {"delivered": 10},
			"data": [
				{"period": "2026-07-01", "delivered": 10}
			]
		}`)
	})

	resp, err := client.Emails.MetricsWithOptions(context.Background(), &MetricsOptions{
		Dimensions: []MetricsDimension{MetricsDimensionPeriod},
	})
	if err != nil {
		t.Fatalf("Emails.MetricsWithOptions returned error: %v", err)
	}
	assert.Equal(t, 1, len(resp.Data))
	assert.NotNil(t, resp.Data[0].Period)
	assert.Equal(t, "2026-07-01", *resp.Data[0].Period)
	assert.NotNil(t, resp.Data[0].Delivered)
	assert.Equal(t, int64(10), *resp.Data[0].Delivered)
}

func TestEmailsMetricsWithDomainDimension(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, "domain", r.URL.Query().Get("dimensions"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "metrics",
			"start_date": "2026-07-01T00:00:00.000Z",
			"end_date": "2026-07-08T00:00:00.000Z",
			"metrics": ["delivered"],
			"dimensions": ["domain"],
			"granularity": "daily",
			"totals": {"delivered": 10},
			"data": [
				{"domain_id": "d1a2b3c4-0000-4000-8000-000000000001", "domain_name": "example.com", "delivered": 10}
			]
		}`)
	})

	resp, err := client.Emails.MetricsWithOptions(context.Background(), &MetricsOptions{
		Dimensions: []MetricsDimension{MetricsDimensionDomain},
	})
	if err != nil {
		t.Fatalf("Emails.MetricsWithOptions returned error: %v", err)
	}
	assert.Equal(t, 1, len(resp.Data))
	assert.NotNil(t, resp.Data[0].DomainId)
	assert.Equal(t, "d1a2b3c4-0000-4000-8000-000000000001", *resp.Data[0].DomainId)
	assert.NotNil(t, resp.Data[0].DomainName)
	assert.Equal(t, "example.com", *resp.Data[0].DomainName)
}

func TestEmailsMetricsWithEmailDimension(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, "email", r.URL.Query().Get("dimensions"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "metrics",
			"start_date": "2026-07-01T00:00:00.000Z",
			"end_date": "2026-07-08T00:00:00.000Z",
			"metrics": ["delivered"],
			"dimensions": ["email"],
			"granularity": "daily",
			"totals": {"delivered": 1},
			"data": [
				{"email_id": "e1a2b3c4-0000-4000-8000-000000000002", "delivered": 1}
			]
		}`)
	})

	resp, err := client.Emails.MetricsWithOptions(context.Background(), &MetricsOptions{
		Dimensions: []MetricsDimension{MetricsDimensionEmail},
	})
	if err != nil {
		t.Fatalf("Emails.MetricsWithOptions returned error: %v", err)
	}
	assert.Equal(t, 1, len(resp.Data))
	assert.NotNil(t, resp.Data[0].EmailId)
	assert.Equal(t, "e1a2b3c4-0000-4000-8000-000000000002", *resp.Data[0].EmailId)
}

func TestEmailsMetricsWithBroadcastDimension(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, "broadcast", r.URL.Query().Get("dimensions"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "metrics",
			"start_date": "2026-07-01T00:00:00.000Z",
			"end_date": "2026-07-08T00:00:00.000Z",
			"metrics": ["delivered", "opened"],
			"dimensions": ["broadcast"],
			"granularity": "daily",
			"totals": {"delivered": 10, "opened": 4},
			"data": [
				{"broadcast_id": "b1a2b3c4-0000-4000-8000-000000000003", "broadcast_name": "July Newsletter", "delivered": 10, "opened": 4}
			]
		}`)
	})

	resp, err := client.Emails.MetricsWithOptions(context.Background(), &MetricsOptions{
		Dimensions: []MetricsDimension{MetricsDimensionBroadcast},
	})
	if err != nil {
		t.Fatalf("Emails.MetricsWithOptions returned error: %v", err)
	}
	assert.Equal(t, 1, len(resp.Data))
	assert.NotNil(t, resp.Data[0].BroadcastId)
	assert.Equal(t, "b1a2b3c4-0000-4000-8000-000000000003", *resp.Data[0].BroadcastId)
	assert.NotNil(t, resp.Data[0].BroadcastName)
	assert.Equal(t, "July Newsletter", *resp.Data[0].BroadcastName)
	assert.Equal(t, int64(4), *resp.Data[0].Opened)
}

func TestEmailsMetricsWithSingleDomainIdFilter(t *testing.T) {
	setup()
	defer teardown()

	domainId := "d1a2b3c4-0000-4000-8000-000000000001"

	mux.HandleFunc("/emails/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, domainId, r.URL.Query().Get("domain_id"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "metrics",
			"start_date": "2026-07-01T00:00:00.000Z",
			"end_date": "2026-07-08T00:00:00.000Z",
			"metrics": ["delivered"],
			"dimensions": [],
			"granularity": "daily",
			"totals": {"delivered": 3}
		}`)
	})

	_, err := client.Emails.MetricsWithOptions(context.Background(), &MetricsOptions{
		DomainId: []string{domainId},
	})
	if err != nil {
		t.Fatalf("Emails.MetricsWithOptions returned error: %v", err)
	}
}

func TestEmailsMetricsWithMultipleDomainIdFilter(t *testing.T) {
	setup()
	defer teardown()

	domainIds := []string{
		"d1a2b3c4-0000-4000-8000-000000000001",
		"d1a2b3c4-0000-4000-8000-000000000002",
	}

	mux.HandleFunc("/emails/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, strings.Join(domainIds, ","), r.URL.Query().Get("domain_id"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "metrics",
			"start_date": "2026-07-01T00:00:00.000Z",
			"end_date": "2026-07-08T00:00:00.000Z",
			"metrics": ["delivered"],
			"dimensions": [],
			"granularity": "daily",
			"totals": {"delivered": 6}
		}`)
	})

	_, err := client.Emails.MetricsWithOptions(context.Background(), &MetricsOptions{
		DomainId: domainIds,
	})
	if err != nil {
		t.Fatalf("Emails.MetricsWithOptions returned error: %v", err)
	}
}

func TestEmailsMetricsWithSingleEmailIdFilter(t *testing.T) {
	setup()
	defer teardown()

	emailId := "e1a2b3c4-0000-4000-8000-000000000001"

	mux.HandleFunc("/emails/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, emailId, r.URL.Query().Get("email_id"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "metrics",
			"start_date": "2026-07-01T00:00:00.000Z",
			"end_date": "2026-07-08T00:00:00.000Z",
			"metrics": ["delivered"],
			"dimensions": [],
			"granularity": "daily",
			"totals": {"delivered": 1}
		}`)
	})

	_, err := client.Emails.MetricsWithOptions(context.Background(), &MetricsOptions{
		EmailId: []string{emailId},
	})
	if err != nil {
		t.Fatalf("Emails.MetricsWithOptions returned error: %v", err)
	}
}

func TestEmailsMetricsWithMultipleEmailIdFilter(t *testing.T) {
	setup()
	defer teardown()

	emailIds := []string{
		"e1a2b3c4-0000-4000-8000-000000000001",
		"e1a2b3c4-0000-4000-8000-000000000002",
	}

	mux.HandleFunc("/emails/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, strings.Join(emailIds, ","), r.URL.Query().Get("email_id"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "metrics",
			"start_date": "2026-07-01T00:00:00.000Z",
			"end_date": "2026-07-08T00:00:00.000Z",
			"metrics": ["delivered"],
			"dimensions": [],
			"granularity": "daily",
			"totals": {"delivered": 2}
		}`)
	})

	_, err := client.Emails.MetricsWithOptions(context.Background(), &MetricsOptions{
		EmailId: emailIds,
	})
	if err != nil {
		t.Fatalf("Emails.MetricsWithOptions returned error: %v", err)
	}
}

func TestEmailsMetricsWithSingleBroadcastIdFilter(t *testing.T) {
	setup()
	defer teardown()

	broadcastId := "b1a2b3c4-0000-4000-8000-000000000001"

	mux.HandleFunc("/emails/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, broadcastId, r.URL.Query().Get("broadcast_id"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "metrics",
			"start_date": "2026-07-01T00:00:00.000Z",
			"end_date": "2026-07-08T00:00:00.000Z",
			"metrics": ["delivered"],
			"dimensions": [],
			"granularity": "daily",
			"totals": {"delivered": 12}
		}`)
	})

	_, err := client.Emails.MetricsWithOptions(context.Background(), &MetricsOptions{
		BroadcastId: []string{broadcastId},
	})
	if err != nil {
		t.Fatalf("Emails.MetricsWithOptions returned error: %v", err)
	}
}

func TestEmailsMetricsWithMultipleBroadcastIdFilter(t *testing.T) {
	setup()
	defer teardown()

	broadcastIds := []string{
		"b1a2b3c4-0000-4000-8000-000000000001",
		"b1a2b3c4-0000-4000-8000-000000000002",
	}

	mux.HandleFunc("/emails/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		assert.Equal(t, strings.Join(broadcastIds, ","), r.URL.Query().Get("broadcast_id"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "metrics",
			"start_date": "2026-07-01T00:00:00.000Z",
			"end_date": "2026-07-08T00:00:00.000Z",
			"metrics": ["delivered"],
			"dimensions": [],
			"granularity": "daily",
			"totals": {"delivered": 20}
		}`)
	})

	_, err := client.Emails.MetricsWithOptions(context.Background(), &MetricsOptions{
		BroadcastId: broadcastIds,
	})
	if err != nil {
		t.Fatalf("Emails.MetricsWithOptions returned error: %v", err)
	}
}

func TestEmailsMetricsRejectsEmailAndBroadcastCombinations(t *testing.T) {
	setup()
	defer teardown()

	const wantErr = "[ERROR]: The `email` dimension/EmailId cannot be combined with the `broadcast` dimension/BroadcastId."

	_, err := client.Emails.MetricsWithOptions(context.Background(), &MetricsOptions{
		Dimensions: []MetricsDimension{MetricsDimensionEmail, MetricsDimensionBroadcast},
	})
	assert.Error(t, err)
	assert.Equal(t, wantErr, err.Error())

	_, err = client.Emails.MetricsWithOptions(context.Background(), &MetricsOptions{
		Dimensions: []MetricsDimension{MetricsDimensionBroadcast},
		EmailId:    []string{"4dd369bc-aa82-4ff3-97de-514ae3000ee0"},
	})
	assert.Error(t, err)
	assert.Equal(t, wantErr, err.Error())

	_, err = client.Emails.MetricsWithOptions(context.Background(), &MetricsOptions{
		Dimensions:  []MetricsDimension{MetricsDimensionEmail},
		BroadcastId: []string{"5a5a3b1e-3b1a-4b1a-8b1a-3b1a4b1a8b1a"},
	})
	assert.Error(t, err)
	assert.Equal(t, wantErr, err.Error())

	_, err = client.Emails.MetricsWithOptions(context.Background(), &MetricsOptions{
		EmailId:     []string{"4dd369bc-aa82-4ff3-97de-514ae3000ee0"},
		BroadcastId: []string{"5a5a3b1e-3b1a-4b1a-8b1a-3b1a4b1a8b1a"},
	})
	assert.Error(t, err)
	assert.Equal(t, wantErr, err.Error())
}

func TestEmailsMetricsWithMetricsGranularityAndTimezone(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/metrics", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		query := r.URL.Query()
		assert.Equal(t, "delivered,opened,clicked", query.Get("metrics"))
		assert.Equal(t, "hourly", query.Get("granularity"))
		assert.Equal(t, "America/New_York", query.Get("timezone"))
		assert.Equal(t, "2026-07-01", query.Get("start_date"))
		assert.Equal(t, "2026-07-08", query.Get("end_date"))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
			"object": "metrics",
			"start_date": "2026-07-01T00:00:00.000Z",
			"end_date": "2026-07-08T00:00:00.000Z",
			"metrics": ["delivered", "opened", "clicked"],
			"dimensions": [],
			"granularity": "hourly",
			"totals": {"delivered": 10, "opened": 4, "clicked": 2}
		}`)
	})

	startDate := "2026-07-01"
	endDate := "2026-07-08"
	timezone := "America/New_York"
	granularity := MetricsGranularityHourly

	resp, err := client.Emails.MetricsWithOptions(context.Background(), &MetricsOptions{
		StartDate:   &startDate,
		EndDate:     &endDate,
		Timezone:    &timezone,
		Granularity: &granularity,
		Metrics:     []MetricName{MetricDelivered, MetricOpened, MetricClicked},
	})
	if err != nil {
		t.Fatalf("Emails.MetricsWithOptions returned error: %v", err)
	}
	assert.Equal(t, "hourly", resp.Granularity)
	assert.Equal(t, []string{"delivered", "opened", "clicked"}, resp.Metrics)
}
