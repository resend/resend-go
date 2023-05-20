package resend

import (
	"encoding/json"
	"fmt"
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

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}
