// Copyright 2025 Gleb Kaplinskiy. All Rights Reserved.
// See LICENSE for licensing terms.

package handlers

type handlerOptions struct {
	methodNames []string
}

// Option configures gRPC proxy.
type Option func(*handlerOptions)

// WithMethodNames configures list of method names to proxy for non-transparent handler.
func WithMethodNames(methodNames ...string) Option {
	return func(o *handlerOptions) {
		o.methodNames = append([]string(nil), methodNames...)
	}
}
