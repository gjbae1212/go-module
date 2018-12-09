package gcp_dialogflow

import (
	"context"

	"fmt"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"github.com/gjbae1212/go-module/gcp"
	"github.com/golang/protobuf/ptypes/struct"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

var (
	defaultLanguage = "ko"
	defaultTimezone = "Asia/Tokyo"
)

type DialogFlow struct {
	projectId string
	client    *dialogflow.SessionsClient
}

type Request struct {
	Message               string // required
	Username              string // required
	Language              string // optional <https://dialogflow.com/docs/languages>
	Timezone              string // optional <https://www.iana.org/time-zones>
	IntentContext         string // optional context
	IntentContextDeadline int32  // optional context lifespan count
	ResetContext          bool   // optional whether delete all contexts in the current session
}

type Response struct {
	Intent     string
	Confidence float32
	Entities   map[string]*structpb.Value
	Response   string
}

func (req *Request) Validate() bool {
	if req.Message == "" {
		return false
	}
	if req.Username == "" {
		return false
	}
	if req.Language == "" {
		req.Language = defaultLanguage
	}
	if req.Timezone == "" {
		req.Timezone = defaultTimezone
	}
	return true
}

func (req *Request) Session(projectId string) string {
	return fmt.Sprintf("projects/%s/agent/sessions/%s", projectId, req.Username)
}

func (req *Request) Contexts(projectId string) []*dialogflowpb.Context {
	var contexts []*dialogflowpb.Context
	if req.IntentContext != "" {
		name := fmt.Sprintf("projects/%s/agent/sessions/%s/contexts/%s", projectId, req.Username, req.IntentContext)
		contexts = append(contexts, &dialogflowpb.Context{Name: name, LifespanCount: req.IntentContextDeadline})
	}
	return contexts
}

func (res *Response) Set(result *dialogflowpb.DetectIntentResponse) error {
	if result == nil {
		return gcp.EmptyError.New("Response Set")
	}

	query := result.GetQueryResult()
	if query.Intent != nil {
		res.Intent = query.Intent.DisplayName
		res.Confidence = float32(query.IntentDetectionConfidence)
	}
	res.Entities = make(map[string]*structpb.Value)
	for name, value := range query.Parameters.GetFields() {
		res.Entities[name] = value
	}
	res.Response = result.QueryResult.GetFulfillmentText()
	return nil
}

func (res *Response) GetEntity(name string) interface{} {
	v, ok := res.Entities[name]
	if !ok {
		return nil
	}

	switch v.GetKind().(type) {
	case *structpb.Value_StringValue:
		return v.GetStringValue()
	case *structpb.Value_NumberValue:
		return v.GetNumberValue()
	case *structpb.Value_BoolValue:
		return v.GetBoolValue()
	case *structpb.Value_StructValue:
		return v.GetStructValue()
	case *structpb.Value_ListValue:
		return v.GetListValue()
	case *structpb.Value_NullValue:
		return v.GetNullValue()
	default:
		return nil
	}
}

func (df *DialogFlow) DetectIntent(ctx context.Context, req *Request) (*Response, error) {
	if !req.Validate() {
		return nil, gcp.InvalidError.New("DialogFlow DetectIntent")
	}

	request := &dialogflowpb.DetectIntentRequest{
		Session: req.Session(df.projectId),
		QueryInput: &dialogflowpb.QueryInput{
			Input: &dialogflowpb.QueryInput_Text{
				Text: &dialogflowpb.TextInput{
					Text:         req.Message,
					LanguageCode: req.Language,
				},
			},
		},
		QueryParams: &dialogflowpb.QueryParameters{
			TimeZone:      req.Timezone,
			ResetContexts: req.ResetContext,
			Contexts:      req.Contexts(df.projectId),
		},
	}

	response, err := df.client.DetectIntent(ctx, request)
	if err != nil {
		return nil, err
	}

	value := &Response{}
	if err := value.Set(response); err != nil {
		return nil, err
	}

	return value, nil
}
