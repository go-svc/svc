// Code generated by protoc-gen-go.
// source: ping.proto
// DO NOT EDIT!

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	ping.proto

It has these top-level messages:
	Void
	Info
*/
package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Void 呈現一個什麼都沒有的資料。
type Void struct {
}

func (m *Void) Reset()                    { *m = Void{} }
func (m *Void) String() string            { return proto.CompactTextString(m) }
func (*Void) ProtoMessage()               {}
func (*Void) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// Info 呈現一個帶有服務版本號的資料。
type Info struct {
	Version string `protobuf:"bytes,1,opt,name=version" json:"version,omitempty"`
}

func (m *Info) Reset()                    { *m = Info{} }
func (m *Info) String() string            { return proto.CompactTextString(m) }
func (*Info) ProtoMessage()               {}
func (*Info) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Info) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func init() {
	proto.RegisterType((*Void)(nil), "pb.Void")
	proto.RegisterType((*Info)(nil), "pb.Info")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Test service

type TestClient interface {
	Ping(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Info, error)
}

type testClient struct {
	cc *grpc.ClientConn
}

func NewTestClient(cc *grpc.ClientConn) TestClient {
	return &testClient{cc}
}

func (c *testClient) Ping(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Info, error) {
	out := new(Info)
	err := grpc.Invoke(ctx, "/pb.Test/Ping", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Test service

type TestServer interface {
	Ping(context.Context, *Void) (*Info, error)
}

func RegisterTestServer(s *grpc.Server, srv TestServer) {
	s.RegisterService(&_Test_serviceDesc, srv)
}

func _Test_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Test/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServer).Ping(ctx, req.(*Void))
	}
	return interceptor(ctx, in, info, handler)
}

var _Test_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Test",
	HandlerType: (*TestServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Test_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ping.proto",
}

func init() { proto.RegisterFile("ping.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 112 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0xc8, 0xcc, 0x4b,
	0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0x62, 0xe3, 0x62, 0x09, 0xcb,
	0xcf, 0x4c, 0x51, 0x52, 0xe0, 0x62, 0xf1, 0xcc, 0x4b, 0xcb, 0x17, 0x92, 0xe0, 0x62, 0x2f, 0x4b,
	0x2d, 0x2a, 0xce, 0xcc, 0xcf, 0x93, 0x60, 0x54, 0x60, 0xd4, 0xe0, 0x0c, 0x82, 0x71, 0x8d, 0x54,
	0xb8, 0x58, 0x42, 0x52, 0x8b, 0x4b, 0x84, 0x64, 0xb8, 0x58, 0x02, 0x32, 0xf3, 0xd2, 0x85, 0x38,
	0xf4, 0x0a, 0x92, 0xf4, 0x40, 0x7a, 0xa5, 0xc0, 0x2c, 0x90, 0x6e, 0x25, 0x86, 0x24, 0x36, 0xb0,
	0xd1, 0xc6, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x4e, 0x0b, 0xd6, 0x02, 0x68, 0x00, 0x00, 0x00,
}
