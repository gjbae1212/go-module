package feedly

import (
	"net/http"
	"testing"

	"os"

	"log"

	"github.com/stretchr/testify/assert"
)

func TestFeedly_Request(t *testing.T) {
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

func TestFeedly_ID(t *testing.T) {
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
	userId, err := fly.UserId()
	assert.NoError(err)
	log.Println(userId)

	feedId, err := fly.FeedId("http://test.com")
	assert.NoError(err)
	log.Println(feedId)

	categoryId, err := fly.CategoryId("GO")
	assert.NoError(err)
	log.Println(categoryId)

	categoryId, err = fly.CategoryAllId()
	assert.NoError(err)
	log.Println(categoryId)

	tagId, err := fly.TagId("GO")
	assert.NoError(err)
	log.Println(tagId)

	tagId, err = fly.TagAllId()
	assert.NoError(err)
	log.Println(tagId)

	tagId, err = fly.TagSavedId()
	assert.NoError(err)
	log.Println(tagId)

	tagId, err = fly.TagReadId()
	assert.NoError(err)
	log.Println(tagId)
}
