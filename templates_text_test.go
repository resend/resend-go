package resend

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type templateTextTransport func(*http.Request) (*http.Response, error)

func (f templateTextTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestTemplateRequestMarshalJSON(t *testing.T) {
	create := CreateTemplateRequest{Name: "welcome", Html: "<p>Welcome</p>"}
	create.SetText("")
	update := UpdateTemplateRequest{Subject: "New subject"}
	update.SetText("")
	tests := []struct {
		name    string
		request any
		want    string
	}{
		{"zero create value", CreateTemplateRequest{}, `{"name":"","html":""}`},
		{"zero update value", UpdateTemplateRequest{}, `{}`},
		{"nil create pointer", (*CreateTemplateRequest)(nil), `null`},
		{"nil update pointer", (*UpdateTemplateRequest)(nil), `null`},
		{"explicit empty create value", create, `{"name":"welcome","html":"<p>Welcome</p>","text":""}`},
		{"explicit empty update value", update, `{"subject":"New subject","text":""}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatal(err)
			}
			var gotValue, wantValue any
			if err := json.Unmarshal(got, &gotValue); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(tt.want), &wantValue); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(gotValue, wantValue) {
				t.Errorf("JSON = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestTemplateTextPresence(t *testing.T) {
	empty, nonempty := "", "Plain text"
	tests := []struct {
		name         string
		text         string
		setText      []string
		directAfter  *string
		allFields    bool
		expectedText *string
	}{
		{name: "omitted"},
		{name: "direct empty remains omitted", text: ""},
		{name: "direct nonempty", text: nonempty, expectedText: &nonempty},
		{name: "explicit empty", setText: []string{empty}, expectedText: &empty},
		{name: "explicit nonempty", setText: []string{nonempty}, expectedText: &nonempty},
		{name: "setter overwrites direct value", text: nonempty, setText: []string{empty}, expectedText: &empty},
		{name: "setter resets to empty", setText: []string{nonempty, empty}, expectedText: &empty},
		{name: "direct value after setter", setText: []string{empty}, directAfter: &nonempty, expectedText: &nonempty},
		{name: "direct empty after setter", setText: []string{nonempty}, directAfter: &empty, expectedText: &empty},
		{name: "explicit empty retains other fields", setText: []string{empty}, allFields: true, expectedText: &empty},
	}
	for _, operation := range []string{"create", "create with context", "update", "update with context"} {
		for _, tt := range tests {
			t.Run(operation+"/"+tt.name, func(t *testing.T) {
				create := &CreateTemplateRequest{Name: "welcome", Html: "<p>Welcome</p>", Text: tt.text}
				update := &UpdateTemplateRequest{Subject: "New subject", Text: tt.text}
				for _, text := range tt.setText {
					create.SetText(text)
					update.SetText(text)
				}
				if tt.directAfter != nil {
					create.Text = *tt.directAfter
					update.Text = *tt.directAfter
				}
				expected := map[string]any{"subject": "New subject"}
				method, path := http.MethodPatch, "/templates/owned-template"
				if strings.HasPrefix(operation, "create") {
					expected = map[string]any{"name": "welcome", "html": "<p>Welcome</p>"}
					method, path = http.MethodPost, "/templates"
				}
				if tt.allFields {
					variables := []*TemplateVariable{{Key: "NAME", Type: VariableTypeString, FallbackValue: "friend"}}
					create.Alias, update.Alias = "welcome-alias", "welcome-alias"
					create.From, update.From = "team@example.com", "team@example.com"
					create.Subject, update.Subject = "New subject", "New subject"
					create.ReplyTo, update.ReplyTo = []string{"help@example.com"}, []string{"help@example.com"}
					update.Name, update.Html = create.Name, create.Html
					create.Variables, update.Variables = variables, variables
					expected = map[string]any{
						"name": "welcome", "html": "<p>Welcome</p>", "alias": "welcome-alias",
						"from": "team@example.com", "subject": "New subject", "reply_to": []any{"help@example.com"},
						"variables": []any{map[string]any{"key": "NAME", "type": "string", "fallback_value": "friend"}},
					}
				}
				if tt.expectedText != nil {
					expected["text"] = *tt.expectedText
				}
				calls := 0
				client := NewCustomClient(&http.Client{Transport: templateTextTransport(func(req *http.Request) (*http.Response, error) {
					calls++
					if req.Method != method || req.URL.Path != path {
						t.Errorf("request = %s %s, want %s %s", req.Method, req.URL.Path, method, path)
					}
					var body map[string]any
					if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
						return nil, err
					}
					if !reflect.DeepEqual(expected, body) {
						t.Errorf("body = %#v, want %#v", body, expected)
					}
					return &http.Response{
						StatusCode: http.StatusOK, Header: make(http.Header), Request: req,
						Body: io.NopCloser(strings.NewReader(`{"object":"template","id":"owned-template"}`)),
					}, nil
				})}, "owned-test-key")
				var err error
				switch operation {
				case "create":
					_, err = client.Templates.Create(create)
				case "create with context":
					_, err = client.Templates.CreateWithContext(context.Background(), create)
				case "update":
					_, err = client.Templates.Update("owned-template", update)
				case "update with context":
					_, err = client.Templates.UpdateWithContext(context.Background(), "owned-template", update)
				}
				if err != nil {
					t.Fatal(err)
				}
				if calls != 1 {
					t.Errorf("transport calls = %d, want 1", calls)
				}
			})
		}
	}
}
