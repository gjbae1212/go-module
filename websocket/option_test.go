package websocket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithErrorHandlerOption(t *testing.T) {
	assert := assert.New(t)

	bk := &breaker{}
	f := WithErrorHandlerOption(func(err error) {})
	f.apply(bk)
	assert.NotNil(bk.errorHandler)
}

func TestWithMaxMessageLength(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		inputs []interface{}
		wants  []interface{}
	}{
		"step1": {inputs: []interface{}{int64(10)}, wants: []interface{}{int64(10)}},
	}
	for _, test := range tests {
		bk := &breaker{}
		f := WithMaxReadLimit(test.inputs[0].(int64))
		f.apply(bk)
		assert.Equal(test.wants[0].(int64), bk.maxReadLimit)
	}

}

func TestWithMaxMessagePoolLength(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		inputs []interface{}
		wants  []interface{}
	}{
		"step1": {inputs: []interface{}{int64(10)}, wants: []interface{}{int(10)}},
	}
	for _, test := range tests {
		bk := &breaker{}
		f := WithMaxMessagePoolLength(test.inputs[0].(int64))
		f.apply(bk)
		assert.Equal(test.wants[0].(int), cap(bk.broadcast))
	}
}
