package gcp_oauth2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	assert := assert.New(t)

	_, err := NewManager("", "", "", []string{})
	assert.Error(err)

	_, err = NewManager("id", "secret", "url", []string{ScopeProfile, ScopeEmail})
	assert.NoError(err)
}
