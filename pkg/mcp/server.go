// Package mcp provides Model Context Protocol server implementation
package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Itx-Psycho0/knative-func-agentic-poc/pkg/prerequisite"
)

// Server represents an MCP server instance
type Server struct {
	readonly bool
	stdin    io.Reader
	stdout   io.Writer
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
		stdin:    os.Stdin,
		stdout:   os.Stdout,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// MCPRequest represents an incoming MCP request
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// MCPResponse represents an outgoing MCP response
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an MCP error
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ToolInfo represents information about an available tool
type ToolInfo struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

// Start starts the MCP server
func (s *Server) Start(ctx context.Context) error {
	mode := "read-only"
	if !s.readonly {
		mode = "read-write"
	}

	// Log to stderr (stdout is for MCP protocol)
	fmt.Fprintf(os.Stderr, "🚀 MCP Server starting in %s mode...\n", mode)
	fmt.Fprintf(os.Stderr, "📡 Listening for MCP requests on stdin/stdout\n")
	fmt.Fprintf(os.Stderr, "🔧 Available tools: check_prerequisites, list_runtimes\n")
	fmt.Fprintf(os.Stderr, "💡 Connect this server to Claude Desktop, Cursor, or any MCP client\n\n")

	scanner := bufio.NewScanner(s.stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req MCPRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.sendError(nil, -32700, "Parse error")
			continue
		}

		s.handleRequest(ctx, &req)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stdin: %w", err)
	}

	return nil
}

func (s *Server) handleRequest(ctx context.Context, req *MCPRequest) {
	switch req.Method {
	case "initialize":
		s.handleInitialize(req)
	case "tools/list":
		s.handleToolsList(req)
	case "tools/call":
		s.handleToolsCall(ctx, req)
	default:
		s.sendError(req.ID, -32601, fmt.Sprintf("Method not found: %s", req.Method))
	}
}

func (s *Server) handleInitialize(req *MCPRequest) {
	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"serverInfo": map[string]string{
			"name":    "func-agentic",
			"version": "0.1.0",
		},
		"capabilities": map[string]interface{}{
			"tools": map[string]bool{
				"listChanged": false,
			},
		},
	}
	s.sendResponse(req.ID, result)
}

func (s *Server) handleToolsList(req *MCPRequest) {
	tools := []ToolInfo{
		{
			Name:        "check_prerequisites",
			Description: "Check if prerequisites are met for Knative Functions development (Docker, Kubernetes, runtime)",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"runtime": map[string]interface{}{
						"type":        "string",
						"description": "Programming language runtime to check (go, python, node, java, rust)",
						"enum":        []string{"go", "python", "node", "java", "rust"},
					},
					"checkCluster": map[string]interface{}{
						"type":        "boolean",
						"description": "Whether to check Kubernetes cluster connectivity",
						"default":     true,
					},
					"checkDocker": map[string]interface{}{
						"type":        "boolean",
						"description": "Whether to check Docker availability",
						"default":     true,
					},
				},
			},
		},
		{
			Name:        "list_runtimes",
			Description: "List all supported programming language runtimes for Knative Functions",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}

	result := map[string]interface{}{
		"tools": tools,
	}
	s.sendResponse(req.ID, result)
}

func (s *Server) handleToolsCall(ctx context.Context, req *MCPRequest) {
	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.sendError(req.ID, -32602, "Invalid params")
		return
	}

	switch params.Name {
	case "check_prerequisites":
		s.handleCheckPrerequisites(ctx, req.ID, params.Arguments)
	case "list_runtimes":
		s.handleListRuntimes(req.ID)
	default:
		s.sendError(req.ID, -32601, fmt.Sprintf("Tool not found: %s", params.Name))
	}
}

func (s *Server) handleCheckPrerequisites(ctx context.Context, id interface{}, args map[string]interface{}) {
	runtime, _ := args["runtime"].(string)
	checkCluster := true
	checkDocker := true

	if val, ok := args["checkCluster"].(bool); ok {
		checkCluster = val
	}
	if val, ok := args["checkDocker"].(bool); ok {
		checkDocker = val
	}

	checker := prerequisite.NewChecker(
		prerequisite.WithRuntime(runtime),
		prerequisite.WithClusterCheck(checkCluster),
		prerequisite.WithDockerCheck(checkDocker),
	)

	results, err := checker.CheckAll(ctx)
	if err != nil {
		s.sendError(id, -32603, fmt.Sprintf("Check failed: %v", err))
		return
	}

	// Format results
	output := "🔍 Prerequisite Check Results:\n\n"
	allPassed := true
	for _, result := range results {
		if result.Passed {
			output += fmt.Sprintf("✅ %s: %s\n", result.Name, result.Message)
		} else {
			output += fmt.Sprintf("❌ %s: %s\n", result.Name, result.Message)
			allPassed = false
			if len(result.Suggestions) > 0 {
				output += "   Suggestions:\n"
				for _, suggestion := range result.Suggestions {
					output += fmt.Sprintf("   - %s\n", suggestion)
				}
			}
		}
	}

	if allPassed {
		output += "\n✨ All prerequisites are met! You're ready to develop Knative Functions."
	} else {
		output += "\n⚠️  Some prerequisites are not met. Please address the issues above."
	}

	result := map[string]interface{}{
		"content": []map[string]string{
			{
				"type": "text",
				"text": output,
			},
		},
	}
	s.sendResponse(id, result)
}

func (s *Server) handleListRuntimes(id interface{}) {
	runtimes := []map[string]string{
		{"name": "go", "description": "Go (Golang) - Compiled, statically typed language"},
		{"name": "python", "description": "Python - Dynamic, interpreted language"},
		{"name": "node", "description": "Node.js - JavaScript runtime"},
		{"name": "java", "description": "Java - Object-oriented, compiled language"},
		{"name": "rust", "description": "Rust - Systems programming language"},
	}

	output := "🚀 Supported Runtimes for Knative Functions:\n\n"
	for _, rt := range runtimes {
		output += fmt.Sprintf("• %s: %s\n", rt["name"], rt["description"])
	}

	result := map[string]interface{}{
		"content": []map[string]string{
			{
				"type": "text",
				"text": output,
			},
		},
	}
	s.sendResponse(id, result)
}

func (s *Server) sendResponse(id interface{}, result interface{}) {
	resp := MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	s.writeJSON(resp)
}

func (s *Server) sendError(id interface{}, code int, message string) {
	resp := MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
	}
	s.writeJSON(resp)
}

func (s *Server) writeJSON(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
		return
	}
	fmt.Fprintf(s.stdout, "%s\n", data)
}
