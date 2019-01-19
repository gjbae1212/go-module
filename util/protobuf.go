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

func StructPBToMap(s *structpb.Struct) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	if s == nil {
		return result, nil
	}
	for name, value := range s.GetFields() {
		result[name] = ValuePBToInterface(value)
	}
	return result, nil
}

func ValuePBToInterface(v *structpb.Value) interface{} {
	switch v.GetKind().(type) {
	case *structpb.Value_StringValue:
		return v.GetStringValue()
	case *structpb.Value_NumberValue:
		return v.GetNumberValue()
	case *structpb.Value_BoolValue:
		return v.GetBoolValue()
	case *structpb.Value_StructValue:
		return v.GetStructValue()
	case *structpb.Value_ListValue:
		return v.GetListValue()
	case *structpb.Value_NullValue:
		return v.GetNullValue()
	default:
		return nil
	}
}
