package feedly

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gjbae1212/go-module/util"
)

// https://developer.feedly.com/cloud
const (
	baseURL      = "https://cloud.feedly.com"
	version      = "v3"
	clientId     = "feedlydev"
	clientSecret = "feedlydev"
)

type Feedly struct {
	sync.RWMutex
	client       *http.Client
	hold         *Threshold // threshold
	accessToken  string     // access token
	refreshToken string     // refresh token
}

type Threshold struct {
	limitCount int // api 콜 제한
	callCount  int // api 콜 횟수
	resetCount int // api 콜 리셋까지 남은시간
}

func (fly *Feedly) newRequest(method, path string, params map[string]interface{}) (*http.Request, error) {
	if method == "" || path == "" {
		return nil, emptyError.New("Feedly newRequest")
	}

	if !util.CheckHttpMethod(method) {
		return nil, invalidError.New("Feedly newRequest")
	}

	if strings.HasPrefix(path, "/") {
		path = path[1:len(path)]
	}
	urlString := baseURL + "/" + version + "/" + path

	buf := new(bytes.Buffer)
	if params != nil && len(params) > 0 {
		data, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		buf.Write(data)
	}

	req, err := http.NewRequest(method, urlString, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", strings.Join([]string{"Bearer", fly.accessToken}, " "))
	return req, nil
}

func (fly *Feedly) do(request *http.Request) ([]byte, error) {
	res, err := fly.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	// change threshold
	if err := fly.changeThreshold(res); err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case http.StatusUnauthorized:
		// renew access token
		fly.renewAccessToken(fly.accessToken)
		return nil, unknownError.New("Feedly do status unauthorized")
	case http.StatusOK:
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return body, nil
	default:
		return nil, unknownError.New("Feedly do status not ok")
	}
}

func (fly *Feedly) renewAccessToken(prevAccessToken string) error {
	fly.Lock()
	defer fly.Unlock()
	if fly.accessToken != prevAccessToken { // already change
		return nil
	}

	urlString := baseURL + "/" + version + "/auth/token"
	params := map[string]interface{}{
		"refresh_token": fly.refreshToken,
		"client_id":     clientId,
		"client_secret": clientSecret,
		"grant_type":    "refresh_token",
	}
	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPost, urlString, buf)
	if err != nil {
		return err
	}
	res, err := fly.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return unknownError.New("Feedly renewAccessToken status not ok")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	token, ok := result["access_token"]
	if !ok {
		return unknownError.New("Feedly renewAccessToken not access token field")
	}
	fly.accessToken = token.(string)
	return nil
}

func (fly *Feedly) changeThreshold(res *http.Response) error {
	v := res.Header.Get("X-Ratelimit-Count")
	if v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		fly.hold.callCount = i
	}
	v = res.Header.Get("X-Ratelimit-Limit")
	if v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		fly.hold.limitCount = i
	}
	v = res.Header.Get("X-Ratelimit-Reset")
	if v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		fly.hold.resetCount = i
	}
	return nil
}
