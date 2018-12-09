package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetModulePath(t *testing.T) {
	assert := assert.New(t)
	assert.NotEmpty(GetModulePath())
}
