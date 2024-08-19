/*
 * Copyright 2021 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package remote

import (
	"context"
	"net"
	"time"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/profiler"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/grpc"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/serviceinfo"
	"github.com/cloudwego/kitex/pkg/streaming"
	"github.com/cloudwego/kitex/pkg/unknownservice/service"
)

// Option is used to pack the inbound and outbound handlers.
type Option struct {
	Outbounds []OutboundHandler

	Inbounds []InboundHandler

	StreamingMetaHandlers []StreamingMetaHandler
}

// PrependBoundHandler adds a BoundHandler to the head.
func (o *Option) PrependBoundHandler(h BoundHandler) {
	switch v := h.(type) {
	case DuplexBoundHandler:
		o.Inbounds = append([]InboundHandler{v}, o.Inbounds...)
		o.Outbounds = append([]OutboundHandler{v}, o.Outbounds...)
	case InboundHandler:
		o.Inbounds = append([]InboundHandler{v}, o.Inbounds...)
	case OutboundHandler:
		o.Outbounds = append([]OutboundHandler{v}, o.Outbounds...)
	default:
		panic("invalid BoundHandler: must implement InboundHandler or OutboundHandler")
	}
}

// AppendBoundHandler adds a BoundHandler to the end.
func (o *Option) AppendBoundHandler(h BoundHandler) {
	switch v := h.(type) {
	case DuplexBoundHandler:
		o.Inbounds = append(o.Inbounds, v)
		o.Outbounds = append(o.Outbounds, v)
	case InboundHandler:
		o.Inbounds = append(o.Inbounds, v)
	case OutboundHandler:
		o.Outbounds = append(o.Outbounds, v)
	default:
		panic("invalid BoundHandler: must implement InboundHandler or OutboundHandler")
	}
}

// ServerOption contains option that is used to init the remote server.
type ServerOption struct {
	TargetSvcInfo *serviceinfo.ServiceInfo

	SvcSearcher ServiceSearcher

	TransServerFactory TransServerFactory

	SvrHandlerFactory ServerTransHandlerFactory

	Codec Codec

	PayloadCodec PayloadCodec

	// Listener is used to specify the server listener, which comes with higher priority than Address below.
	Listener net.Listener

	// Address is the listener addr
	Address net.Addr

	ReusePort bool

	// Duration that server waits for to allow any existing connection to be closed gracefully.
	ExitWaitTime time.Duration

	// Duration that server waits for after error occurs during connection accepting.
	AcceptFailedDelayTime time.Duration

	// Duration that the accepted connection waits for to read or write data.
	MaxConnectionIdleTime time.Duration

	ReadWriteTimeout time.Duration

	InitOrResetRPCInfoFunc func(rpcinfo.RPCInfo, net.Addr) rpcinfo.RPCInfo

	TracerCtl *rpcinfo.TraceController

	Profiler                 profiler.Profiler
	ProfilerTransInfoTagging TransInfoTagging
	ProfilerMessageTagging   MessageTagging

	GRPCCfg *grpc.ServerConfig

	GRPCUnknownServiceHandler func(ctx context.Context, method string, stream streaming.Stream) error

	UnknownServiceHandler service.UnknownServiceHandler

	Option

	// invoking chain with recv/send middlewares for streaming APIs
	RecvEndpoint endpoint.RecvEndpoint
	SendEndpoint endpoint.SendEndpoint

	// for thrift streaming, this is enabled by default
	// for grpc(protobuf) streaming, it's disabled by default, enable with server.WithCompatibleMiddlewareForUnary
	CompatibleMiddlewareForUnary bool
}

// ClientOption is used to init the remote client.
type ClientOption struct {
	SvcInfo *serviceinfo.ServiceInfo

	CliHandlerFactory ClientTransHandlerFactory

	Codec Codec

	PayloadCodec PayloadCodec

	ConnPool ConnPool

	Dialer Dialer

	Option

	EnableConnPoolReporter bool
}
