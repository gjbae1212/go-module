package gcp_dialogflow

import (
	"testing"

	"encoding/json"

	"strings"

	"sort"

	"github.com/stretchr/testify/assert"
	dialogflow "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

func TestGenerateAndParseSession(t *testing.T) {
	assert := assert.New(t)
	projectId := "allan"
	sessionId := "3483478324vcxe/////sdf//nfkj(*#^$*(!!@#()@&$($&T*%#$*(_)@#*$"

	_, err := GenerateSession("", "")
	assert.Error(err)

	session, err := GenerateSession(projectId, sessionId)
	assert.NoError(err)

	_, _, err = ParseSession("dsfsf34u390dsf/3praf/")
	assert.Error(err)

	pId, sId, err := ParseSession(session)
	assert.NoError(err)
	assert.Equal(pId, projectId)
	assert.Equal(sId, sessionId)

	test := "projects/allan/agent/sdfsdfsd/324sd-dsfdsf/--/sessions/abc"
	pId, sId, err = ParseSession(test)
	assert.NoError(err)
	assert.Equal(sId, "sdfsdfsd/324sd-dsfdsf/--/sessions/abc")
	assert.Equal(pId, "allan")
}

func TestGenerateAndParseContext(t *testing.T) {
	assert := assert.New(t)
	projectId := "allan"
	sessionId := "3483478324vcxe/////sdf//nfkj(*#^$*(!!@#()@&$($&T*%#$*(_)@#*$"
	contextId := "dsfjh9sdo30ujdvioj324"

	_, err := GenerateContext("", "", "")
	assert.Error(err)

	context, err := GenerateContext(projectId, sessionId, contextId)
	assert.NoError(err)

	_, _, _, err = ParseContext("dsfsf34u390dsf/3praf/")
	assert.Error(err)

	pId, sId, cId, err := ParseContext(context)
	assert.NoError(err)
	assert.Equal(pId, projectId)
	assert.Equal(sId, sessionId)
	assert.Equal(cId, contextId)
}

func TestGenerateAndParseIntent(t *testing.T) {
	assert := assert.New(t)
	projectId := "allan"
	intentId := "esfoij23j4xcfd$$dfkl"

	_, err := GenerateIntent("", "")
	assert.Error(err)

	intent, err := GenerateIntent(projectId, intentId)
	assert.NoError(err)

	_, _, err = ParseIntent("dsfsf34u390dsf/3praf/")
	assert.Error(err)

	pId, iId, err := ParseIntent(intent)
	assert.NoError(err)
	assert.Equal(pId, projectId)
	assert.Equal(iId, intentId)
}

func TestWebhookRequest(t *testing.T) {
	assert := assert.New(t)
	req := &dialogflow.WebhookRequest{
		Session:    "session1",
		ResponseId: "response1",
		QueryResult: &dialogflow.QueryResult{
			QueryText:       "abc",
			FulfillmentText: "abc",
		},
	}
	result, err := WebhookRequestToJson(req)
	assert.NoError(err)
	req2, err := JsonToWebhookRequest(result)
	assert.NoError(err)
	assert.Equal(req.Session, req2.Session)
	assert.Equal(req.ResponseId, req2.ResponseId)
	assert.Equal(req.QueryResult.QueryText, req2.QueryResult.QueryText)
	sample := `{
  "responseId": "ea3d77e8-ae27-41a4-9e1d-174bd461b68c",
  "session": "projects/your-agents-project-id/agent/sessions/88d13aa8-2999-4f71-b233-39cbf3a824a0",
  "queryResult": {
    "queryText": "user's original query to your agent",
    "parameters": {
      "param": "param value"
    },
    "allRequiredParamsPresent": true,
    "fulfillmentText": "Text defined in Dialogflow's console for the intent that was matched",
    "fulfillmentMessages": [
      {
        "text": {
          "text": [
            "Text defined in Dialogflow's console for the intent that was matched"
          ]
        }
      }
    ],
    "outputContexts": [
      {
        "name": "projects/your-agents-project-id/agent/sessions/88d13aa8-2999-4f71-b233-39cbf3a824a0/contexts/generic",
        "lifespanCount": 5,
        "parameters": {
          "param": "param value"
        }
      }
    ],
    "intent": {
      "name": "projects/your-agents-project-id/agent/intents/29bcd7f8-f717-4261-a8fd-2d3e451b8af8",
      "displayName": "Matched Intent Name"
    },
    "intentDetectionConfidence": 1,
    "diagnosticInfo": {},
    "languageCode": "en"
  },
  "originalDetectIntentRequest": {}
}`

	var data map[string]interface{}
	err = json.Unmarshal([]byte(sample), &data)
	assert.NoError(err)

	req3, err := JsonToWebhookRequest([]byte(sample))
	assert.NoError(err)
	assert.Equal("ea3d77e8-ae27-41a4-9e1d-174bd461b68c", req3.ResponseId)
	assert.Equal("projects/your-agents-project-id/agent/sessions/88d13aa8-2999-4f71-b233-39cbf3a824a0",
		req3.Session)
	assert.Equal("param value", req3.QueryResult.Parameters.Fields["param"].GetStringValue())
}

func TestWebhookResponse(t *testing.T) {
	assert := assert.New(t)

	sample := `{
  "fulfillmentText": "This is a text response",
  "fulfillmentMessages": [
    {
      "card": {
        "title": "card title",
        "subtitle": "card text",
        "imageUri": "https://assistant.google.com/static/images/molecule/Molecule-Formation-stop.png",
        "buttons": [
          {
            "text": "button text",
            "postback": "https://assistant.google.com/"
          }
        ]
      }
    }
  ],
  "source": "example.com",
  "payload": {
    "google": {
      "expectUserResponse": true,
      "richResponse": {
        "items": [
          {
            "simpleResponse": {
              "textToSpeech": "this is a simple response"
            }
          }
        ]
      }
    },
    "facebook": {
      "text": "Hello, Facebook!"
    },
    "slack": {
      "text": "This is a text response for Slack."
    }
  },
  "outputContexts": [
    {
      "name": "projects/${PROJECT_ID}/agent/sessions/${SESSION_ID}/contexts/context name",
      "lifespanCount": 5,
      "parameters": {
        "param": "param value"
      }
    }
  ],
  "followupEventInput": {
    "name": "event name",
    "languageCode": "en-US",
    "parameters": {
      "param": "param value"
    }
  }
}`
	var data map[string]interface{}
	err := json.Unmarshal([]byte(sample), &data)
	assert.NoError(err)

	res, err := JsonToWebhookResponse([]byte(sample))
	assert.NoError(err)
	assert.Equal("This is a text response", res.FulfillmentText)
	assert.Equal("This is a text response for Slack.",
		res.Payload.Fields["slack"].GetStructValue().Fields["text"].GetStringValue())

	jsonValue, err := WebhookResponseToJson(res)
	assert.NoError(err)
	s := strings.Split(string(jsonValue), "")
	sort.Strings(s)
	e1 := strings.TrimSpace(strings.Join(s, ""))
	s1 := strings.Split(string(sample), "")
	sort.Strings(s1)
	e2 := strings.TrimSpace(strings.Join(s1, ""))
	assert.Equal(e1, e2)
}
