package actions_on_google

import (
	"testing"

	"encoding/json"
	"log"

	"github.com/stretchr/testify/assert"
)

func TestSystemIntentData(t *testing.T) {
	assert := assert.New(t)

	data := &SignInData{
		SignInValueSpec: SignInValueSpec{OptContext: "embeding"},
	}
	result, err := json.Marshal(data)
	assert.NoError(err)
	log.Println(string(result))
}
