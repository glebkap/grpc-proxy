// Copyright 2017 Michal Witkowski. All Rights Reserved.
// Copyright 2025 Gleb Kaplinskiy. All Rights Reserved.
// See LICENSE for licensing terms.

package streamdirector

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// StreamDirector returns a gRPC ClientConn to be used to forward the call to.
//
// The presence of the `Context` allows for rich filtering, e.g. based on Metadata (headers).
// If no handling is meant to be done, a `codes.NotImplemented` gRPC error should be returned.
//
// The context returned from this function should be the context for the *outgoing* (to backend) call. In case you want
// to forward any Metadata between the inbound request and outbound requests, you should do it manually. However, you
// *must* propagate the cancel function (`context.WithCancel`) of the inbound context to the one returned.
//
// It is worth noting that the StreamDirector will be fired *after* all server-side stream interceptors
// are invoked. So decisions around authorization, monitoring etc. are better to be handled there.
//
// See the rather rich example.
type StreamDirector func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error)

// DefaultDirector returns a very simple forwarding StreamDirector that forwards all
// calls.
func DefaultDirector(cc *grpc.ClientConn) StreamDirector {
	return func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		ctx = metadata.NewOutgoingContext(ctx, md.Copy())
		return ctx, cc, nil
	}
}
