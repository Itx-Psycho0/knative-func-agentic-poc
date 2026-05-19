# 🧪 Testing Guide - Knative Functions Agentic POC

## Hinglish Explanation (Kya Hai Ye?)

Ye tool **AI agents** ko **Knative Functions** ke saath kaam karne deta hai. Jaise tum ChatGPT se baat karte ho, waise hi AI automatically functions deploy kar sakta hai.

### Kaise Kaam Karta Hai?

```
You → AI (Claude/ChatGPT) → MCP Server → Knative Functions
```

1. **Tum AI se bolte ho**: "Check karo Docker installed hai ya nahi"
2. **AI MCP server ko call karta hai**: `check_prerequisites`
3. **MCP server check karta hai** aur result deta hai
4. **AI tumhe batata hai**: "Haan bhai, Docker installed hai version 29.4.3"

---

## 📋 Phase 1: Basic CLI Testing (Abhi Karo)

### Test 1: Version Check
```bash
./bin/func-agentic version
```

**Expected Output:**
```
func-agentic v0.1.0-dev
Git Commit: unknown
Build Date: unknown

LFX 2026 Mentorship Proposal POC
```

### Test 2: Prerequisite Check (Go)
```bash
./bin/func-agentic prerequisite check --runtime go
```

**Expected Output:**
```
🔍 Checking prerequisites...

✅ Docker: Docker is running (version 29.4.3)
❌ Kubernetes: Kubernetes cluster checking not yet implemented
✅ Runtime (go): Go is installed: go version go1.26.1 linux/amd64
```

### Test 3: Prerequisite Check (Python)
```bash
./bin/func-agentic prerequisite check --runtime python
```

### Test 4: Prerequisite Check (Node.js)
```bash
./bin/func-agentic prerequisite check --runtime node
```

### Test 5: MCP Server Help
```bash
./bin/func-agentic mcp start --help
```

---

## 🤖 Phase 2: MCP Server Testing (Manual)

### Start MCP Server
```bash
./bin/func-agentic mcp start
```

**Expected Output:**
```
🚀 MCP Server starting in read-only mode...
📡 Listening for MCP requests on stdin/stdout
🔧 Available tools: check_prerequisites, list_runtimes
💡 Connect this server to Claude Desktop, Cursor, or any MCP client
```

### Test MCP Protocol (Manual JSON-RPC)

**Terminal 1:** Start server
```bash
./bin/func-agentic mcp start
```

**Terminal 2:** Send test requests
```bash
# Test 1: Initialize
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./bin/func-agentic mcp start

# Test 2: List tools
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' | ./bin/func-agentic mcp start

# Test 3: Check prerequisites
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"check_prerequisites","arguments":{"runtime":"go"}}}' | ./bin/func-agentic mcp start

# Test 4: List runtimes
echo '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"list_runtimes","arguments":{}}}' | ./bin/func-agentic mcp start
```

---

## 🎯 Phase 3: Connect with AI (Claude Desktop)

### Step 1: Create MCP Config

**File:** `~/.config/Claude/claude_desktop_config.json`

```json
{
  "mcpServers": {
    "knative-func-agentic": {
      "command": "/home/psycho/Downloads/Knative_Ai-agent/bin/func-agentic",
      "args": ["mcp", "start"],
      "env": {
        "FUNC_ENABLE_MCP_WRITE": "false"
      }
    }
  }
}
```

### Step 2: Restart Claude Desktop

Close and reopen Claude Desktop app.

### Step 3: Test with AI

**Prompt to Claude:**
```
Can you check if I have the prerequisites installed for Knative Functions development with Go?
```

**Claude will:**
1. See the `check_prerequisites` tool
2. Call it with `runtime: "go"`
3. Get the results
4. Tell you in natural language

**Example Conversation:**
```
You: Check if Docker is installed

Claude: Let me check that for you.
[Calls check_prerequisites tool]

Claude: Yes! Docker is installed and running. You have version 29.4.3.
```

---

## 🔌 Phase 4: Connect with Cursor IDE

### Step 1: Create MCP Config

**File:** `~/.cursor/mcp.json`

```json
{
  "mcpServers": {
    "knative-func": {
      "command": "/home/psycho/Downloads/Knative_Ai-agent/bin/func-agentic",
      "args": ["mcp", "start"]
    }
  }
}
```

### Step 2: Restart Cursor

### Step 3: Use in Cursor Chat

```
You: @knative-func check prerequisites for Python development
```

---

## 🧪 Phase 5: Automated Testing Script

Create a test script:

```bash
#!/bin/bash
# test-mcp.sh

echo "🧪 Testing MCP Server..."

# Start server in background
./bin/func-agentic mcp start > /tmp/mcp-output.log 2>&1 &
MCP_PID=$!

sleep 2

# Test 1: Initialize
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | nc localhost 3000

# Test 2: List tools
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' | nc localhost 3000

# Cleanup
kill $MCP_PID

echo "✅ Tests complete!"
```

---

## 📊 Available MCP Tools

### 1. `check_prerequisites`

**Description:** Check if prerequisites are met for Knative Functions development

**Parameters:**
- `runtime` (string): Programming language (go, python, node, java, rust)
- `checkCluster` (boolean): Check Kubernetes cluster (default: true)
- `checkDocker` (boolean): Check Docker (default: true)

**Example:**
```json
{
  "name": "check_prerequisites",
  "arguments": {
    "runtime": "go",
    "checkDocker": true,
    "checkCluster": true
  }
}
```

### 2. `list_runtimes`

**Description:** List all supported programming language runtimes

**Parameters:** None

**Example:**
```json
{
  "name": "list_runtimes",
  "arguments": {}
}
```

---

## 🎬 Demo Video Script

### What to Show:

1. **CLI Demo (2 min)**
   ```bash
   ./bin/func-agentic version
   ./bin/func-agentic prerequisite check --runtime go
   ./bin/func-agentic prerequisite check --runtime python
   ```

2. **MCP Server Demo (2 min)**
   ```bash
   ./bin/func-agentic mcp start
   # Show it's listening
   # Send a test JSON-RPC request
   ```

3. **AI Integration Demo (3 min)**
   - Open Claude Desktop
   - Show MCP config
   - Ask Claude: "Check my prerequisites for Go development"
   - Show Claude using the tool
   - Show the response

4. **Explain the Vision (2 min)**
   - "Right now it checks prerequisites"
   - "In future, AI can deploy functions, test them, debug them"
   - "Teaching + Automation = Perfect for LFX mentorship"

---

## 🐛 Troubleshooting

### MCP Server Not Starting
```bash
# Check if binary exists
ls -la bin/func-agentic

# Check permissions
chmod +x bin/func-agentic

# Run with verbose output
./bin/func-agentic mcp start 2>&1 | tee mcp.log
```

### Claude Not Seeing Tools
1. Check config file location: `~/.config/Claude/claude_desktop_config.json`
2. Check binary path is absolute
3. Restart Claude Desktop completely
4. Check Claude logs: `~/.config/Claude/logs/`

### Tool Calls Failing
```bash
# Test manually first
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | ./bin/func-agentic mcp start

# Check stderr for errors
./bin/func-agentic mcp start 2> errors.log
```

---

## 📝 Next Steps for Full Implementation

### Phase 1 (Current - POC) ✅
- [x] CLI framework
- [x] Prerequisite checking
- [x] MCP server skeleton
- [x] Basic tools (check_prerequisites, list_runtimes)

### Phase 2 (LFX Mentorship)
- [ ] Full Kubernetes integration
- [ ] Function creation tool
- [ ] Function deployment tool
- [ ] CI/CD pipeline generation
- [ ] Testing and debugging tools

### Phase 3 (Advanced)
- [ ] Multi-cluster support
- [ ] GitOps integration
- [ ] Observability tools
- [ ] Cost optimization suggestions

---

## 🎓 Educational Value (For LFX Proposal)

### What Mentees Will Learn:

1. **Go Programming**
   - CLI development with Cobra
   - Context handling
   - Error handling patterns
   - Testing strategies

2. **Kubernetes & Knative**
   - Function deployment
   - Service mesh concepts
   - Serverless architecture
   - Resource management

3. **AI Integration**
   - Model Context Protocol (MCP)
   - JSON-RPC communication
   - Tool design for AI agents
   - Prompt engineering

4. **DevOps**
   - CI/CD pipelines
   - GitHub Actions
   - Container orchestration
   - GitOps workflows

---

## 📧 Questions?

For LFX proposal or technical questions:
- GitHub: https://github.com/Itx-Psycho0/knative-func-agentic-poc
- Issue: https://github.com/knative/func/issues/3646
