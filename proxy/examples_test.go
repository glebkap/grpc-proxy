// Copyright 2017 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package proxy_test

import (
	"context"
	"log"
	"strings"

	"github.com/glebkap/grpc-proxy/proxy/codec"
	"github.com/glebkap/grpc-proxy/proxy/handlers"
	sd "github.com/glebkap/grpc-proxy/proxy/stream_director"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	director sd.StreamDirector
)

// ExampleRegisterService is a simple example of registering a service with the proxy.
func ExampleRegisterService() {
	// A gRPC server with the proxying codec enabled.
	server := grpc.NewServer(grpc.ForceServerCodecV2(codec.Codec()))

	// Register a TestService with 4 of its methods explicitly.
	handlers.RegisterService(server, director,
		"gk.testproto.TestService",
		handlers.WithMethodNames("PingEmpty", "Ping", "PingError", "PingList"),
	)
}

// ExampleTransparentHandler is an example of redirecting all requests to the proxy.
func ExampleTransparentHandler() {
	grpc.NewServer(
		grpc.ForceServerCodecV2(codec.Codec()),
		grpc.UnknownServiceHandler(handlers.TransparentHandler(director)),
	)
}

// Provides a simple example of a director that shields internal services and dials a staging or production backend.
// This is a *very naive* implementation that creates a new connection on every request. Consider using pooling.
func ExampleStreamDirector() {
	connGen := func(hostname string) (*grpc.ClientConn, error) {
		conn, err := grpc.NewClient(
			hostname,
			grpc.WithDefaultCallOptions(grpc.ForceCodecV2(codec.Codec())),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, err
		}

		return conn, nil
	}

	stagingConn, err := connGen("api-service.staging.svc.local")
	if err != nil {
		log.Fatal("failed to create staging backend:", err)
	}

	prodConn, err := connGen("api-service.prod.svc.local")
	if err != nil {
		log.Fatal("failed to create production backend:", err)
	}

	director = func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {

		// Make sure we never forward internal services.
		if strings.HasPrefix(fullMethodName, "/com.example.internal.") {
			return nil, nil, status.Errorf(codes.Unimplemented, "Unknown method")
		}

		md, ok := metadata.FromIncomingContext(ctx)
		// Copy the inbound metadata explicitly.
		outCtx := metadata.NewOutgoingContext(ctx, md.Copy())

		if ok {
			// Decide on which backend to dial
			if val, exists := md[":authority"]; exists && val[0] == "staging.api.example.com" {
				return outCtx, stagingConn, nil
			} else if val, exists := md[":authority"]; exists && val[0] == "api.example.com" {
				return outCtx, prodConn, nil
			}
		}

		return nil, nil, status.Errorf(codes.Unimplemented, "Unknown method")
	}
}
