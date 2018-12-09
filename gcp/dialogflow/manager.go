package gcp_dialogflow

import (
	"context"

	"encoding/json"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"github.com/gjbae1212/go-module/gcp"
	"google.golang.org/api/option"
)

type Manager interface {
	DetectIntent(ctx context.Context, req *Request) (*Response, error)
}

func NewManager(gcpjwt []byte) (Manager, error) {
	if len(gcpjwt) == 0 {
		return nil, gcp.EmptyError.New("dialogflow NewManager")
	}

	var jwtJson map[string]interface{}
	if err := json.Unmarshal(gcpjwt, &jwtJson); err != nil {
		return nil, err
	}

	client, err := dialogflow.NewSessionsClient(context.Background(), option.WithCredentialsJSON(gcpjwt))
	if err != nil {
		return nil, err
	}

	df := &DialogFlow{
		projectId: jwtJson["project_id"].(string),
		client:    client,
	}

	return Manager(df), nil
}
