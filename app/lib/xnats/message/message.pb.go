// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package message

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import any "github.com/golang/protobuf/ptypes/any"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Message struct {
	Body                 *any.Any          `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
	Meta                 map[string]string `protobuf:"bytes,2,rep,name=meta,proto3" json:"meta,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_message_d85672dedb117d15, []int{0}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Message.Unmarshal(m, b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Message.Marshal(b, m, deterministic)
}
func (dst *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(dst, src)
}
func (m *Message) XXX_Size() int {
	return xxx_messageInfo_Message.Size(m)
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetBody() *any.Any {
	if m != nil {
		return m.Body
	}
	return nil
}

func (m *Message) GetMeta() map[string]string {
	if m != nil {
		return m.Meta
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "xnats.message.Message")
	proto.RegisterMapType((map[string]string)(nil), "xnats.message.Message.MetaEntry")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_message_d85672dedb117d15) }

var fileDescriptor_message_d85672dedb117d15 = []byte{
	// 187 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x4d, 0x2d, 0x2e,
	0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0xad, 0xc8, 0x4b, 0x2c, 0x29,
	0xd6, 0x83, 0x0a, 0x4a, 0x49, 0xa6, 0xe7, 0xe7, 0xa7, 0xe7, 0xa4, 0xea, 0x83, 0x25, 0x93, 0x4a,
	0xd3, 0xf4, 0x13, 0xf3, 0x2a, 0x21, 0x2a, 0x95, 0x16, 0x31, 0x72, 0xb1, 0xfb, 0x42, 0x94, 0x09,
	0x69, 0x70, 0xb1, 0x24, 0xe5, 0xa7, 0x54, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x1b, 0x89, 0xe8,
	0x41, 0x74, 0xe9, 0xc1, 0x74, 0xe9, 0x39, 0xe6, 0x55, 0x06, 0x81, 0x55, 0x08, 0x99, 0x70, 0xb1,
	0xe4, 0xa6, 0x96, 0x24, 0x4a, 0x30, 0x29, 0x30, 0x6b, 0x70, 0x1b, 0x29, 0xe8, 0xa1, 0x58, 0xa7,
	0xe7, 0x0b, 0xa7, 0x4b, 0x12, 0x5d, 0xf3, 0x4a, 0x8a, 0x2a, 0x83, 0xc0, 0xaa, 0xa5, 0xcc, 0xb9,
	0x38, 0xe1, 0x42, 0x42, 0x02, 0x5c, 0xcc, 0xd9, 0xa9, 0x10, 0xbb, 0x38, 0x83, 0x40, 0x4c, 0x21,
	0x11, 0x2e, 0xd6, 0xb2, 0xc4, 0x9c, 0xd2, 0x54, 0x09, 0x26, 0xb0, 0x18, 0x84, 0x63, 0xc5, 0x64,
	0xc1, 0xe8, 0xc4, 0x19, 0xc5, 0x0e, 0x35, 0x3b, 0x89, 0x0d, 0xec, 0x1a, 0x63, 0x40, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xbb, 0x2d, 0xef, 0x5d, 0xf1, 0x00, 0x00, 0x00,
}