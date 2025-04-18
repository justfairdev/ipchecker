// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: ipchecker.proto

package ipchecker

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
	IPChecker_CheckIP_FullMethodName = "/ipchecker.v1.IPChecker/CheckIP"
)

// IPCheckerClient is the client API for IPChecker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// IPChecker service for checking an IP against allowed countries.
type IPCheckerClient interface {
	// CheckIP returns whether the IP is in the allowed list.
	CheckIP(ctx context.Context, in *IPCheckRequest, opts ...grpc.CallOption) (*IPCheckResponse, error)
}

type iPCheckerClient struct {
	cc grpc.ClientConnInterface
}

func NewIPCheckerClient(cc grpc.ClientConnInterface) IPCheckerClient {
	return &iPCheckerClient{cc}
}

func (c *iPCheckerClient) CheckIP(ctx context.Context, in *IPCheckRequest, opts ...grpc.CallOption) (*IPCheckResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IPCheckResponse)
	err := c.cc.Invoke(ctx, IPChecker_CheckIP_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IPCheckerServer is the server API for IPChecker service.
// All implementations must embed UnimplementedIPCheckerServer
// for forward compatibility.
//
// IPChecker service for checking an IP against allowed countries.
type IPCheckerServer interface {
	// CheckIP returns whether the IP is in the allowed list.
	CheckIP(context.Context, *IPCheckRequest) (*IPCheckResponse, error)
	mustEmbedUnimplementedIPCheckerServer()
}

// UnimplementedIPCheckerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedIPCheckerServer struct{}

func (UnimplementedIPCheckerServer) CheckIP(context.Context, *IPCheckRequest) (*IPCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckIP not implemented")
}
func (UnimplementedIPCheckerServer) mustEmbedUnimplementedIPCheckerServer() {}
func (UnimplementedIPCheckerServer) testEmbeddedByValue()                   {}

// UnsafeIPCheckerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IPCheckerServer will
// result in compilation errors.
type UnsafeIPCheckerServer interface {
	mustEmbedUnimplementedIPCheckerServer()
}

func RegisterIPCheckerServer(s grpc.ServiceRegistrar, srv IPCheckerServer) {
	// If the following call pancis, it indicates UnimplementedIPCheckerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&IPChecker_ServiceDesc, srv)
}

func _IPChecker_CheckIP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IPCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IPCheckerServer).CheckIP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IPChecker_CheckIP_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IPCheckerServer).CheckIP(ctx, req.(*IPCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// IPChecker_ServiceDesc is the grpc.ServiceDesc for IPChecker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IPChecker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ipchecker.v1.IPChecker",
	HandlerType: (*IPCheckerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckIP",
			Handler:    _IPChecker_CheckIP_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ipchecker.proto",
}
