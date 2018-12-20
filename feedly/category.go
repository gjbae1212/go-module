package feedly

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

type Category struct {
	Id           string          `json:"id,omitempty"`
	Label        string          `json:"label,omitempty"`
	Description  string          `json:"description,omitempty"`
	Customizable bool            `json:"customizable,omitempty"`
	Enterprise   bool            `json:"enterprise,omitempty"`
	Cover        string          `json:"cover,omitempty"`
	Created      int             `json:"created,omitempty"`
	NumFeeds     int             `json:"numFeeds,omitempty"`
	Feeds        []*Subscription `json:"feeds,omitempty"`
}

type CategoryService interface {
	GetCategories(withStats bool) ([]*Category, error)
}

func (fly *Feedly) GetCategories(withStats bool) ([]*Category, error) {
	req, err := fly.newRequest(http.MethodGet, "/collections?withStats="+strconv.FormatBool(withStats), nil)
	if err != nil {
		return nil, err
	}
	body, err := fly.do(req)
	if err != nil {
		return nil, err
	}
	var result []*Category
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fly *Feedly) GetCategory(categoryId string) (*Category, error) {
	if categoryId == "" {
		return nil, emptyError.New("Feedly GetCategory")
	}
	req, err := fly.newRequest(http.MethodGet, "/collections/"+url.PathEscape(categoryId), nil)
	if err != nil {
		return nil, err
	}
	body, err := fly.do(req)
	if err != nil {
		return nil, err
	}

	var result []*Category
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, nil
}
