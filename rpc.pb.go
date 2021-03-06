// Code generated by protoc-gen-go.
// source: rpc.proto
// DO NOT EDIT!

/*
Package chordstore is a generated protocol buffer package.

It is generated from these files:
	rpc.proto

It has these top-level messages:
	DHTKeyValue
	DHTHashKeyValue
	DHTBytes
	DHTBytesErr
	SnapshotOptions
	DataStream
*/
package chordstore

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import chord "github.com/euforia/go-chord"

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

type DHTKeyValue struct {
	Vn    *chord.Vnode `protobuf:"bytes,1,opt,name=vn" json:"vn,omitempty"`
	Key   []byte       `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte       `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *DHTKeyValue) Reset()                    { *m = DHTKeyValue{} }
func (m *DHTKeyValue) String() string            { return proto.CompactTextString(m) }
func (*DHTKeyValue) ProtoMessage()               {}
func (*DHTKeyValue) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *DHTKeyValue) GetVn() *chord.Vnode {
	if m != nil {
		return m.Vn
	}
	return nil
}

func (m *DHTKeyValue) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *DHTKeyValue) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type DHTHashKeyValue struct {
	Vn       *chord.Vnode `protobuf:"bytes,1,opt,name=vn" json:"vn,omitempty"`
	PrevHash []byte       `protobuf:"bytes,2,opt,name=prevHash,proto3" json:"prevHash,omitempty"`
	Key      []byte       `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	Value    []byte       `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
}

func (m *DHTHashKeyValue) Reset()                    { *m = DHTHashKeyValue{} }
func (m *DHTHashKeyValue) String() string            { return proto.CompactTextString(m) }
func (*DHTHashKeyValue) ProtoMessage()               {}
func (*DHTHashKeyValue) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *DHTHashKeyValue) GetVn() *chord.Vnode {
	if m != nil {
		return m.Vn
	}
	return nil
}

func (m *DHTHashKeyValue) GetPrevHash() []byte {
	if m != nil {
		return m.PrevHash
	}
	return nil
}

func (m *DHTHashKeyValue) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *DHTHashKeyValue) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type DHTBytes struct {
	Vn *chord.Vnode `protobuf:"bytes,1,opt,name=vn" json:"vn,omitempty"`
	B  []byte       `protobuf:"bytes,2,opt,name=b,proto3" json:"b,omitempty"`
}

func (m *DHTBytes) Reset()                    { *m = DHTBytes{} }
func (m *DHTBytes) String() string            { return proto.CompactTextString(m) }
func (*DHTBytes) ProtoMessage()               {}
func (*DHTBytes) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *DHTBytes) GetVn() *chord.Vnode {
	if m != nil {
		return m.Vn
	}
	return nil
}

func (m *DHTBytes) GetB() []byte {
	if m != nil {
		return m.B
	}
	return nil
}

type DHTBytesErr struct {
	B   []byte `protobuf:"bytes,1,opt,name=b,proto3" json:"b,omitempty"`
	Err string `protobuf:"bytes,2,opt,name=err" json:"err,omitempty"`
}

func (m *DHTBytesErr) Reset()                    { *m = DHTBytesErr{} }
func (m *DHTBytesErr) String() string            { return proto.CompactTextString(m) }
func (*DHTBytesErr) ProtoMessage()               {}
func (*DHTBytesErr) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *DHTBytesErr) GetB() []byte {
	if m != nil {
		return m.B
	}
	return nil
}

func (m *DHTBytesErr) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type SnapshotOptions struct {
	Vn *chord.Vnode `protobuf:"bytes,1,opt,name=vn" json:"vn,omitempty"`
}

func (m *SnapshotOptions) Reset()                    { *m = SnapshotOptions{} }
func (m *SnapshotOptions) String() string            { return proto.CompactTextString(m) }
func (*SnapshotOptions) ProtoMessage()               {}
func (*SnapshotOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *SnapshotOptions) GetVn() *chord.Vnode {
	if m != nil {
		return m.Vn
	}
	return nil
}

type DataStream struct {
	Data []byte `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (m *DataStream) Reset()                    { *m = DataStream{} }
func (m *DataStream) String() string            { return proto.CompactTextString(m) }
func (*DataStream) ProtoMessage()               {}
func (*DataStream) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *DataStream) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*DHTKeyValue)(nil), "chordstore.DHTKeyValue")
	proto.RegisterType((*DHTHashKeyValue)(nil), "chordstore.DHTHashKeyValue")
	proto.RegisterType((*DHTBytes)(nil), "chordstore.DHTBytes")
	proto.RegisterType((*DHTBytesErr)(nil), "chordstore.DHTBytesErr")
	proto.RegisterType((*SnapshotOptions)(nil), "chordstore.SnapshotOptions")
	proto.RegisterType((*DataStream)(nil), "chordstore.DataStream")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for DHT service

type DHTClient interface {
	PutKeyRPC(ctx context.Context, in *DHTKeyValue, opts ...grpc.CallOption) (*chord.ErrResponse, error)
	GetKeyRPC(ctx context.Context, in *DHTBytes, opts ...grpc.CallOption) (*DHTBytesErr, error)
	UpdateKeyRPC(ctx context.Context, in *DHTHashKeyValue, opts ...grpc.CallOption) (*chord.ErrResponse, error)
	RemoveKeyRPC(ctx context.Context, in *DHTBytes, opts ...grpc.CallOption) (*chord.ErrResponse, error)
	PutObjectRPC(ctx context.Context, opts ...grpc.CallOption) (DHT_PutObjectRPCClient, error)
	GetObjectRPC(ctx context.Context, in *DHTBytes, opts ...grpc.CallOption) (DHT_GetObjectRPCClient, error)
	RemoveObjectRPC(ctx context.Context, in *DHTBytes, opts ...grpc.CallOption) (*chord.ErrResponse, error)
	SnapshotRPC(ctx context.Context, in *chord.Vnode, opts ...grpc.CallOption) (DHT_SnapshotRPCClient, error)
	RestoreRPC(ctx context.Context, opts ...grpc.CallOption) (DHT_RestoreRPCClient, error)
}

type dHTClient struct {
	cc *grpc.ClientConn
}

func NewDHTClient(cc *grpc.ClientConn) DHTClient {
	return &dHTClient{cc}
}

func (c *dHTClient) PutKeyRPC(ctx context.Context, in *DHTKeyValue, opts ...grpc.CallOption) (*chord.ErrResponse, error) {
	out := new(chord.ErrResponse)
	err := grpc.Invoke(ctx, "/chordstore.DHT/PutKeyRPC", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dHTClient) GetKeyRPC(ctx context.Context, in *DHTBytes, opts ...grpc.CallOption) (*DHTBytesErr, error) {
	out := new(DHTBytesErr)
	err := grpc.Invoke(ctx, "/chordstore.DHT/GetKeyRPC", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dHTClient) UpdateKeyRPC(ctx context.Context, in *DHTHashKeyValue, opts ...grpc.CallOption) (*chord.ErrResponse, error) {
	out := new(chord.ErrResponse)
	err := grpc.Invoke(ctx, "/chordstore.DHT/UpdateKeyRPC", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dHTClient) RemoveKeyRPC(ctx context.Context, in *DHTBytes, opts ...grpc.CallOption) (*chord.ErrResponse, error) {
	out := new(chord.ErrResponse)
	err := grpc.Invoke(ctx, "/chordstore.DHT/RemoveKeyRPC", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dHTClient) PutObjectRPC(ctx context.Context, opts ...grpc.CallOption) (DHT_PutObjectRPCClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_DHT_serviceDesc.Streams[0], c.cc, "/chordstore.DHT/PutObjectRPC", opts...)
	if err != nil {
		return nil, err
	}
	x := &dHTPutObjectRPCClient{stream}
	return x, nil
}

type DHT_PutObjectRPCClient interface {
	Send(*DataStream) error
	CloseAndRecv() (*chord.ErrResponse, error)
	grpc.ClientStream
}

type dHTPutObjectRPCClient struct {
	grpc.ClientStream
}

func (x *dHTPutObjectRPCClient) Send(m *DataStream) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dHTPutObjectRPCClient) CloseAndRecv() (*chord.ErrResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(chord.ErrResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dHTClient) GetObjectRPC(ctx context.Context, in *DHTBytes, opts ...grpc.CallOption) (DHT_GetObjectRPCClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_DHT_serviceDesc.Streams[1], c.cc, "/chordstore.DHT/GetObjectRPC", opts...)
	if err != nil {
		return nil, err
	}
	x := &dHTGetObjectRPCClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DHT_GetObjectRPCClient interface {
	Recv() (*DataStream, error)
	grpc.ClientStream
}

type dHTGetObjectRPCClient struct {
	grpc.ClientStream
}

func (x *dHTGetObjectRPCClient) Recv() (*DataStream, error) {
	m := new(DataStream)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dHTClient) RemoveObjectRPC(ctx context.Context, in *DHTBytes, opts ...grpc.CallOption) (*chord.ErrResponse, error) {
	out := new(chord.ErrResponse)
	err := grpc.Invoke(ctx, "/chordstore.DHT/RemoveObjectRPC", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dHTClient) SnapshotRPC(ctx context.Context, in *chord.Vnode, opts ...grpc.CallOption) (DHT_SnapshotRPCClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_DHT_serviceDesc.Streams[2], c.cc, "/chordstore.DHT/SnapshotRPC", opts...)
	if err != nil {
		return nil, err
	}
	x := &dHTSnapshotRPCClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DHT_SnapshotRPCClient interface {
	Recv() (*DataStream, error)
	grpc.ClientStream
}

type dHTSnapshotRPCClient struct {
	grpc.ClientStream
}

func (x *dHTSnapshotRPCClient) Recv() (*DataStream, error) {
	m := new(DataStream)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *dHTClient) RestoreRPC(ctx context.Context, opts ...grpc.CallOption) (DHT_RestoreRPCClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_DHT_serviceDesc.Streams[3], c.cc, "/chordstore.DHT/RestoreRPC", opts...)
	if err != nil {
		return nil, err
	}
	x := &dHTRestoreRPCClient{stream}
	return x, nil
}

type DHT_RestoreRPCClient interface {
	Send(*DataStream) error
	CloseAndRecv() (*chord.ErrResponse, error)
	grpc.ClientStream
}

type dHTRestoreRPCClient struct {
	grpc.ClientStream
}

func (x *dHTRestoreRPCClient) Send(m *DataStream) error {
	return x.ClientStream.SendMsg(m)
}

func (x *dHTRestoreRPCClient) CloseAndRecv() (*chord.ErrResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(chord.ErrResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for DHT service

type DHTServer interface {
	PutKeyRPC(context.Context, *DHTKeyValue) (*chord.ErrResponse, error)
	GetKeyRPC(context.Context, *DHTBytes) (*DHTBytesErr, error)
	UpdateKeyRPC(context.Context, *DHTHashKeyValue) (*chord.ErrResponse, error)
	RemoveKeyRPC(context.Context, *DHTBytes) (*chord.ErrResponse, error)
	PutObjectRPC(DHT_PutObjectRPCServer) error
	GetObjectRPC(*DHTBytes, DHT_GetObjectRPCServer) error
	RemoveObjectRPC(context.Context, *DHTBytes) (*chord.ErrResponse, error)
	SnapshotRPC(*chord.Vnode, DHT_SnapshotRPCServer) error
	RestoreRPC(DHT_RestoreRPCServer) error
}

func RegisterDHTServer(s *grpc.Server, srv DHTServer) {
	s.RegisterService(&_DHT_serviceDesc, srv)
}

func _DHT_PutKeyRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DHTKeyValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DHTServer).PutKeyRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chordstore.DHT/PutKeyRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DHTServer).PutKeyRPC(ctx, req.(*DHTKeyValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _DHT_GetKeyRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DHTBytes)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DHTServer).GetKeyRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chordstore.DHT/GetKeyRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DHTServer).GetKeyRPC(ctx, req.(*DHTBytes))
	}
	return interceptor(ctx, in, info, handler)
}

func _DHT_UpdateKeyRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DHTHashKeyValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DHTServer).UpdateKeyRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chordstore.DHT/UpdateKeyRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DHTServer).UpdateKeyRPC(ctx, req.(*DHTHashKeyValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _DHT_RemoveKeyRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DHTBytes)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DHTServer).RemoveKeyRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chordstore.DHT/RemoveKeyRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DHTServer).RemoveKeyRPC(ctx, req.(*DHTBytes))
	}
	return interceptor(ctx, in, info, handler)
}

func _DHT_PutObjectRPC_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DHTServer).PutObjectRPC(&dHTPutObjectRPCServer{stream})
}

type DHT_PutObjectRPCServer interface {
	SendAndClose(*chord.ErrResponse) error
	Recv() (*DataStream, error)
	grpc.ServerStream
}

type dHTPutObjectRPCServer struct {
	grpc.ServerStream
}

func (x *dHTPutObjectRPCServer) SendAndClose(m *chord.ErrResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dHTPutObjectRPCServer) Recv() (*DataStream, error) {
	m := new(DataStream)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _DHT_GetObjectRPC_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DHTBytes)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DHTServer).GetObjectRPC(m, &dHTGetObjectRPCServer{stream})
}

type DHT_GetObjectRPCServer interface {
	Send(*DataStream) error
	grpc.ServerStream
}

type dHTGetObjectRPCServer struct {
	grpc.ServerStream
}

func (x *dHTGetObjectRPCServer) Send(m *DataStream) error {
	return x.ServerStream.SendMsg(m)
}

func _DHT_RemoveObjectRPC_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DHTBytes)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DHTServer).RemoveObjectRPC(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chordstore.DHT/RemoveObjectRPC",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DHTServer).RemoveObjectRPC(ctx, req.(*DHTBytes))
	}
	return interceptor(ctx, in, info, handler)
}

func _DHT_SnapshotRPC_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(chord.Vnode)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DHTServer).SnapshotRPC(m, &dHTSnapshotRPCServer{stream})
}

type DHT_SnapshotRPCServer interface {
	Send(*DataStream) error
	grpc.ServerStream
}

type dHTSnapshotRPCServer struct {
	grpc.ServerStream
}

func (x *dHTSnapshotRPCServer) Send(m *DataStream) error {
	return x.ServerStream.SendMsg(m)
}

func _DHT_RestoreRPC_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DHTServer).RestoreRPC(&dHTRestoreRPCServer{stream})
}

type DHT_RestoreRPCServer interface {
	SendAndClose(*chord.ErrResponse) error
	Recv() (*DataStream, error)
	grpc.ServerStream
}

type dHTRestoreRPCServer struct {
	grpc.ServerStream
}

func (x *dHTRestoreRPCServer) SendAndClose(m *chord.ErrResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *dHTRestoreRPCServer) Recv() (*DataStream, error) {
	m := new(DataStream)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _DHT_serviceDesc = grpc.ServiceDesc{
	ServiceName: "chordstore.DHT",
	HandlerType: (*DHTServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PutKeyRPC",
			Handler:    _DHT_PutKeyRPC_Handler,
		},
		{
			MethodName: "GetKeyRPC",
			Handler:    _DHT_GetKeyRPC_Handler,
		},
		{
			MethodName: "UpdateKeyRPC",
			Handler:    _DHT_UpdateKeyRPC_Handler,
		},
		{
			MethodName: "RemoveKeyRPC",
			Handler:    _DHT_RemoveKeyRPC_Handler,
		},
		{
			MethodName: "RemoveObjectRPC",
			Handler:    _DHT_RemoveObjectRPC_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "PutObjectRPC",
			Handler:       _DHT_PutObjectRPC_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "GetObjectRPC",
			Handler:       _DHT_GetObjectRPC_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SnapshotRPC",
			Handler:       _DHT_SnapshotRPC_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "RestoreRPC",
			Handler:       _DHT_RestoreRPC_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "rpc.proto",
}

func init() { proto.RegisterFile("rpc.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 431 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x93, 0x5f, 0x6b, 0xdb, 0x30,
	0x10, 0xc0, 0xa3, 0x24, 0x1b, 0xf5, 0xc5, 0x90, 0x21, 0xca, 0x1a, 0xbc, 0x3d, 0x04, 0xc1, 0x20,
	0x2f, 0xb5, 0x47, 0x07, 0x1b, 0x94, 0x6e, 0xec, 0x8f, 0x43, 0x0d, 0x7d, 0x68, 0x70, 0xbc, 0xbe,
	0xcb, 0xf6, 0xad, 0xce, 0xb6, 0x58, 0x46, 0x96, 0x0d, 0xfe, 0x62, 0xfb, 0x7c, 0xc3, 0xb2, 0xd3,
	0xb9, 0xc5, 0x4d, 0x46, 0xdf, 0x64, 0xe9, 0x7e, 0xf7, 0xbb, 0xd3, 0xc9, 0x60, 0xc8, 0x2c, 0xb2,
	0x33, 0x29, 0x94, 0xa0, 0x10, 0x25, 0x42, 0xc6, 0xb9, 0x12, 0x12, 0xad, 0x37, 0xb7, 0x1b, 0x95,
	0x14, 0xa1, 0x1d, 0x89, 0xad, 0x83, 0xc5, 0x0f, 0x21, 0x37, 0xdc, 0xb9, 0x15, 0xa7, 0x3a, 0xc2,
	0x49, 0x51, 0x35, 0x08, 0x5b, 0xc3, 0xc4, 0xf5, 0x82, 0x2b, 0xac, 0x6e, 0xf8, 0xef, 0x02, 0xe9,
	0x6b, 0x18, 0x96, 0xe9, 0x8c, 0xcc, 0xc9, 0x62, 0x72, 0x66, 0xda, 0x3a, 0xd8, 0xbe, 0x49, 0x45,
	0x8c, 0xfe, 0xb0, 0x4c, 0xe9, 0x0b, 0x18, 0xfd, 0xc2, 0x6a, 0x36, 0x9c, 0x93, 0x85, 0xe9, 0xd7,
	0x4b, 0x7a, 0x0c, 0xcf, 0xca, 0x1a, 0x9c, 0x8d, 0xf4, 0x5e, 0xf3, 0xc1, 0x72, 0x98, 0xba, 0x5e,
	0xe0, 0xf1, 0x3c, 0xf9, 0xcf, 0xc4, 0x16, 0x1c, 0x65, 0x12, 0xcb, 0x9a, 0x68, 0xb3, 0xdf, 0x7d,
	0xef, 0xa4, 0xa3, 0x1e, 0xe9, 0xb8, 0x2b, 0x7d, 0x0f, 0x47, 0xae, 0x17, 0x7c, 0xad, 0x14, 0xe6,
	0x07, 0x6c, 0x26, 0x90, 0xb0, 0xd5, 0x90, 0x90, 0x9d, 0xea, 0x1b, 0xd0, 0xdc, 0x52, 0xca, 0xe6,
	0x90, 0xb4, 0x87, 0xb5, 0x1c, 0xa5, 0xd4, 0xc1, 0x86, 0x5f, 0x2f, 0x99, 0x03, 0xd3, 0x75, 0xca,
	0xb3, 0x3c, 0x11, 0xea, 0x3a, 0x53, 0x1b, 0x91, 0x1e, 0xb0, 0xb1, 0x39, 0x80, 0xcb, 0x15, 0x5f,
	0x2b, 0x89, 0x7c, 0x4b, 0x29, 0x8c, 0x63, 0xae, 0x78, 0x6b, 0xd0, 0xeb, 0xb3, 0x3f, 0x63, 0x18,
	0xb9, 0x5e, 0x40, 0xcf, 0xc1, 0x58, 0x15, 0xea, 0x0a, 0x2b, 0x7f, 0xf5, 0x8d, 0x9e, 0xd8, 0xff,
	0x86, 0x69, 0x77, 0x46, 0x64, 0xd1, 0xd6, 0xb0, 0x94, 0xd2, 0xc7, 0x3c, 0x13, 0x69, 0x8e, 0x6c,
	0x40, 0x2f, 0xc0, 0xb8, 0xc4, 0x1d, 0x7b, 0xfc, 0x80, 0xd5, 0xcd, 0x59, 0x27, 0x7d, 0xbb, 0x4b,
	0x29, 0xd9, 0x80, 0x7e, 0x01, 0xf3, 0x7b, 0x16, 0x73, 0x85, 0x6d, 0x82, 0x57, 0x0f, 0x42, 0xbb,
	0xa3, 0x7c, 0xa4, 0x80, 0x73, 0x30, 0x7d, 0xdc, 0x8a, 0x12, 0xf7, 0xd6, 0xd0, 0xcf, 0x7e, 0x02,
	0x73, 0x55, 0xa8, 0xeb, 0xf0, 0x27, 0x46, 0xaa, 0x66, 0x5f, 0xde, 0x63, 0xef, 0x2e, 0xaf, 0x9f,
	0x5e, 0x10, 0xfa, 0x19, 0xcc, 0x4b, 0xec, 0xf0, 0xfd, 0xee, 0x47, 0xb2, 0xb2, 0xc1, 0x5b, 0x42,
	0x3f, 0xc2, 0xb4, 0xa9, 0xfe, 0x50, 0x92, 0xfe, 0x06, 0x3e, 0xc0, 0x64, 0xf7, 0x28, 0x6a, 0xf4,
	0xde, 0x23, 0xd8, 0xeb, 0xbd, 0x00, 0xf0, 0x51, 0x9f, 0x3c, 0xa1, 0xef, 0xf0, 0xb9, 0xfe, 0x87,
	0xdf, 0xfd, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x4b, 0xc1, 0xfe, 0x56, 0x03, 0x04, 0x00, 0x00,
}
