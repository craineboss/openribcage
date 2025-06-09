![openribcage Banner](openribcage-banner.png)

# openribcage - A2A Protocol Client for Avatar Interfaces

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/craine-io/openribcage)](https://goreportcard.com/report/github.com/craine-io/openribcage)
[![A2A Protocol](https://img.shields.io/badge/A2A-Protocol%20Compliant-green.svg)](https://github.com/google-a2a/A2A)

> **Standards-Based A2A Protocol Client for Avatar-Based Agent Coordination**

openribcage is an A2A (Agent2Agent) protocol client designed to enable avatar-based interfaces for AI agent coordination. It serves as the backend communication layer for avatar management interfaces, providing standardized agent discovery, communication, and orchestration capabilities through the proven A2A protocol standard.

## 🎯 The Problem We're Solving

Enterprise teams need to coordinate AI agents naturally through avatar interfaces, but agents from different frameworks (kagent, LangGraph, CrewAI) speak different protocols. What if there was a standards-based approach?

**openribcage implements the Google-backed A2A protocol to enable universal agent communication.**

## 🔍 What openribcage Does

openribcage provides a complete A2A protocol client that enables avatar interfaces to communicate with any A2A-compliant agent framework:

- **🌐 A2A Protocol Client**: Full JSON-RPC 2.0 implementation with Agent2Agent standard compliance
- **🔍 Agent Discovery**: Automatic AgentCard discovery via `.well-known/agent.json` endpoints
- **⚡ Real-time Streaming**: Server-Sent Events (SSE) for live agent communication and status updates
- **🎭 Avatar Integration**: AgentCard to avatar persona mapping for natural conversation interfaces
- **🔒 Enterprise Security**: Standard HTTP authentication with Agent Gateway integration patterns

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                 Avatar Interface Layer                          │
│          (Natural conversation with AI agents)                  │
├─────────────────────────────────────────────────────────────────┤
│              openribcage - A2A Protocol Client                  │
│           (Agent discovery, communication, orchestration)       │
├─────────────────────────────────────────────────────────────────┤
│              A2A Protocol (JSON-RPC 2.0 / HTTP)                 │
│                  (Standard agent communication)                 │
├─────────────────────────────────────────────────────────────────┤
│              A2A-Compliant Agent Frameworks                     │
│           (kagent, LangGraph, CrewAI, Google ADK)               │
└─────────────────────────────────────────────────────────────────┘
```

## 📁 Project Structure

```
openribcage/
├── cmd/                          # Command-line applications
│   ├── openribcage/             # Main A2A client application
│   └── discovery/               # Agent discovery tools
├── pkg/                          # Core library packages
│   ├── a2a/                     # A2A protocol implementation
│   │   ├── client/              # JSON-RPC 2.0 A2A client
│   │   ├── types/               # A2A protocol types and schemas
│   │   └── streaming/           # Server-Sent Events streaming
│   ├── agentcard/               # AgentCard discovery and parsing
│   ├── registry/                # Agent registry and management
│   └── avatar/                  # Avatar interface integration
├── internal/                     # Private application code
│   ├── config/                  # Configuration management
│   ├── auth/                    # Authentication handlers
│   └── logging/                 # Structured logging
├── api/                          # API definitions and schemas
│   ├── openapi/                 # OpenAPI specifications
│   └── grpc/                    # gRPC definitions (if needed)
├── deployments/                  # Deployment configurations
│   ├── kubernetes/              # K8s manifests
│   ├── docker/                  # Docker configurations
│   └── helm/                    # Helm charts
├── scripts/                      # Build and development scripts
│   ├── test-a2a-client.sh       # A2A client testing
│   └── setup-dev.sh             # Development environment setup
├── docs/                         # Documentation
│   ├── a2a-protocol.md          # A2A protocol implementation guide
│   ├── agent-discovery.md       # AgentCard discovery documentation
│   ├── avatar-integration.md    # Avatar interface integration
│   └── examples/                # Usage examples and tutorials
├── test/                         # Integration and e2e tests
│   ├── integration/             # A2A protocol integration tests
│   └── fixtures/                # Test data and mock agents
├── examples/                     # Example implementations
│   ├── kagent/                  # kagent A2A integration example
│   └── avatar-client/           # Avatar interface example
└── tools/                        # Development and build tools
    ├── agentcard-validator/     # AgentCard format validator
    └── a2a-tester/              # A2A protocol compliance tester
```

## 🎪 Core A2A Protocol Components

### JSON-RPC 2.0 Client
Standards-compliant JSON-RPC 2.0 client implementing all A2A protocol methods:
- `message/send` - Send messages to agents
- `message/stream` - Real-time streaming communication
- `task/create` - Create and manage agent tasks
- `task/status` - Monitor task execution status

### AgentCard Discovery
Automatic agent discovery via A2A standard `.well-known/agent.json` endpoints:
- JSON schema validation for AgentCard format compliance
- Capability parsing and indexing
- Health monitoring and status tracking
- Dynamic agent registration and deregistration

### Avatar Interface Integration
Bridge between A2A protocol data and avatar-based interfaces:
- AgentCard to avatar persona mapping
- Real-time avatar updates via A2A streaming
- Natural conversation flow through A2A message handling

## 🛠️ Supported A2A Frameworks

openribcage connects to any A2A-compliant agent framework:

### Current Development Priority
1. **kagent** - Kubernetes-native agent framework with native A2A endpoints
2. **LangGraph** - Graph-based agent workflows with A2A support  
3. **CrewAI** - Role-based team coordination via A2A protocol
4. **Google ADK** - Agent Development Kit with native A2A implementation

### A2A Protocol Compliance
All supported frameworks implement the standard A2A protocol:
- `.well-known/agent.json` AgentCard discovery
- JSON-RPC 2.0 method calls over HTTP(S)
- Server-Sent Events for real-time streaming
- Standard authentication schemes

## 🚀 Getting Started

> **Note**: openribcage is currently in active development. We're building the A2A protocol client foundation and kagent integration.

### Prerequisites
- **Go 1.21+** for core development
- **Docker** for containerization
- **kubectl** for kagent testing
- **kagent sandbox** for A2A endpoint testing

### Development Setup
```bash
# Clone the repository
git clone https://github.com/craine-io/openribcage.git
cd openribcage

# Set up development environment
./scripts/setup-dev.sh

# Test A2A client with kagent endpoints (coming in Phase 2)
./scripts/test-a2a-client.sh

# Run A2A protocol compliance tests
go test ./test/integration/...
```

### kagent A2A Testing
Use the [kagent sandbox](https://github.com/craine-io/istio-envoy-sandboxes/tree/main/k3d-sandboxes/kagent-sandbox) for A2A protocol testing:

```bash
# Set up kagent with A2A endpoints
cd kagent-sandbox
./scripts/cluster-setup-k3d-kagent-everything.sh

# Test openribcage A2A client
# kagent A2A endpoint: http://localhost:8083/api/a2a/kagent/
```

## 🤝 Contributing

We're actively seeking contributors for A2A protocol implementation:

**Scope**: A2A protocol compliance, agent communication, discovery mechanisms  
**Out of Scope**: Avatar interface implementations, UI components, proprietary integrations

### Development Focus Areas
- **A2A Protocol Implementation** - JSON-RPC 2.0 client and method handlers
- **AgentCard Discovery** - Agent discovery and capability parsing
- **Streaming Integration** - Server-Sent Events for real-time updates
- **kagent Integration** - Reference A2A implementation testing

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## 🗺️ Development Roadmap

### Phase 1: A2A Protocol Foundation (Weeks 1-2) 🔄 *Current Phase*
- A2A protocol specification analysis and documentation
- AgentCard discovery and parsing system implementation
- JSON-RPC 2.0 client with basic A2A method support
- Agent lifecycle management via A2A protocol

### Phase 2: Core A2A Client (Weeks 3-5)
- Complete A2A method implementation (`message/send`, `message/stream`, etc.)
- Real-time streaming via Server-Sent Events
- kagent A2A endpoint integration and testing
- Agent registry and health monitoring

### Phase 3: Multi-Agent Orchestration (Weeks 6-8)
- Multi-agent discovery and coordination via A2A protocol
- Agent capability mapping and routing
- Cross-agent communication patterns
- A2A protocol compliance testing across frameworks

### Phase 4: Avatar Interface Integration (Weeks 9-12)
- AgentCard to avatar persona mapping system
- Dynamic UI adaptation based on agent capabilities
- Real-time avatar updates from A2A streams
- Avatar interface API and integration patterns

### Phase 5: Production Deployment (Weeks 13-15)
- End-to-end A2A protocol integration testing
- Performance optimization and connection pooling
- Kubernetes packaging and deployment patterns
- Production observability and monitoring

### Phase 6: Enterprise Scale (Future)
- [Agent Gateway](https://agentgateway.dev/) integration for enterprise data plane
- Federated agent discovery across environments
- Enterprise security with RBAC and audit logging
- Advanced multi-agent coordination patterns

## 🔒 A2A Protocol Security

openribcage implements standard A2A protocol security patterns:
- **HTTP Authentication**: Bearer tokens, API keys, OAuth2
- **Agent Gateway Integration**: Enterprise-grade security for agent communication
- **TLS/HTTPS**: Secure transport for all A2A protocol communication
- **AgentCard Validation**: Schema validation and security scanning

## 🌟 Why A2A Protocol?

- **Industry Standard**: Google-backed, mature protocol specification
- **Framework Agnostic**: Works with any A2A-compliant agent framework  
- **Future-Proof**: New frameworks adopting A2A automatically compatible
- **Enterprise Ready**: Built-in authentication, streaming, task management
- **Avatar Optimized**: Perfect foundation for avatar-based agent interfaces

## 📖 Documentation

- **[A2A Protocol Guide](docs/a2a-protocol.md)** - Complete A2A implementation guide
- **[Agent Discovery](docs/agent-discovery.md)** - AgentCard discovery and parsing
- **[Avatar Integration](docs/avatar-integration.md)** - Avatar interface patterns
- **[kagent Integration](docs/examples/kagent.md)** - kagent A2A testing guide

## 🎯 Community & Support

- **💬 Discord**: [Join our community](https://discord.gg/craine-io)
- **🐛 Issues**: [Report bugs or request features](https://github.com/craine-io/openribcage/issues)
- **📚 Documentation**: [A2A Protocol Docs](https://github.com/google-a2a/A2A)
- **🤝 Contributing**: [Development Guidelines](CONTRIBUTING.md)

## 🔮 The Vision

Transform complex multi-agent coordination into natural conversation through A2A protocol compliance and avatar-based interfaces. When your infrastructure agents (kagent) can seamlessly collaborate with your business process agents (CrewAI) and workflow agents (LangGraph)—all through standardized A2A communication and natural avatar interfaces—that's the future openribcage enables.

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Google** for the Agent2Agent (A2A) protocol specification
- **Solo.io** for kagent framework and A2A implementation
- **Agent Gateway** for enterprise security pattern inspiration
- **The open-source AI agent community** for A2A protocol adoption

---

## 🚀 Ready to Build A2A Protocol Client?

**A2A Protocol Client Foundation**: JSON-RPC 2.0 + AgentCard Discovery + Avatar Integration

**[Get Started →](CONTRIBUTING.md)** | **[Join Discord →](https://discord.gg/craine-io)** | **[View Issues →](https://github.com/craine-io/project-openribcage/issues)**
