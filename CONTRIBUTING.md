# Contributing to openribcage

🎉 Thank you for your interest in contributing to openribcage! We're building an A2A (Agent2Agent) protocol client for avatar-based agent coordination, and we need passionate developers, protocol experts, and enterprise practitioners to make this vision a reality.

## 🎯 What We're Building

openribcage is a standards-based A2A protocol client that enables avatar interfaces to communicate with any A2A-compliant agent framework. We're focused on implementing the Google-backed Agent2Agent protocol specification to provide universal agent communication through natural avatar interfaces.

## 🤝 How You Can Contribute

### 🌐 A2A Protocol Implementation (Primary Need)
**Build the core A2A protocol client infrastructure**

We're actively seeking developers to work on:
- **JSON-RPC 2.0 Client** - Core A2A method implementation
- **AgentCard Discovery** - Agent discovery via `.well-known/agent.json` endpoints
- **Server-Sent Events Streaming** - Real-time A2A communication
- **Agent Registry** - A2A agent management and lifecycle

### 🧪 Framework Integration Testing
**Validate A2A protocol compliance across frameworks**

- **kagent Integration** - Test with kagent A2A endpoints (Priority #1)
- **LangGraph A2A Testing** - Validate LangGraph A2A compliance
- **CrewAI Protocol Testing** - Test CrewAI A2A implementation
- **Google ADK Integration** - Validate Google Agent Development Kit A2A support

### 🎭 Avatar Interface Integration
**Bridge A2A protocol data to avatar-based interfaces**

- AgentCard to avatar persona mapping systems
- Real-time avatar updates via A2A streaming
- Avatar interface APIs and integration patterns
- Natural conversation flow through A2A messaging

### 🛡️ Security & Enterprise Features
**Build enterprise-grade A2A protocol security**

- A2A authentication and authorization patterns
- Agent Gateway integration for enterprise data plane
- A2A protocol security validation and testing
- Compliance framework support (SOX, GDPR, HIPAA)

### 📚 Documentation & Community
**Help others succeed with A2A protocol implementation**

- A2A protocol implementation guides
- AgentCard format documentation
- Avatar integration tutorials
- Community support and mentorship

## 🚀 Getting Started

### Prerequisites

- **Go 1.21+** for core A2A client development
- **Docker** for containerization and testing
- **kubectl** for kagent A2A endpoint testing
- Basic understanding of JSON-RPC 2.0 and HTTP protocols

### Development Environment Setup

```bash
# Clone the repository
git clone https://github.com/craine-io/openribcage.git
cd openribcage

# Set up development environment
./scripts/setup-dev.sh

# Test A2A protocol client (coming in Phase 2)
./scripts/test-a2a-client.sh

# Run A2A protocol compliance tests
go test ./test/integration/...
```

### kagent A2A Testing Environment

Set up the [kagent sandbox](https://github.com/craine-io/istio-envoy-sandboxes/tree/main/k3d-sandboxes/kagent-sandbox) for real A2A endpoint testing:

```bash
# Set up kagent with A2A endpoints
cd kagent-sandbox
./scripts/cluster-setup-k3d-kagent-everything.sh

# Test openribcage A2A client against kagent
# kagent A2A endpoint: http://localhost:8083/api/a2a/kagent/
```

### Your First Contribution

1. **Join our Discord**: [https://discord.gg/craine-io](https://discord.gg/craine-io)
2. **Browse A2A protocol issues**: Look for `a2a-protocol`, `agentcard`, and `json-rpc` labels
3. **Read the A2A specification**: Study [Google's A2A Protocol](https://github.com/google-a2a/A2A)
4. **Start with Issue #2**: "Study A2A Protocol Specification" is perfect for newcomers
5. **Test with kagent**: Validate your A2A understanding with real endpoints

## 🌐 A2A Protocol Development Guide

### Core A2A Components

openribcage implements the complete A2A protocol specification:

```go
// A2A Protocol Client Interface
type A2AClient interface {
    // Agent Discovery
    DiscoverAgent(baseURL string) (*AgentCard, error)
    RegisterAgent(agent *AgentCard) error
    
    // A2A Method Calls (JSON-RPC 2.0)
    SendMessage(agentID string, message *A2AMessage) (*A2AResponse, error)
    CreateTask(agentID string, task *A2ATask) (*A2ATaskResponse, error)
    GetTaskStatus(taskID string) (*A2ATaskStatus, error)
    
    // Real-time Streaming (Server-Sent Events)
    StreamMessages(agentID string) (<-chan *A2AMessage, error)
    StreamTaskUpdates(taskID string) (<-chan *A2ATaskUpdate, error)
    
    // Health and Diagnostics
    Ping(agentID string) error
    GetAgentStatus(agentID string) (*A2AAgentStatus, error)
}
```

### A2A Development Priorities

#### 1. JSON-RPC 2.0 Client (`pkg/a2a/client/`)
- HTTP transport layer with connection pooling
- Request/response correlation and timeout handling
- A2A method call implementation (`message/send`, `task/create`, etc.)
- Authentication integration (Bearer tokens, API keys)

#### 2. AgentCard Discovery (`pkg/agentcard/`)
- `.well-known/agent.json` HTTP endpoint discovery
- JSON schema validation for AgentCard format compliance
- Agent capability parsing and indexing
- Dynamic agent registration and health monitoring

#### 3. Streaming Integration (`pkg/a2a/streaming/`)
- Server-Sent Events (SSE) client implementation
- Real-time message and task update streaming
- Connection management and recovery
- Event parsing and routing

#### 4. Avatar Integration (`pkg/avatar/`)
- AgentCard to avatar persona mapping
- Real-time avatar updates from A2A streams
- Avatar interface API patterns
- Natural conversation flow management

### Framework-Specific A2A Testing

#### kagent A2A Integration
- Test AgentCard discovery from kagent endpoints
- Validate JSON-RPC 2.0 method calls with kagent
- Test SSE streaming with kagent real-time responses
- Performance testing with kagent A2A endpoints

#### LangGraph A2A Compliance
- Validate LangGraph A2A protocol implementation
- Test workflow-based agent coordination via A2A
- Document LangGraph-specific A2A patterns
- Integration testing with openribcage client

#### CrewAI A2A Validation
- Test role-based team coordination via A2A protocol
- Validate CrewAI AgentCard format compliance
- Test crew communication patterns through A2A
- Document CrewAI-specific capabilities mapping

## 🏗️ Development Workflow

### 1. Issue Creation and Discussion
- **Check A2A protocol issues** before creating new ones
- **Use A2A-specific labels** for protocol-related issues
- **Join Discord #a2a-protocol** channel for technical discussions
- **Reference A2A specification** in issue descriptions

### 2. Development Process
- **Follow A2A protocol standards** for all implementations
- **Test against kagent endpoints** for validation
- **Write comprehensive A2A compliance tests**
- **Document A2A patterns and integration points**
- **Validate with A2A specification requirements**

### 3. Pull Request Guidelines
- **Focus on A2A protocol compliance** in all changes
- **Include A2A protocol testing** for new functionality
- **Reference A2A specification sections** in documentation
- **Test against multiple A2A-compliant frameworks** when possible
- **Document avatar interface integration patterns**

### Code Style and Standards

#### Go Code Standards for A2A Implementation
- Follow A2A protocol naming conventions
- Implement proper JSON-RPC 2.0 error handling
- Use structured logging for A2A protocol debugging
- Include comprehensive A2A method testing
- Document all A2A protocol compliance decisions

#### A2A Protocol Testing Requirements
- **AgentCard Validation**: Test `.well-known/agent.json` discovery
- **JSON-RPC 2.0 Compliance**: Validate all A2A method implementations
- **SSE Streaming**: Test real-time communication patterns
- **Multi-Framework Testing**: Validate against kagent, LangGraph, CrewAI
- **Avatar Integration**: Test AgentCard to avatar persona mapping

## 🛡️ A2A Protocol Security

### Secure A2A Implementation
- **Follow A2A authentication standards** for all agent communication
- **Validate AgentCard schemas** to prevent malicious agent registration
- **Implement proper TLS/HTTPS** for all A2A protocol communication
- **Use secure JSON-RPC 2.0 patterns** with proper error handling
- **Follow Agent Gateway security patterns** for enterprise integration

### A2A Protocol Compliance Testing
- **Security testing**: Validate A2A authentication mechanisms
- **Schema validation**: Test AgentCard format compliance
- **Protocol fuzzing**: Test JSON-RPC 2.0 implementation robustness
- **Integration security**: Validate cross-framework A2A security

## 📋 A2A Development Phases

### Phase 1: A2A Protocol Foundation (Weeks 1-2) 🔄 *Current*
**Priority contributions:**
- A2A protocol specification analysis and documentation
- AgentCard discovery system implementation
- JSON-RPC 2.0 client foundation
- kagent A2A endpoint testing and validation

### Phase 2: Core A2A Client (Weeks 3-5)
**Priority contributions:**
- Complete A2A method implementation
- Server-Sent Events streaming client
- Agent registry and lifecycle management
- A2A protocol compliance testing

### Phase 3: Multi-Agent Orchestration (Weeks 6-8)
**Priority contributions:**
- Multi-framework A2A protocol testing
- Cross-agent communication via A2A
- A2A protocol performance optimization
- Framework compatibility validation

### Phase 4: Avatar Interface Integration (Weeks 9-12)
**Priority contributions:**
- AgentCard to avatar persona mapping
- Real-time avatar updates via A2A streaming
- Avatar interface API development
- Natural conversation flow implementation

## 🏆 A2A Protocol Contributor Recognition

### A2A Protocol Specialist Recognition
Contributors who build A2A protocol components become **A2A Protocol Specialists** with:
- **A2A protocol expertise recognition** in community
- **Direct input on A2A compliance decisions**
- **Framework integration authority** for A2A testing
- **Conference speaking opportunities** on A2A protocol implementation

### Framework A2A Integration Maintainership
Contributors who implement A2A integration with specific frameworks become **Framework A2A Maintainers** with:
- **Framework-specific A2A testing authority**
- **A2A compliance validation responsibilities**
- **Community leadership** for framework A2A integration
- **Protocol specification input** for framework-specific patterns

## 📞 Community and Communication

### Discord Server
**Primary A2A development hub**: [https://discord.gg/craine-io](https://discord.gg/craine-io)

**A2A-Focused Channels:**
- `#a2a-protocol` - A2A protocol implementation discussions
- `#agentcard-discovery` - Agent discovery and AgentCard format
- `#json-rpc-client` - JSON-RPC 2.0 client development
- `#avatar-integration` - Avatar interface A2A integration
- `#kagent-testing` - kagent A2A endpoint testing
- `#framework-a2a` - Multi-framework A2A compliance

### A2A Protocol Office Hours
**Weekly A2A protocol focus**: Every Wednesday 3-4 PM PT
- **A2A specification discussions** and implementation guidance
- **Framework A2A compliance testing** and validation
- **Avatar integration planning** and pattern review
- **Community A2A protocol showcase** and feedback

## ❓ A2A Protocol FAQ

### Q: I'm new to JSON-RPC 2.0. Can I still contribute?
**A:** Absolutely! We provide A2A protocol learning resources and mentorship. The A2A specification is well-documented and approachable.

### Q: How do I test A2A protocol compliance?
**A:** Use the kagent sandbox for real A2A endpoint testing. We also provide A2A protocol compliance testing tools.

### Q: What if a framework doesn't support A2A protocol?
**A:** We focus exclusively on A2A-compliant frameworks. Frameworks without A2A support are out of scope for openribcage.

### Q: How do I contribute to AgentCard format validation?
**A:** Study the A2A specification AgentCard schema and implement validation in `pkg/agentcard/`. Test against kagent AgentCard examples.

### Q: Can I propose changes to A2A protocol implementation?
**A:** Implementation improvements are welcome, but protocol changes should go through the official A2A specification process.

## 🎯 Contributor License Agreement

For substantial A2A protocol contributions, we require a Contributor License Agreement (CLA) that ensures open source licensing and protects all contributors.

## 📧 Contact and Support

### A2A Protocol Maintainer Team
- **Jason T Clark** - Project Lead - [@craineboss](https://github.com/craineboss)
- **A2A Protocol Team** - [Contact via Discord #a2a-protocol](https://discord.gg/craine-io)

### Getting A2A Protocol Help
- **Discord #a2a-protocol**: [https://discord.gg/craine-io](https://discord.gg/craine-io)
- **A2A Protocol Issues**: [GitHub Issues with a2a-protocol label](https://github.com/craine-io/openribcage/issues)
- **A2A Specification**: [Official A2A Protocol Documentation](https://github.com/google-a2a/A2A)

## 🌟 Thank You

openribcage exists because of contributors who believe in standards-based agent communication and avatar-based interfaces. Every A2A protocol contribution—whether it's JSON-RPC 2.0 implementation, AgentCard discovery, or avatar integration—helps build the future of natural AI agent coordination.

Together, we're implementing the A2A protocol standard that enables universal agent communication through beautiful avatar interfaces.

**Ready to build the A2A protocol client that powers avatar-based agent coordination?** We can't wait to see what you implement!

---

*This contributing guide focuses on A2A protocol implementation. Join our community discussions to help improve openribcage's A2A compliance and avatar integration capabilities.*