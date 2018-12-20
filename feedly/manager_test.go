package feedly

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	assert := assert.New(t)
	_, err := NewManager("accesstoken", "refreshtoken")
	assert.NoError(err)
}
