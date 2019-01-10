package actions_on_google

import (
	"testing"

	"github.com/golang/protobuf/ptypes/struct"
	"github.com/stretchr/testify/assert"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

func TestGetUserByOriginalDetectIntentRequest(t *testing.T) {
	assert := assert.New(t)
	req := &dialogflowpb.OriginalDetectIntentRequest{}

	user := GetUserByOriginalDetectIntentRequest(req)
	assert.True(user.IsEmpty())

	req.Payload = &structpb.Struct{}
	user = GetUserByOriginalDetectIntentRequest(req)
	assert.True(user.IsEmpty())

	inner := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"userId":      &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "allan"}},
			"userStorage": &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "data"}},
		},
	}

	payload := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"user": &structpb.Value{Kind: &structpb.Value_StructValue{
				StructValue: inner,
			}},
		},
	}

	req.Payload = payload
	user = GetUserByOriginalDetectIntentRequest(req)
	assert.False(user.IsEmpty())
	assert.Equal(user.UserId, "allan")
	assert.Equal(user.UserStorage, "data")
}
