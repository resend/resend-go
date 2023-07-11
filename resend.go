package resend

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	version        = "1.7.0"
	defaultBaseURL = "https://api.resend.com/"
	userAgent      = "resend-go/" + version
	contentType    = "application/json"
)

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
	Emails  EmailsSvc
	ApiKeys ApiKeysSvc
	Domains DomainsSvc
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
	c.ApiKeys = &ApiKeysSvcImpl{client: c}
	c.Domains = &DomainsSvcImpl{client: c}

	c.ApiKey = apiKey
	c.headers = make(map[string]string)
	return c
}

// NewRequest builds and returns a new HTTP request object
// based on the given arguments
func (c *Client) NewRequest(method, path string, params interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	var req *http.Request
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		req, err = http.NewRequest(method, u.String(), nil)
		if err != nil {
			return nil, err
		}

	default:
		buf := new(bytes.Buffer)
		if params != nil {
			err = json.NewEncoder(buf).Encode(params)
			if err != nil {
				return nil, err
			}
		}

		req, err = http.NewRequest(method, u.String(), buf)
		if err != nil {
			return nil, err
		}
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
