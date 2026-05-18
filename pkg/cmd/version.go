package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is the current version of the application
	Version = "v0.1.0-dev"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
	// BuildDate is the build date
	BuildDate = "unknown"
)

// NewVersionCommand creates the version command
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  `Print version information for func-agentic.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("func-agentic %s\n", Version)
			fmt.Printf("Git Commit: %s\n", GitCommit)
			fmt.Printf("Build Date: %s\n", BuildDate)
			fmt.Println()
			fmt.Println("LFX 2026 Mentorship Proposal POC")
			fmt.Println("Knative Functions: End-to-End Agentic Workflow")
		},
	}

	return cmd
}
