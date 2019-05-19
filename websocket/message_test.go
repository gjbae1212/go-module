package websocket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInternalMessage_GetMessage(t *testing.T) {
	assert := assert.New(t)
	msg := &internalMessage{payload: []byte("hello")}
	assert.Equal("hello", string(msg.GetMessage()))
}
