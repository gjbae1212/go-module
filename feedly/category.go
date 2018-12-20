package feedly

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

type CategoryService interface {
	GetCategories(withStats bool) ([]*Category, error)
	GetCategory(categoryId string) (*Category, error)
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
