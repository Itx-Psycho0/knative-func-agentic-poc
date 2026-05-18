package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/Itx-Psycho0/knative-func-agentic-poc/pkg/mcp"
	"github.com/spf13/cobra"
)

// NewMCPCommand creates the MCP command
func NewMCPCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Model Context Protocol (MCP) server operations",
		Long: `The MCP server enables AI agents to interact with Knative Functions
through the Model Context Protocol.

By default, the server operates in read-only mode. To enable write operations
(deploy, delete), set the environment variable FUNC_ENABLE_MCP_WRITE=true.

Example configuration for Claude Desktop:
  {
    "mcpServers": {
      "func-agentic": {
        "command": "func-agentic",
        "args": ["mcp", "start"]
      }
    }
  }`,
	}

	cmd.AddCommand(NewMCPStartCommand())

	return cmd
}

// NewMCPStartCommand creates the MCP start command
func NewMCPStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the MCP server",
		Long: `Start the Model Context Protocol (MCP) server.

This command is designed to be invoked by MCP clients (Claude, Cursor, VS Code, etc.)
and should not typically be run directly by users.

The server provides the following capabilities:
  - Prerequisite checking and validation
  - CI/CD pipeline generation
  - Workflow orchestration
  - Function lifecycle management

Write mode can be enabled by setting FUNC_ENABLE_MCP_WRITE=true.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPStart(cmd.Context())
		},
	}

	return cmd
}

func runMCPStart(ctx context.Context) error {
	// Check if write mode is enabled
	writeEnabled := false
	if val := os.Getenv("FUNC_ENABLE_MCP_WRITE"); val != "" {
		parsed, err := strconv.ParseBool(val)
		if err != nil {
			return fmt.Errorf("FUNC_ENABLE_MCP_WRITE must be a boolean (true/false, 1/0). Got: %q", val)
		}
		writeEnabled = parsed
	}

	// Create and configure MCP server
	server := mcp.NewServer(
		mcp.WithReadonly(!writeEnabled),
	)

	// Log startup information
	mode := "read-only"
	if writeEnabled {
		mode = "read-write"
	}
	fmt.Fprintf(os.Stderr, "Starting MCP server in %s mode...\n", mode)

	// Start the server
	if err := server.Start(ctx); err != nil {
		return fmt.Errorf("failed to start MCP server: %w", err)
	}

	return nil
}
