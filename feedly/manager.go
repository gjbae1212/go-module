package feedly

import (
	"net/http"
	"time"
)

type Manager interface {
	ProfileService
	CategoryService
	FeedService
	SubscriptionService
}

func NewManager(accestoken, refreshtoken string) (Manager, error) {
	if accestoken == "" || refreshtoken == "" {
		return nil, emptyError.New("NewManager")
	}
	tr := &http.Transport{
		IdleConnTimeout: time.Duration(10) * time.Minute,
	}
	fly := &Feedly{
		client: &http.Client{Transport: tr, Timeout: time.Duration(10) * time.Second},
		hold:   &Threshold{},
	}

	return Manager(fly), nil
}
