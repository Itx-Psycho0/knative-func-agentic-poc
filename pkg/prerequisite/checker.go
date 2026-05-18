// Package prerequisite provides environment validation for Knative Functions.
//
// This package includes checkers for:
// - Kubernetes cluster connectivity
// - Docker availability
// - Runtime dependencies (Go, Python, Node.js, Java, Rust)
// - Knative Serving installation
//
// Example usage:
//
//	checker := prerequisite.NewChecker(
//	    prerequisite.WithRuntime("go"),
//	    prerequisite.WithClusterCheck(true),
//	)
//	results, err := checker.CheckAll(ctx)
//	if err != nil {
//	    return err
//	}
//	for _, result := range results {
//	    fmt.Printf("%s: %s\n", result.Name, result.Message)
//	}
package prerequisite

import (
	"context"
	"fmt"
)

// Severity represents the severity level of a check result
type Severity int

const (
	// SeverityInfo indicates informational message
	SeverityInfo Severity = iota
	// SeverityWarning indicates a warning
	SeverityWarning
	// SeverityError indicates an error
	SeverityError
)

// CheckResult represents the result of a prerequisite check
type CheckResult struct {
	Name        string
	Passed      bool
	Message     string
	Suggestions []string
	Severity    Severity
}

// Checker interface defines a prerequisite checker
type Checker interface {
	Check(ctx context.Context) (*CheckResult, error)
	Name() string
	Description() string
}

// PrerequisiteChecker orchestrates multiple prerequisite checks
type PrerequisiteChecker struct {
	runtime      string
	checkCluster bool
	checkDocker  bool
	checkKnative bool
	checkers     []Checker
}

// Option is a functional option for PrerequisiteChecker
type Option func(*PrerequisiteChecker)

// WithRuntime sets the target runtime
func WithRuntime(runtime string) Option {
	return func(c *PrerequisiteChecker) {
		c.runtime = runtime
	}
}

// WithClusterCheck enables cluster checking
func WithClusterCheck(enabled bool) Option {
	return func(c *PrerequisiteChecker) {
		c.checkCluster = enabled
	}
}

// WithDockerCheck enables Docker checking
func WithDockerCheck(enabled bool) Option {
	return func(c *PrerequisiteChecker) {
		c.checkDocker = enabled
	}
}

// WithKnativeCheck enables Knative checking
func WithKnativeCheck(enabled bool) Option {
	return func(c *PrerequisiteChecker) {
		c.checkKnative = enabled
	}
}

// NewChecker creates a new PrerequisiteChecker
func NewChecker(opts ...Option) *PrerequisiteChecker {
	c := &PrerequisiteChecker{
		checkCluster: true,
		checkDocker:  true,
		checkKnative: false,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.initializeCheckers()
	return c
}

func (c *PrerequisiteChecker) initializeCheckers() {
	c.checkers = []Checker{}

	// Add Docker checker
	if c.checkDocker {
		c.checkers = append(c.checkers, NewDockerChecker())
	}

	// Add Kubernetes checker
	if c.checkCluster {
		c.checkers = append(c.checkers, NewKubernetesChecker(c.checkKnative))
	}

	// Add runtime checker
	if c.runtime != "" {
		c.checkers = append(c.checkers, NewRuntimeChecker(c.runtime))
	}
}

// CheckAll runs all configured prerequisite checks
func (c *PrerequisiteChecker) CheckAll(ctx context.Context) ([]*CheckResult, error) {
	results := make([]*CheckResult, 0, len(c.checkers))

	for _, checker := range c.checkers {
		result, err := checker.Check(ctx)
		if err != nil {
			return nil, fmt.Errorf("check %s failed: %w", checker.Name(), err)
		}
		results = append(results, result)
	}

	return results, nil
}

// GetSuggestions returns remediation suggestions for failed checks
func GetSuggestions(results []*CheckResult) []string {
	suggestions := []string{}

	for _, result := range results {
		if !result.Passed {
			suggestions = append(suggestions, result.Suggestions...)
		}
	}

	return suggestions
}

// AllPassed returns true if all checks passed
func AllPassed(results []*CheckResult) bool {
	for _, result := range results {
		if !result.Passed {
			return false
		}
	}
	return true
}
