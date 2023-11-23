package resend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytesToIntArray(t *testing.T) {
	assert.Equal(t, BytesToIntArray([]byte{44, 45, 46}), []int{44, 45, 46})
}
