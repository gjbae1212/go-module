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

	err := Load(sentryDSN)
	assert.NoError(err)
}

func TestRaiseSentry(t *testing.T) {
	assert := assert.New(t)
	sentryDSN := os.Getenv("SENTRY_DSN")
	if sentryDSN == "" {
		return
	}

	err := Load(sentryDSN)
	assert.NoError(err)
	RaiseSentry(fmt.Errorf("test error without request"), true)
	RaiseSentry(fmt.Errorf("test error without request"), false)
	req := httptest.NewRequest("GET", "http://localhost", nil)
	RaiseSentryWithRequest(fmt.Errorf("test error with request"), req, true)
	RaiseSentryWithRequest(fmt.Errorf("test error with request"), req, false)
	time.Sleep(3 * time.Second)
}
