// Package mcp provides Model Context Protocol server implementation
package mcp

import (
	"context"
	"fmt"
)

// Server represents an MCP server instance
type Server struct {
	readonly bool
}

// Option is a functional option for configuring the Server
type Option func(*Server)

// WithReadonly sets the server to read-only mode
func WithReadonly(readonly bool) Option {
	return func(s *Server) {
		s.readonly = readonly
	}
}

// NewServer creates a new MCP server with the given options
func NewServer(opts ...Option) *Server {
	s := &Server{
		readonly: true, // Default to read-only for safety
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// Start starts the MCP server
func (s *Server) Start(ctx context.Context) error {
	// TODO: Implement MCP server startup
	// This is a placeholder for the actual MCP server implementation

	mode := "read-only"
	if !s.readonly {
		mode = "read-write"
	}

	return fmt.Errorf("MCP server implementation coming soon (mode: %s)", mode)
}
