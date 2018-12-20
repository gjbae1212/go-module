package feedly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type StreamService interface {
	GetEntryIdsOfStream(streamId, continuation string, orderDesc, unreadOnly bool,
		count, newerThan int64) (*Stream, error)
	GetContentsOfStream(streamId, continuation string, orderDesc, unreadOnly bool,
		count, newerThan int) (*Stream, error)
	SearchContentsOfStream(streamId, continuation, query string, count, newerThan int) (*Stream, error)
}

// A feedId, a categoryId, a tagId or a system category ids can be used as stream ids.
func (fly *Feedly) GetEntryIdsOfStream(streamId, continuation string, orderDesc, unreadOnly bool,
	count, newerThan int64) (*Stream, error) {
	if streamId == "" {
		return nil, emptyError.New("Feedly GetEntryIdsOfStream")
	}
	ranked := "oldest"
	if orderDesc {
		ranked = "newest"
	}
	if count < 0 {
		count = 20
	} else if count > 10000 {
		count = 10000
	}
	var params []string
	params = append(params, "ranked="+ranked)
	params = append(params, "unreadOnly="+strconv.FormatBool(unreadOnly))
	params = append(params, "count="+fmt.Sprintf("%d", count))
	if newerThan > 0 {
		params = append(params, "newThan="+fmt.Sprintf("%d", newerThan))
	}
	if continuation != "" {
		params = append(params, "continuation="+continuation)
	}

	req, err := fly.newRequest(http.MethodGet, "/streams/ids?streamId="+url.PathEscape(streamId)+"&"+strings.Join(params, "&"),
		nil)
	if err != nil {
		return nil, err
	}
	body, err := fly.do(req)
	if err != nil {
		return nil, err
	}
	result := &Stream{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fly *Feedly) GetContentsOfStream(streamId, continuation string, orderDesc, unreadOnly bool,
	count, newerThan int) (*Stream, error) {
	if streamId == "" {
		return nil, emptyError.New("Feedly GetContentOfStream")
	}
	ranked := "oldest"
	if orderDesc {
		ranked = "newest"
	}
	if count < 0 {
		count = 20
	} else if count > 10000 {
		count = 10000
	}
	var params []string
	params = append(params, "ranked="+ranked)
	params = append(params, "unreadOnly="+strconv.FormatBool(unreadOnly))
	params = append(params, "count="+fmt.Sprintf("%d", count))
	if newerThan > 0 {
		params = append(params, "newThan="+fmt.Sprintf("%d", newerThan))
	}
	if continuation != "" {
		params = append(params, "continuation="+continuation)
	}

	req, err := fly.newRequest(http.MethodGet, "/streams/contents?streamId="+url.PathEscape(streamId)+"&"+strings.Join(params, "&"),
		nil)
	if err != nil {
		return nil, err
	}
	body, err := fly.do(req)
	if err != nil {
		return nil, err
	}
	result := &Stream{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (fly *Feedly) SearchContentsOfStream(streamId, continuation, query string, count, newerThan int) (*Stream, error) {
	if streamId == "" || query == "" {
		return nil, emptyError.New("Feedly SearchContentsOfStream")
	}

	var params []string
	params = append(params, "query="+url.PathEscape(query))
	if count < 0 {
		count = 10
	} else if count > 20 {
		count = 20
	}
	params = append(params, "count="+fmt.Sprintf("%d", count))
	if continuation != "" {
		params = append(params, "continuation="+continuation)
	}
	if newerThan > 0 {
		params = append(params, "newThan="+fmt.Sprintf("%d", newerThan))
	}

	req, err := fly.newRequest(http.MethodGet, "/search/contents?streamId="+url.PathEscape(streamId)+"&"+strings.Join(params, "&"),
		nil)
	if err != nil {
		return nil, err
	}
	body, err := fly.do(req)
	if err != nil {
		return nil, err
	}
	result := &Stream{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}
