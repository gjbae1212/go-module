package feedly

import (
	"net/http"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestFeedly_GetEntryIdsOfStream(t *testing.T) {
	assert := assert.New(t)

	accessToken := os.Getenv("FEEDLY_ACCESS_TOKEN")
	refreshToken := os.Getenv("FEEDLY_REFRESH_TOKEN")
	fly := &Feedly{
		client:       http.DefaultClient,
		accessToken:  accessToken,
		refreshToken: refreshToken,
		hold:         &Threshold{},
	}

	stream, err := fly.GetEntryIdsOfStream("feed/http://platum.kr/feed",
		"", true, false, 5, 0)
	assert.NoError(err)
	spew.Dump(stream)

}

func TestFeedly_GetContentOfStream(t *testing.T) {
	assert := assert.New(t)

	accessToken := os.Getenv("FEEDLY_ACCESS_TOKEN")
	refreshToken := os.Getenv("FEEDLY_REFRESH_TOKEN")
	fly := &Feedly{
		client:       http.DefaultClient,
		accessToken:  accessToken,
		refreshToken: refreshToken,
		hold:         &Threshold{},
	}

	stream, err := fly.GetContentsOfStream("feed/http://platum.kr/feed",
		"", true, false, 1, 0)
	assert.NoError(err)
	spew.Dump(stream)

	categoryId, err := fly.CategoryId("GO")
	assert.NoError(err)

	stream, err = fly.GetContentsOfStream(categoryId,
		"", true, false, 1, 0)
	assert.NoError(err)
	spew.Dump(stream)

	savedId, err := fly.TagSavedId()
	assert.NoError(err)

	stream, err = fly.GetContentsOfStream(savedId,
		"", true, false, 2, 0)
	assert.NoError(err)
	spew.Dump(stream)
}

func TestFeedly_SearchContentsOfStream(t *testing.T) {
	assert := assert.New(t)

	accessToken := os.Getenv("FEEDLY_ACCESS_TOKEN")
	refreshToken := os.Getenv("FEEDLY_REFRESH_TOKEN")
	fly := &Feedly{
		client:       http.DefaultClient,
		accessToken:  accessToken,
		refreshToken: refreshToken,
		hold:         &Threshold{},
	}

	stream, err := fly.SearchContentsOfStream("feed/http://platum.kr/feed", "", "체인", 2, 0)
	assert.NoError(err)
	spew.Dump(stream)

	stream, err = fly.SearchContentsOfStream("feed/http://platum.kr/feed", stream.Continuation, "체인", 2, 0)
	assert.NoError(err)
	spew.Dump(stream)
}
