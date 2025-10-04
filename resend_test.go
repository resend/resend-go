package resend

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
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
