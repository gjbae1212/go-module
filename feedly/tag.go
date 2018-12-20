package feedly

import (
	"encoding/json"
	"net/http"
)

type TagService interface {
	GetTags() ([]*Tag, error)
}

func (fly *Feedly) GetTags() ([]*Tag, error) {
	req, err := fly.newRequest(http.MethodGet, "/tags", nil)
	if err != nil {
		return nil, err
	}
	body, err := fly.do(req)
	if err != nil {
		return nil, err
	}
	var result []*Tag
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}
