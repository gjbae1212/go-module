package redis

import (
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	assert := assert.New(t)
	addr := "127.0.0.1:6378"
	s := miniredis.NewMiniRedis()
	if err := miniredis.NewMiniRedis().StartAddr(addr); err != nil {
		panic(err)
	}
	defer s.Close()

	_, err := NewManager([]string{})
	assert.Error(err)

	_, err = NewManager([]string{addr})
	assert.NoError(err)
}
