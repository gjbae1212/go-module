package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetModulePath(t *testing.T) {
	assert := assert.New(t)
	assert.NotEmpty(GetModulePath())
}


func TestCheckHttpMethod(t *testing.T) {
	assert := assert.New(t)
	assert.True(CheckHttpMethod("get"))
	assert.True(CheckHttpMethod("GET"))
	assert.False(CheckHttpMethod("GE"))
}
