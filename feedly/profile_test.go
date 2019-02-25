package feedly

import (
	"net/http"
	"testing"

	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestFeedly_GetProfile(t *testing.T) {
	assert := assert.New(t)

	accessToken := os.Getenv("FEEDLY_ACCESS_TOKEN")
	refreshToken := os.Getenv("FEEDLY_REFRESH_TOKEN")
	if accessToken == "" || refreshToken == "" {
		return
	}

	fly := &Feedly{
		client:       http.DefaultClient,
		accessToken:  accessToken,
		refreshToken: refreshToken,
		hold:         &Threshold{},
	}

	p, err := fly.GetProfile()
	assert.NoError(err)
	spew.Dump(p)
}
