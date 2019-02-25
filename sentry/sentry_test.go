package sentry

import (
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	assert := assert.New(t)
	sentryDSN := os.Getenv("SENTRY_DSN")
	if sentryDSN == "" {
		return
	}

	err := Load(sentryDSN, 10)
	assert.NoError(err)
}

func TestRaiseSentry(t *testing.T) {
	assert := assert.New(t)
	sentryDSN := os.Getenv("SENTRY_DSN")
	if sentryDSN == "" {
		return
	}

	err := Load(sentryDSN, defaultStackLength)
	assert.NoError(err)
	packet1 := MakePacket(fmt.Errorf("test error without request"), true)
	Raise(packet1, nil)

	req := httptest.NewRequest("GET", "http://localhost", nil)
	packet2 := MakePacketWithRequest(fmt.Errorf("test error with request"), req, true)
	Raise(packet2, nil)
	time.Sleep(3 * time.Second)
}
