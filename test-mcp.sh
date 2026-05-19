#!/bin/bash
# Quick MCP Server Test Script

echo "🧪 Testing MCP Server..."
echo ""

# Test 1: Initialize
echo "📝 Test 1: Initialize"
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./bin/func-agentic mcp start 2>/dev/null | head -1
echo ""

# Test 2: List tools
echo "📝 Test 2: List Tools"
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' | ./bin/func-agentic mcp start 2>/dev/null | head -1 | jq '.result.tools[].name' 2>/dev/null || echo "Install jq for pretty output"
echo ""

# Test 3: Check prerequisites (Go)
echo "📝 Test 3: Check Prerequisites (Go)"
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"check_prerequisites","arguments":{"runtime":"go"}}}' | ./bin/func-agentic mcp start 2>/dev/null | head -1
echo ""

# Test 4: List runtimes
echo "📝 Test 4: List Runtimes"
echo '{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"list_runtimes","arguments":{}}}' | ./bin/func-agentic mcp start 2>/dev/null | head -1
echo ""

echo "✅ All tests complete!"
echo ""
echo "💡 To use with AI:"
echo "   1. Add to Claude Desktop config: ~/.config/Claude/claude_desktop_config.json"
echo "   2. Restart Claude Desktop"
echo "   3. Ask: 'Check my prerequisites for Go development'"
