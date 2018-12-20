package feedly

import (
	"encoding/json"
	"net/http"
)

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
