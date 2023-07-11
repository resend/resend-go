package resend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteArrayToArrayString(t *testing.T) {
	assert.Equal(t, ByteArrayToStringArray([]byte{44, 45, 46}), []string{"44", "45", "46"})
}
