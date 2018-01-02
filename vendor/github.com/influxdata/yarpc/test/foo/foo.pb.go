// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: test/foo/foo.proto

/*
Package foo is a generated protocol buffer package.

It is generated from these files:
	test/foo/foo.proto

It has these top-level messages:
	Request
	Response
*/
package foo

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/influxdata/yarpc/yarpcproto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type Request struct {
	In string `protobuf:"bytes,1,opt,name=in,proto3" json:"in,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptorFoo, []int{0} }

func (m *Request) GetIn() string {
	if m != nil {
		return m.In
	}
	return ""
}

type Response struct {
	Out string `protobuf:"bytes,1,opt,name=out,proto3" json:"out,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptorFoo, []int{1} }

func (m *Response) GetOut() string {
	if m != nil {
		return m.Out
	}
	return ""
}

func init() {
	proto.RegisterType((*Request)(nil), "foo.Request")
	proto.RegisterType((*Response)(nil), "foo.Response")
}

func init() { proto.RegisterFile("test/foo/foo.proto", fileDescriptorFoo) }

var fileDescriptorFoo = []byte{
	// 206 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2a, 0x49, 0x2d, 0x2e,
	0xd1, 0x4f, 0xcb, 0xcf, 0x07, 0x61, 0xbd, 0x82, 0xa2, 0xfc, 0x92, 0x7c, 0x21, 0xe6, 0xb4, 0xfc,
	0x7c, 0x29, 0x93, 0xf4, 0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0xfd, 0xe2, 0x92,
	0xd2, 0xc4, 0xa2, 0x92, 0xe4, 0xc4, 0xa2, 0xbc, 0xcc, 0x54, 0xfd, 0xca, 0xc4, 0xa2, 0x82, 0x64,
	0x08, 0x09, 0x56, 0x0e, 0x61, 0x42, 0xb4, 0x2a, 0x49, 0x72, 0xb1, 0x07, 0xa5, 0x16, 0x96, 0xa6,
	0x16, 0x97, 0x08, 0xf1, 0x71, 0x31, 0x65, 0xe6, 0x49, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x31,
	0x65, 0xe6, 0x29, 0xc9, 0x70, 0x71, 0x04, 0xa5, 0x16, 0x17, 0xe4, 0xe7, 0x15, 0xa7, 0x0a, 0x09,
	0x70, 0x31, 0xe7, 0x97, 0x96, 0x40, 0x25, 0x41, 0x4c, 0xa3, 0x0a, 0x2e, 0x66, 0xb7, 0xfc, 0x7c,
	0x21, 0x03, 0x2e, 0xee, 0xd0, 0xbc, 0xc4, 0xa2, 0x4a, 0xdf, 0xd4, 0x92, 0x8c, 0xfc, 0x14, 0x21,
	0x1e, 0x3d, 0x90, 0xab, 0xa0, 0x26, 0x4a, 0xf1, 0x42, 0x79, 0x10, 0x43, 0x94, 0x58, 0x1a, 0xb6,
	0x4a, 0x30, 0x08, 0x59, 0x72, 0x09, 0x05, 0xa7, 0x16, 0x95, 0xa5, 0x16, 0x05, 0x97, 0x14, 0xa5,
	0x26, 0xe6, 0x12, 0xab, 0x91, 0xd1, 0x80, 0x51, 0x0a, 0x6c, 0x40, 0x12, 0x1b, 0xd8, 0xe5, 0xc6,
	0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x29, 0xba, 0xb4, 0xf3, 0x0a, 0x01, 0x00, 0x00,
}