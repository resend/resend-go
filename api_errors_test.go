package resend_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	resend "github.com/resend/resend-go/v4"
)

type errorResponseTransport func(*http.Request) (*http.Response, error)

func (f errorResponseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestEmailsGetReturnsStructuredAPIError(t *testing.T) {
	cases := []struct {
		name        string
		status      int
		contentType string
		body        string
		wantName    string
		wantMessage string
	}{
		{"bad request", 400, "application/json", `{"name":"validation_error","message":"Invalid input"}`, "validation_error", "Invalid input"},
		{"unauthorized", 401, "application/json", `{"name":"validation_error","message":"Invalid API key"}`, "validation_error", "Invalid API key"},
		{"forbidden", 403, "application/json", `{"name":"restricted_api_key","message":"Not allowed"}`, "restricted_api_key", "Not allowed"},
		{"not found", 404, "application/json", `{"name":"not_found","message":"Email not found"}`, "not_found", "Email not found"},
		{"conflict", 409, "application/json", `{"name":"invalid_idempotent_request","message":"Key already used"}`, "invalid_idempotent_request", "Key already used"},
		{"validation", 422, "application/json; charset=utf-8", `{"name":"validation_error","message":"Invalid input"}`, "validation_error", "Invalid input"},
		{"server error", 500, "application/json", `{"name":"internal_server_error","message":"Try again later"}`, "internal_server_error", "Try again later"},
		{"HTTP status overrides body", 503, "application/json", `{"statusCode":200,"name":"internal_server_error","message":"Unavailable"}`, "internal_server_error", "Unavailable"},
		{"non-JSON body", 502, "text/html", `<html>Bad Gateway</html>`, "", "502 Bad Gateway"},
		{"missing content type", 503, "", `Unavailable`, "", "503 Service Unavailable"},
		{"malformed JSON", 422, "application/json", `{`, "", "422 Unprocessable Entity"},
		{"invalid field type", 500, "application/json", `{"name":"internal_server_error","message":42}`, "", "500 Internal Server Error"},
		{"empty body", 503, "application/json", ``, "", "503 Service Unavailable"},
		{"missing message", 500, "application/json", `{"name":"internal_server_error"}`, "internal_server_error", "Unknown Error"},
		{"null body", 500, "application/json", `null`, "", "Unknown Error"},
		{"empty bad request message", 400, "application/json", `{}`, "", ""},
		{"empty validation message", 422, "application/json", `{}`, "", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := resend.NewCustomClient(&http.Client{
				Transport: errorResponseTransport(func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: tc.status,
						Status:     fmt.Sprintf("%d %s", tc.status, http.StatusText(tc.status)),
						Header:     http.Header{"Content-Type": {tc.contentType}},
						Body:       io.NopCloser(strings.NewReader(tc.body)),
						Request:    req,
					}, nil
				}),
			}, "re_test")

			_, err := client.Emails.Get("test-email")
			if err == nil {
				t.Fatal("expected an API error")
			}
			var apiError *resend.ResendError
			if !errors.As(fmt.Errorf("get email: %w", err), &apiError) {
				t.Fatalf("expected *resend.ResendError, got %T: %v", err, err)
			}
			if apiError.StatusCode != tc.status || apiError.Name != tc.wantName || apiError.Message != tc.wantMessage {
				t.Errorf("got %+v, want status=%d name=%q message=%q", apiError, tc.status, tc.wantName, tc.wantMessage)
			}
			if want := "[ERROR]: " + tc.wantMessage; err.Error() != want {
				t.Errorf("Error() = %q, want %q", err.Error(), want)
			}
			if errors.Is(err, resend.ErrRateLimit) {
				t.Error("non-429 error must not match ErrRateLimit")
			}
		})
	}
}

func TestEmailsGetPreservesRateLimitError(t *testing.T) {
	client := resend.NewCustomClient(&http.Client{
		Transport: errorResponseTransport(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Status:     "429 Too Many Requests",
				Header: http.Header{
					"Content-Type":        {"application/json"},
					"Ratelimit-Limit":     {"10"},
					"Ratelimit-Remaining": {"0"},
					"Ratelimit-Reset":     {"60"},
					"Retry-After":         {"60"},
				},
				Body:    io.NopCloser(strings.NewReader(`{"message":"Too many requests"}`)),
				Request: req,
			}, nil
		}),
	}, "re_test")

	_, err := client.Emails.Get("test-email")
	wrapped := fmt.Errorf("get email: %w", err)
	if !errors.Is(wrapped, resend.ErrRateLimit) {
		t.Fatalf("expected ErrRateLimit, got %v", err)
	}
	var rateLimit *resend.RateLimitError
	if !errors.As(wrapped, &rateLimit) {
		t.Fatalf("expected *resend.RateLimitError, got %T", err)
	}
	if rateLimit.Message != "Too many requests" || rateLimit.Limit != "10" || rateLimit.Remaining != "0" || rateLimit.Reset != "60" || rateLimit.RetryAfter != "60" {
		t.Errorf("unexpected rate-limit metadata: %+v", rateLimit)
	}
	var apiError *resend.ResendError
	if errors.As(wrapped, &apiError) {
		t.Error("429 must continue to return RateLimitError")
	}
}
