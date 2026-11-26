package resend

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVerifyWebhookSignatureVersionPrefix(t *testing.T) {
	secretBytes := []byte("test-secret-bytes-1234567890")
	secret := "whsec_" + base64.StdEncoding.EncodeToString(secretBytes)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	msgId := "msg_123"
	payload := `{"type":"email.sent"}`
	signedContent := fmt.Sprintf("%s.%s.%s", msgId, timestamp, payload)
	sig := generateSignature(secretBytes, []byte(signedContent))

	client := NewClient("123")

	err := client.Webhooks.Verify(&VerifyWebhookOptions{
		Payload: payload,
		Headers: WebhookHeaders{
			Id:        msgId,
			Timestamp: timestamp,
			Signature: "v2," + sig,
		},
		WebhookSecret: secret,
	})
	assert.Error(t, err)

	err = client.Webhooks.Verify(&VerifyWebhookOptions{
		Payload: payload,
		Headers: WebhookHeaders{
			Id:        msgId,
			Timestamp: timestamp,
			Signature: "v2,bad v1," + sig,
		},
		WebhookSecret: secret,
	})
	assert.NoError(t, err)
}
