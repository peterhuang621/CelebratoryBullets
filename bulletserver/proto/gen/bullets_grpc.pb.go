// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: bullets.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	BulletService_DirectDrawBullets_FullMethodName = "/BulletService/DirectDrawBullets"
)

// BulletServiceClient is the client API for BulletService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BulletServiceClient interface {
	DirectDrawBullets(ctx context.Context, in *BulletList, opts ...grpc.CallOption) (*Ack, error)
}

type bulletServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBulletServiceClient(cc grpc.ClientConnInterface) BulletServiceClient {
	return &bulletServiceClient{cc}
}

func (c *bulletServiceClient) DirectDrawBullets(ctx context.Context, in *BulletList, opts ...grpc.CallOption) (*Ack, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Ack)
	err := c.cc.Invoke(ctx, BulletService_DirectDrawBullets_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BulletServiceServer is the server API for BulletService service.
// All implementations must embed UnimplementedBulletServiceServer
// for forward compatibility.
type BulletServiceServer interface {
	DirectDrawBullets(context.Context, *BulletList) (*Ack, error)
	mustEmbedUnimplementedBulletServiceServer()
}

// UnimplementedBulletServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedBulletServiceServer struct{}

func (UnimplementedBulletServiceServer) DirectDrawBullets(context.Context, *BulletList) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DirectDrawBullets not implemented")
}
func (UnimplementedBulletServiceServer) mustEmbedUnimplementedBulletServiceServer() {}
func (UnimplementedBulletServiceServer) testEmbeddedByValue()                       {}

// UnsafeBulletServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BulletServiceServer will
// result in compilation errors.
type UnsafeBulletServiceServer interface {
	mustEmbedUnimplementedBulletServiceServer()
}

func RegisterBulletServiceServer(s grpc.ServiceRegistrar, srv BulletServiceServer) {
	// If the following call pancis, it indicates UnimplementedBulletServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BulletService_ServiceDesc, srv)
}

func _BulletService_DirectDrawBullets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BulletList)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BulletServiceServer).DirectDrawBullets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BulletService_DirectDrawBullets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BulletServiceServer).DirectDrawBullets(ctx, req.(*BulletList))
	}
	return interceptor(ctx, in, info, handler)
}

// BulletService_ServiceDesc is the grpc.ServiceDesc for BulletService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BulletService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "BulletService",
	HandlerType: (*BulletServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DirectDrawBullets",
			Handler:    _BulletService_DirectDrawBullets_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "bullets.proto",
}
