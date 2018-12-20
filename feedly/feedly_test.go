package feedly

import (
	"net/http"
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
)

func TestFeedly_Request(t *testing.T) {
	assert := assert.New(t)

	accessToken := os.Getenv("FEEDLY_ACCESS_TOKEN")
	refreshToken := os.Getenv("FEEDLY_REFRESH_TOKEN")
	fly := &Feedly{
		client:       http.DefaultClient,
		accessToken:  accessToken,
		refreshToken: refreshToken,
		hold:         &Threshold{},
	}
	req, err := fly.newRequest(http.MethodGet, "/profile", nil)
	assert.NoError(err)
	_, err = fly.do(req)
	assert.NoError(err)
	assert.NotEqual(fly.hold.resetCount, 0)
	assert.NotEqual(fly.hold.callCount, 0)
	assert.NotEqual(fly.hold.limitCount, 0)

	err = fly.renewAccessToken(accessToken)
	assert.NoError(err)
	assert.NotEqual(accessToken, fly.accessToken)
}
