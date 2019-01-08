package actions_on_google

import (
	"log"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

func TestComponent(t *testing.T) {
	assert := assert.New(t)
	mr := &MediaResponse{
		MediaType: MT_AUDIO,
		MediaObjects: []*MediaObject{
			{Name: "hi"},
		},
	}
	data, err := json.Marshal(mr)
	assert.NoError(err)
	log.Println(string(data))
}
