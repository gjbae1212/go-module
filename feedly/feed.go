package feedly

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Feed struct {
	Id          string   `json:"id,omitempty"`
	FeedId      string   `json:"feedId,omitempty"`
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Language    string   `json:"language,omitempty"`
	Website     string   `json:"website,omitempty"`
	Topics      []string `json:"topics,omitempty"`
	Velocity    float64  `json:"velocity,omitempty"`
	Subscribers int      `json:"subscribers,omitempty"`
	State       string   `json:"state,omitempty"`
}

type FeedService interface {
	GetFeed(feedId string) (*Feed, error)
}

func (fly *Feedly) GetFeed(feedId string) (*Feed, error) {
	if feedId == "" {
		return nil, emptyError.New("Feedly GetFeed")
	}
	req, err := fly.newRequest(http.MethodGet, "/feeds/"+url.PathEscape(feedId), nil)
	if err != nil {
		return nil, err
	}
	body, err := fly.do(req)
	if err != nil {
		return nil, err
	}
	result := &Feed{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}
	return result, nil
}
