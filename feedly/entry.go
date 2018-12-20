package feedly

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type EntryService interface {
	GetEntry(entryId string) (*Entry, error)
}

func (fly *Feedly) GetEntry(entryId string) (*Entry, error) {
	if entryId == "" {
		return nil, emptyError.New("Feedly GetEntry")
	}

	req, err := fly.newRequest(http.MethodGet, "/entries/"+url.PathEscape(entryId), nil)
	if err != nil {
		return nil, err
	}
	body, err := fly.do(req)
	if err != nil {
		return nil, err
	}

	var result []*Entry
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, nil
}
