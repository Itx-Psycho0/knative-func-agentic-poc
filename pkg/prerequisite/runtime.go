package prerequisite

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// RuntimeChecker checks if the specified runtime is installed
type RuntimeChecker struct {
	runtime string
}

// NewRuntimeChecker creates a new runtime checker
func NewRuntimeChecker(runtime string) *RuntimeChecker {
	return &RuntimeChecker{
		runtime: runtime,
	}
}

// Name returns the checker name
func (r *RuntimeChecker) Name() string {
	return fmt.Sprintf("Runtime (%s)", r.runtime)
}

// Description returns the checker description
func (r *RuntimeChecker) Description() string {
	return fmt.Sprintf("Checks if %s runtime is installed and available", r.runtime)
}

// Check performs the runtime prerequisite check
func (r *RuntimeChecker) Check(ctx context.Context) (*CheckResult, error) {
	switch strings.ToLower(r.runtime) {
	case "go":
		return r.checkGo(ctx)
	case "python":
		return r.checkPython(ctx)
	case "node", "nodejs":
		return r.checkNode(ctx)
	case "java":
		return r.checkJava(ctx)
	case "rust":
		return r.checkRust(ctx)
	default:
		return &CheckResult{
			Name:     r.Name(),
			Passed:   false,
			Message:  fmt.Sprintf("Unknown runtime: %s", r.runtime),
			Severity: SeverityError,
			Suggestions: []string{
				"Supported runtimes: go, python, node, java, rust",
			},
		}, nil
	}
}

func (r *RuntimeChecker) checkGo(ctx context.Context) (*CheckResult, error) {
	if _, err := exec.LookPath("go"); err != nil {
		return &CheckResult{
			Name:     r.Name(),
			Passed:   false,
			Message:  "Go is not installed",
			Severity: SeverityError,
			Suggestions: []string{
				"Install Go from https://golang.org/dl/",
				"Recommended version: 1.21 or later",
			},
		}, nil
	}

	cmd := exec.CommandContext(ctx, "go", "version")
	output, err := cmd.Output()
	if err != nil {
		return &CheckResult{
			Name:     r.Name(),
			Passed:   false,
			Message:  "Go is installed but version check failed",
			Severity: SeverityWarning,
		}, nil
	}

	version := strings.TrimSpace(string(output))
	return &CheckResult{
		Name:     r.Name(),
		Passed:   true,
		Message:  fmt.Sprintf("Go is installed: %s", version),
		Severity: SeverityInfo,
	}, nil
}

func (r *RuntimeChecker) checkPython(ctx context.Context) (*CheckResult, error) {
	// Try python3 first, then python
	pythonCmd := "python3"
	if _, err := exec.LookPath("python3"); err != nil {
		pythonCmd = "python"
		if _, err := exec.LookPath("python"); err != nil {
			return &CheckResult{
				Name:     r.Name(),
				Passed:   false,
				Message:  "Python is not installed",
				Severity: SeverityError,
				Suggestions: []string{
					"Install Python from https://www.python.org/downloads/",
					"Recommended version: 3.9 or later",
				},
			}, nil
		}
	}

	cmd := exec.CommandContext(ctx, pythonCmd, "--version")
	output, err := cmd.Output()
	if err != nil {
		return &CheckResult{
			Name:     r.Name(),
			Passed:   false,
			Message:  "Python is installed but version check failed",
			Severity: SeverityWarning,
		}, nil
	}

	version := strings.TrimSpace(string(output))
	return &CheckResult{
		Name:     r.Name(),
		Passed:   true,
		Message:  fmt.Sprintf("Python is installed: %s", version),
		Severity: SeverityInfo,
	}, nil
}

func (r *RuntimeChecker) checkNode(ctx context.Context) (*CheckResult, error) {
	if _, err := exec.LookPath("node"); err != nil {
		return &CheckResult{
			Name:     r.Name(),
			Passed:   false,
			Message:  "Node.js is not installed",
			Severity: SeverityError,
			Suggestions: []string{
				"Install Node.js from https://nodejs.org/",
				"Recommended version: 18 LTS or later",
			},
		}, nil
	}

	cmd := exec.CommandContext(ctx, "node", "--version")
	output, err := cmd.Output()
	if err != nil {
		return &CheckResult{
			Name:     r.Name(),
			Passed:   false,
			Message:  "Node.js is installed but version check failed",
			Severity: SeverityWarning,
		}, nil
	}

	version := strings.TrimSpace(string(output))
	return &CheckResult{
		Name:     r.Name(),
		Passed:   true,
		Message:  fmt.Sprintf("Node.js is installed: %s", version),
		Severity: SeverityInfo,
	}, nil
}

func (r *RuntimeChecker) checkJava(ctx context.Context) (*CheckResult, error) {
	if _, err := exec.LookPath("java"); err != nil {
		return &CheckResult{
			Name:     r.Name(),
			Passed:   false,
			Message:  "Java is not installed",
			Severity: SeverityError,
			Suggestions: []string{
				"Install Java from https://adoptium.net/",
				"Recommended version: Java 17 or later",
			},
		}, nil
	}

	cmd := exec.CommandContext(ctx, "java", "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &CheckResult{
			Name:     r.Name(),
			Passed:   false,
			Message:  "Java is installed but version check failed",
			Severity: SeverityWarning,
		}, nil
	}

	// Java version output goes to stderr, get first line
	lines := strings.Split(string(output), "\n")
	version := "unknown"
	if len(lines) > 0 {
		version = strings.TrimSpace(lines[0])
	}

	return &CheckResult{
		Name:     r.Name(),
		Passed:   true,
		Message:  fmt.Sprintf("Java is installed: %s", version),
		Severity: SeverityInfo,
	}, nil
}

func (r *RuntimeChecker) checkRust(ctx context.Context) (*CheckResult, error) {
	if _, err := exec.LookPath("rustc"); err != nil {
		return &CheckResult{
			Name:     r.Name(),
			Passed:   false,
			Message:  "Rust is not installed",
			Severity: SeverityError,
			Suggestions: []string{
				"Install Rust from https://rustup.rs/",
				"Run: curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh",
			},
		}, nil
	}

	cmd := exec.CommandContext(ctx, "rustc", "--version")
	output, err := cmd.Output()
	if err != nil {
		return &CheckResult{
			Name:     r.Name(),
			Passed:   false,
			Message:  "Rust is installed but version check failed",
			Severity: SeverityWarning,
		}, nil
	}

	version := strings.TrimSpace(string(output))
	return &CheckResult{
		Name:     r.Name(),
		Passed:   true,
		Message:  fmt.Sprintf("Rust is installed: %s", version),
		Severity: SeverityInfo,
	}, nil
}
