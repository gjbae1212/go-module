package sentry

import (
	"fmt"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/stretchr/testify/assert"
)

func TestInitSentry(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		inputSentryDSN   string
		inputEnvironment string
		inputRelease     string
		inputHostname    string
		inputStack       bool
		inputDebug       bool
		outputError      bool
	}{
		"empty": {outputError: true},
		"success": {inputSentryDSN: "https://aabb@blahblah.com/1111", inputEnvironment: "local", inputRelease: "local",
			inputHostname: "local"},
	}

	for _, t := range tests {
		err := InitSentry(t.inputSentryDSN, t.inputEnvironment, t.inputRelease, t.inputHostname, false, false)
		assert.Equal(t.outputError, err != nil)
	}
}

func TestError(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		inputErr error
	}{
		"success": {inputErr: fmt.Errorf("[err] test")},
	}

	for _, t := range tests {
		Error(t.inputErr)
	}
	_ = assert
}

func TestErrorWithEcho(t *testing.T) {
	assert := assert.New(t)

	e := echo.New()

	tests := map[string]struct {
		inputErr     error
		inputContext echo.Context
	}{
		"success": {inputErr: fmt.Errorf("[err] test"), inputContext: e.NewContext(nil, nil)},
	}

	for _, t := range tests {
		ErrorWithEcho(t.inputErr, t.inputContext, nil)
	}
	_ = assert
}
