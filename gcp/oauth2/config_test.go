package gcp_oauth2

import (
	"context"
	"testing"

	"encoding/base64"
	"strings"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type mockconfig struct {
	*oauth2.Config
}

func (c *mockconfig) ExchangeWithCode(ctx context.Context, code string) (client *OAuthClient, err error) {
	client = &OAuthClient{Client: c.Client(ctx, &oauth2.Token{}), Token: &oauth2.Token{}}
	return
}

func (c *mockconfig) GetLoginUrlWithState(state string) string {
	return c.AuthCodeURL(state)
}

func (c *mockconfig) GetLoginUrl() (state string, url string) {
	state = c.GetRandomState()
	url = c.AuthCodeURL(state)
	return
}

func (c *mockconfig) GetRandomState() string {
	// generate random string
	uid := uuid.NewV4()
	ba64 := base64.StdEncoding.EncodeToString(uid.Bytes())
	// remove base64 characters(/, +)
	return strings.Replace(strings.Replace(ba64, "/", "", -1), "+", "", -1)
}

func TestConfig_ExchangeWithCode(t *testing.T) {
	assert := assert.New(t)
	m := mockManager()
	_, err := m.ExchangeWithCode(context.Background(), "")
	assert.NoError(err)
}

func TestConfig_GetRandomState(t *testing.T) {
	assert := assert.New(t)
	c := mockManager()
	for i := 0; i < 100; i++ {
		assert.NotContains(c.GetRandomState(), "/")
		assert.NotContains(c.GetRandomState(), "+")
	}
}

func TestConfig_GetLoginUrl(t *testing.T) {
	assert := assert.New(t)
	m := mockManager()
	state, url := m.GetLoginUrl()
	assert.NotEmpty(state)
	assert.NotEmpty(url)
}

func TestConfig_GetLoginUrlWithState(t *testing.T) {
	assert := assert.New(t)
	c := mockManager()
	state := c.GetRandomState()
	url := c.GetLoginUrlWithState(state)
	assert.NotEmpty(state)
	assert.NotEmpty(url)
}

func mockManager() Manager {
	oauth := &oauth2.Config{
		ClientID:     "id",
		ClientSecret: "secret",
		RedirectURL:  "http://localhost:8080/oauth2",
		Scopes:       []string{ScopeProfile, ScopeEmail},
		Endpoint:     google.Endpoint,
	}

	mockConfg := &mockconfig{Config: oauth}
	return Manager(mockConfg)
}
