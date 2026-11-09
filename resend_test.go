package resend

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResend(t *testing.T) {
	client := NewClient("123")
	assert.NotNil(t, client)
}

func TestResendRequestHeaders(t *testing.T) {
	ctx := context.TODO()
	client := NewClient("123")
	params := &SendEmailRequest{
		To: []string{"email@example.com", "email2@example.com"},
	}
	req, err := client.NewRequest(ctx, "POST", "/emails/", params)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, req.Header["Accept"][0], "application/json")
	assert.Equal(t, req.Header["Content-Type"][0], "application/json")
	assert.Equal(t, req.Method, http.MethodPost)
	assert.Equal(t, req.URL.String(), "https://api.resend.com/emails/")
	assert.Equal(t, req.Header["Authorization"][0], "Bearer 123")

	_, ok := req.Header["Idempotency-Key"]
	assert.False(t, ok, "expected 'Idempotency-Key' header to be absent")
}

func TestResendRequestShouldReturnErrorIfContextIsCancelled(t *testing.T) {
	client := NewClient("123")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req, err := client.NewRequest(ctx, "POST", "/", nil)
	if err != nil {
		t.Error(err)
	}

	res, err := client.Perform(req, nil)
	assert.True(t, errors.Unwrap(err) == context.Canceled)
	assert.Nil(t, res)
}

func errorResponse(status int, contentType string, body string) *http.Response {
	header := http.Header{}
	if contentType != "" {
		header.Set("Content-Type", contentType)
	}
	return &http.Response{
		StatusCode: status,
		Status:     fmt.Sprintf("%d %s", status, http.StatusText(status)),
		Header:     header,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}

func TestHandleError(t *testing.T) {
	const jsonType = "application/json; charset=utf-8"

	cases := []struct {
		desc string
		resp *http.Response
		want error
	}{
		{
			desc: "rate_limit_error",
			resp: &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Status:     fmt.Sprintf("%d %s", http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests)),
				Header: http.Header{
					"Content-Type":        {"application/json; charset=utf-8"},
					"Ratelimit-Limit":     {"2"},
					"Ratelimit-Remaining": {"0"},
					"Ratelimit-Reset":     {"1"},
					"Retry-After":         {"1"},
				},
				Body: io.NopCloser(bytes.NewBufferString(`{"message":"Rate limit exceeded"}`)),
			},
			want: &RateLimitError{
				Message:    "Rate limit exceeded",
				Limit:      "2",
				Remaining:  "0",
				Reset:      "1",
				RetryAfter: "1",
			},
		},
		{
			desc: "validation_error",
			resp: errorResponse(http.StatusUnprocessableEntity, jsonType, `{"name":"validation_error","message":"Validation error"}`),
			want: &ResendError{StatusCode: http.StatusUnprocessableEntity, Name: "validation_error", Message: "Validation error"},
		},
		{
			desc: "validation_error_no_json",
			resp: errorResponse(http.StatusUnprocessableEntity, "", `Validation error`),
			want: &ResendError{StatusCode: http.StatusUnprocessableEntity, Message: "422 Unprocessable Entity"},
		},
		{
			desc: "validation_error_no_message",
			resp: errorResponse(http.StatusUnprocessableEntity, jsonType, `{}`),
			want: &ResendError{StatusCode: http.StatusUnprocessableEntity, Message: "Unknown Error"},
		},
		{
			desc: "bad_request",
			resp: errorResponse(http.StatusBadRequest, jsonType, `{"message":"Validation error"}`),
			want: &ResendError{StatusCode: http.StatusBadRequest, Message: "Validation error"},
		},
		{
			desc: "bad_request_no_json",
			resp: errorResponse(http.StatusBadRequest, "", `Validation error`),
			want: &ResendError{StatusCode: http.StatusBadRequest, Message: "400 Bad Request"},
		},
		{
			desc: "bad_request_invalid_json",
			resp: errorResponse(http.StatusBadRequest, jsonType, `{`),
			want: &ResendError{StatusCode: http.StatusBadRequest, Message: "400 Bad Request"},
		},
		{
			desc: "bad_request_no_message",
			resp: errorResponse(http.StatusBadRequest, jsonType, `{}`),
			want: &ResendError{StatusCode: http.StatusBadRequest, Message: "Unknown Error"},
		},
		{
			desc: "forbidden",
			resp: errorResponse(http.StatusForbidden, jsonType, `{"name":"restricted_api_key","message":"Not allowed"}`),
			want: &ResendError{StatusCode: http.StatusForbidden, Name: "restricted_api_key", Message: "Not allowed"},
		},
		{
			desc: "server_error",
			resp: errorResponse(http.StatusInternalServerError, jsonType, `{"message":"Server error"}`),
			want: &ResendError{StatusCode: http.StatusInternalServerError, Message: "Server error"},
		},
		{
			desc: "server_error_no_json",
			resp: errorResponse(http.StatusInternalServerError, "", `Server error`),
			want: &ResendError{StatusCode: http.StatusInternalServerError, Message: "500 Internal Server Error"},
		},
		{
			desc: "server_error_invalid_json",
			resp: errorResponse(http.StatusInternalServerError, jsonType, `{`),
			want: &ResendError{StatusCode: http.StatusInternalServerError, Message: "500 Internal Server Error"},
		},
		{
			desc: "server_error_no_message",
			resp: errorResponse(http.StatusInternalServerError, jsonType, `{}`),
			want: &ResendError{StatusCode: http.StatusInternalServerError, Message: "Unknown Error"},
		},
		{
			desc: "server_error_partially_decoded_body",
			resp: errorResponse(http.StatusInternalServerError, jsonType, `{"name":"internal_server_error","message":42}`),
			want: &ResendError{StatusCode: http.StatusInternalServerError, Message: "500 Internal Server Error"},
		},
		{
			desc: "server_error_null_body",
			resp: errorResponse(http.StatusInternalServerError, jsonType, `null`),
			want: &ResendError{StatusCode: http.StatusInternalServerError, Message: "Unknown Error"},
		},
		{
			desc: "server_error_empty_body",
			resp: errorResponse(http.StatusServiceUnavailable, jsonType, ``),
			want: &ResendError{StatusCode: http.StatusServiceUnavailable, Message: "503 Service Unavailable"},
		},
		{
			desc: "status_code_from_response_not_body",
			resp: errorResponse(http.StatusServiceUnavailable, jsonType, `{"statusCode":200,"name":"internal_server_error","message":"Unavailable"}`),
			want: &ResendError{StatusCode: http.StatusServiceUnavailable, Name: "internal_server_error", Message: "Unavailable"},
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			err := handleError(c.resp)
			assert.Equal(t, c.want, err)
		})
	}
}

func TestResendErrorPreservesMessageFormat(t *testing.T) {
	err := &ResendError{StatusCode: http.StatusNotFound, Name: "not_found", Message: "Email not found"}
	assert.EqualError(t, err, "[ERROR]: Email not found")
}

func TestPerformReturnsResendErrorForNonRateLimitResponses(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emails/missing", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"statusCode":404,"name":"not_found","message":"Email not found"}`)
	})

	_, err := client.Emails.Get("missing")

	var apiErr *ResendError
	assert.True(t, errors.As(fmt.Errorf("get email: %w", err), &apiErr))
	assert.Equal(t, &ResendError{StatusCode: http.StatusNotFound, Name: "not_found", Message: "Email not found"}, apiErr)
	assert.False(t, errors.Is(err, ErrRateLimit))
}

func TestRateLimitErrorIs(t *testing.T) {
	// Create a rate limit error
	rateLimitErr := &RateLimitError{
		Message:    "Rate limit exceeded",
		Limit:      "2",
		Remaining:  "0",
		Reset:      "1",
		RetryAfter: "1",
	}

	// Test that errors.Is correctly identifies RateLimitError
	assert.True(t, errors.Is(rateLimitErr, ErrRateLimit))

	// Test that a regular error is not identified as a rate limit error
	regularErr := errors.New("some other error")
	assert.False(t, errors.Is(regularErr, ErrRateLimit))
}

func TestRateLimitErrorHandling(t *testing.T) {
	// Simulate a 429 response
	resp := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Status:     fmt.Sprintf("%d %s", http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests)),
		Header: http.Header{
			"Content-Type":        {"application/json; charset=utf-8"},
			"Ratelimit-Limit":     {"10"},
			"Ratelimit-Remaining": {"0"},
			"Ratelimit-Reset":     {"60"},
			"Retry-After":         {"60"},
		},
		Body: io.NopCloser(bytes.NewBufferString(`{"message":"Too many requests"}`)),
	}

	err := handleError(resp)

	// Verify it's a RateLimitError
	assert.True(t, errors.Is(err, ErrRateLimit))

	// Verify we can type assert to access fields
	var rateLimitErr *RateLimitError
	assert.True(t, errors.As(err, &rateLimitErr))
	assert.Equal(t, "Too many requests", rateLimitErr.Message)
	assert.Equal(t, "10", rateLimitErr.Limit)
	assert.Equal(t, "0", rateLimitErr.Remaining)
	assert.Equal(t, "60", rateLimitErr.Reset)
	assert.Equal(t, "60", rateLimitErr.RetryAfter)
}

func TestNewRequestSetsReplayableBody(t *testing.T) {
	client := NewClient("123")
	req, err := client.NewRequest(context.Background(), http.MethodPost, "/emails", map[string]string{"a": "b"})
	assert.NoError(t, err)
	want := "{\"a\":\"b\"}\n"
	assert.Equal(t, int64(len(want)), req.ContentLength)
	assert.NotNil(t, req.GetBody)
	body, err := req.GetBody()
	assert.NoError(t, err)
	got, err := io.ReadAll(body)
	assert.NoError(t, err)
	assert.Equal(t, want, string(got))
}

func TestNewRequestWithoutParamsHasNoBody(t *testing.T) {
	client := NewClient("123")
	req, err := client.NewRequest(context.Background(), http.MethodGet, "/emails", nil)
	assert.NoError(t, err)
	assert.Nil(t, req.Body)
	assert.Equal(t, int64(0), req.ContentLength)
	_, ok := req.Header["Content-Type"]
	assert.False(t, ok)
}

func TestNewRequestReturnsEncodeError(t *testing.T) {
	client := NewClient("123")
	req, err := client.NewRequest(context.Background(), http.MethodPost, "/emails", make(chan int))
	assert.Error(t, err)
	assert.Nil(t, req)
}

func TestPerformFollowsTemporaryAndPermanentRedirects(t *testing.T) {
	for _, status := range []int{http.StatusTemporaryRedirect, http.StatusPermanentRedirect} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			var gotMethod, gotBody, gotKey, gotType string
			dest := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				b, _ := io.ReadAll(r.Body)
				gotMethod, gotBody = r.Method, string(b)
				gotKey, gotType = r.Header.Get("Idempotency-Key"), r.Header.Get("Content-Type")
				w.Write([]byte(`{"id":"1"}`))
			}))
			defer dest.Close()
			src := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, dest.URL+"/emails", status)
			}))
			defer src.Close()

			client := NewClient("123")
			base, err := url.Parse(src.URL + "/")
			assert.NoError(t, err)
			client.BaseURL = base

			opts := &SendEmailOptions{IdempotencyKey: "key-1"}
			req, err := client.NewRequestWithOptions(context.Background(), http.MethodPost, "emails", map[string]string{"a": "b"}, opts)
			assert.NoError(t, err)

			var out map[string]string
			_, err = client.Perform(req, &out)
			assert.NoError(t, err)
			assert.Equal(t, "1", out["id"])
			assert.Equal(t, http.MethodPost, gotMethod)
			assert.Equal(t, "{\"a\":\"b\"}\n", gotBody)
			assert.Equal(t, "key-1", gotKey)
			assert.Equal(t, "application/json", gotType)
		})
	}
}

func TestNewRequestReturnsErrorForInvalidMethod(t *testing.T) {
	client := NewClient("123")
	for _, params := range []any{nil, map[string]string{"a": "b"}} {
		req, err := client.NewRequest(context.Background(), "BAD METHOD", "/emails", params)
		assert.Error(t, err)
		assert.Nil(t, req)
	}
}
