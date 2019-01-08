package actions_on_google

import (
	"log"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	assert := assert.New(t)

	response := &WebhookResponse{
		ExpectUserResponse: false,
		RichResponse:       &RichResponse{},
		SystemIntent:       &SystemIntent{},
	}
	bys, err := json.Marshal(response)
	assert.NoError(err)
	log.Println(string(bys))
}
