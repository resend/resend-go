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
