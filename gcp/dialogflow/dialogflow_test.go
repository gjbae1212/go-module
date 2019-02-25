package gcp_dialogflow

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"fmt"

	"github.com/gjbae1212/go-module/util"
	"github.com/stretchr/testify/assert"
)

func TestDialogFlow_DetectIntent(t *testing.T) {
	assert := assert.New(t)

	jwtpath := filepath.Join(util.GetModulePath(), "asset", "gcp_jwt.json")
	_, err := os.Stat(jwtpath)
	if os.IsNotExist(err) {
		return
	}

	jwt, err := ioutil.ReadFile(jwtpath)
	assert.NoError(err)
	m, err := NewManager(jwt)
	assert.NoError(err)
	_ = m

	//req := &Request{Message: "맛집 알려줘", Username: "allan"}
	//response, err := m.(*DialogFlow).DetectIntent(context.Background(), req)
	//assert.NoError(err)
	//assert.Equal("foodmenu", response.Intent)
	//assert.Equal("", response.GetEntity("city").(string))
	//
	//req = &Request{Message: "머긴머야", Username: "allan"}
	//response, err = m.(*DialogFlow).DetectIntent(context.Background(), req)
	//assert.NoError(err)
	//assert.Equal("foodmenu", response.Intent)
	//assert.Equal("", response.GetEntity("city").(string))
	//
	//req = &Request{Message: "취소", Username: "allan"}
	//response, err = m.(*DialogFlow).DetectIntent(context.Background(), req)
	//assert.NoError(err)
	//assert.Equal("foodmenu_cancel", response.Intent)
	//assert.Nil(response.GetEntity("city"))
	//
	//req = &Request{Message: "강남", Username: "allan"}
	//response, err = m.(*DialogFlow).DetectIntent(context.Background(), req)
	//assert.NoError(err)
	//assert.Equal("fallback", response.Intent)
	//
	//req = &Request{Message: "강남 맛집 알려줘", Username: "allan"}
	//response, err = m.(*DialogFlow).DetectIntent(context.Background(), req)
	//assert.NoError(err)
	//assert.Equal("foodmenu", response.Intent)
	//assert.Equal("강남", response.GetEntity("city").(string))
	//
	//req = &Request{Message: "강남", Username: "allan"}
	//response, err = m.(*DialogFlow).DetectIntent(context.Background(), req)
	//assert.NoError(err)
	//assert.Equal("foodmenu_city", response.Intent)
	//assert.Equal("강남", response.GetEntity("city").(string))
	//
	//req = &Request{Message: "분당 맛집 알려줘", Username: "allan"}
	//response, err = m.(*DialogFlow).DetectIntent(context.Background(), req)
	//assert.NoError(err)
	//assert.Equal("foodmenu", response.Intent)
	//assert.Equal("분당", response.GetEntity("city").(string))
	//
	//// context reset 하기
	//req = &Request{Message: "context reset", Username: "allan", ResetContext: true}
	//response, err = m.(*DialogFlow).DetectIntent(context.Background(), req)
	//assert.Equal("context_reset", response.Intent)
	//assert.NoError(err)
	//
	//req = &Request{Message: "분당 맛집 알려줘", Username: "allan"}
	//response, err = m.(*DialogFlow).DetectIntent(context.Background(), req)
	//assert.NoError(err)
	//assert.Equal("foodmenu", response.Intent)
	//assert.Equal("분당", response.GetEntity("city").(string))
	//
	//// context reset 하기
	//req = &Request{Message: "Context Reset", Username: "allan", IntentContext: "foodmenu", IntentContextDeadline: 0}
	//response, err = m.(*DialogFlow).DetectIntent(context.Background(), req)
	//assert.NoError(err)
	//assert.Equal("context_reset", response.Intent)
	//assert.NoError(err)
}

func TestRequest_Session(t *testing.T) {
	assert := assert.New(t)
	req := &Request{Message: "분당 맛집 알려줘", Username: "allan"}
	assert.Equal(fmt.Sprintf("projects/%s/agent/sessions/%s", "aaa", req.Username), req.Session("aaa"))
}

func TestRequest_Validate(t *testing.T) {
	assert := assert.New(t)
	req := &Request{Message: "분당 맛집 알려줘", Username: "allan"}
	assert.True(req.Validate())
	assert.Equal(req.Language, defaultLanguage)
	assert.Equal(req.Timezone, defaultTimezone)

	req = &Request{Message: "분당 맛집 알려줘", Username: "allan", Timezone: "default", Language: "default"}
	assert.True(req.Validate())
	assert.Equal(req.Language, "default")
	assert.Equal(req.Timezone, "default")
}

func TestResponse_Set(t *testing.T) {
	assert := assert.New(t)
	res := &Response{}
	assert.Error(res.Set(nil))
}

func TestResponse_GetEntity(t *testing.T) {
	assert := assert.New(t)
	res := &Response{}
	assert.Equal(res.GetEntity("allan"), nil)
}
