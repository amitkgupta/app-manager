// Code generated by protoc-gen-gogo.
// source: limit_bandwidth.proto
// DO NOT EDIT!

package warden

import proto "code.google.com/p/gogoprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type LimitBandwidthRequest struct {
	Handle           *string `protobuf:"bytes,1,req,name=handle" json:"handle,omitempty"`
	Rate             *uint64 `protobuf:"varint,2,req,name=rate" json:"rate,omitempty"`
	Burst            *uint64 `protobuf:"varint,3,req,name=burst" json:"burst,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *LimitBandwidthRequest) Reset()         { *m = LimitBandwidthRequest{} }
func (m *LimitBandwidthRequest) String() string { return proto.CompactTextString(m) }
func (*LimitBandwidthRequest) ProtoMessage()    {}

func (m *LimitBandwidthRequest) GetHandle() string {
	if m != nil && m.Handle != nil {
		return *m.Handle
	}
	return ""
}

func (m *LimitBandwidthRequest) GetRate() uint64 {
	if m != nil && m.Rate != nil {
		return *m.Rate
	}
	return 0
}

func (m *LimitBandwidthRequest) GetBurst() uint64 {
	if m != nil && m.Burst != nil {
		return *m.Burst
	}
	return 0
}

type LimitBandwidthResponse struct {
	Rate             *uint64 `protobuf:"varint,1,req,name=rate" json:"rate,omitempty"`
	Burst            *uint64 `protobuf:"varint,2,req,name=burst" json:"burst,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *LimitBandwidthResponse) Reset()         { *m = LimitBandwidthResponse{} }
func (m *LimitBandwidthResponse) String() string { return proto.CompactTextString(m) }
func (*LimitBandwidthResponse) ProtoMessage()    {}

func (m *LimitBandwidthResponse) GetRate() uint64 {
	if m != nil && m.Rate != nil {
		return *m.Rate
	}
	return 0
}

func (m *LimitBandwidthResponse) GetBurst() uint64 {
	if m != nil && m.Burst != nil {
		return *m.Burst
	}
	return 0
}

func init() {
}
