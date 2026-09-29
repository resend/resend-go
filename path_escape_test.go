package resend

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentifiersAreEscapedInRequestPaths(t *testing.T) {
	setup()
	defer teardown()

	const raw = "a/b?c#d%e@example.com"
	const escaped = "a%2Fb%3Fc%23d%25e@example.com"

	var got []string
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		got = append(got, r.Method+" "+r.URL.EscapedPath())
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{}`))
	})

	client.Contacts.Get(&GetContactOptions{Id: raw})
	client.Contacts.Update(&UpdateContactRequest{Email: raw})
	client.Contacts.Remove(&RemoveContactOptions{Id: raw})
	client.Contacts.Topics.List(raw)
	client.Contacts.Topics.Update(&UpdateContactTopicsRequest{
		Email:  raw,
		Topics: []TopicSubscriptionUpdate{{Id: "t1", Subscription: "opt_in"}},
	})
	client.Contacts.Segments.List(&ListContactSegmentsRequest{Email: raw})
	client.Contacts.Segments.Add(&AddContactSegmentRequest{Email: raw, SegmentId: "s1"})
	client.Contacts.Segments.Remove(&RemoveContactSegmentRequest{Email: raw, SegmentId: "s1"})
	client.Events.Get(raw)
	client.Events.Remove(raw)
	client.Templates.Get(raw)
	client.Templates.Remove(raw)

	assert.Equal(t, []string{
		"GET /contacts/" + escaped,
		"PATCH /contacts/" + escaped,
		"DELETE /contacts/" + escaped,
		"GET /contacts/" + escaped + "/topics",
		"PATCH /contacts/" + escaped + "/topics",
		"GET /contacts/" + escaped + "/segments",
		"POST /contacts/" + escaped + "/segments/s1",
		"DELETE /contacts/" + escaped + "/segments/s1",
		"GET /events/" + escaped,
		"DELETE /events/" + escaped,
		"GET /templates/" + escaped,
		"DELETE /templates/" + escaped,
	}, got)
}
