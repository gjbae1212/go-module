package gcp_oauth2

import (
	"context"
	"encoding/base64"

	"strings"

	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
)

type config struct {
	*oauth2.Config
}

func newConfig(conf *oauth2.Config) *config {
	return &config{Config: conf}
}

func (c *config) ExchangeWithCode(ctx context.Context, code string) (client *OAuthClient, err error) {
	token, suberr := c.Exchange(ctx, code)
	if suberr != nil {
		err = suberr
		return
	}
	client = &OAuthClient{Client: c.Client(ctx, token), Token: token}
	return
}

func (c *config) GetLoginUrlWithState(state string) string {
	return c.AuthCodeURL(state)
}

func (c *config) GetLoginUrl() (state string, url string) {
	state = c.GetRandomState()
	url = c.AuthCodeURL(state)
	return
}

func (c *config) GetRandomState() string {
	// generate random string
	uid := uuid.NewV4()
	ba64 := base64.StdEncoding.EncodeToString(uid.Bytes())
	// remove base64 characters(/, +)
	return strings.Replace(strings.Replace(ba64, "/", "", -1), "+", "", -1)
}
