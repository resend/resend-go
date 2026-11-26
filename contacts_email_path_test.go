package resend

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// TestContactEmailIdentifierIsPathEscaped checks that an email used as a
// contact identifier is sent as a single, escaped path segment.
func TestContactEmailIdentifierIsPathEscaped(t *testing.T) {
	emails := []string{"sales#eu@example.com", "sales?eu@example.com", "sales%eu@example.com"}

	calls := map[string]func(c *Client, email string) error{
		"get": func(c *Client, email string) error {
			_, err := c.Contacts.GetWithContext(context.Background(), &GetContactOptions{Id: email})
			return err
		},
		"remove": func(c *Client, email string) error {
			_, err := c.Contacts.RemoveWithContext(context.Background(), &RemoveContactOptions{Id: email})
			return err
		},
		"update": func(c *Client, email string) error {
			_, err := c.Contacts.UpdateWithContext(context.Background(), &UpdateContactRequest{Email: email})
			return err
		},
		"topics list": func(c *Client, email string) error {
			_, err := c.Contacts.Topics.ListWithContext(context.Background(), email)
			return err
		},
		"topics update": func(c *Client, email string) error {
			_, err := c.Contacts.Topics.UpdateWithContext(context.Background(), &UpdateContactTopicsRequest{
				Email:  email,
				Topics: []TopicSubscriptionUpdate{{Id: "t", Subscription: "opt_in"}},
			})
			return err
		},
		"segments list": func(c *Client, email string) error {
			_, err := c.Contacts.Segments.ListWithContext(context.Background(), &ListContactSegmentsRequest{Email: email})
			return err
		},
	}

	for name, call := range calls {
		for _, email := range emails {
			t.Run(name+"/"+email, func(t *testing.T) {
				var got *http.Request
				c := NewClient("test")
				c.client = &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					got = r
					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     make(http.Header),
						Body:       io.NopCloser(strings.NewReader(`{}`)),
						Request:    r,
					}, nil
				})}

				err := call(c, email)
				assert.NoError(t, err)
				if got == nil {
					t.Fatal("no request was sent")
				}
				assert.Empty(t, got.URL.RawQuery)
				assert.Empty(t, got.URL.Fragment)
				assert.True(t, strings.Contains(got.URL.EscapedPath(), url.PathEscape(email)),
					"escaped path %q does not contain %q", got.URL.EscapedPath(), url.PathEscape(email))
				assert.Contains(t, got.URL.Path, email)
			})
		}
	}
}
