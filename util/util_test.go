package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetModulePath(t *testing.T) {
	assert := assert.New(t)
	assert.NotEmpty(GetModulePath())
}

func TestStringInSlice(t *testing.T) {
	assert := assert.New(t)

	a := []string{"a", "b", "c"}
	assert.True(StringInSlice("a", a), "String in Slice Not Match")

	b := []string{"a", "b", "c"}
	assert.False(StringInSlice("d", b), "String in Slice Not Match")
}

func TestCheckHttpMethod(t *testing.T) {
	assert := assert.New(t)
	assert.True(CheckHttpMethod("get"))
	assert.True(CheckHttpMethod("GET"))
	assert.False(CheckHttpMethod("GE"))
}
