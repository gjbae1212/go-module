package util

import (
	"bytes"
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/struct"
)

func MapToStructPB(m map[string]interface{}) (*structpb.Struct, error) {
	jb, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	unmarshaler := &jsonpb.Unmarshaler{}
	result := &structpb.Struct{}
	if err = unmarshaler.Unmarshal(bytes.NewReader(jb), result); err != nil {
		return nil, err
	}
	return result, nil
}

func ObjectToStructPB(o interface{}) (*structpb.Struct, error) {
	jb, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	unmarshaler := &jsonpb.Unmarshaler{}
	result := &structpb.Struct{}
	if err = unmarshaler.Unmarshal(bytes.NewReader(jb), result); err != nil {
		return nil, err
	}
	return result, nil
}
