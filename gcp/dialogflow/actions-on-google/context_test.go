package actions_on_google

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/gjbae1212/go-module/util"
	"github.com/stretchr/testify/assert"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

func TestContext(t *testing.T) {
	assert := assert.New(t)

	m1 := map[string]interface{}{
		"allan1": 10,
	}
	m2 := map[string]interface{}{
		"allan2": 10,
	}
	data1, err := util.MapToStructPB(m1)
	assert.NoError(err)
	data2, err := util.MapToStructPB(m2)
	assert.NoError(err)

	req := &dialogflowpb.WebhookRequest{
		QueryResult: &dialogflowpb.QueryResult{
			OutputContexts: []*dialogflowpb.Context{
				&dialogflowpb.Context{
					Name:          "a1",
					LifespanCount: 2,
					Parameters:    data1,
				},
				&dialogflowpb.Context{
					Name:          "a2",
					LifespanCount: 3,
					Parameters:    data2,
				},
			},
		},
	}

	contexts, err := GetContextsFromRequest(req)
	assert.NoError(err)
	assert.Equal(contexts[0].Name, "a1")
	assert.Equal(contexts[0].LifespanCount, 2)
	assert.Equal(int(contexts[0].Params["allan1"].(float64)), int(10))
	spew.Dump(contexts)
}
