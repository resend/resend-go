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

func TestHandleError(t *testing.T) {
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
			resp: &http.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Status:     fmt.Sprintf("%d %s", http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity)),
				Header:     http.Header{"Content-Type": {"application/json; charset=utf-8"}},
				Body:       io.NopCloser(bytes.NewBufferString(`{"message":"Validation error"}`)),
			},
			want: errors.New("[ERROR]: Validation error"),
		},
		{
			desc: "validation_error_no_json",
			resp: &http.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Status:     fmt.Sprintf("%d %s", http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity)),
				Body:       io.NopCloser(bytes.NewBufferString(`Validation error`)),
			},
			want: errors.New("[ERROR]: 422 Unprocessable Entity"),
		},
		{
			desc: "bad_request",
			resp: &http.Response{
				StatusCode: http.StatusBadRequest,
				Status:     fmt.Sprintf("%d %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest)),
				Header:     http.Header{"Content-Type": {"application/json; charset=utf-8"}},
				Body:       io.NopCloser(bytes.NewBufferString(`{"message":"Validation error"}`)),
			},
			want: errors.New("[ERROR]: Validation error"),
		},
		{
			desc: "bad_request_no_json",
			resp: &http.Response{
				StatusCode: http.StatusBadRequest,
				Status:     fmt.Sprintf("%d %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest)),
				Body:       io.NopCloser(bytes.NewBufferString(`Validation error`)),
			},
			want: errors.New("[ERROR]: 400 Bad Request"),
		},
		{
			desc: "bad_request_invalid_json",
			resp: &http.Response{
				StatusCode: http.StatusBadRequest,
				Status:     fmt.Sprintf("%d %s", http.StatusBadRequest, http.StatusText(http.StatusBadRequest)),
				Header:     http.Header{"Content-Type": {"application/json; charset=utf-8"}},
				Body:       io.NopCloser(bytes.NewBufferString(`{`)),
			},
			want: errors.New("[ERROR]: 400 Bad Request"),
		},
		{
			desc: "server_error",
			resp: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Status:     fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)),
				Header:     http.Header{"Content-Type": {"application/json; charset=utf-8"}},
				Body:       io.NopCloser(bytes.NewBufferString(`{"message":"Server error"}`)),
			},
			want: errors.New("[ERROR]: Server error"),
		},
		{
			desc: "server_error_no_json",
			resp: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Status:     fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)),
				Body:       io.NopCloser(bytes.NewBufferString(`Server error`)),
			},
			want: errors.New("[ERROR]: 500 Internal Server Error"),
		},
		{
			desc: "server_error_invalid_json",
			resp: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Status:     fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)),
				Header:     http.Header{"Content-Type": {"application/json; charset=utf-8"}},
				Body:       io.NopCloser(bytes.NewBufferString(`{`)),
			},
			want: errors.New("[ERROR]: 500 Internal Server Error"),
		},
		{
			desc: "server_error_no_message",
			resp: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Status:     fmt.Sprintf("%d %s", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)),
				Header:     http.Header{"Content-Type": {"application/json; charset=utf-8"}},
				Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
			},
			want: errors.New("[ERROR]: Unknown Error"),
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			err := handleError(c.resp)
			assert.Equal(t, c.want, err)
		})
	}
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

type failingMarshaler struct{ called *bool }

func (f failingMarshaler) MarshalJSON() ([]byte, error) {
	*f.called = true
	return nil, errors.New("marshal failed")
}

func TestNewRequestValidatesBeforeEncoding(t *testing.T) {
	client := NewClient("123")

	called := false
	req, err := client.NewRequest(context.Background(), "BAD METHOD", "/emails", failingMarshaler{&called})
	assert.Error(t, err)
	assert.Nil(t, req)
	assert.False(t, called, "params must not be marshalled when the method is invalid")

	called = false
	//lint:ignore SA1012 intentionally passing a nil context
	req, err = client.NewRequest(nil, http.MethodPost, "/emails", failingMarshaler{&called}) //nolint:staticcheck
	assert.Error(t, err)
	assert.Nil(t, req)
	assert.False(t, called, "params must not be marshalled when the context is nil")

	called = false
	req, err = client.NewRequest(context.Background(), http.MethodPost, "/emails", failingMarshaler{&called})
	assert.Error(t, err)
	assert.Nil(t, req)
	assert.True(t, called)
}
