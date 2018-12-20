package feedly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type FeedService interface {
	GetFeed(feedId string) (*Feed, error)
	SearchFeeds(query, locale string, count int) ([]*Feed, error)
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

func (fly *Feedly) SearchFeeds(query, locale string, count int) ([]*Feed, error) {
	if query == "" {
		return nil, emptyError.New("Feedly GetSearchFeed")
	}

	var params []string
	params = append(params, "query="+url.PathEscape(query))
	if locale != "" {
		params = append(params, "locale="+locale)
	}
	if count < 0 {
		count = 20
	}
	params = append(params, "count="+fmt.Sprintf("%d", count))

	req, err := fly.newRequest(http.MethodGet, "/search/feeds?"+strings.Join(params, "&"), nil)
	if err != nil {
		return nil, err
	}
	body, err := fly.do(req)
	if err != nil {
		return nil, err
	}
	result := &SearchResult{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}
	return result.Feeds, nil
}
