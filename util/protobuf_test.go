package util

import (
	"testing"

	"github.com/golang/protobuf/ptypes/struct"
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

func TestStructPBToMap(t *testing.T) {
	assert := assert.New(t)
	s := &structpb.Struct{
		Fields: map[string]*structpb.Value{},
	}
	s.Fields["allan"] = &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: "test"}}
	s.Fields["allan2"] = &structpb.Value{Kind: &structpb.Value_ListValue{
		ListValue: &structpb.ListValue{
			Values: []*structpb.Value{
				&structpb.Value{
					Kind: &structpb.Value_StringValue{StringValue: "test2"},
				},
			},
		},
	}}
	s.Fields["allan3"] = &structpb.Value{Kind: &structpb.Value_StructValue{
		StructValue: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"test3": &structpb.Value{
					Kind: &structpb.Value_StringValue{StringValue: "test4"},
				},
			},
		},
	}}

	m, err := StructPBToMap(s)
	assert.NoError(err)
	_, ok := m["allan"]
	assert.True(ok)
	_, ok = m["allan2"]
	assert.True(ok)
	_, ok = m["allan3"]
	assert.True(ok)
}
