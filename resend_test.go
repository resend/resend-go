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
