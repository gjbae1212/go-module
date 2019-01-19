package actions_on_google

import (
	"github.com/gjbae1212/go-module/util"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type Context struct {
	Name          string
	LifespanCount int
	Params        map[string]interface{}
}

func (u *Context) IsEmpty() bool {
	return u.Name == ""
}

func GetContextsFromRequest(req *dialogflowpb.WebhookRequest) ([]*Context, error) {
	var contexts []*Context
	if req.QueryResult.OutputContexts == nil || len(req.QueryResult.OutputContexts) == 0 {
		return contexts, nil
	}

	for _, opt := range req.QueryResult.OutputContexts {
		ctx := &Context{}
		ctx.Name = opt.Name
		ctx.LifespanCount = int(opt.LifespanCount)
		m, err := util.StructPBToMap(opt.Parameters)
		if err != nil {
			return nil, err
		}

		ctx.Params = m
		contexts = append(contexts, ctx)
	}

	return contexts, nil
}
