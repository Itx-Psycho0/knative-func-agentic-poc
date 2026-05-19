#!/usr/bin/env python3
"""
Simple MCP Client for Testing
Usage: python3 test-mcp-client.py
"""

import subprocess
import json
import sys

class MCPClient:
    def __init__(self, command):
        self.command = command
        
    def call(self, method, params=None):
        """Send a JSON-RPC request to the MCP server"""
        request = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": method,
        }
        if params:
            request["params"] = params
            
        # Start the MCP server process
        process = subprocess.Popen(
            self.command,
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        
        # Send request
        request_json = json.dumps(request) + "\n"
        stdout, stderr = process.communicate(input=request_json, timeout=5)
        
        # Parse response
        if stdout.strip():
            response = json.loads(stdout.strip().split('\n')[0])
            return response
        else:
            return {"error": "No response", "stderr": stderr}

def main():
    print("🤖 MCP Client - Knative Functions Agentic POC")
    print("=" * 60)
    print()
    
    # Initialize client
    client = MCPClient([
        "/home/psycho/Downloads/Knative_Ai-agent/bin/func-agentic",
        "mcp",
        "start"
    ])
    
    # Test 1: Initialize
    print("📝 Test 1: Initialize Server")
    response = client.call("initialize", {})
    if "result" in response:
        print(f"✅ Server: {response['result']['serverInfo']['name']} v{response['result']['serverInfo']['version']}")
        print(f"✅ Protocol: {response['result']['protocolVersion']}")
    else:
        print(f"❌ Error: {response}")
    print()
    
    # Test 2: List Tools
    print("📝 Test 2: List Available Tools")
    response = client.call("tools/list")
    if "result" in response:
        tools = response['result']['tools']
        print(f"✅ Found {len(tools)} tools:")
        for tool in tools:
            print(f"   • {tool['name']}: {tool['description']}")
    else:
        print(f"❌ Error: {response}")
    print()
    
    # Test 3: Check Prerequisites (Go)
    print("📝 Test 3: Check Prerequisites (Go)")
    response = client.call("tools/call", {
        "name": "check_prerequisites",
        "arguments": {
            "runtime": "go",
            "checkDocker": True,
            "checkCluster": True
        }
    })
    if "result" in response:
        text = response['result']['content'][0]['text']
        print(text)
    else:
        print(f"❌ Error: {response}")
    print()
    
    # Test 4: List Runtimes
    print("📝 Test 4: List Supported Runtimes")
    response = client.call("tools/call", {
        "name": "list_runtimes",
        "arguments": {}
    })
    if "result" in response:
        text = response['result']['content'][0]['text']
        print(text)
    else:
        print(f"❌ Error: {response}")
    print()
    
    print("=" * 60)
    print("✅ All tests completed!")
    print()
    print("💡 Next Steps:")
    print("   1. Install Cursor IDE or MCP Inspector")
    print("   2. Configure MCP server")
    print("   3. Use AI to interact with your tools")

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\n\n👋 Bye!")
        sys.exit(0)
    except Exception as e:
        print(f"\n❌ Error: {e}")
        sys.exit(1)
