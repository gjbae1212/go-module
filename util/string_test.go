package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	assert := assert.New(t)

	a := []string{"a", "b", "c"}
	assert.True(StringInSlice("a", a), "String in Slice Not Match")

	b := []string{"a", "b", "c"}
	assert.False(StringInSlice("d", b), "String in Slice Not Match")
}

func TestSplitStringByN(t *testing.T) {
	assert := assert.New(t)
	s := "안녕하세요 저는바보 입니다."
	values := SplitStringByN(s, 2)
	assert.Len(values, 8)
}
