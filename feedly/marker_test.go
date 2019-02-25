package feedly

import (
	"testing"

	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestFeedly_GetUnreads(t *testing.T) {
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

	unread, err := fly.GetUnreads("", false, 0)
	assert.NoError(err)
	spew.Dump(unread)

}

func TestFeedly_MarkEntriesAsAction(t *testing.T) {
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

	categories, err := fly.GetCategories(false)
	assert.NoError(err)
	for _, category := range categories {
		if category.Label == "TEST" {
			stream, err := fly.GetEntryIdsOfStream(category.Id, "", true, false, 20, 0)
			assert.NoError(err)
			err = fly.MarkEntriesAsAction(stream.Ids, AsRead)
			assert.NoError(err)
			err = fly.MarkEntriesAsAction(stream.Ids, AsUnRead)
			assert.NoError(err)
			err = fly.MarkEntriesAsAction(stream.Ids, AsSaved)
			assert.NoError(err)
			err = fly.MarkEntriesAsAction(stream.Ids, AsUnSaved)
			assert.NoError(err)
		}
	}

}
