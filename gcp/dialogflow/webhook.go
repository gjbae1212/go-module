package gcp_dialogflow

import (
	"bytes"

	"regexp"

	"fmt"

	"github.com/gjbae1212/go-module/gcp"
	"github.com/golang/protobuf/jsonpb"
	dialogflow "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

var (
	sessionRegex = regexp.MustCompile("projects/(.+)/agent/sessions/(.+)")
	contextRegex = regexp.MustCompile("projects/(.+)/agent/sessions/(.+)/contexts/(.+)")
	intentRegex  = regexp.MustCompile("projects/(.+)/agent/intents/(.+)")
)

// https://dialogflow.com/docs/fulfillment/how-it-works
// https://developers.google.com/actions/build/json/

func GenerateSession(projectId, sessionId string) (string, error) {
	if projectId == "" || sessionId == "" {
		return "", gcp.EmptyError.New("webhook GenerateSession")
	}
	return fmt.Sprintf("projects/%s/agent/sessions/%s", projectId, sessionId), nil
}

func GenerateContext(projectId, sessionId, contextId string) (string, error) {
	if projectId == "" || sessionId == "" || contextId == "" {
		return "", gcp.EmptyError.New("webhook GenerateContext")
	}
	return fmt.Sprintf("projects/%s/agent/sessions/%s/contexts/%s", projectId, sessionId, contextId), nil
}

func GenerateIntent(projectId, intentId string) (string, error) {
	if projectId == "" || intentId == "" {
		return "", gcp.EmptyError.New("webhook GenerateIntent")
	}
	return fmt.Sprintf("projects/%s/agent/intents/%s", projectId, intentId), nil
}

func ParseSession(session string) (projectId, sessionId string, err error) {
	if session == "" {
		err = gcp.EmptyError.New("webhook ParseSession")
		return
	}

	matches := sessionRegex.FindStringSubmatch(session)
	if len(matches) == 0 || len(matches) != 3 {
		err = gcp.InvalidError.New("webhook ParseSession invalid format")
		return
	}

	projectId = matches[1]
	sessionId = matches[2]
	return
}

func ParseContext(context string) (projectId, sessionId string, contextId string, err error) {
	if context == "" {
		err = gcp.EmptyError.New("webhook ParseContext")
		return
	}

	matches := contextRegex.FindStringSubmatch(context)
	if len(matches) == 0 || len(matches) != 4 {
		err = gcp.InvalidError.New("webhook ParseContext invalid format")
		return
	}

	projectId = matches[1]
	sessionId = matches[2]
	contextId = matches[3]
	return
}

func ParseIntent(intent string) (projectId, intentId string, err error) {
	if intent == "" {
		err = gcp.EmptyError.New("webhook ParseIntent")
		return
	}

	matches := intentRegex.FindStringSubmatch(intent)
	if len(matches) == 0 || len(matches) != 3 {
		err = gcp.InvalidError.New("webhook ParseIntent invalid format")
		return
	}

	projectId = matches[1]
	intentId = matches[2]
	return
}

func JsonToWebhookRequest(json []byte) (*dialogflow.WebhookRequest, error) {
	if len(json) == 0 {
		return nil, gcp.EmptyError.New("webhook JsonToWebhookRequest")
	}

	unmarshaler := &jsonpb.Unmarshaler{}
	unmarshaler.AllowUnknownFields = true
	wr := &dialogflow.WebhookRequest{}
	if err := unmarshaler.Unmarshal(bytes.NewReader(json), wr); err != nil {
		return nil, err
	}
	return wr, nil
}

func WebhookRequestToJson(req *dialogflow.WebhookRequest) ([]byte, error) {
	if req == nil {
		return nil, gcp.EmptyError.New("webhook WebhookRequestToJson")
	}

	marshaler := &jsonpb.Marshaler{}
	json, err := marshaler.MarshalToString(req)
	if err != nil {
		return nil, err
	}
	return []byte(json), err
}

func JsonToWebhookResponse(json []byte) (*dialogflow.WebhookResponse, error) {
	if len(json) == 0 {
		return nil, gcp.EmptyError.New("webhook JsonToWebhookResponse")
	}
	unmarshaler := &jsonpb.Unmarshaler{}
	unmarshaler.AllowUnknownFields = true
	wr := &dialogflow.WebhookResponse{}
	if err := unmarshaler.Unmarshal(bytes.NewReader(json), wr); err != nil {
		return nil, err
	}
	return wr, nil
}

func WebhookResponseToJson(res *dialogflow.WebhookResponse) ([]byte, error) {
	if res == nil {
		return nil, gcp.EmptyError.New("webhook WebhookResponseToJson")
	}
	marshaler := &jsonpb.Marshaler{}
	json, err := marshaler.MarshalToString(res)
	if err != nil {
		return nil, err
	}
	return []byte(json), err
}