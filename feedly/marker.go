package feedly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type MarkAction int

const (
	AsRead = iota
	AsUnRead
	AsSaved
	AsUnSaved
)

type MarkerService interface {
	GetUnreads(streamId string, autorefresh bool, newerThan int) ([]*Unread, error)
	MarkEntriesAsAction(entryIds []string, action MarkAction) error
}

func (fly *Feedly) GetUnreads(streamId string, autoRefresh bool, newerThan int) ([]*Unread, error) {
	var params []string
	if streamId != "" {
		params = append(params, "streamId="+url.PathEscape(streamId))
	}
	params = append(params, "autorefresh="+strconv.FormatBool(autoRefresh))
	if newerThan > 0 {
		params = append(params, "newThan="+fmt.Sprintf("%d", newerThan))
	}

	req, err := fly.newRequest(http.MethodGet, "/markers/counts?"+strings.Join(params, "&"),
		nil)
	if err != nil {
		return nil, err
	}
	body, err := fly.do(req)
	if err != nil {
		return nil, err
	}
	result := &Marker{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}
	return result.UnreadCounts, nil
}

func (fly *Feedly) MarkEntriesAsAction(entryIds []string, action MarkAction) error {
	if len(entryIds) == 0 {
		return emptyError.New("Feedly MarkAsRead")
	}
	params := make(map[string]interface{})
	params["type"] = "entries"
	params["entryIds"] = entryIds
	switch action {
	case AsRead:
		params["action"] = "markAsRead"
	case AsUnRead:
		params["action"] = "keepUnread"
	case AsSaved:
		params["action"] = "markAsSaved"
	case AsUnSaved:
		params["action"] = "markAsUnsaved"
	}

	req, err := fly.newRequest(http.MethodPost, "/markers", params)
	if err != nil {
		return err
	}
	if _, err = fly.do(req); err != nil {
		return err
	}
	return nil
}
