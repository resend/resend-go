package resend

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	version     = "2.28.0"
	userAgent   = "resend-go/" + version
	contentType = "application/json"
)

var defaultBaseURL = getEnv("RESEND_BASE_URL", "https://api.resend.com/")

var defaultHTTPClient = &http.Client{
	Timeout: time.Minute,
}

// Options interface is used to define additional options that can be passed
// to the API methods.
type Options interface {
	// GetIdempotencyKey returns the idempotency key
	GetIdempotencyKey() string
}

// Client handles communication with Resend API.
type Client struct {
	// HTTP client
	client *http.Client

	// Api Key
	ApiKey string

	// Base URL
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	// HTTP headers
	headers map[string]string

	// Services
	Emails     EmailsSvc
	Batch      BatchSvc
	ApiKeys    ApiKeysSvc
	Domains    DomainsSvc
	Audiences  AudiencesSvc
	Contacts   ContactsSvc
	Broadcasts BroadcastsSvc
	Webhooks   WebhooksSvc
}

// NewClient is the default client constructor
func NewClient(apiKey string) *Client {
	key := strings.Trim(strings.TrimSpace(apiKey), "'")
	return NewCustomClient(defaultHTTPClient, key)
}

// NewCustomClient builds a new Resend API client, using a provided Http client.
func NewCustomClient(httpClient *http.Client, apiKey string) *Client {
	if httpClient == nil {
		httpClient = defaultHTTPClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}

	c.Emails = &EmailsSvcImpl{client: c}
	c.Batch = &BatchSvcImpl{client: c}
	c.ApiKeys = &ApiKeysSvcImpl{client: c}
	c.Domains = &DomainsSvcImpl{client: c}
	c.Audiences = &AudiencesSvcImpl{client: c}
	c.Contacts = &ContactsSvcImpl{client: c}
	c.Broadcasts = &BroadcastsSvcImpl{client: c}
	c.Webhooks = &WebhooksSvcImpl{client: c}

	c.ApiKey = apiKey
	c.headers = make(map[string]string)
	return c
}

// NewRequestWithOptions builds and returns a new HTTP request object
// based on the given arguments and options
// It is used to set additional options like idempotency key
func (c *Client) NewRequestWithOptions(ctx context.Context, method, path string, params interface{}, options Options) (*http.Request, error) {
	req, err := c.NewRequest(ctx, method, path, params)

	if err != nil {
		return nil, err
	}

	// Set the idempotency key if provided and the method is POST for now.
	if options != nil {
		if options.GetIdempotencyKey() != "" && method == http.MethodPost {
			req.Header.Set("Idempotency-Key", options.GetIdempotencyKey())
		}

		// Handle batch-specific options
		if batchOptions, ok := options.(*BatchSendEmailOptions); ok {
			if batchOptions.GetBatchValidation() != "" {
				req.Header.Set("x-batch-validation", batchOptions.GetBatchValidation())
			}
		}
	}

	return req, nil
}

// NewRequest builds and returns a new HTTP request object
// based on the given arguments
func (c *Client) NewRequest(ctx context.Context, method, path string, params interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if params != nil {
		buf := new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(params)
		if err != nil {
			return nil, err
		}

		req.Body = io.NopCloser(buf)
		req.Header.Set("Content-Type", contentType)
	}

	for k, v := range c.headers {
		req.Header.Add(k, v)
	}

	req.Header.Set("Accept", contentType)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)

	return req, nil
}

// Perform sends the request to the Resend API
func (c *Client) Perform(req *http.Request, ret interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle possible errors.
	// Any 2xx status code is considered success
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, handleError(resp)
	}

	if resp.StatusCode != http.StatusNoContent && ret != nil {
		if w, ok := ret.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			if resp.Body != nil {
				err = json.NewDecoder(resp.Body).Decode(ret)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return resp, err
}

// handleError tries to handle errors based on HTTP status codes
func handleError(resp *http.Response) error {
	switch resp.StatusCode {

	// Handle rate limit errors (429)
	case http.StatusTooManyRequests:
		r := &DefaultError{}
		if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
			err := json.NewDecoder(resp.Body).Decode(r)
			if err != nil {
				r.Message = resp.Status
			}
		} else {
			r.Message = resp.Status
		}

		return &RateLimitError{
			Message:    r.Message,
			Limit:      resp.Header.Get("ratelimit-limit"),
			Remaining:  resp.Header.Get("ratelimit-remaining"),
			Reset:      resp.Header.Get("ratelimit-reset"),
			RetryAfter: resp.Header.Get("retry-after"),
		}

	// Handles errors most likely caused by the client
	case http.StatusUnprocessableEntity, http.StatusBadRequest:
		r := &InvalidRequestError{}
		if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
			err := json.NewDecoder(resp.Body).Decode(r)
			if err != nil {
				r.Message = resp.Status
			}
		} else {
			r.Message = resp.Status
		}

		// TODO: replace this with a new ResendError type
		return errors.New("[ERROR]: " + r.Message)
	default:
		// Tries to parse `message` attr from error
		r := &DefaultError{}

		if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
			err := json.NewDecoder(resp.Body).Decode(r)
			if err != nil {
				r.Message = resp.Status
			}
		} else {
			r.Message = resp.Status
		}

		if r.Message != "" {
			// TODO: replace this with a new ResendError type
			return errors.New("[ERROR]: " + r.Message)
		}
		return errors.New("[ERROR]: Unknown Error")
	}
}

type InvalidRequestError struct {
	StatusCode int    `json:"statusCode"`
	Name       string `json:"name"`
	Message    string `json:"message"`
}

type DefaultError struct {
	Message string `json:"message"`
}

// ListOptions contains pagination parameters for list methods
type ListOptions struct {
	Limit  *int    `json:"limit,omitempty"`
	After  *string `json:"after,omitempty"`
	Before *string `json:"before,omitempty"`
}

// buildPaginationQuery constructs query parameters for pagination
func buildPaginationQuery(options *ListOptions) string {
	if options == nil {
		return ""
	}

	query := make(url.Values)
	if options.Limit != nil {
		query.Set("limit", fmt.Sprintf("%d", *options.Limit))
	}
	if options.After != nil {
		query.Set("after", *options.After)
	}
	if options.Before != nil {
		query.Set("before", *options.Before)
	}

	if len(query) > 0 {
		return "?" + query.Encode()
	}
	return ""
}
