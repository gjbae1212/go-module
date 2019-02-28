package gcp_dialogflow

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"io/ioutil"
)

func TestNewManager(t *testing.T) {
	assert := assert.New(t)
	jwtpath := os.Getenv("GCP_JWT")
	_, err := os.Stat(jwtpath)
	if os.IsNotExist(err) {
		return
	}

	jwt, err := ioutil.ReadFile(jwtpath)
	assert.NoError(err)
	m, err := NewManager(jwt)
	assert.NoError(err)
	assert.NotNil(m.(*DialogFlow).client)
	assert.NotEmpty(m.(*DialogFlow).projectId)
}
