#!/bin/bash

# openribcage Development Environment Setup Script
# A2A Protocol Client for Avatar Interfaces
#
# This script sets up the complete development environment for openribcage
# with one command as referenced in the PR and README documentation.

set -e

echo "🏗️  Setting up openribcage development environment..."
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

# Check if Go is installed and version requirement
check_go_version() {
    print_status "Checking Go installation and version..."
    
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21+ from https://golang.org/dl/"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    REQUIRED_VERSION="1.21"
    
    if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
        print_error "Go version $GO_VERSION found, but Go 1.21+ is required"
        print_error "Please upgrade Go from https://golang.org/dl/"
        exit 1
    fi
    
    print_success "Go $GO_VERSION found (meets requirement: 1.21+)"
}

# Setup Go module and download dependencies
setup_go_module() {
    print_status "Setting up Go module and downloading dependencies..."
    
    if [ ! -f "go.mod" ]; then
        print_error "go.mod not found. Are you in the openribcage project root?"
        exit 1
    fi
    
    # Download dependencies
    go mod download
    
    # Tidy up module
    go mod tidy
    
    print_success "Go dependencies downloaded and module tidied"
}

# Install golangci-lint if not present
install_golangci_lint() {
    print_status "Checking golangci-lint installation..."
    
    if command -v golangci-lint &> /dev/null; then
        LINT_VERSION=$(golangci-lint version 2>/dev/null | awk '{print $4}' | head -n1 || echo "unknown")
        print_success "golangci-lint found (version: $LINT_VERSION)"
        return
    fi
    
    print_status "Installing golangci-lint..."
    
    # Install latest version
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    
    if command -v golangci-lint &> /dev/null; then
        print_success "golangci-lint installed successfully"
    else
        print_error "Failed to install golangci-lint"
        print_error "Please ensure \$GOPATH/bin is in your PATH"
        exit 1
    fi
}

# Verify development tools
verify_dev_tools() {
    print_status "Verifying development tools..."
    
    # Check make
    if ! command -v make &> /dev/null; then
        print_warning "make not found - Makefile targets may not work"
        print_warning "Install build-essential (Ubuntu/Debian) or Command Line Tools (macOS)"
    else
        print_success "make found"
    fi
    
    # Check docker (optional)
    if command -v docker &> /dev/null; then
        print_success "Docker found (for containerization)"
    else
        print_warning "Docker not found (optional - needed for containerization)"
    fi
    
    # Check kubectl (optional for kagent testing)
    if command -v kubectl &> /dev/null; then
        print_success "kubectl found (for kagent A2A testing)"
    else
        print_warning "kubectl not found (optional - needed for kagent sandbox testing)"
    fi
}

# Test the build process
test_build() {
    print_status "Testing build process..."
    
    # Create bin directory if it doesn't exist
    mkdir -p bin
    
    # Test build
    if go build -o bin/openribcage ./cmd/openribcage; then
        print_success "Main application builds successfully"
    else
        print_error "Failed to build main application"
        exit 1
    fi
    
    if go build -o bin/discovery ./cmd/discovery; then
        print_success "Discovery tool builds successfully"
    else
        print_error "Failed to build discovery tool"
        exit 1
    fi
    
    # Clean up test binaries
    rm -f bin/openribcage bin/discovery
}

# Run basic tests
run_tests() {
    print_status "Running basic tests..."
    
    if go test ./...; then
        print_success "All tests pass"
    else
        print_warning "Some tests failed - this is expected in scaffolding phase"
    fi
}

# Run linter
run_linter() {
    print_status "Running code linter..."
    
    if golangci-lint run; then
        print_success "Code passes linting checks"
    else
        print_warning "Linter found issues - address these for clean code"
    fi
}

# Print setup completion summary
print_setup_summary() {
    echo
    echo "🎉 Development environment setup complete!"
    echo
    echo -e "${GREEN}Next steps:${NC}"
    echo "  1. Build applications:     ${BLUE}make build${NC}"
    echo "  2. Run tests:             ${BLUE}make test${NC}"
    echo "  3. Run linter:            ${BLUE}make lint${NC}"
    echo "  4. Development workflow:   ${BLUE}make dev${NC}"
    echo
    echo -e "${YELLOW}For A2A testing with kagent:${NC}"
    echo "  • Set up kagent sandbox from: https://github.com/craine-io/istio-envoy-sandboxes"
    echo "  • Run A2A tests:          ${BLUE}./scripts/test-a2a-client.sh${NC}"
    echo
    echo -e "${BLUE}Happy coding! 🚀${NC}"
}

# Main execution
main() {
    # Verify we're in the right directory
    if [ ! -f "go.mod" ] || [ ! -f "Makefile" ]; then
        print_error "This script must be run from the openribcage project root directory"
        exit 1
    fi
    
    check_go_version
    setup_go_module
    install_golangci_lint
    verify_dev_tools
    test_build
    run_tests
    run_linter
    print_setup_summary
}

# Execute main function
main "$@"
