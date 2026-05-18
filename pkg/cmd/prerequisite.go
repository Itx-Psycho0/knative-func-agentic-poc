package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/Itx-Psycho0/knative-func-agentic-poc/pkg/prerequisite"
	"github.com/spf13/cobra"
)

// NewPrerequisiteCommand creates the prerequisite command
func NewPrerequisiteCommand() *cobra.Command {
	var runtime string
	var checkCluster bool
	var checkDocker bool
	var checkKnative bool

	cmd := &cobra.Command{
		Use:   "prerequisite",
		Short: "Check prerequisites for function development",
		Long: `Check that all prerequisites are met for developing and deploying
Knative Functions.

This command validates:
  - Kubernetes cluster connectivity
  - Docker availability
  - Runtime dependencies (Go, Python, Node.js, etc.)
  - Knative Serving installation (optional)

Example:
  func-agentic prerequisite --runtime go --check-cluster --check-knative`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPrerequisiteCheck(cmd.Context(), runtime, checkCluster, checkDocker, checkKnative)
		},
	}

	cmd.Flags().StringVar(&runtime, "runtime", "", "Target runtime (go, python, node, java, rust)")
	cmd.Flags().BoolVar(&checkCluster, "check-cluster", true, "Check Kubernetes cluster connectivity")
	cmd.Flags().BoolVar(&checkDocker, "check-docker", true, "Check Docker availability")
	cmd.Flags().BoolVar(&checkKnative, "check-knative", false, "Check Knative Serving installation")

	return cmd
}

func runPrerequisiteCheck(ctx context.Context, runtime string, checkCluster, checkDocker, checkKnative bool) error {
	fmt.Println("🔍 Checking prerequisites...")
	fmt.Println()

	checker := prerequisite.NewChecker(
		prerequisite.WithRuntime(runtime),
		prerequisite.WithClusterCheck(checkCluster),
		prerequisite.WithDockerCheck(checkDocker),
		prerequisite.WithKnativeCheck(checkKnative),
	)

	results, err := checker.CheckAll(ctx)
	if err != nil {
		return fmt.Errorf("prerequisite check failed: %w", err)
	}

	// Display results
	allPassed := true
	for _, result := range results {
		icon := "✅"
		if !result.Passed {
			icon = "❌"
			allPassed = false
		}

		fmt.Printf("%s %s: %s\n", icon, result.Name, result.Message)

		if !result.Passed && len(result.Suggestions) > 0 {
			fmt.Println("   Suggestions:")
			for _, suggestion := range result.Suggestions {
				fmt.Printf("   - %s\n", suggestion)
			}
		}
		fmt.Println()
	}

	if allPassed {
		fmt.Println("✨ All prerequisites met! You're ready to create functions.")
		return nil
	}

	fmt.Fprintln(os.Stderr, "⚠️  Some prerequisites are not met. Please address the issues above.")
	os.Exit(1)
	return nil
}
