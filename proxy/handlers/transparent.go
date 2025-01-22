// Copyright 2025 Gleb Kaplinskiy. All Rights Reserved.
// See LICENSE for licensing terms.

package handlers

import (
	sd "github.com/glebkap/grpc-proxy/proxy/stream_director"
	"google.golang.org/grpc"
)


// TransparentHandler returns a handler that attempts to proxy all requests that are not registered in the server.
// The indented use here is as a transparent proxy, where the server doesn't know about the services implemented by the
// backends. It should be used as a `grpc.UnknownServiceHandler`.
func TransparentHandler(director sd.StreamDirector) grpc.StreamHandler {
	streamer := &handler{director: director}
	return streamer.handler
}
