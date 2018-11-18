package gcp_oauth2

import (
	"context"

	"github.com/gjbae1212/go-module/gcp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	ScopeEmail       = "https://www.googleapis.com/auth/userinfo.email"
	ScopeProfile     = "https://www.googleapis.com/auth/userinfo.profile"
)

type Manager interface {
	ExchangeWithCode(ctx context.Context, code string) (client *OAuthClient, err error)
	GetLoginUrlWithState(state string) (url string)
	GetLoginUrl() (state string, url string)
	GetRandomState() string
}

func NewManager(oauthId, oauthSecret, redirectURL string, scopes []string) (Manager, error) {
	if oauthId == "" || oauthSecret == "" || redirectURL == "" || len(scopes) == 0 {
		return nil, gcp.EmptyError.New("oauth2 NewManager")
	}
	oauth := &oauth2.Config{
		ClientID:     oauthId,
		ClientSecret: oauthSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}
	config := newConfig(oauth)
	return Manager(config), nil
}
