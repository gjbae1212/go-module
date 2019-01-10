package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToStructPB(t *testing.T) {
	assert := assert.New(t)

	m := make(map[string]interface{})
	m["allan"] = "dong"
	m["power"] = 10
	result, err := MapToStructPB(m)
	assert.NoError(err)
	assert.Equal(result.Fields["allan"].GetStringValue(), m["allan"])
	assert.Equal(result.Fields["power"].GetNumberValue(), float64(m["power"].(int)))
}

func TestObjectToStructPB(t *testing.T) {
	assert := assert.New(t)

	s := struct {
		Text string `json:"text,omitempty"`
	}{
		Text: "hihi",
	}

	result, err := ObjectToStructPB(s)
	assert.NoError(err)
	assert.Equal(result.Fields["text"].GetStringValue(), s.Text)
}
