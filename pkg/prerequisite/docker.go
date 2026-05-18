package prerequisite

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// DockerChecker checks Docker availability and configuration
type DockerChecker struct{}

// NewDockerChecker creates a new Docker checker
func NewDockerChecker() *DockerChecker {
	return &DockerChecker{}
}

// Name returns the checker name
func (d *DockerChecker) Name() string {
	return "Docker"
}

// Description returns the checker description
func (d *DockerChecker) Description() string {
	return "Checks if Docker is installed and running"
}

// Check performs the Docker prerequisite check
func (d *DockerChecker) Check(ctx context.Context) (*CheckResult, error) {
	// Check if docker command exists
	if _, err := exec.LookPath("docker"); err != nil {
		return &CheckResult{
			Name:     d.Name(),
			Passed:   false,
			Message:  "Docker is not installed",
			Severity: SeverityError,
			Suggestions: []string{
				"Install Docker Desktop from https://www.docker.com/products/docker-desktop",
				"Or install Docker Engine: https://docs.docker.com/engine/install/",
			},
		}, nil
	}

	// Check if Docker daemon is running
	cmd := exec.CommandContext(ctx, "docker", "info")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &CheckResult{
			Name:     d.Name(),
			Passed:   false,
			Message:  "Docker is installed but not running",
			Severity: SeverityError,
			Suggestions: []string{
				"Start Docker Desktop",
				"Or start Docker daemon: sudo systemctl start docker",
			},
		}, nil
	}

	// Check Docker version
	versionCmd := exec.CommandContext(ctx, "docker", "version", "--format", "{{.Server.Version}}")
	versionOutput, err := versionCmd.Output()
	if err != nil {
		return &CheckResult{
			Name:     d.Name(),
			Passed:   true,
			Message:  "Docker is running",
			Severity: SeverityInfo,
		}, nil
	}

	version := strings.TrimSpace(string(versionOutput))
	message := fmt.Sprintf("Docker is running (version %s)", version)

	// Check if we can access Docker socket
	if strings.Contains(string(output), "permission denied") {
		return &CheckResult{
			Name:     d.Name(),
			Passed:   false,
			Message:  "Docker is running but current user lacks permissions",
			Severity: SeverityError,
			Suggestions: []string{
				"Add your user to the docker group: sudo usermod -aG docker $USER",
				"Then log out and log back in",
				"Or run commands with sudo",
			},
		}, nil
	}

	return &CheckResult{
		Name:     d.Name(),
		Passed:   true,
		Message:  message,
		Severity: SeverityInfo,
	}, nil
}
