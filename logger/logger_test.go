package logger

import (
	"path"
	"path/filepath"
	"testing"

	"runtime"

	"os"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger(t *testing.T) {
	assert := assert.New(t)

	logger, err := NewLogger("", "")
	assert.NoError(err)
	logger.Print("unit test step 1")
	_, filename, _, ok := runtime.Caller(0)
	assert.True(ok)
	_, err = NewLogger(filepath.Join(path.Dir(filename), "../", "empty"), "test.log")
	assert.Error(err)
	logPath := filepath.Join(path.Dir(filename), "../", "log")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		os.Mkdir(logPath, os.ModePerm)
	}
	logger, err = NewLogger(logPath, "test.log")
	assert.NoError(err)
	logger.Printf("unit test step 2")
}
