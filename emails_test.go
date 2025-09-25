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
			"from":"from@example.com",
			"to":["james@bond.com"],
			"created_at":"2023-04-03T22:13:42.674981+00:00",
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
					To:        []string{"recipient@example.com"},
					From:      "sender@example.com",
					CreatedAt: "2024-01-01T00:00:00Z",
					Subject:   "Test Email 1",
					Html:      "<p>Test content</p>",
					Text:      "Test content",
					LastEvent: "delivered",
				},
				{
					Id:        "2",
					Object:    "email",
					To:        []string{"recipient2@example.com"},
					From:      "sender@example.com",
					CreatedAt: "2024-01-02T00:00:00Z",
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
					CreatedAt: "2024-01-03T00:00:00Z",
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
