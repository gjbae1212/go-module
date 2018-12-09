package gcp_dialogflow

import (
	"testing"

	"io/ioutil"
	"path/filepath"

	"github.com/gjbae1212/go-module/util"
	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	assert := assert.New(t)

	// TODO: required jwt path
	jwtpath := filepath.Join(util.GetModulePath(), "asset", "gcp_jwt.json")
	jwt, err := ioutil.ReadFile(jwtpath)
	assert.NoError(err)
	m, err := NewManager(jwt)
	assert.NoError(err)
	assert.NotNil(m.(*DialogFlow).client)
	assert.NotEmpty(m.(*DialogFlow).projectId)
}
