package cmd

import (
	"github.com/spf13/cobra"
)

// NewRootCommand creates the root command for the func-agentic CLI
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "func-agentic",
		Short: "Knative Functions with End-to-End Agentic Workflow",
		Long: `func-agentic extends Knative Functions with comprehensive MCP server
capabilities for AI agents to autonomously manage the complete lifecycle
of serverless functions.

Features:
  - Comprehensive prerequisite checking
  - Multi-platform CI/CD integration
  - Guided workflow orchestration
  - GitOps support

This is a proof of concept for LFX 2026 Mentorship Program.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add subcommands
	cmd.AddCommand(NewMCPCommand())
	cmd.AddCommand(NewPrerequisiteCommand())
	cmd.AddCommand(NewVersionCommand())

	return cmd
}
