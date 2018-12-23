package util

import (
	"testing"

	"time"

	"github.com/stretchr/testify/assert"
)

func TestBackoff(t *testing.T) {
	assert := assert.New(t)

	b := NewBackoff(2 * time.Second)
	assert.NotEmpty(b)
	_ = b.Wait()
	b.Success()
}
