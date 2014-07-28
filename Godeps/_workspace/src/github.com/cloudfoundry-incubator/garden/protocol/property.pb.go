// Code generated by protoc-gen-gogo.
// source: property.proto
// DO NOT EDIT!

package warden

import proto "code.google.com/p/gogoprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type Property struct {
	Key              *string `protobuf:"bytes,1,req" json:"Key,omitempty"`
	Value            *string `protobuf:"bytes,2,req" json:"Value,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Property) Reset()         { *m = Property{} }
func (m *Property) String() string { return proto.CompactTextString(m) }
func (*Property) ProtoMessage()    {}

func (m *Property) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *Property) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

func init() {
}
