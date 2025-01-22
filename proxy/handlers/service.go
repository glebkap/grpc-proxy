// Copyright 2025 Gleb Kaplinskiy. All Rights Reserved.
// See LICENSE for licensing terms.

package handlers

import (
	sd "github.com/glebkap/grpc-proxy/proxy/stream_director"
	"google.golang.org/grpc"
)

// RegisterService sets up a proxy handler for a particular gRPC service and method.
// The behaviour is the same as if you were registering a handler method, e.g. from a generated pb.go file.
func RegisterService(server *grpc.Server, director sd.StreamDirector, serviceName string, options ...Option) {
	h := &handler{
		director: director,
		options:  handlerOptions{},
	}

	for _, o := range options {
		o(&h.options)
	}

	fakeDesc := &grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*interface{})(nil),
	}

	for _, m := range h.options.methodNames {
		streamDesc := grpc.StreamDesc{
			StreamName:    m,
			Handler:       h.handler,
			ServerStreams: true,
			ClientStreams: true,
		}
		fakeDesc.Streams = append(fakeDesc.Streams, streamDesc)
	}
	server.RegisterService(fakeDesc, h)
}
