package gcp_dialogflow

import (
	"bytes"

	"github.com/gjbae1212/go-module/gcp"
	"github.com/golang/protobuf/jsonpb"
	dialogflow "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// https://dialogflow.com/docs/fulfillment/how-it-works
func JsonToWebhookRequest(json []byte) (*dialogflow.WebhookRequest, error) {
	if len(json) == 0 {
		return nil, gcp.EmptyError.New("webhook JsonToWebhookRequest")
	}
	wr := &dialogflow.WebhookRequest{}
	if err := jsonpb.Unmarshal(bytes.NewReader(json), wr); err != nil {
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
	wr := &dialogflow.WebhookResponse{}
	if err := jsonpb.Unmarshal(bytes.NewReader(json), wr); err != nil {
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
