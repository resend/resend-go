package resend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTemplate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify it
		var req CreateTemplateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		assert.Equal(t, "welcome-email", req.Name)
		assert.Equal(t, "<strong>Welcome!</strong>", req.Html)

		ret := `
		{
			"id": "49a3999c-0ce1-4ea6-ab68-afcd6dc2e794",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Create(&CreateTemplateRequest{
		Name: "welcome-email",
		Html: "<strong>Welcome!</strong>",
	})
	if err != nil {
		t.Errorf("Templates.Create returned error: %v", err)
	}
	assert.Equal(t, "49a3999c-0ce1-4ea6-ab68-afcd6dc2e794", resp.Id)
	assert.Equal(t, "template", resp.Object)
}

func TestCreateTemplateWithVariables(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify variables
		var req CreateTemplateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		assert.Equal(t, "user-notification", req.Name)
		assert.Equal(t, "<strong>Hey, {{{NAME}}}, you are {{{AGE}}} years old.</strong>", req.Html)
		assert.Equal(t, 3, len(req.Variables))
		assert.Equal(t, "NAME", req.Variables[0].Key)
		assert.Equal(t, VariableTypeString, req.Variables[0].Type)
		assert.Equal(t, "user", req.Variables[0].FallbackValue)
		assert.Equal(t, "AGE", req.Variables[1].Key)
		assert.Equal(t, VariableTypeNumber, req.Variables[1].Type)
		// JSON numbers are decoded as float64
		assert.Equal(t, float64(25), req.Variables[1].FallbackValue)
		assert.Equal(t, "OPTIONAL_VARIABLE", req.Variables[2].Key)
		assert.Equal(t, VariableTypeString, req.Variables[2].Type)

		ret := `
		{
			"id": "template-with-vars-id",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Create(&CreateTemplateRequest{
		Name: "user-notification",
		Html: "<strong>Hey, {{{NAME}}}, you are {{{AGE}}} years old.</strong>",
		Variables: []*TemplateVariable{
			{
				Key:           "NAME",
				Type:          VariableTypeString,
				FallbackValue: "user",
			},
			{
				Key:           "AGE",
				Type:          VariableTypeNumber,
				FallbackValue: 25,
			},
			{
				Key:  "OPTIONAL_VARIABLE",
				Type: VariableTypeString,
			},
		},
	})
	if err != nil {
		t.Errorf("Templates.Create returned error: %v", err)
	}
	assert.Equal(t, "template-with-vars-id", resp.Id)
	assert.Equal(t, "template", resp.Object)
}

func TestCreateTemplateWithAllFields(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify all fields
		var req CreateTemplateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		assert.Equal(t, "full-template", req.Name)
		assert.Equal(t, "full-alias", req.Alias)
		assert.Equal(t, "Team <team@example.com>", req.From)
		assert.Equal(t, "Important Update", req.Subject)
		assert.Equal(t, "<h1>Hello</h1>", req.Html)
		assert.Equal(t, "Hello", req.Text)

		// ReplyTo can be string or []string
		replyTo, ok := req.ReplyTo.([]interface{})
		assert.True(t, ok)
		assert.Equal(t, 2, len(replyTo))
		assert.Equal(t, "support@example.com", replyTo[0].(string))
		assert.Equal(t, "help@example.com", replyTo[1].(string))

		ret := `
		{
			"id": "full-template-id",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Create(&CreateTemplateRequest{
		Name:    "full-template",
		Alias:   "full-alias",
		From:    "Team <team@example.com>",
		Subject: "Important Update",
		ReplyTo: []string{"support@example.com", "help@example.com"},
		Html:    "<h1>Hello</h1>",
		Text:    "Hello",
	})
	if err != nil {
		t.Errorf("Templates.Create returned error: %v", err)
	}
	assert.Equal(t, "full-template-id", resp.Id)
	assert.Equal(t, "template", resp.Object)
}

func TestCreateTemplateWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "context-template-id",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	resp, err := client.Templates.CreateWithContext(ctx, &CreateTemplateRequest{
		Name: "context-template",
		Html: "<p>Content</p>",
	})
	if err != nil {
		t.Errorf("Templates.CreateWithContext returned error: %v", err)
	}
	assert.Equal(t, "context-template-id", resp.Id)
	assert.Equal(t, "template", resp.Object)
}

func TestCreateTemplateWithAllVariableTypes(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify all variable types
		var req CreateTemplateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		assert.Equal(t, 2, len(req.Variables))

		// String variable
		assert.Equal(t, "STRING_VAR", req.Variables[0].Key)
		assert.Equal(t, VariableTypeString, req.Variables[0].Type)
		assert.Equal(t, "default", req.Variables[0].FallbackValue)

		// Number variable
		assert.Equal(t, "NUMBER_VAR", req.Variables[1].Key)
		assert.Equal(t, VariableTypeNumber, req.Variables[1].Type)
		assert.Equal(t, float64(42), req.Variables[1].FallbackValue)

		ret := `
		{
			"id": "all-types-template-id",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Create(&CreateTemplateRequest{
		Name: "all-variable-types",
		Html: "<div>Test</div>",
		Variables: []*TemplateVariable{
			{
				Key:           "STRING_VAR",
				Type:          VariableTypeString,
				FallbackValue: "default",
			},
			{
				Key:           "NUMBER_VAR",
				Type:          VariableTypeNumber,
				FallbackValue: 42,
			},
		},
	})
	if err != nil {
		t.Errorf("Templates.Create returned error: %v", err)
	}
	assert.Equal(t, "all-types-template-id", resp.Id)
	assert.Equal(t, "template", resp.Object)
}

func TestCreateTemplateWithSingleReplyTo(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify ReplyTo as string
		var req CreateTemplateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		replyTo, ok := req.ReplyTo.(string)
		assert.True(t, ok)
		assert.Equal(t, "reply@example.com", replyTo)

		ret := `
		{
			"id": "single-reply-to-id",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Create(&CreateTemplateRequest{
		Name:    "single-reply-to",
		Html:    "<p>Test</p>",
		ReplyTo: "reply@example.com",
	})
	if err != nil {
		t.Errorf("Templates.Create returned error: %v", err)
	}
	assert.Equal(t, "single-reply-to-id", resp.Id)
}

func TestGetTemplate(t *testing.T) {
	setup()
	defer teardown()

	templateId := "34a080c9-b17d-4187-ad80-5af20266e535"

	mux.HandleFunc(fmt.Sprintf("/templates/%s", templateId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "template",
			"id": "34a080c9-b17d-4187-ad80-5af20266e535",
			"alias": "reset-password",
			"name": "reset-password",
			"created_at": "2023-10-06T23:47:56.678Z",
			"updated_at": "2023-10-06T23:47:56.678Z",
			"status": "published",
			"published_at": "2023-10-06T23:47:56.678Z",
			"from": "John Doe <john.doe@example.com>",
			"subject": "Hello, world!",
			"reply_to": null,
			"html": "<h1>Hello, world!</h1>",
			"text": "Hello, world!",
			"variables": [
				{
					"id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
					"key": "user_name",
					"type": "string",
					"fallback_value": "John Doe",
					"created_at": "2023-10-06T23:47:56.678Z",
					"updated_at": "2023-10-06T23:47:56.678Z"
				}
			]
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Get(templateId)
	if err != nil {
		t.Errorf("Templates.Get returned error: %v", err)
	}
	assert.Equal(t, "template", resp.Object)
	assert.Equal(t, "34a080c9-b17d-4187-ad80-5af20266e535", resp.Id)
	assert.Equal(t, "reset-password", resp.Alias)
	assert.Equal(t, "reset-password", resp.Name)
	assert.Equal(t, "published", resp.Status)
	assert.Equal(t, "2023-10-06T23:47:56.678Z", resp.CreatedAt)
	assert.Equal(t, "2023-10-06T23:47:56.678Z", resp.UpdatedAt)
	assert.Equal(t, "2023-10-06T23:47:56.678Z", resp.PublishedAt)
	assert.Equal(t, "John Doe <john.doe@example.com>", resp.From)
	assert.Equal(t, "Hello, world!", resp.Subject)
	assert.Nil(t, resp.ReplyTo)
	assert.Equal(t, "<h1>Hello, world!</h1>", resp.Html)
	assert.Equal(t, "Hello, world!", resp.Text)
	assert.Equal(t, 1, len(resp.Variables))
	assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", resp.Variables[0].Id)
	assert.Equal(t, "user_name", resp.Variables[0].Key)
	assert.Equal(t, VariableTypeString, resp.Variables[0].Type)
	assert.Equal(t, "John Doe", resp.Variables[0].FallbackValue)
	assert.Equal(t, "2023-10-06T23:47:56.678Z", resp.Variables[0].CreatedAt)
	assert.Equal(t, "2023-10-06T23:47:56.678Z", resp.Variables[0].UpdatedAt)
}

func TestGetTemplateByAlias(t *testing.T) {
	setup()
	defer teardown()

	templateAlias := "welcome-email"

	mux.HandleFunc(fmt.Sprintf("/templates/%s", templateAlias), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "template",
			"id": "template-id-123",
			"alias": "welcome-email",
			"name": "Welcome Email",
			"created_at": "2023-10-06T23:47:56.678Z",
			"updated_at": "2023-10-06T23:47:56.678Z",
			"status": "draft",
			"published_at": "",
			"from": "support@example.com",
			"subject": "Welcome!",
			"reply_to": "noreply@example.com",
			"html": "<p>Welcome!</p>",
			"text": "Welcome!",
			"variables": []
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Get(templateAlias)
	if err != nil {
		t.Errorf("Templates.Get returned error: %v", err)
	}
	assert.Equal(t, "template", resp.Object)
	assert.Equal(t, "template-id-123", resp.Id)
	assert.Equal(t, "welcome-email", resp.Alias)
	assert.Equal(t, "Welcome Email", resp.Name)
	assert.Equal(t, "draft", resp.Status)
	assert.Equal(t, "", resp.PublishedAt)
	assert.Equal(t, "support@example.com", resp.From)
	assert.Equal(t, "Welcome!", resp.Subject)
	assert.Equal(t, "noreply@example.com", resp.ReplyTo)
	assert.Equal(t, 0, len(resp.Variables))
}

func TestGetTemplateWithContext(t *testing.T) {
	setup()
	defer teardown()

	templateId := "context-test-id"

	mux.HandleFunc(fmt.Sprintf("/templates/%s", templateId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "template",
			"id": "context-test-id",
			"alias": "",
			"name": "Context Test",
			"created_at": "2023-10-06T23:47:56.678Z",
			"updated_at": "2023-10-06T23:47:56.678Z",
			"status": "published",
			"published_at": "2023-10-06T23:47:56.678Z",
			"from": "",
			"subject": "",
			"reply_to": null,
			"html": "<p>Test</p>",
			"text": "Test",
			"variables": []
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	resp, err := client.Templates.GetWithContext(ctx, templateId)
	if err != nil {
		t.Errorf("Templates.GetWithContext returned error: %v", err)
	}
	assert.Equal(t, "template", resp.Object)
	assert.Equal(t, "context-test-id", resp.Id)
	assert.Equal(t, "Context Test", resp.Name)
	assert.Equal(t, "published", resp.Status)
}

func TestListTemplates(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "2", query.Get("limit"))

		ret := `
		{
			"object": "list",
			"data": [
				{
					"id": "e169aa45-1ecf-4183-9955-b1499d5701d3",
					"name": "reset-password",
					"status": "draft",
					"published_at": null,
					"created_at": "2023-10-06T23:47:56.678Z",
					"updated_at": "2023-10-06T23:47:56.678Z",
					"alias": "reset-password"
				},
				{
					"id": "b7f9c2e1-1234-4abc-9def-567890abcdef",
					"name": "welcome-message",
					"status": "published",
					"published_at": "2023-10-06T23:47:56.678Z",
					"created_at": "2023-10-06T23:47:56.678Z",
					"updated_at": "2023-10-06T23:47:56.678Z",
					"alias": "welcome-message"
				}
			],
			"has_more": false
		}`
		fmt.Fprintf(w, ret)
	})

	limit := 2
	resp, err := client.Templates.List(&ListOptions{
		Limit: &limit,
	})
	if err != nil {
		t.Errorf("Templates.List returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.False(t, resp.HasMore)
	assert.Equal(t, 2, len(resp.Data))
	assert.Equal(t, "e169aa45-1ecf-4183-9955-b1499d5701d3", resp.Data[0].Id)
	assert.Equal(t, "reset-password", resp.Data[0].Name)
	assert.Equal(t, "draft", resp.Data[0].Status)
	assert.Nil(t, resp.Data[0].PublishedAt)
	assert.Equal(t, "reset-password", resp.Data[0].Alias)
	assert.Equal(t, "b7f9c2e1-1234-4abc-9def-567890abcdef", resp.Data[1].Id)
	assert.Equal(t, "welcome-message", resp.Data[1].Name)
	assert.Equal(t, "published", resp.Data[1].Status)
	assert.NotNil(t, resp.Data[1].PublishedAt)
	assert.Equal(t, "2023-10-06T23:47:56.678Z", *resp.Data[1].PublishedAt)
}

func TestListTemplatesWithAfter(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "2", query.Get("limit"))
		assert.Equal(t, "34a080c9-b17d-4187-ad80-5af20266e535", query.Get("after"))

		ret := `
		{
			"object": "list",
			"data": [
				{
					"id": "next-template-id",
					"name": "next-template",
					"status": "published",
					"published_at": "2023-10-07T23:47:56.678Z",
					"created_at": "2023-10-07T23:47:56.678Z",
					"updated_at": "2023-10-07T23:47:56.678Z",
					"alias": "next-template"
				}
			],
			"has_more": true
		}`
		fmt.Fprintf(w, ret)
	})

	limit := 2
	after := "34a080c9-b17d-4187-ad80-5af20266e535"
	resp, err := client.Templates.List(&ListOptions{
		Limit: &limit,
		After: &after,
	})
	if err != nil {
		t.Errorf("Templates.List returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.True(t, resp.HasMore)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "next-template-id", resp.Data[0].Id)
}

func TestListTemplatesWithBefore(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Check query parameters
		query := r.URL.Query()
		assert.Equal(t, "2", query.Get("limit"))
		assert.Equal(t, "34a080c9-b17d-4187-ad80-5af20266e535", query.Get("before"))

		ret := `
		{
			"object": "list",
			"data": [
				{
					"id": "previous-template-id",
					"name": "previous-template",
					"status": "draft",
					"published_at": null,
					"created_at": "2023-10-05T23:47:56.678Z",
					"updated_at": "2023-10-05T23:47:56.678Z",
					"alias": "previous-template"
				}
			],
			"has_more": false
		}`
		fmt.Fprintf(w, ret)
	})

	limit := 2
	before := "34a080c9-b17d-4187-ad80-5af20266e535"
	resp, err := client.Templates.List(&ListOptions{
		Limit:  &limit,
		Before: &before,
	})
	if err != nil {
		t.Errorf("Templates.List returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.False(t, resp.HasMore)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "previous-template-id", resp.Data[0].Id)
}

func TestListTemplatesWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "list",
			"data": [
				{
					"id": "context-template-id",
					"name": "context-template",
					"status": "published",
					"published_at": "2023-10-06T23:47:56.678Z",
					"created_at": "2023-10-06T23:47:56.678Z",
					"updated_at": "2023-10-06T23:47:56.678Z",
					"alias": "context-template"
				}
			],
			"has_more": false
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	resp, err := client.Templates.ListWithContext(ctx, &ListOptions{})
	if err != nil {
		t.Errorf("Templates.ListWithContext returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.False(t, resp.HasMore)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, "context-template-id", resp.Data[0].Id)
}

func TestListTemplatesWithoutOptions(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/templates", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Check that there are no query parameters
		query := r.URL.Query()
		assert.Equal(t, "", query.Get("limit"))
		assert.Equal(t, "", query.Get("after"))
		assert.Equal(t, "", query.Get("before"))

		ret := `
		{
			"object": "list",
			"data": [
				{
					"id": "template-1",
					"name": "template-1",
					"status": "published",
					"published_at": "2023-10-06T23:47:56.678Z",
					"created_at": "2023-10-06T23:47:56.678Z",
					"updated_at": "2023-10-06T23:47:56.678Z",
					"alias": "template-1"
				},
				{
					"id": "template-2",
					"name": "template-2",
					"status": "draft",
					"published_at": null,
					"created_at": "2023-10-06T23:47:56.678Z",
					"updated_at": "2023-10-06T23:47:56.678Z",
					"alias": "template-2"
				}
			],
			"has_more": false
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.List(nil)
	if err != nil {
		t.Errorf("Templates.List returned error: %v", err)
	}
	assert.Equal(t, "list", resp.Object)
	assert.False(t, resp.HasMore)
	assert.Equal(t, 2, len(resp.Data))
}

func TestGetTemplateWithMultipleReplyTo(t *testing.T) {
	setup()
	defer teardown()

	templateId := "multi-reply-to-id"

	mux.HandleFunc(fmt.Sprintf("/templates/%s", templateId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "template",
			"id": "multi-reply-to-id",
			"alias": "",
			"name": "Multi Reply To",
			"created_at": "2023-10-06T23:47:56.678Z",
			"updated_at": "2023-10-06T23:47:56.678Z",
			"status": "published",
			"published_at": "2023-10-06T23:47:56.678Z",
			"from": "",
			"subject": "",
			"reply_to": ["support@example.com", "help@example.com"],
			"html": "<p>Test</p>",
			"text": "Test",
			"variables": []
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Get(templateId)
	if err != nil {
		t.Errorf("Templates.Get returned error: %v", err)
	}
	assert.Equal(t, "multi-reply-to-id", resp.Id)
	assert.NotNil(t, resp.ReplyTo)

	// ReplyTo is []interface{} when decoded from JSON
	replyTo, ok := resp.ReplyTo.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, 2, len(replyTo))
	assert.Equal(t, "support@example.com", replyTo[0].(string))
	assert.Equal(t, "help@example.com", replyTo[1].(string))
}

func TestUpdateTemplate(t *testing.T) {
	setup()
	defer teardown()

	templateId := "34a080c9-b17d-4187-ad80-5af20266e535"

	mux.HandleFunc(fmt.Sprintf("/templates/%s", templateId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify it
		var req UpdateTemplateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		assert.Equal(t, "welcome-email-updated", req.Name)
		assert.Equal(t, "<strong>Updated content</strong>", req.Html)

		ret := `
		{
			"id": "34a080c9-b17d-4187-ad80-5af20266e535",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Update(templateId, &UpdateTemplateRequest{
		Name: "welcome-email-updated",
		Html: "<strong>Updated content</strong>",
	})
	if err != nil {
		t.Errorf("Templates.Update returned error: %v", err)
	}
	assert.Equal(t, "34a080c9-b17d-4187-ad80-5af20266e535", resp.Id)
	assert.Equal(t, "template", resp.Object)
}

func TestUpdateTemplateWithVariables(t *testing.T) {
	setup()
	defer teardown()

	templateId := "template-with-vars"

	mux.HandleFunc(fmt.Sprintf("/templates/%s", templateId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify variables
		var req UpdateTemplateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		assert.Equal(t, "updated-template", req.Name)
		assert.Equal(t, "<p>Hello {{{NAME}}}, you have {{{COUNT}}} messages</p>", req.Html)
		assert.Equal(t, 2, len(req.Variables))
		assert.Equal(t, "NAME", req.Variables[0].Key)
		assert.Equal(t, VariableTypeString, req.Variables[0].Type)
		assert.Equal(t, "COUNT", req.Variables[1].Key)
		assert.Equal(t, VariableTypeNumber, req.Variables[1].Type)

		ret := `
		{
			"id": "template-with-vars",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Update(templateId, &UpdateTemplateRequest{
		Name: "updated-template",
		Html: "<p>Hello {{{NAME}}}, you have {{{COUNT}}} messages</p>",
		Variables: []*TemplateVariable{
			{
				Key:           "NAME",
				Type:          VariableTypeString,
				FallbackValue: "User",
			},
			{
				Key:           "COUNT",
				Type:          VariableTypeNumber,
				FallbackValue: 0,
			},
		},
	})
	if err != nil {
		t.Errorf("Templates.Update returned error: %v", err)
	}
	assert.Equal(t, "template-with-vars", resp.Id)
	assert.Equal(t, "template", resp.Object)
}

func TestUpdateTemplateByAlias(t *testing.T) {
	setup()
	defer teardown()

	templateAlias := "my-template-alias"

	mux.HandleFunc(fmt.Sprintf("/templates/%s", templateAlias), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "updated-by-alias-id",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Update(templateAlias, &UpdateTemplateRequest{
		Name: "updated-name",
		Html: "<p>Updated by alias</p>",
	})
	if err != nil {
		t.Errorf("Templates.Update returned error: %v", err)
	}
	assert.Equal(t, "updated-by-alias-id", resp.Id)
	assert.Equal(t, "template", resp.Object)
}

func TestUpdateTemplateWithContext(t *testing.T) {
	setup()
	defer teardown()

	templateId := "context-update-id"

	mux.HandleFunc(fmt.Sprintf("/templates/%s", templateId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "context-update-id",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	resp, err := client.Templates.UpdateWithContext(ctx, templateId, &UpdateTemplateRequest{
		Name: "context-updated",
		Html: "<p>Context update</p>",
	})
	if err != nil {
		t.Errorf("Templates.UpdateWithContext returned error: %v", err)
	}
	assert.Equal(t, "context-update-id", resp.Id)
	assert.Equal(t, "template", resp.Object)
}

func TestUpdateTemplateWithAllFields(t *testing.T) {
	setup()
	defer teardown()

	templateId := "full-update-id"

	mux.HandleFunc(fmt.Sprintf("/templates/%s", templateId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Decode request body to verify all fields
		var req UpdateTemplateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		assert.Equal(t, "full-update", req.Name)
		assert.Equal(t, "updated-alias", req.Alias)
		assert.Equal(t, "Updated <updated@example.com>", req.From)
		assert.Equal(t, "Updated Subject", req.Subject)
		assert.Equal(t, "<p>Updated HTML</p>", req.Html)
		assert.Equal(t, "Updated Text", req.Text)

		// ReplyTo can be string or []string
		replyTo, ok := req.ReplyTo.([]interface{})
		assert.True(t, ok)
		assert.Equal(t, 1, len(replyTo))
		assert.Equal(t, "updated@example.com", replyTo[0].(string))

		ret := `
		{
			"id": "full-update-id",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Update(templateId, &UpdateTemplateRequest{
		Name:    "full-update",
		Alias:   "updated-alias",
		From:    "Updated <updated@example.com>",
		Subject: "Updated Subject",
		ReplyTo: []string{"updated@example.com"},
		Html:    "<p>Updated HTML</p>",
		Text:    "Updated Text",
	})
	if err != nil {
		t.Errorf("Templates.Update returned error: %v", err)
	}
	assert.Equal(t, "full-update-id", resp.Id)
	assert.Equal(t, "template", resp.Object)
}

func TestPublishTemplate(t *testing.T) {
	setup()
	defer teardown()

	templateId := "34a080c9-b17d-4187-ad80-5af20266e535"

	mux.HandleFunc(fmt.Sprintf("/templates/%s/publish", templateId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "34a080c9-b17d-4187-ad80-5af20266e535",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Publish(templateId)
	if err != nil {
		t.Errorf("Templates.Publish returned error: %v", err)
	}
	assert.Equal(t, "34a080c9-b17d-4187-ad80-5af20266e535", resp.Id)
	assert.Equal(t, "template", resp.Object)
}

func TestPublishTemplateByAlias(t *testing.T) {
	setup()
	defer teardown()

	templateAlias := "my-template"

	mux.HandleFunc(fmt.Sprintf("/templates/%s/publish", templateAlias), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "published-by-alias-id",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Publish(templateAlias)
	if err != nil {
		t.Errorf("Templates.Publish returned error: %v", err)
	}
	assert.Equal(t, "published-by-alias-id", resp.Id)
	assert.Equal(t, "template", resp.Object)
}

func TestPublishTemplateWithContext(t *testing.T) {
	setup()
	defer teardown()

	templateId := "context-publish-id"

	mux.HandleFunc(fmt.Sprintf("/templates/%s/publish", templateId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "context-publish-id",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	resp, err := client.Templates.PublishWithContext(ctx, templateId)
	if err != nil {
		t.Errorf("Templates.PublishWithContext returned error: %v", err)
	}
	assert.Equal(t, "context-publish-id", resp.Id)
	assert.Equal(t, "template", resp.Object)
}

func TestDuplicateTemplate(t *testing.T) {
	setup()
	defer teardown()

	templateId := "34a080c9-b17d-4187-ad80-5af20266e535"

	mux.HandleFunc(fmt.Sprintf("/templates/%s/duplicate", templateId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "duplicated-template-id-789",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Duplicate(templateId)
	if err != nil {
		t.Errorf("Templates.Duplicate returned error: %v", err)
	}
	assert.Equal(t, "duplicated-template-id-789", resp.Id)
	assert.Equal(t, "template", resp.Object)
}

func TestDuplicateTemplateByAlias(t *testing.T) {
	setup()
	defer teardown()

	templateAlias := "my-template"

	mux.HandleFunc(fmt.Sprintf("/templates/%s/duplicate", templateAlias), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "duplicated-by-alias-id",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Duplicate(templateAlias)
	if err != nil {
		t.Errorf("Templates.Duplicate returned error: %v", err)
	}
	assert.Equal(t, "duplicated-by-alias-id", resp.Id)
	assert.Equal(t, "template", resp.Object)
}

func TestDuplicateTemplateWithContext(t *testing.T) {
	setup()
	defer teardown()

	templateId := "context-duplicate-id"

	mux.HandleFunc(fmt.Sprintf("/templates/%s/duplicate", templateId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"id": "context-duplicate-id-result",
			"object": "template"
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	resp, err := client.Templates.DuplicateWithContext(ctx, templateId)
	if err != nil {
		t.Errorf("Templates.DuplicateWithContext returned error: %v", err)
	}
	assert.Equal(t, "context-duplicate-id-result", resp.Id)
	assert.Equal(t, "template", resp.Object)
}

func TestRemoveTemplate(t *testing.T) {
	setup()
	defer teardown()

	templateId := "34a080c9-b17d-4187-ad80-5af20266e535"

	mux.HandleFunc(fmt.Sprintf("/templates/%s", templateId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "template",
			"id": "34a080c9-b17d-4187-ad80-5af20266e535",
			"deleted": true
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Remove(templateId)
	if err != nil {
		t.Errorf("Templates.Remove returned error: %v", err)
	}
	assert.Equal(t, "template", resp.Object)
	assert.Equal(t, "34a080c9-b17d-4187-ad80-5af20266e535", resp.Id)
	assert.True(t, resp.Deleted)
}

func TestRemoveTemplateByAlias(t *testing.T) {
	setup()
	defer teardown()

	templateAlias := "my-template"

	mux.HandleFunc(fmt.Sprintf("/templates/%s", templateAlias), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "template",
			"id": "removed-by-alias-id",
			"deleted": true
		}`
		fmt.Fprintf(w, ret)
	})

	resp, err := client.Templates.Remove(templateAlias)
	if err != nil {
		t.Errorf("Templates.Remove returned error: %v", err)
	}
	assert.Equal(t, "template", resp.Object)
	assert.Equal(t, "removed-by-alias-id", resp.Id)
	assert.True(t, resp.Deleted)
}

func TestRemoveTemplateWithContext(t *testing.T) {
	setup()
	defer teardown()

	templateId := "context-remove-id"

	mux.HandleFunc(fmt.Sprintf("/templates/%s", templateId), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ret := `
		{
			"object": "template",
			"id": "context-remove-id",
			"deleted": true
		}`
		fmt.Fprintf(w, ret)
	})

	ctx := context.Background()
	resp, err := client.Templates.RemoveWithContext(ctx, templateId)
	if err != nil {
		t.Errorf("Templates.RemoveWithContext returned error: %v", err)
	}
	assert.Equal(t, "template", resp.Object)
	assert.Equal(t, "context-remove-id", resp.Id)
	assert.True(t, resp.Deleted)
}
