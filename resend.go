package resend

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	version     = "2.11.0"
	userAgent   = "resend-go/" + version
	contentType = "application/json"
)

var defaultBaseURL = getEnv("RESEND_BASE_URL", "https://api.resend.com/")

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
	Emails    EmailsSvc
	Batch     BatchSvc
	ApiKeys   ApiKeysSvc
	Domains   DomainsSvc
	Audiences AudiencesSvc
	Contacts  ContactsSvc
}

// NewClient is the default client constructor
func NewClient(apiKey string) *Client {
	key := strings.Trim(strings.TrimSpace(apiKey), "'")
	return NewCustomClient(http.DefaultClient, key)
}

// NewCustomClient builds a new Resend API client, using a provided Http client.
func NewCustomClient(httpClient *http.Client, apiKey string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}

	c.Emails = &EmailsSvcImpl{client: c}
	c.Batch = &BatchSvcImpl{client: c}
	c.ApiKeys = &ApiKeysSvcImpl{client: c}
	c.Domains = &DomainsSvcImpl{client: c}
	c.Audiences = &AudiencesSvcImpl{client: c}
	c.Contacts = &ContactsSvcImpl{client: c}

	c.ApiKey = apiKey
	c.headers = make(map[string]string)
	return c
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
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
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

	// Handles errors most likely caused by the client
	case http.StatusUnprocessableEntity, http.StatusBadRequest:
		r := &InvalidRequestError{}
		err := json.NewDecoder(resp.Body).Decode(r)
		if err != nil {
			return err
		}
		return errors.New("[ERROR]: " + r.Message)
	default:
		// Tries to parse `message` attr from error
		r := &DefaultError{}
		err := json.NewDecoder(resp.Body).Decode(r)
		if err != nil {
			return err
		}
		if r.Message != "" {
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
