# Agent Discovery Guide

This document explains how openribcage discovers A2A-compliant agents using the AgentCard system and `.well-known/agent.json` endpoints.

## Overview

Agent discovery in openribcage follows the A2A protocol standard for automatic agent registration and capability detection. The system uses AgentCard endpoints to discover agents, parse their capabilities, and maintain an active agent registry.

## AgentCard Discovery

### Discovery Process

1. **Endpoint Scanning** - Query `.well-known/agent.json` endpoints
2. **AgentCard Parsing** - Validate and parse agent capabilities
3. **Registry Registration** - Add agents to the openribcage registry
4. **Health Monitoring** - Track agent availability and status
5. **Capability Indexing** - Enable capability-based agent selection

### AgentCard Format

AgentCards follow the standard A2A format:

```json
{
  "name": "k8s-agent",
  "description": "Kubernetes cluster management agent",
  "version": "1.0.0",
  "capabilities": [
    "kubernetes-management",
    "cluster-monitoring",
    "resource-scaling"
  ],
  "endpoints": [
    {
      "type": "a2a",
      "url": "http://localhost:8083/api/a2a/kagent/k8s-agent",
      "methods": [
        "tasks/send",
        "tasks/sendSubscribe",
        "tasks/status",
        "tasks/cancel"
      ],
      "description": "A2A protocol endpoint for Kubernetes operations"
    },
    {
      "type": "streaming",
      "url": "http://localhost:8083/api/a2a/kagent/k8s-agent/stream",
      "methods": [
        "tasks/sendSubscribe"
      ],
      "description": "Server-Sent Events streaming endpoint"
    }
  ],
  "metadata": {
    "framework": "kagent",
    "language": "go",
    "deployment": "kubernetes"
  }
}
```

## Implementation in openribcage

### Discoverer Architecture

```go
// pkg/agentcard/agentcard.go
type Discoverer struct {
    client   *http.Client
    logger   *logrus.Logger
    timeout  time.Duration
    registry *registry.Registry
}

// Core discovery methods
func (d *Discoverer) Discover(ctx context.Context, agentURL string) (*types.AgentCard, error)
func (d *Discoverer) DiscoverFromHost(ctx context.Context, baseURL string) ([]*types.AgentCard, error)
func (d *Discoverer) Validate(card *types.AgentCard) error
func (d *Discoverer) Parse(data []byte) (*types.AgentCard, error)
```

### Discovery Patterns

#### Single Agent Discovery

```go
// Discover a specific agent
discoverer := agentcard.NewDiscoverer(30 * time.Second)
card, err := discoverer.Discover(ctx, "http://localhost:8083/api/a2a/kagent/k8s-agent")
if err != nil {
    log.Errorf("Discovery failed: %v", err)
    return
}

log.Infof("Discovered agent: %s with capabilities: %v", card.Name, card.Capabilities)
```

#### Bulk Discovery

```go
// Discover all agents from a namespace
cards, err := discoverer.DiscoverFromHost(ctx, "http://localhost:8083/api/a2a/kagent")
for _, card := range cards {
    registry.Register(&types.Agent{
        ID:           generateAgentID(card),
        Name:         card.Name,
        URL:          card.Endpoints[0].URL,
        Card:         card,
        Status:       types.AgentStatusOnline,
        DiscoveredAt: time.Now(),
    })
}
```

### AgentCard Validation

#### Schema Validation

```go
func (d *Discoverer) Validate(card *types.AgentCard) error {
    // Required fields
    if card.Name == "" {
        return fmt.Errorf("agent name is required")
    }
    if card.Version == "" {
        return fmt.Errorf("agent version is required")
    }
    if len(card.Endpoints) == 0 {
        return fmt.Errorf("at least one endpoint is required")
    }

    // Validate endpoints
    for _, endpoint := range card.Endpoints {
        if err := d.validateEndpoint(&endpoint); err != nil {
            return fmt.Errorf("invalid endpoint: %w", err)
        }
    }

    // Validate capabilities
    for _, capability := range card.Capabilities {
        if !isValidCapability(capability) {
            return fmt.Errorf("invalid capability: %s", capability)
        }
    }

    return nil
}
```

#### Endpoint Validation

```go
func (d *Discoverer) validateEndpoint(endpoint *types.Endpoint) error {
    // Validate URL
    if _, err := url.Parse(endpoint.URL); err != nil {
        return fmt.Errorf("invalid URL: %w", err)
    }

    // Validate endpoint type
    validTypes := []string{"a2a", "streaming", "webhook"}
    if !contains(validTypes, endpoint.Type) {
        return fmt.Errorf("unsupported endpoint type: %s", endpoint.Type)
    }

    // Validate methods
    for _, method := range endpoint.Methods {
        if !isValidA2AMethod(method) {
            return fmt.Errorf("invalid A2A method: %s", method)
        }
    }

    return nil
}
```

## Agent Registry

### Registry Management

```go
// pkg/registry/registry.go
type Registry struct {
    mu      sync.RWMutex
    agents  map[string]*types.Agent
    logger  *logrus.Logger
    cleanup time.Duration
}

// Registry operations
func (r *Registry) Register(agent *types.Agent) error
func (r *Registry) Unregister(agentID string) error
func (r *Registry) Get(agentID string) (*types.Agent, error)
func (r *Registry) List() []*types.Agent
func (r *Registry) FindByCapability(capability string) []*types.Agent
```

### Agent Lifecycle

#### Registration Process

```go
// Register discovered agent
agent := &types.Agent{
    ID:           generateAgentID(card.Name, card.Version),
    Name:         card.Name,
    URL:          primaryEndpoint.URL,
    Card:         card,
    Status:       types.AgentStatusDiscovering,
    DiscoveredAt: time.Now(),
    LastSeen:     time.Now(),
}

// Validate agent connectivity
if err := validateAgentConnectivity(agent); err != nil {
    agent.Status = types.AgentStatusError
} else {
    agent.Status = types.AgentStatusOnline
}

registry.Register(agent)
```

#### Health Monitoring

```go
func (r *Registry) StartHealthMonitoring(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            r.checkAgentHealth()
        }
    }
}

func (r *Registry) checkAgentHealth() {
    agents := r.List()
    for _, agent := range agents {
        if err := pingAgent(agent); err != nil {
            r.UpdateStatus(agent.ID, types.AgentStatusOffline)
        } else {
            r.UpdateStatus(agent.ID, types.AgentStatusOnline)
        }
    }
}
```

## Discovery Strategies

### Proactive Discovery

```go
// Continuous discovery from known hosts
func (d *Discoverer) StartContinuousDiscovery(ctx context.Context, hosts []string) {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            for _, host := range hosts {
                go d.discoverFromHost(ctx, host)
            }
        }
    }
}
```

### Event-Driven Discovery

```go
// Discovery triggered by external events
func (d *Discoverer) HandleDiscoveryEvent(event *DiscoveryEvent) {
    switch event.Type {
    case "agent_added":
        d.Discover(context.Background(), event.AgentURL)
    case "agent_removed":
        d.registry.Unregister(event.AgentID)
    case "capability_changed":
        d.refreshAgent(event.AgentID)
    }
}
```

## Capability-Based Selection

### Capability Indexing

```go
func (r *Registry) FindByCapability(capability string) []*types.Agent {
    r.mu.RLock()
    defer r.mu.RUnlock()

    var matches []*types.Agent
    for _, agent := range r.agents {
        if agent.Card != nil && agent.Status == types.AgentStatusOnline {
            for _, cap := range agent.Card.Capabilities {
                if matchesCapability(cap, capability) {
                    matches = append(matches, agent)
                    break
                }
            }
        }
    }

    return matches
}
```

### Smart Agent Selection

```go
func (r *Registry) SelectBestAgent(capability string, criteria SelectionCriteria) (*types.Agent, error) {
    candidates := r.FindByCapability(capability)
    if len(candidates) == 0 {
        return nil, fmt.Errorf("no agents found with capability: %s", capability)
    }

    // Score agents based on criteria
    scored := make([]ScoredAgent, len(candidates))
    for i, agent := range candidates {
        scored[i] = ScoredAgent{
            Agent: agent,
            Score: calculateAgentScore(agent, criteria),
        }
    }

    // Sort by score (highest first)
    sort.Slice(scored, func(i, j int) bool {
        return scored[i].Score > scored[j].Score
    })

    return scored[0].Agent, nil
}
```

## Testing Agent Discovery

### Discovery Testing with kagent

```bash
# Test AgentCard discovery
curl http://localhost:8083/api/a2a/kagent/k8s-agent/.well-known/agent.json

# Expected response format
{
  "name": "k8s-agent",
  "description": "Kubernetes cluster management agent",
  "version": "1.0.0",
  "capabilities": ["kubernetes-management"],
  "endpoints": [...]
}
```

### Integration Testing

```go
func TestAgentDiscovery(t *testing.T) {
    discoverer := agentcard.NewDiscoverer(30 * time.Second)
    
    // Test valid agent discovery
    card, err := discoverer.Discover(context.Background(), testAgentURL)
    assert.NoError(t, err)
    assert.NotNil(t, card)
    assert.Equal(t, "k8s-agent", card.Name)
    
    // Test validation
    err = discoverer.Validate(card)
    assert.NoError(t, err)
    
    // Test registry registration
    registry := registry.NewRegistry(5 * time.Minute)
    agent := &types.Agent{
        ID:   "test-agent",
        Name: card.Name,
        Card: card,
    }
    
    err = registry.Register(agent)
    assert.NoError(t, err)
}
```

## Error Handling

### Discovery Errors

```go
// Handle common discovery errors
func (d *Discoverer) handleDiscoveryError(err error, agentURL string) {
    switch {
    case isNetworkError(err):
        d.logger.Warnf("Network error discovering %s: %v", agentURL, err)
        // Schedule retry
        d.scheduleRetry(agentURL)
    
    case isParseError(err):
        d.logger.Errorf("Invalid AgentCard format at %s: %v", agentURL, err)
        // Don't retry parse errors
    
    case isTimeoutError(err):
        d.logger.Warnf("Timeout discovering %s: %v", agentURL, err)
        // Retry with longer timeout
        d.scheduleRetryWithTimeout(agentURL, d.timeout*2)
    
    default:
        d.logger.Errorf("Unknown discovery error for %s: %v", agentURL, err)
    }
}
```

## Configuration

### Discovery Configuration

```yaml
# discovery configuration
discovery:
  timeout: "30s"
  retry_attempts: 3
  retry_delay: "5s"
  health_check_interval: "30s"
  discovery_interval: "5m"
  
  # Known discovery hosts
  hosts:
    - "http://localhost:8083/api/a2a/kagent"
    - "http://production.example.com/api/a2a"
  
  # Capability filters
  capability_filters:
    - "kubernetes-*"
    - "monitoring-*"
    - "security-*"
```

## Implementation Roadmap

### Issue #3: AgentCard Discovery System
- [ ] Basic AgentCard discovery implementation
- [ ] Schema validation and parsing
- [ ] Agent registry foundation
- [ ] Health monitoring basics

### Future Enhancements
- [ ] Advanced capability matching
- [ ] Federated discovery across environments
- [ ] Discovery event webhooks
- [ ] Performance optimization for large agent populations

## Best Practices

1. **Always validate AgentCards** before registration
2. **Implement proper error handling** for network failures
3. **Use timeouts** to prevent hanging discovery requests
4. **Monitor agent health** continuously
5. **Index capabilities** for efficient agent selection
6. **Handle agent lifecycle events** gracefully

## References

- [A2A Protocol Specification](https://github.com/google-a2a/A2A)
- [AgentCard Standard Format](https://github.com/google-a2a/A2A/blob/main/agentcard.md)
- [kagent Discovery Examples](https://github.com/solo-io/kagent)

---

*This document will be updated as AgentCard discovery implementation progresses through Issue #3.*
