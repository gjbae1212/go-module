package feedly

import (
	"testing"

	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestFeedly_GetFeed(t *testing.T) {
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

	feed, err := fly.GetFeed("feed/http://feeds.feedburner.com/TechCrunch/startups")
	assert.NoError(err)
	spew.Dump(feed)

}

func TestFeedly_SearchFeeds(t *testing.T) {
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

	feed, err := fly.SearchFeeds("golang", "en", 5)
	assert.NoError(err)
	spew.Dump(feed)

}
