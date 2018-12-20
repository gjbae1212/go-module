package feedly

import (
	"encoding/json"
	"net/http"
)

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
