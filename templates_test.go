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

		assert.Equal(t, 5, len(req.Variables))

		// String variable
		assert.Equal(t, "STRING_VAR", req.Variables[0].Key)
		assert.Equal(t, VariableTypeString, req.Variables[0].Type)
		assert.Equal(t, "default", req.Variables[0].FallbackValue)

		// Number variable
		assert.Equal(t, "NUMBER_VAR", req.Variables[1].Key)
		assert.Equal(t, VariableTypeNumber, req.Variables[1].Type)
		assert.Equal(t, float64(42), req.Variables[1].FallbackValue)

		// Boolean variable
		assert.Equal(t, "BOOLEAN_VAR", req.Variables[2].Key)
		assert.Equal(t, VariableTypeBoolean, req.Variables[2].Type)
		assert.Equal(t, true, req.Variables[2].FallbackValue)

		// Object variable
		assert.Equal(t, "OBJECT_VAR", req.Variables[3].Key)
		assert.Equal(t, VariableTypeObject, req.Variables[3].Type)

		// List variable
		assert.Equal(t, "LIST_VAR", req.Variables[4].Key)
		assert.Equal(t, VariableTypeList, req.Variables[4].Type)
		assert.NotNil(t, req.Variables[4].FallbackValue)

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
			{
				Key:           "BOOLEAN_VAR",
				Type:          VariableTypeBoolean,
				FallbackValue: true,
			},
			{
				Key:           "OBJECT_VAR",
				Type:          VariableTypeObject,
				FallbackValue: map[string]interface{}{"key": "value"},
			},
			{
				Key:           "LIST_VAR",
				Type:          VariableTypeList,
				FallbackValue: []interface{}{"item1", "item2"},
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
