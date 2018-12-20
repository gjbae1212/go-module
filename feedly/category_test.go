package feedly

import (
	"testing"

	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestFeedly_GetCategories(t *testing.T) {
	assert := assert.New(t)

	accessToken := os.Getenv("FEEDLY_ACCESS_TOKEN")
	refreshToken := os.Getenv("FEEDLY_REFRESH_TOKEN")
	fly := &Feedly{
		client:       http.DefaultClient,
		accessToken:  accessToken,
		refreshToken: refreshToken,
		hold:         &Threshold{},
	}

	ca, err := fly.GetCategories(false)
	assert.NoError(err)
	spew.Dump(ca)
}

func TestFeedly_GetCategory(t *testing.T) {
	assert := assert.New(t)

	accessToken := os.Getenv("FEEDLY_ACCESS_TOKEN")
	refreshToken := os.Getenv("FEEDLY_REFRESH_TOKEN")
	fly := &Feedly{
		client:       http.DefaultClient,
		accessToken:  accessToken,
		refreshToken: refreshToken,
		hold:         &Threshold{},
	}

	categories, err := fly.GetCategories(false)
	assert.NoError(err)
	if len(categories) > 0 {
		ca, err := fly.GetCategory(categories[0].Id)
		assert.NoError(err)
		spew.Dump(ca)
	}
}
