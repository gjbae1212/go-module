package gcp_oauth2

import (
	"testing"

	"net/http"

	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{Transport: RoundTripFunc(fn)}
}

func TestOAuthClient_UserInfo(t *testing.T) {
	assert := assert.New(t)
	// http mock
	hc := NewTestClient(func(req *http.Request) *http.Response {
		data, _ := json.Marshal(map[string]string{"success": "ok"})
		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBuffer(data)),
		}
	})

	client := &OAuthClient{Client: hc}
	result, err := client.UserInfo()
	assert.NoError(err)
	assert.Equal(result["success"].(string), "ok")
}
