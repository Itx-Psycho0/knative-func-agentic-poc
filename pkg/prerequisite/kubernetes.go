package prerequisite

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// KubernetesChecker checks Kubernetes cluster availability and configuration
type KubernetesChecker struct {
	checkKnative bool
}

// NewKubernetesChecker creates a new Kubernetes checker
func NewKubernetesChecker(checkKnative bool) *KubernetesChecker {
	return &KubernetesChecker{
		checkKnative: checkKnative,
	}
}

// Name returns the checker name
func (k *KubernetesChecker) Name() string {
	return "Kubernetes"
}

// Description returns the checker description
func (k *KubernetesChecker) Description() string {
	return "Checks if Kubernetes cluster is accessible and configured"
}

// Check performs the Kubernetes prerequisite check
func (k *KubernetesChecker) Check(ctx context.Context) (*CheckResult, error) {
	// Check if kubectl is installed
	if _, err := exec.LookPath("kubectl"); err != nil {
		return &CheckResult{
			Name:     k.Name(),
			Passed:   false,
			Message:  "kubectl is not installed",
			Severity: SeverityError,
			Suggestions: []string{
				"Install kubectl: https://kubernetes.io/docs/tasks/tools/",
				"For local development, consider using minikube or kind",
			},
		}, nil
	}

	// Check if cluster is accessible
	cmd := exec.CommandContext(ctx, "kubectl", "cluster-info")
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Check if it's a config issue
		if strings.Contains(string(output), "connection refused") ||
			strings.Contains(string(output), "Unable to connect") {
			return &CheckResult{
				Name:     k.Name(),
				Passed:   false,
				Message:  "kubectl is installed but cannot connect to cluster",
				Severity: SeverityError,
				Suggestions: []string{
					"Start your Kubernetes cluster (minikube start, kind create cluster, etc.)",
					"Check your kubeconfig: kubectl config view",
					"Verify cluster context: kubectl config current-context",
				},
			}, nil
		}

		if strings.Contains(string(output), "no configuration") ||
			strings.Contains(string(output), "does not exist") {
			return &CheckResult{
				Name:     k.Name(),
				Passed:   false,
				Message:  "kubectl is installed but not configured",
				Severity: SeverityError,
				Suggestions: []string{
					"Configure kubectl with a cluster",
					"For local development: minikube start",
					"Or create a kind cluster: kind create cluster",
				},
			}, nil
		}

		return &CheckResult{
			Name:     k.Name(),
			Passed:   false,
			Message:  fmt.Sprintf("kubectl check failed: %v", err),
			Severity: SeverityError,
			Suggestions: []string{
				"Run 'kubectl cluster-info' to diagnose the issue",
			},
		}, nil
	}

	// Get cluster version
	versionCmd := exec.CommandContext(ctx, "kubectl", "version", "--short")
	versionOutput, _ := versionCmd.Output()
	versionInfo := strings.TrimSpace(string(versionOutput))

	// Get current context
	contextCmd := exec.CommandContext(ctx, "kubectl", "config", "current-context")
	contextOutput, _ := contextCmd.Output()
	currentContext := strings.TrimSpace(string(contextOutput))

	message := fmt.Sprintf("Kubernetes cluster is accessible (context: %s)", currentContext)
	if versionInfo != "" {
		// Extract just the server version line if available
		lines := strings.Split(versionInfo, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Server") {
				message = fmt.Sprintf("Kubernetes cluster is accessible - %s (context: %s)",
					strings.TrimSpace(line), currentContext)
				break
			}
		}
	}

	result := &CheckResult{
		Name:     k.Name(),
		Passed:   true,
		Message:  message,
		Severity: SeverityInfo,
	}

	// Check Knative if requested
	if k.checkKnative {
		knativeResult := k.checkKnativeServing(ctx)
		if !knativeResult.Passed {
			result.Suggestions = append(result.Suggestions, knativeResult.Suggestions...)
		}
	}

	return result, nil
}

// checkKnativeServing checks if Knative Serving is installed
func (k *KubernetesChecker) checkKnativeServing(ctx context.Context) *CheckResult {
	// Check if knative-serving namespace exists
	cmd := exec.CommandContext(ctx, "kubectl", "get", "namespace", "knative-serving")
	if err := cmd.Run(); err != nil {
		return &CheckResult{
			Name:     "Knative Serving",
			Passed:   false,
			Message:  "Knative Serving is not installed",
			Severity: SeverityWarning,
			Suggestions: []string{
				"Install Knative Serving: https://knative.dev/docs/install/",
				"Quick install: kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.12.0/serving-crds.yaml",
				"Then: kubectl apply -f https://github.com/knative/serving/releases/download/knative-v1.12.0/serving-core.yaml",
			},
		}
	}

	// Check if Knative Serving pods are running
	cmd = exec.CommandContext(ctx, "kubectl", "get", "pods", "-n", "knative-serving",
		"--field-selector=status.phase!=Running", "--no-headers")
	output, _ := cmd.Output()

	if len(strings.TrimSpace(string(output))) > 0 {
		return &CheckResult{
			Name:     "Knative Serving",
			Passed:   false,
			Message:  "Knative Serving is installed but some pods are not running",
			Severity: SeverityWarning,
			Suggestions: []string{
				"Check pod status: kubectl get pods -n knative-serving",
				"Check logs: kubectl logs -n knative-serving -l app=controller",
			},
		}
	}

	return &CheckResult{
		Name:     "Knative Serving",
		Passed:   true,
		Message:  "Knative Serving is installed and running",
		Severity: SeverityInfo,
	}
}
