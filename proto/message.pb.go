// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

package message

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type JoinRequest struct {
	IpAddr               int32    `protobuf:"varint,1,opt,name=ipAddr,proto3" json:"ipAddr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JoinRequest) Reset()         { *m = JoinRequest{} }
func (m *JoinRequest) String() string { return proto.CompactTextString(m) }
func (*JoinRequest) ProtoMessage()    {}
func (*JoinRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}

func (m *JoinRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinRequest.Unmarshal(m, b)
}
func (m *JoinRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinRequest.Marshal(b, m, deterministic)
}
func (m *JoinRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinRequest.Merge(m, src)
}
func (m *JoinRequest) XXX_Size() int {
	return xxx_messageInfo_JoinRequest.Size(m)
}
func (m *JoinRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinRequest.DiscardUnknown(m)
}

var xxx_messageInfo_JoinRequest proto.InternalMessageInfo

func (m *JoinRequest) GetIpAddr() int32 {
	if m != nil {
		return m.IpAddr
	}
	return 0
}

type JoinReply struct {
	Reply                bool     `protobuf:"varint,1,opt,name=reply,proto3" json:"reply,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JoinReply) Reset()         { *m = JoinReply{} }
func (m *JoinReply) String() string { return proto.CompactTextString(m) }
func (*JoinReply) ProtoMessage()    {}
func (*JoinReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{1}
}

func (m *JoinReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinReply.Unmarshal(m, b)
}
func (m *JoinReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinReply.Marshal(b, m, deterministic)
}
func (m *JoinReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinReply.Merge(m, src)
}
func (m *JoinReply) XXX_Size() int {
	return xxx_messageInfo_JoinReply.Size(m)
}
func (m *JoinReply) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinReply.DiscardUnknown(m)
}

var xxx_messageInfo_JoinReply proto.InternalMessageInfo

func (m *JoinReply) GetReply() bool {
	if m != nil {
		return m.Reply
	}
	return false
}

type ExitRequest struct {
	IpAddr               int32    `protobuf:"varint,1,opt,name=ipAddr,proto3" json:"ipAddr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExitRequest) Reset()         { *m = ExitRequest{} }
func (m *ExitRequest) String() string { return proto.CompactTextString(m) }
func (*ExitRequest) ProtoMessage()    {}
func (*ExitRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{2}
}

func (m *ExitRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExitRequest.Unmarshal(m, b)
}
func (m *ExitRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExitRequest.Marshal(b, m, deterministic)
}
func (m *ExitRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExitRequest.Merge(m, src)
}
func (m *ExitRequest) XXX_Size() int {
	return xxx_messageInfo_ExitRequest.Size(m)
}
func (m *ExitRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ExitRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ExitRequest proto.InternalMessageInfo

func (m *ExitRequest) GetIpAddr() int32 {
	if m != nil {
		return m.IpAddr
	}
	return 0
}

type ExitReply struct {
	Reply                bool     `protobuf:"varint,1,opt,name=reply,proto3" json:"reply,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExitReply) Reset()         { *m = ExitReply{} }
func (m *ExitReply) String() string { return proto.CompactTextString(m) }
func (*ExitReply) ProtoMessage()    {}
func (*ExitReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{3}
}

func (m *ExitReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExitReply.Unmarshal(m, b)
}
func (m *ExitReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExitReply.Marshal(b, m, deterministic)
}
func (m *ExitReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExitReply.Merge(m, src)
}
func (m *ExitReply) XXX_Size() int {
	return xxx_messageInfo_ExitReply.Size(m)
}
func (m *ExitReply) XXX_DiscardUnknown() {
	xxx_messageInfo_ExitReply.DiscardUnknown(m)
}

var xxx_messageInfo_ExitReply proto.InternalMessageInfo

func (m *ExitReply) GetReply() bool {
	if m != nil {
		return m.Reply
	}
	return false
}

type HeartBeatRequest struct {
	Time                 int64    `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartBeatRequest) Reset()         { *m = HeartBeatRequest{} }
func (m *HeartBeatRequest) String() string { return proto.CompactTextString(m) }
func (*HeartBeatRequest) ProtoMessage()    {}
func (*HeartBeatRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{4}
}

func (m *HeartBeatRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartBeatRequest.Unmarshal(m, b)
}
func (m *HeartBeatRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartBeatRequest.Marshal(b, m, deterministic)
}
func (m *HeartBeatRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartBeatRequest.Merge(m, src)
}
func (m *HeartBeatRequest) XXX_Size() int {
	return xxx_messageInfo_HeartBeatRequest.Size(m)
}
func (m *HeartBeatRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartBeatRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HeartBeatRequest proto.InternalMessageInfo

func (m *HeartBeatRequest) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

type HeartBeatReply struct {
	Time                 int64    `protobuf:"varint,1,opt,name=time,proto3" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartBeatReply) Reset()         { *m = HeartBeatReply{} }
func (m *HeartBeatReply) String() string { return proto.CompactTextString(m) }
func (*HeartBeatReply) ProtoMessage()    {}
func (*HeartBeatReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{5}
}

func (m *HeartBeatReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartBeatReply.Unmarshal(m, b)
}
func (m *HeartBeatReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartBeatReply.Marshal(b, m, deterministic)
}
func (m *HeartBeatReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartBeatReply.Merge(m, src)
}
func (m *HeartBeatReply) XXX_Size() int {
	return xxx_messageInfo_HeartBeatReply.Size(m)
}
func (m *HeartBeatReply) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartBeatReply.DiscardUnknown(m)
}

var xxx_messageInfo_HeartBeatReply proto.InternalMessageInfo

func (m *HeartBeatReply) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func init() {
	proto.RegisterType((*JoinRequest)(nil), "message.JoinRequest")
	proto.RegisterType((*JoinReply)(nil), "message.JoinReply")
	proto.RegisterType((*ExitRequest)(nil), "message.ExitRequest")
	proto.RegisterType((*ExitReply)(nil), "message.ExitReply")
	proto.RegisterType((*HeartBeatRequest)(nil), "message.HeartBeatRequest")
	proto.RegisterType((*HeartBeatReply)(nil), "message.HeartBeatReply")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 226 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcd, 0x4d, 0x2d, 0x2e,
	0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x72, 0x95, 0x54, 0xb9,
	0xb8, 0xbd, 0xf2, 0x33, 0xf3, 0x82, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0xc4, 0xb8, 0xd8,
	0x32, 0x0b, 0x1c, 0x53, 0x52, 0x8a, 0x24, 0x18, 0x15, 0x18, 0x35, 0x58, 0x83, 0xa0, 0x3c, 0x25,
	0x45, 0x2e, 0x4e, 0x88, 0xb2, 0x82, 0x9c, 0x4a, 0x21, 0x11, 0x2e, 0xd6, 0x22, 0x10, 0x03, 0xac,
	0x86, 0x23, 0x08, 0xc2, 0x01, 0x99, 0xe4, 0x5a, 0x91, 0x59, 0x42, 0x84, 0x49, 0x10, 0x65, 0xb8,
	0x4d, 0x52, 0xe3, 0x12, 0xf0, 0x48, 0x4d, 0x2c, 0x2a, 0x71, 0x4a, 0x4d, 0x84, 0x1b, 0x27, 0xc4,
	0xc5, 0x52, 0x92, 0x99, 0x9b, 0x0a, 0x56, 0xc8, 0x1c, 0x04, 0x66, 0x2b, 0xa9, 0x70, 0xf1, 0x21,
	0xa9, 0x03, 0x99, 0x87, 0x45, 0x95, 0xd1, 0x41, 0x46, 0x2e, 0x2e, 0xe7, 0xfc, 0xbc, 0xbc, 0xd4,
	0xe4, 0x92, 0xcc, 0xfc, 0x3c, 0x21, 0x73, 0x88, 0x4f, 0xdc, 0x8b, 0xf2, 0x4b, 0x0b, 0x84, 0x44,
	0xf4, 0x60, 0xc1, 0x82, 0x14, 0x08, 0x52, 0x42, 0x68, 0xa2, 0x20, 0x37, 0x31, 0x80, 0x34, 0x82,
	0x1c, 0x8e, 0xae, 0x11, 0xc9, 0xcf, 0x48, 0x1a, 0xe1, 0x5e, 0x54, 0x62, 0x10, 0x72, 0xe4, 0xe2,
	0x84, 0x3b, 0x53, 0x48, 0x12, 0xae, 0x04, 0xdd, 0x8b, 0x52, 0xe2, 0xd8, 0xa4, 0xc0, 0x46, 0x24,
	0xb1, 0x81, 0x63, 0xcd, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x62, 0x00, 0x19, 0xb8, 0xc6, 0x01,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ConnectionClient is the client API for Connection service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ConnectionClient interface {
	JoinGroup(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinReply, error)
	ExitGroup(ctx context.Context, in *ExitRequest, opts ...grpc.CallOption) (*ExitReply, error)
	HeartBeat(ctx context.Context, in *HeartBeatRequest, opts ...grpc.CallOption) (*HeartBeatReply, error)
}

type connectionClient struct {
	cc *grpc.ClientConn
}

func NewConnectionClient(cc *grpc.ClientConn) ConnectionClient {
	return &connectionClient{cc}
}

func (c *connectionClient) JoinGroup(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinReply, error) {
	out := new(JoinReply)
	err := c.cc.Invoke(ctx, "/message.Connection/JoinGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectionClient) ExitGroup(ctx context.Context, in *ExitRequest, opts ...grpc.CallOption) (*ExitReply, error) {
	out := new(ExitReply)
	err := c.cc.Invoke(ctx, "/message.Connection/ExitGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *connectionClient) HeartBeat(ctx context.Context, in *HeartBeatRequest, opts ...grpc.CallOption) (*HeartBeatReply, error) {
	out := new(HeartBeatReply)
	err := c.cc.Invoke(ctx, "/message.Connection/HeartBeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ConnectionServer is the server API for Connection service.
type ConnectionServer interface {
	JoinGroup(context.Context, *JoinRequest) (*JoinReply, error)
	ExitGroup(context.Context, *ExitRequest) (*ExitReply, error)
	HeartBeat(context.Context, *HeartBeatRequest) (*HeartBeatReply, error)
}

// UnimplementedConnectionServer can be embedded to have forward compatible implementations.
type UnimplementedConnectionServer struct {
}

func (*UnimplementedConnectionServer) JoinGroup(ctx context.Context, req *JoinRequest) (*JoinReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinGroup not implemented")
}
func (*UnimplementedConnectionServer) ExitGroup(ctx context.Context, req *ExitRequest) (*ExitReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExitGroup not implemented")
}
func (*UnimplementedConnectionServer) HeartBeat(ctx context.Context, req *HeartBeatRequest) (*HeartBeatReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HeartBeat not implemented")
}

func RegisterConnectionServer(s *grpc.Server, srv ConnectionServer) {
	s.RegisterService(&_Connection_serviceDesc, srv)
}

func _Connection_JoinGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectionServer).JoinGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.Connection/JoinGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectionServer).JoinGroup(ctx, req.(*JoinRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Connection_ExitGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectionServer).ExitGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.Connection/ExitGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectionServer).ExitGroup(ctx, req.(*ExitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Connection_HeartBeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeartBeatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConnectionServer).HeartBeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/message.Connection/HeartBeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConnectionServer).HeartBeat(ctx, req.(*HeartBeatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Connection_serviceDesc = grpc.ServiceDesc{
	ServiceName: "message.Connection",
	HandlerType: (*ConnectionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "JoinGroup",
			Handler:    _Connection_JoinGroup_Handler,
		},
		{
			MethodName: "ExitGroup",
			Handler:    _Connection_ExitGroup_Handler,
		},
		{
			MethodName: "HeartBeat",
			Handler:    _Connection_HeartBeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "message.proto",
}
