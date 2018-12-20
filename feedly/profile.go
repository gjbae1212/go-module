package feedly

import (
	"encoding/json"
	"net/http"
)

type Profile struct {
	Id                          string           `json:"id,omitempty"`
	Email                       string           `json:"email,omitempty"`
	GivenName                   string           `json:"givenName,omitempty"`
	FamilyName                  string           `json:"familyName,omitempty"`
	FullName                    string           `json:"fullName,omitempty"`
	Picture                     string           `json:"picture,omitempty"`
	Gender                      string           `json:"gender,omitempty"`
	Locale                      string           `json:"locale,omitempty"`
	Google                      string           `json:"google,omitempty"`
	Reader                      string           `json:"reader,omitempty"`
	Wave                        string           `json:"wave,omitempty"`
	Client                      string           `json:"client,omitempty"`
	Source                      string           `json:"source,omitempty"`
	Created                     int              `json:"created,omitempty"`
	Product                     string           `json:"product,omitempty"`
	ProductExpiration           int              `json:"productExpiration,omitempty"`
	ProductRenewalAmount        int              `json:"productRenewalAmount,omitempty"`
	UpgradeDate                 int              `json:"upgradeDate,omitempty"`
	SubscriptionRenewalDate     int              `json:"subscriptionRenewalDate,omitempty"`
	SubscriptionPaymentProvider string           `json:"subscriptionPaymentProvider,omitempty"`
	SubscriptionStatus          string           `json:"subscriptionStatus,omitempty"`
	EvernoteStoreUrl            string           `json:"evernoteStoreUrl,omitempty"`
	EvernoteWebApiPrefix        string           `json:"evernoteWebApiPrefix,omitempty"`
	EvernotePartialOAuth        bool             `json:"evernotePartialOAuth,omitempty"`
	RefPage                     string           `json:"refPage,omitempty"`
	LandingPage                 string           `json:"landingPage,omitempty"`
	LoginProviders              []*LoginProvider `json:"logins,omitempty"`
	CardDetails                 *CardDetails     `json:"cardDetails,omitempty"`
	TwitterUserId               string           `json:"twitterUserId,omitempty"`
	FacebookUserId              string           `json:"facebookUserId,omitempty"`
	WordPressId                 string           `json:"wordPressId,omitempty"`
	WindowsLiveId               string           `json:"windowsLiveId,omitempty"`
	EvernoteUserId              string           `json:"evernoteUserId,omitempty"`
	EvernoteConnected           bool             `json:"evernoteConnected,omitempty"`
	PocketConnected             bool             `json:"pocketConnected,omitempty"`
	DropboxConnected            bool             `json:"dropboxConnected,omitempty"`
	TwitterConnected            bool             `json:"twitterConnected,omitempty"`
	FacebookConnected           bool             `json:"facebookConnected,omitempty"`
	WordPressConnected          bool             `json:"wordPressConnected,omitempty"`
	WindowsLiveConnected        bool             `json:"windowsLiveConnected,omitempty"`
	InstapaperConnected         bool             `json:"instapaperConnected,omitempty"`
}

type LoginProvider struct {
	Id         string `json:"id,omitempty"`
	Verified   bool   `json:"verified,omitempty"`
	Picture    string `json:"picture,omitempty"`
	Provider   string `json:"provider,omitempty"`
	ProviderId string `json:"providerId,omitempty"`
	FullName   string `json:"fullName,omitempty"`
}

type CardDetails struct {
	Brand           string `json:"brand,omitempty"`
	ExpirationMonth int    `json:"expirationMonth,omitempty"`
	ExpirationYear  int    `json:"expirationYear,omitempty"`
	Last4           string `json:"last4,omitempty"`
	Country         string `json:"country,omitempty"`
}

type ProfileService interface {
	GetProfile() (*Profile, error)
}

func (fly *Feedly) GetProfile() (*Profile, error) {
	req, err := fly.newRequest(http.MethodGet, "/profile", nil)
	if err != nil {
		return nil, err
	}
	body, err := fly.do(req)
	if err != nil {
		return nil, err
	}
	result := &Profile{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}
	return result, nil
}
