// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AnalysisInterfaceClient is the client API for AnalysisInterface service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AnalysisInterfaceClient interface {
	RunAnalysis(ctx context.Context, in *RunAnalysisRequest, opts ...grpc.CallOption) (*RunAnalysisReply, error)
	Top10Category(ctx context.Context, in *Top10CategoryRequest, opts ...grpc.CallOption) (*Top10Reply, error)
	Top10Area(ctx context.Context, in *Top10AreaRequest, opts ...grpc.CallOption) (*Top10Reply, error)
	Top10Keywords(ctx context.Context, in *Top10KeywordsRequest, opts ...grpc.CallOption) (*Top10Reply, error)
	EmotionDistribute(ctx context.Context, in *EmotionDistributeRequest, opts ...grpc.CallOption) (*EmotionDistributeReply, error)
	EmotionNotice(ctx context.Context, in *EmotionNoticeRequest, opts ...grpc.CallOption) (*EmotionNoticeReply, error)
	GetAnalysis(ctx context.Context, in *GetAnalysisRequest, opts ...grpc.CallOption) (*GetAnalysisReply, error)
	ListAnalysis(ctx context.Context, in *ListAnalysisRequest, opts ...grpc.CallOption) (*ListAnalysisReply, error)
}

type analysisInterfaceClient struct {
	cc grpc.ClientConnInterface
}

func NewAnalysisInterfaceClient(cc grpc.ClientConnInterface) AnalysisInterfaceClient {
	return &analysisInterfaceClient{cc}
}

func (c *analysisInterfaceClient) RunAnalysis(ctx context.Context, in *RunAnalysisRequest, opts ...grpc.CallOption) (*RunAnalysisReply, error) {
	out := new(RunAnalysisReply)
	err := c.cc.Invoke(ctx, "/api.analysis.v1.AnalysisInterface/RunAnalysis", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analysisInterfaceClient) Top10Category(ctx context.Context, in *Top10CategoryRequest, opts ...grpc.CallOption) (*Top10Reply, error) {
	out := new(Top10Reply)
	err := c.cc.Invoke(ctx, "/api.analysis.v1.AnalysisInterface/Top10Category", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analysisInterfaceClient) Top10Area(ctx context.Context, in *Top10AreaRequest, opts ...grpc.CallOption) (*Top10Reply, error) {
	out := new(Top10Reply)
	err := c.cc.Invoke(ctx, "/api.analysis.v1.AnalysisInterface/Top10Area", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analysisInterfaceClient) Top10Keywords(ctx context.Context, in *Top10KeywordsRequest, opts ...grpc.CallOption) (*Top10Reply, error) {
	out := new(Top10Reply)
	err := c.cc.Invoke(ctx, "/api.analysis.v1.AnalysisInterface/Top10Keywords", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analysisInterfaceClient) EmotionDistribute(ctx context.Context, in *EmotionDistributeRequest, opts ...grpc.CallOption) (*EmotionDistributeReply, error) {
	out := new(EmotionDistributeReply)
	err := c.cc.Invoke(ctx, "/api.analysis.v1.AnalysisInterface/EmotionDistribute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analysisInterfaceClient) EmotionNotice(ctx context.Context, in *EmotionNoticeRequest, opts ...grpc.CallOption) (*EmotionNoticeReply, error) {
	out := new(EmotionNoticeReply)
	err := c.cc.Invoke(ctx, "/api.analysis.v1.AnalysisInterface/EmotionNotice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analysisInterfaceClient) GetAnalysis(ctx context.Context, in *GetAnalysisRequest, opts ...grpc.CallOption) (*GetAnalysisReply, error) {
	out := new(GetAnalysisReply)
	err := c.cc.Invoke(ctx, "/api.analysis.v1.AnalysisInterface/GetAnalysis", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *analysisInterfaceClient) ListAnalysis(ctx context.Context, in *ListAnalysisRequest, opts ...grpc.CallOption) (*ListAnalysisReply, error) {
	out := new(ListAnalysisReply)
	err := c.cc.Invoke(ctx, "/api.analysis.v1.AnalysisInterface/ListAnalysis", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AnalysisInterfaceServer is the server API for AnalysisInterface service.
// All implementations must embed UnimplementedAnalysisInterfaceServer
// for forward compatibility
type AnalysisInterfaceServer interface {
	RunAnalysis(context.Context, *RunAnalysisRequest) (*RunAnalysisReply, error)
	Top10Category(context.Context, *Top10CategoryRequest) (*Top10Reply, error)
	Top10Area(context.Context, *Top10AreaRequest) (*Top10Reply, error)
	Top10Keywords(context.Context, *Top10KeywordsRequest) (*Top10Reply, error)
	EmotionDistribute(context.Context, *EmotionDistributeRequest) (*EmotionDistributeReply, error)
	EmotionNotice(context.Context, *EmotionNoticeRequest) (*EmotionNoticeReply, error)
	GetAnalysis(context.Context, *GetAnalysisRequest) (*GetAnalysisReply, error)
	ListAnalysis(context.Context, *ListAnalysisRequest) (*ListAnalysisReply, error)
	mustEmbedUnimplementedAnalysisInterfaceServer()
}

// UnimplementedAnalysisInterfaceServer must be embedded to have forward compatible implementations.
type UnimplementedAnalysisInterfaceServer struct {
}

func (UnimplementedAnalysisInterfaceServer) RunAnalysis(context.Context, *RunAnalysisRequest) (*RunAnalysisReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunAnalysis not implemented")
}
func (UnimplementedAnalysisInterfaceServer) Top10Category(context.Context, *Top10CategoryRequest) (*Top10Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Top10Category not implemented")
}
func (UnimplementedAnalysisInterfaceServer) Top10Area(context.Context, *Top10AreaRequest) (*Top10Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Top10Area not implemented")
}
func (UnimplementedAnalysisInterfaceServer) Top10Keywords(context.Context, *Top10KeywordsRequest) (*Top10Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Top10Keywords not implemented")
}
func (UnimplementedAnalysisInterfaceServer) EmotionDistribute(context.Context, *EmotionDistributeRequest) (*EmotionDistributeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmotionDistribute not implemented")
}
func (UnimplementedAnalysisInterfaceServer) EmotionNotice(context.Context, *EmotionNoticeRequest) (*EmotionNoticeReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EmotionNotice not implemented")
}
func (UnimplementedAnalysisInterfaceServer) GetAnalysis(context.Context, *GetAnalysisRequest) (*GetAnalysisReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAnalysis not implemented")
}
func (UnimplementedAnalysisInterfaceServer) ListAnalysis(context.Context, *ListAnalysisRequest) (*ListAnalysisReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAnalysis not implemented")
}
func (UnimplementedAnalysisInterfaceServer) mustEmbedUnimplementedAnalysisInterfaceServer() {}

// UnsafeAnalysisInterfaceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AnalysisInterfaceServer will
// result in compilation errors.
type UnsafeAnalysisInterfaceServer interface {
	mustEmbedUnimplementedAnalysisInterfaceServer()
}

func RegisterAnalysisInterfaceServer(s grpc.ServiceRegistrar, srv AnalysisInterfaceServer) {
	s.RegisterService(&AnalysisInterface_ServiceDesc, srv)
}

func _AnalysisInterface_RunAnalysis_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunAnalysisRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalysisInterfaceServer).RunAnalysis(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.analysis.v1.AnalysisInterface/RunAnalysis",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalysisInterfaceServer).RunAnalysis(ctx, req.(*RunAnalysisRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalysisInterface_Top10Category_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Top10CategoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalysisInterfaceServer).Top10Category(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.analysis.v1.AnalysisInterface/Top10Category",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalysisInterfaceServer).Top10Category(ctx, req.(*Top10CategoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalysisInterface_Top10Area_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Top10AreaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalysisInterfaceServer).Top10Area(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.analysis.v1.AnalysisInterface/Top10Area",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalysisInterfaceServer).Top10Area(ctx, req.(*Top10AreaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalysisInterface_Top10Keywords_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Top10KeywordsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalysisInterfaceServer).Top10Keywords(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.analysis.v1.AnalysisInterface/Top10Keywords",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalysisInterfaceServer).Top10Keywords(ctx, req.(*Top10KeywordsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalysisInterface_EmotionDistribute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmotionDistributeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalysisInterfaceServer).EmotionDistribute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.analysis.v1.AnalysisInterface/EmotionDistribute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalysisInterfaceServer).EmotionDistribute(ctx, req.(*EmotionDistributeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalysisInterface_EmotionNotice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmotionNoticeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalysisInterfaceServer).EmotionNotice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.analysis.v1.AnalysisInterface/EmotionNotice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalysisInterfaceServer).EmotionNotice(ctx, req.(*EmotionNoticeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalysisInterface_GetAnalysis_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAnalysisRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalysisInterfaceServer).GetAnalysis(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.analysis.v1.AnalysisInterface/GetAnalysis",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalysisInterfaceServer).GetAnalysis(ctx, req.(*GetAnalysisRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AnalysisInterface_ListAnalysis_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAnalysisRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AnalysisInterfaceServer).ListAnalysis(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.analysis.v1.AnalysisInterface/ListAnalysis",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AnalysisInterfaceServer).ListAnalysis(ctx, req.(*ListAnalysisRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AnalysisInterface_ServiceDesc is the grpc.ServiceDesc for AnalysisInterface service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AnalysisInterface_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.analysis.v1.AnalysisInterface",
	HandlerType: (*AnalysisInterfaceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RunAnalysis",
			Handler:    _AnalysisInterface_RunAnalysis_Handler,
		},
		{
			MethodName: "Top10Category",
			Handler:    _AnalysisInterface_Top10Category_Handler,
		},
		{
			MethodName: "Top10Area",
			Handler:    _AnalysisInterface_Top10Area_Handler,
		},
		{
			MethodName: "Top10Keywords",
			Handler:    _AnalysisInterface_Top10Keywords_Handler,
		},
		{
			MethodName: "EmotionDistribute",
			Handler:    _AnalysisInterface_EmotionDistribute_Handler,
		},
		{
			MethodName: "EmotionNotice",
			Handler:    _AnalysisInterface_EmotionNotice_Handler,
		},
		{
			MethodName: "GetAnalysis",
			Handler:    _AnalysisInterface_GetAnalysis_Handler,
		},
		{
			MethodName: "ListAnalysis",
			Handler:    _AnalysisInterface_ListAnalysis_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/analysis/v1/analysis_interface.proto",
}
