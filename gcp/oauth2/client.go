package gcp_oauth2

import (
	"io/ioutil"
	"net/http"

	"encoding/json"

	"golang.org/x/oauth2"
)

var (
	EndpointUserInfo = "https://www.googleapis.com/oauth2/v3/userinfo"
)

type OAuthClient struct {
	Client *http.Client
	Token  *oauth2.Token
}

func (oc *OAuthClient) UserInfo() (map[string]interface{}, error) {
	res, err := oc.Client.Get(EndpointUserInfo)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}
