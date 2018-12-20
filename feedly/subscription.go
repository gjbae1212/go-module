package feedly

import (
	"encoding/json"
	"net/http"
)

type Subscription struct {
	Id          string      `json:"id,omitempty"`
	FeedId      string      `json:"feedId,omitempty"`
	SortId      string      `json:"sortid,omitempty"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	ContentType string      `json:"contentType,omitempty"`
	Language    string      `json:"language,omitempty"`
	Website     string      `json:"website,omitempty"`
	IconUrl     string      `json:"iconUrl,omitempty"`
	CoverUrl    string      `json:"coverUrl,omitempty"`
	VisualUrl   string      `json:"visualUrl,omitempty"`
	CoverColor  string      `json:"coverColor,omitempty"`
	Subscribers int         `json:"subscribers,omitempty"`
	Added       int         `json:"added,omitempty"`
	Updated     int         `json:"updated,omitempty"`
	Velocity    float64     `json:"velocity,omitempty"`
	Partial     bool        `json:"partial,omitempty"`
	Categories  []*Category `json:"categories,omitempty"`
	Topics      []string    `json:"topics,omitempty"`
}

type SubscriptionService interface {
	GetSubscriptions() ([]*Subscription, error)
	Subscribe(feedId string, label string, categoriesIds []string) (*Subscription, error)
}

func (fly *Feedly) GetSubscriptions() ([]*Subscription, error) {
	req, err := fly.newRequest(http.MethodGet, "/subscriptions", nil)
	if err != nil {
		return nil, err
	}
	body, err := fly.do(req)
	if err != nil {
		return nil, err
	}
	var result []*Subscription
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fly *Feedly) Subscribe(feedId string, label string, categoriesIds []string) (*Subscription, error) {
	if feedId == "" {
		return nil, emptyError.New("Feedly Subscribe")
	}
	params := make(map[string]interface{})
	params["id"] = feedId
	params["label"] = label
	if len(categoriesIds) > 0 {
		params["categories"] = []interface{}{}
		for _, id := range categoriesIds {
			params["categories"] = append(params["categories"].([]interface{}), map[string]string{"id": id})
		}
	}

	req, err := fly.newRequest(http.MethodPost, "/subscriptions", params)
	if err != nil {
		return nil, err
	}
	body, err := fly.do(req)
	if err != nil {
		return nil, err
	}
	var result []*Subscription
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, nil
}
