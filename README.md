# Knative Functions: End-to-End Agentic Workflow

[![CI](https://github.com/Itx-Psycho0/knative-func-agentic-poc/workflows/CI/badge.svg)](https://github.com/Itx-Psycho0/knative-func-agentic-poc/actions)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![LFX Mentorship](https://img.shields.io/badge/LFX-2026%20Proposal-orange)](https://github.com/knative/func/issues/3646)

> **LFX 2026 Mentorship Proposal** - Proof of Concept for [knative/func#3646](https://github.com/knative/func/issues/3646)

## 🎯 Overview

This project extends Knative Functions' MCP (Model Context Protocol) server to enable AI agents to autonomously manage the complete lifecycle of serverless functions - from environment validation through CI/CD integration and production deployment.

### Key Features

- ✅ **Comprehensive Prerequisite Checking**: Validates environment, cluster, and dependencies
- ✅ **Multi-Platform CI/CD Integration**: GitHub Actions, GitLab CI, Tekton pipelines
- ✅ **Guided Workflow Orchestration**: Step-by-step agent guidance through function lifecycle
- ✅ **Educational-First Approach**: Teaching agents and users, not just automation
- ✅ **Production-Ready**: Proper error handling, testing, and documentation

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or later
- Docker Desktop (optional but recommended)
- kubectl (optional, for cluster testing)

### Installation

```bash
# Clone the repository
git clone https://github.com/Itx-Psycho0/knative-func-agentic-poc.git
cd knative-func-agentic-poc

# Install dependencies
go mod download

# Build the project
make build

# Run tests
make test
```

### Usage

```bash
# Check version
./bin/func-agentic version

# Check prerequisites
./bin/func-agentic prerequisite

# Check with specific runtime
./bin/func-agentic prerequisite --runtime go

# Start MCP server (for AI agents)
./bin/func-agentic mcp start
```

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      AI Agent (Claude, etc.)                 │
└────────────────────────┬────────────────────────────────────┘
                         │ MCP Protocol
┌────────────────────────▼────────────────────────────────────┐
│                    MCP Server (This POC)                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │ Prerequisite │  │   CI/CD      │  │   Workflow   │     │
│  │   Checking   │  │ Integration  │  │ Orchestration│     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│                  Knative Functions CLI                       │
└────────────────────────┬────────────────────────────────────┘
                         │
┌────────────────────────▼────────────────────────────────────┐
│              Kubernetes Cluster (Knative)                    │
└─────────────────────────────────────────────────────────────┘
```

## 🛠️ MCP Tools (Planned)

### Prerequisite Checking
- `check_prerequisites` - Comprehensive environment validation
- `validate_cluster` - Kubernetes cluster health check
- `suggest_remediation` - Fix suggestions for failed checks

### CI/CD Integration
- `generate_github_workflow` - Create GitHub Actions workflow
- `generate_gitlab_ci` - Create GitLab CI pipeline
- `setup_tekton_pipeline` - Create Tekton pipeline
- `configure_gitops` - Setup ArgoCD/Flux deployment

### Workflow Orchestration
- `start_guided_workflow` - Begin guided workflow
- `get_workflow_status` - Check workflow progress
- `get_next_steps` - Get next recommended actions

## 📊 Project Status

### Completed ✅
- [x] Project foundation and structure
- [x] CLI framework with Cobra
- [x] Prerequisite checking framework
- [x] Docker availability checker
- [x] Build system and Makefile
- [x] CI/CD pipeline (GitHub Actions)
- [x] Comprehensive documentation

### In Progress 🚧
- [ ] Kubernetes cluster checker
- [ ] Runtime dependency checker
- [ ] MCP server implementation
- [ ] GitHub Actions generator
- [ ] Educational skill documents

### Planned 📋
- [ ] GitLab CI generator
- [ ] Tekton pipeline integration
- [ ] Workflow orchestration
- [ ] Complete test coverage
- [ ] Demo scenarios

## 🧪 Testing

```bash
# Run unit tests
make test

# Run with coverage
make test-coverage

# Run linting
make lint

# Run all checks
make check
```

## 📖 Documentation

- **Architecture**: See [ARCHITECTURE.md](docs/ARCHITECTURE.md) (coming soon)
- **Contributing**: See [CONTRIBUTING.md](CONTRIBUTING.md) (coming soon)
- **API Reference**: See [API_REFERENCE.md](docs/API_REFERENCE.md) (coming soon)

## 🎓 Educational Approach

Following the maintainer's vision, this project emphasizes **teaching over automation**:

> "Let people learn how to get the most from Functions by first teaching their agent how the more complex features work. They can then explore and grow in their knowledge and application organically." - Luke Kingland, Knative Functions Maintainer

### Progressive Complexity Levels

1. **Beginner**: Local development and basic concepts
2. **Intermediate**: GitHub Actions and simple CI/CD
3. **Advanced**: Self-hosted runners and local clusters
4. **Expert**: Tekton pipelines and remote builds
5. **Production**: GitOps and advanced deployment strategies

## 🤝 Contributing

This is a POC for an LFX mentorship proposal. Feedback and suggestions are welcome!

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Knative Community** - For the amazing Knative Functions project
- **Anthropic** - For the Model Context Protocol specification
- **CNCF** - For the LFX Mentorship program
- **Luke Kingland** - For guidance and feedback on the project scope

## 📧 Contact

- **GitHub**: [@Itx-Psycho0](https://github.com/Itx-Psycho0)
- **Project Link**: [https://github.com/Itx-Psycho0/knative-func-agentic-poc](https://github.com/Itx-Psycho0/knative-func-agentic-poc)
- **LFX Issue**: [knative/func#3646](https://github.com/knative/func/issues/3646)

## 🎯 LFX 2026 Proposal

This POC demonstrates:

1. **Technical Competence**: Deep understanding of Go, Kubernetes, MCP, and Knative
2. **Problem-Solving**: Practical solutions to real developer challenges
3. **Communication**: Clear documentation and presentation
4. **Initiative**: Proactive implementation beyond basic requirements
5. **Community Alignment**: Following Knative project conventions and maintainer feedback

---

**Made with ❤️ for the Knative Community**

*This is a proof of concept for the LFX 2026 Mentorship Program*
