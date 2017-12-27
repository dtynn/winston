// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	common.proto
	meta.proto
	raft.proto

It has these top-level messages:
	Result
	Node
	MetaAddMetaNodeReq
	MetaAddMetaNodeResp
	DataNode
	RaftGroup
	RaftAddGroup
	RaftAddGroupResp
	RaftMessage
	RaftGroupMessage
	RaftMessageResponse
	RaftGroupMessageResponse
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ResultCode int32

const (
	ResultCode_ResultCodeOK                   ResultCode = 0
	ResultCode_ResultCodeUnknown              ResultCode = 1
	ResultCode_ResultCodeInternal             ResultCode = 2
	ResultCode_ResultCodeRaftDuplicateGroupID ResultCode = 100001
)

var ResultCode_name = map[int32]string{
	0:      "ResultCodeOK",
	1:      "ResultCodeUnknown",
	2:      "ResultCodeInternal",
	100001: "ResultCodeRaftDuplicateGroupID",
}
var ResultCode_value = map[string]int32{
	"ResultCodeOK":                   0,
	"ResultCodeUnknown":              1,
	"ResultCodeInternal":             2,
	"ResultCodeRaftDuplicateGroupID": 100001,
}

func (x ResultCode) String() string {
	return proto.EnumName(ResultCode_name, int32(x))
}
func (ResultCode) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Result struct {
	Code    ResultCode `protobuf:"varint,1,opt,name=code,enum=pb.ResultCode" json:"code,omitempty"`
	Message string     `protobuf:"bytes,2,opt,name=Message" json:"Message,omitempty"`
}

func (m *Result) Reset()                    { *m = Result{} }
func (m *Result) String() string            { return proto.CompactTextString(m) }
func (*Result) ProtoMessage()               {}
func (*Result) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Result) GetCode() ResultCode {
	if m != nil {
		return m.Code
	}
	return ResultCode_ResultCodeOK
}

func (m *Result) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type Node struct {
	Id      uint64 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Context []byte `protobuf:"bytes,2,opt,name=context,proto3" json:"context,omitempty"`
}

func (m *Node) Reset()                    { *m = Node{} }
func (m *Node) String() string            { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()               {}
func (*Node) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Node) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Node) GetContext() []byte {
	if m != nil {
		return m.Context
	}
	return nil
}

func init() {
	proto.RegisterType((*Result)(nil), "pb.Result")
	proto.RegisterType((*Node)(nil), "pb.Node")
	proto.RegisterEnum("pb.ResultCode", ResultCode_name, ResultCode_value)
}

func init() { proto.RegisterFile("common.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 216 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0xce, 0xcf, 0xcd,
	0xcd, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0x72, 0xe3, 0x62,
	0x0b, 0x4a, 0x2d, 0x2e, 0xcd, 0x29, 0x11, 0x52, 0xe2, 0x62, 0x49, 0xce, 0x4f, 0x49, 0x95, 0x60,
	0x54, 0x60, 0xd4, 0xe0, 0x33, 0xe2, 0xd3, 0x2b, 0x48, 0xd2, 0x83, 0xc8, 0x38, 0xe7, 0xa7, 0xa4,
	0x06, 0x81, 0xe5, 0x84, 0x24, 0xb8, 0xd8, 0x7d, 0x53, 0x8b, 0x8b, 0x13, 0xd3, 0x53, 0x25, 0x98,
	0x14, 0x18, 0x35, 0x38, 0x83, 0x60, 0x5c, 0x25, 0x03, 0x2e, 0x16, 0x3f, 0x90, 0x0a, 0x3e, 0x2e,
	0xa6, 0xcc, 0x14, 0xb0, 0x19, 0x2c, 0x41, 0x4c, 0x99, 0x29, 0x20, 0x1d, 0xc9, 0xf9, 0x79, 0x25,
	0xa9, 0x15, 0x25, 0x60, 0x1d, 0x3c, 0x41, 0x30, 0xae, 0x56, 0x31, 0x17, 0x17, 0xc2, 0x7c, 0x21,
	0x01, 0x2e, 0x1e, 0x04, 0xcf, 0xdf, 0x5b, 0x80, 0x41, 0x48, 0x94, 0x4b, 0x10, 0x21, 0x12, 0x9a,
	0x97, 0x9d, 0x97, 0x5f, 0x9e, 0x27, 0xc0, 0x28, 0x24, 0xc6, 0x25, 0x84, 0x10, 0xf6, 0xcc, 0x2b,
	0x49, 0x2d, 0xca, 0x4b, 0xcc, 0x11, 0x60, 0x12, 0x52, 0xe1, 0x92, 0x43, 0x72, 0x6e, 0x62, 0x5a,
	0x89, 0x4b, 0x69, 0x41, 0x4e, 0x66, 0x72, 0x62, 0x49, 0xaa, 0x7b, 0x51, 0x7e, 0x69, 0x81, 0xa7,
	0x8b, 0xc0, 0xc2, 0x5e, 0xb6, 0x24, 0x36, 0xb0, 0xcf, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff,
	0x0f, 0x5e, 0x4e, 0xb2, 0x09, 0x01, 0x00, 0x00,
}