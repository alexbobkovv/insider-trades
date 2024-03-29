// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: trade.proto

package api

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

// TradeServiceClient is the client API for TradeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TradeServiceClient interface {
	ListTrades(ctx context.Context, in *TradeRequest, opts ...grpc.CallOption) (TradeService_ListTradesClient, error)
	ListTransactions(ctx context.Context, in *TradeRequest, opts ...grpc.CallOption) (TradeService_ListTransactionsClient, error)
	ListViews(ctx context.Context, in *TradeViewRequest, opts ...grpc.CallOption) (TradeService_ListViewsClient, error)
}

type tradeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTradeServiceClient(cc grpc.ClientConnInterface) TradeServiceClient {
	return &tradeServiceClient{cc}
}

func (c *tradeServiceClient) ListTrades(ctx context.Context, in *TradeRequest, opts ...grpc.CallOption) (TradeService_ListTradesClient, error) {
	stream, err := c.cc.NewStream(ctx, &TradeService_ServiceDesc.Streams[0], "/api.TradeService/ListTrades", opts...)
	if err != nil {
		return nil, err
	}
	x := &tradeServiceListTradesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TradeService_ListTradesClient interface {
	Recv() (*Trade, error)
	grpc.ClientStream
}

type tradeServiceListTradesClient struct {
	grpc.ClientStream
}

func (x *tradeServiceListTradesClient) Recv() (*Trade, error) {
	m := new(Trade)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *tradeServiceClient) ListTransactions(ctx context.Context, in *TradeRequest, opts ...grpc.CallOption) (TradeService_ListTransactionsClient, error) {
	stream, err := c.cc.NewStream(ctx, &TradeService_ServiceDesc.Streams[1], "/api.TradeService/ListTransactions", opts...)
	if err != nil {
		return nil, err
	}
	x := &tradeServiceListTransactionsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TradeService_ListTransactionsClient interface {
	Recv() (*Transaction, error)
	grpc.ClientStream
}

type tradeServiceListTransactionsClient struct {
	grpc.ClientStream
}

func (x *tradeServiceListTransactionsClient) Recv() (*Transaction, error) {
	m := new(Transaction)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *tradeServiceClient) ListViews(ctx context.Context, in *TradeViewRequest, opts ...grpc.CallOption) (TradeService_ListViewsClient, error) {
	stream, err := c.cc.NewStream(ctx, &TradeService_ServiceDesc.Streams[2], "/api.TradeService/ListViews", opts...)
	if err != nil {
		return nil, err
	}
	x := &tradeServiceListViewsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TradeService_ListViewsClient interface {
	Recv() (*TradeViewResponse, error)
	grpc.ClientStream
}

type tradeServiceListViewsClient struct {
	grpc.ClientStream
}

func (x *tradeServiceListViewsClient) Recv() (*TradeViewResponse, error) {
	m := new(TradeViewResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TradeServiceServer is the server API for TradeService service.
// All implementations must embed UnimplementedTradeServiceServer
// for forward compatibility
type TradeServiceServer interface {
	ListTrades(*TradeRequest, TradeService_ListTradesServer) error
	ListTransactions(*TradeRequest, TradeService_ListTransactionsServer) error
	ListViews(*TradeViewRequest, TradeService_ListViewsServer) error
	mustEmbedUnimplementedTradeServiceServer()
}

// UnimplementedTradeServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTradeServiceServer struct {
}

func (UnimplementedTradeServiceServer) ListTrades(*TradeRequest, TradeService_ListTradesServer) error {
	return status.Errorf(codes.Unimplemented, "method ListTrades not implemented")
}
func (UnimplementedTradeServiceServer) ListTransactions(*TradeRequest, TradeService_ListTransactionsServer) error {
	return status.Errorf(codes.Unimplemented, "method ListTransactions not implemented")
}
func (UnimplementedTradeServiceServer) ListViews(*TradeViewRequest, TradeService_ListViewsServer) error {
	return status.Errorf(codes.Unimplemented, "method ListViews not implemented")
}
func (UnimplementedTradeServiceServer) mustEmbedUnimplementedTradeServiceServer() {}

// UnsafeTradeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TradeServiceServer will
// result in compilation errors.
type UnsafeTradeServiceServer interface {
	mustEmbedUnimplementedTradeServiceServer()
}

func RegisterTradeServiceServer(s grpc.ServiceRegistrar, srv TradeServiceServer) {
	s.RegisterService(&TradeService_ServiceDesc, srv)
}

func _TradeService_ListTrades_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TradeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TradeServiceServer).ListTrades(m, &tradeServiceListTradesServer{stream})
}

type TradeService_ListTradesServer interface {
	Send(*Trade) error
	grpc.ServerStream
}

type tradeServiceListTradesServer struct {
	grpc.ServerStream
}

func (x *tradeServiceListTradesServer) Send(m *Trade) error {
	return x.ServerStream.SendMsg(m)
}

func _TradeService_ListTransactions_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TradeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TradeServiceServer).ListTransactions(m, &tradeServiceListTransactionsServer{stream})
}

type TradeService_ListTransactionsServer interface {
	Send(*Transaction) error
	grpc.ServerStream
}

type tradeServiceListTransactionsServer struct {
	grpc.ServerStream
}

func (x *tradeServiceListTransactionsServer) Send(m *Transaction) error {
	return x.ServerStream.SendMsg(m)
}

func _TradeService_ListViews_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TradeViewRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TradeServiceServer).ListViews(m, &tradeServiceListViewsServer{stream})
}

type TradeService_ListViewsServer interface {
	Send(*TradeViewResponse) error
	grpc.ServerStream
}

type tradeServiceListViewsServer struct {
	grpc.ServerStream
}

func (x *tradeServiceListViewsServer) Send(m *TradeViewResponse) error {
	return x.ServerStream.SendMsg(m)
}

// TradeService_ServiceDesc is the grpc.ServiceDesc for TradeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TradeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.TradeService",
	HandlerType: (*TradeServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListTrades",
			Handler:       _TradeService_ListTrades_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ListTransactions",
			Handler:       _TradeService_ListTransactions_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ListViews",
			Handler:       _TradeService_ListViews_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "trade.proto",
}
