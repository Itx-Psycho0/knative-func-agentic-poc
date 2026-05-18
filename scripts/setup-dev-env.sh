#!/bin/bash

# Development Environment Setup Script
# For Knative Functions Agentic Workflow POC

set -e

echo "🚀 Setting up development environment..."
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check Go installation
echo -e "${BLUE}Checking Go installation...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed${NC}"
    echo "Please install Go 1.21 or later from https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "${GREEN}✓ Go ${GO_VERSION} installed${NC}"
echo ""

# Check Docker installation
echo -e "${BLUE}Checking Docker installation...${NC}"
if ! command -v docker &> /dev/null; then
    echo -e "${YELLOW}⚠ Docker is not installed${NC}"
    echo "Docker is recommended for building and testing functions"
    echo "Install from https://www.docker.com/products/docker-desktop"
else
    DOCKER_VERSION=$(docker --version | awk '{print $3}' | sed 's/,//')
    echo -e "${GREEN}✓ Docker ${DOCKER_VERSION} installed${NC}"
fi
echo ""

# Check kubectl installation
echo -e "${BLUE}Checking kubectl installation...${NC}"
if ! command -v kubectl &> /dev/null; then
    echo -e "${YELLOW}⚠ kubectl is not installed${NC}"
    echo "kubectl is needed for Kubernetes cluster interaction"
    echo "Install from https://kubernetes.io/docs/tasks/tools/"
else
    KUBECTL_VERSION=$(kubectl version --client --short 2>/dev/null | awk '{print $3}')
    echo -e "${GREEN}✓ kubectl ${KUBECTL_VERSION} installed${NC}"
fi
echo ""

# Install Go development tools
echo -e "${BLUE}Installing Go development tools...${NC}"

# golangci-lint
if ! command -v golangci-lint &> /dev/null; then
    echo "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    echo -e "${GREEN}✓ golangci-lint installed${NC}"
else
    echo -e "${GREEN}✓ golangci-lint already installed${NC}"
fi

# goimports
if ! command -v goimports &> /dev/null; then
    echo "Installing goimports..."
    go install golang.org/x/tools/cmd/goimports@latest
    echo -e "${GREEN}✓ goimports installed${NC}"
else
    echo -e "${GREEN}✓ goimports already installed${NC}"
fi

# godoc
if ! command -v godoc &> /dev/null; then
    echo "Installing godoc..."
    go install golang.org/x/tools/cmd/godoc@latest
    echo -e "${GREEN}✓ godoc installed${NC}"
else
    echo -e "${GREEN}✓ godoc already installed${NC}"
fi

echo ""

# Download Go dependencies
echo -e "${BLUE}Downloading Go dependencies...${NC}"
go mod download
go mod tidy
echo -e "${GREEN}✓ Dependencies downloaded${NC}"
echo ""

# Create necessary directories
echo -e "${BLUE}Creating project directories...${NC}"
mkdir -p bin
mkdir -p docs
mkdir -p examples
mkdir -p tests/integration
mkdir -p tests/e2e
mkdir -p skills
echo -e "${GREEN}✓ Directories created${NC}"
echo ""

# Run initial build
echo -e "${BLUE}Running initial build...${NC}"
make build
echo -e "${GREEN}✓ Build successful${NC}"
echo ""

# Run tests
echo -e "${BLUE}Running tests...${NC}"
make test
echo -e "${GREEN}✓ Tests passed${NC}"
echo ""

# Summary
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✨ Development environment setup complete!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo "Next steps:"
echo "  1. Review the documentation:"
echo "     - README.md"
echo "     - LFX_PROPOSAL_ANALYSIS.md"
echo "     - IMPLEMENTATION_PLAN.md"
echo ""
echo "  2. Start development:"
echo "     make build    # Build the project"
echo "     make test     # Run tests"
echo "     make lint     # Run linter"
echo "     make help     # See all available commands"
echo ""
echo "  3. Run the CLI:"
echo "     ./bin/func-agentic --help"
echo "     ./bin/func-agentic prerequisite --help"
echo "     ./bin/func-agentic mcp --help"
echo ""
echo -e "${BLUE}Happy coding! 🚀${NC}"
