package feedly

import (
	"net/http"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestFeedly_GetEntry(t *testing.T) {
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

	stream, err := fly.GetEntryIdsOfStream("feed/http://platum.kr/feed",
		"", true, false, 1, 0)
	assert.NoError(err)
	entry, err := fly.GetEntry(stream.Ids[0])
	assert.NoError(err)
	spew.Dump(entry)

}
