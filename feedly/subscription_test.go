package feedly

import (
	"testing"

	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestFeedly_GetSubscriptions(t *testing.T) {
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

	ca, err := fly.GetSubscriptions()
	assert.NoError(err)
	spew.Dump(ca)
}

func TestFeedly_Subscribe(t *testing.T) {
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
			sub, err := fly.Subscribe("feed/http://www.engadget.com/rss-full.xml", category.Label, []string{category.Id})
			assert.NoError(err)
			assert.NotEmpty(sub)
			spew.Dump(sub)
		}
	}
}
