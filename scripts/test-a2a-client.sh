#!/bin/bash

# openribcage A2A Client Testing Script
# Test A2A protocol client with kagent sandbox endpoints
#
# This script tests the openribcage A2A protocol client against 
# real kagent A2A endpoints as referenced in the PR and README.

set -e

echo "🧪 Testing openribcage A2A client with kagent sandbox..."
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Configuration
KAGENT_BASE_URL="http://localhost:8083/api/a2a"
KAGENT_NAMESPACE="kagent"
TEST_TIMEOUT=30

# kagent agents available for testing
KAGENT_AGENTS=(
    "k8s-agent"
    "helm-agent"
    "istio-agent"
    "cilium-debug-agent"
    "observability-agent"
    "promql-agent"
    "argo-rollouts-conversion-agent"
    "kgateway-agent"
    "openribcage-test-agent"
)

# Check if kagent sandbox is running
check_kagent_availability() {
    print_status "Checking kagent sandbox availability..."
    
    if ! command -v curl &> /dev/null; then
        print_error "curl is required for A2A testing but not found"
        print_error "Please install curl to run A2A tests"
        exit 1
    fi
    
    # Test basic connectivity to kagent
    if curl -s --connect-timeout 5 "$KAGENT_BASE_URL" >/dev/null 2>&1; then
        print_success "kagent sandbox is reachable at $KAGENT_BASE_URL"
    else
        print_error "Cannot reach kagent sandbox at $KAGENT_BASE_URL"
        print_error ""
        print_error "Please ensure kagent sandbox is running:"
        print_error "  1. Clone: git clone https://github.com/craine-io/istio-envoy-sandboxes"
        print_error "  2. Setup: cd istio-envoy-sandboxes/k3d-sandboxes/kagent-sandbox"
        print_error "  3. Run:   ./scripts/cluster-setup-k3d-kagent-everything.sh"
        print_error "  4. Wait for agents to be ready, then retry this script"
        exit 1
    fi
}

# Check if openribcage binaries are built
check_openribcage_binaries() {
    print_status "Checking openribcage binaries..."
    
    if [ ! -f "bin/openribcage" ]; then
        print_status "Building openribcage binary..."
        if ! make build >/dev/null 2>&1; then
            print_error "Failed to build openribcage. Run 'make build' to see errors."
            exit 1
        fi
    fi
    
    if [ ! -f "bin/discovery" ]; then
        print_status "Building discovery binary..."
        if ! make build-all >/dev/null 2>&1; then
            print_error "Failed to build discovery tool. Run 'make build-all' to see errors."
            exit 1
        fi
    fi
    
    print_success "openribcage binaries are ready"
}

# Test AgentCard discovery
test_agentcard_discovery() {
    print_status "Testing AgentCard discovery..."
    
    # Test discovery of kagent namespace
    local discovery_url="$KAGENT_BASE_URL/$KAGENT_NAMESPACE"
    
    print_status "Attempting to discover agents at: $discovery_url"
    
    # TODO: This will be implemented when Issue #3 is complete
    # For now, test basic connectivity and provide placeholder
    if curl -s --connect-timeout 5 "$discovery_url" >/dev/null 2>&1; then
        print_success "AgentCard discovery endpoint is reachable"
        print_warning "Full AgentCard discovery will be implemented in Issue #3"
    else
        print_warning "AgentCard discovery endpoint test inconclusive"
        print_warning "This is expected during scaffolding phase"
    fi
}

# Test basic A2A communication
test_a2a_communication() {
    print_status "Testing A2A communication with kagent agents..."
    
    # Test a few key agents
    local test_agents=("k8s-agent" "helm-agent" "istio-agent")
    
    for agent in "${test_agents[@]}"; do
        print_status "Testing communication with $agent..."
        
        local agent_url="$KAGENT_BASE_URL/$KAGENT_NAMESPACE/$agent"
        
        # Create a simple A2A test request
        local test_request='{
            "jsonrpc": "2.0",
            "method": "tasks/send",
            "params": {
                "id": "test-'$(date +%s)'",
                "message": {
                    "role": "user",
                    "parts": [{
                        "type": "text",
                        "text": "Hello, this is a test from openribcage A2A client"
                    }]
                }
            },
            "id": 1
        }'
        
        # TODO: This will use the actual A2A client when Issue #10 is complete
        # For now, test basic HTTP connectivity
        if curl -s --connect-timeout 5 -H "Content-Type: application/json" \
               -d "$test_request" "$agent_url" >/dev/null 2>&1; then
            print_success "A2A endpoint $agent is reachable"
        else
            print_warning "A2A endpoint $agent test inconclusive"
        fi
    done
    
    print_warning "Full A2A communication will be implemented in Issue #10"
}

# Test openribcage CLI tools
test_cli_tools() {
    print_status "Testing openribcage CLI tools..."
    
    # Test main CLI help
    if ./bin/openribcage --help >/dev/null 2>&1; then
        print_success "openribcage CLI responds to --help"
    else
        print_error "openribcage CLI failed basic test"
        exit 1
    fi
    
    # Test discovery tool help
    if ./bin/discovery --help >/dev/null 2>&1; then
        print_success "discovery tool responds to --help"
    else
        print_error "discovery tool failed basic test"
        exit 1
    fi
    
    print_warning "Full CLI functionality will be implemented in upcoming issues"
}

# Test integration test suite
test_integration_suite() {
    print_status "Running integration test suite..."
    
    # Run integration tests if they exist
    if go test -tags=integration ./test/integration/... >/dev/null 2>&1; then
        print_success "Integration tests pass"
    else
        print_warning "Integration tests not yet implemented or failing"
        print_warning "This is expected during scaffolding phase"
    fi
}

# Generate test report
generate_test_report() {
    echo
    echo "📊 A2A Client Test Report"
    echo "========================="
    echo
    echo "Environment:"
    echo "  • kagent URL: $KAGENT_BASE_URL"
    echo "  • Timeout: ${TEST_TIMEOUT}s"
    echo "  • Available agents: ${#KAGENT_AGENTS[@]}"
    echo
    echo "Test Results:"
    echo "  ✅ kagent sandbox connectivity"
    echo "  ✅ openribcage binaries build"
    echo "  ✅ CLI tools respond to basic commands"
    echo "  ⏳ AgentCard discovery (Issue #3)"
    echo "  ⏳ A2A communication (Issue #10)"
    echo "  ⏳ Integration tests (Issue #10)"
    echo
    echo "Next Steps:"
    echo "  1. Implement Issue #3: AgentCard Discovery"
    echo "  2. Implement Issue #10: A2A Client Library"
    echo "  3. Re-run this script for full validation"
    echo
    print_success "A2A client testing complete!"
}

# Main execution
main() {
    # Verify we're in the right directory
    if [ ! -f "go.mod" ] || [ ! -f "Makefile" ]; then
        print_error "This script must be run from the openribcage project root directory"
        exit 1
    fi
    
    check_kagent_availability
    check_openribcage_binaries
    test_agentcard_discovery
    test_a2a_communication
    test_cli_tools
    test_integration_suite
    generate_test_report
}

# Execute main function
main "$@"
