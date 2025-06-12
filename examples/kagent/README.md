# kagent A2A Integration Example

This example demonstrates how to use openribcage with kagent A2A endpoints.

## Prerequisites

1. **kagent sandbox running**:
   ```bash
   cd kagent-sandbox
   ./scripts/cluster-setup-k3d-kagent-everything.sh
   ```

2. **openribcage built**:
   ```bash
   make build
   ```

## Basic Usage

### Discover kagent Agents

```bash
# Discover all kagent agents
./bin/discovery scan http://localhost:8083/api/a2a/kagent

# Validate specific agent
./bin/discovery validate http://localhost:8083/api/a2a/kagent/k8s-agent
```

### Communicate with Agents

```bash
# Send message to k8s-agent
./bin/openribcage communicate \
  http://localhost:8083/api/a2a/kagent/k8s-agent \
  "What is the status of my cluster?"

# Send message to helm-agent
./bin/openribcage communicate \
  http://localhost:8083/api/a2a/kagent/helm-agent \
  "List all installed Helm releases"
```

## Available kagent Agents

### Infrastructure Management
- **k8s-agent** - Kubernetes cluster operations
- **helm-agent** - Helm chart management
- **istio-agent** - Service mesh configuration

### Network Management
- **cilium-debug-agent** - Network debugging
- **cilium-manager-agent** - Network policy management
- **cilium-policy-agent** - Policy enforcement

### Observability
- **observability-agent** - Monitoring and metrics
- **promql-agent** - Prometheus query execution

### Deployment
- **argo-rollouts-conversion-agent** - Deployment strategies
- **kgateway-agent** - Gateway management

### Testing
- **openribcage-test-agent** - Testing and validation

## Configuration Example

```yaml
# openribcage.yaml
a2a:
  discovery_hosts:
    - "http://localhost:8083/api/a2a/kagent"
  timeout: "30s"
  retry_attempts: 3
  default_headers:
    User-Agent: "openribcage/1.0.0"
    
logging:
  level: "debug"
  format: "text"
```

## Testing A2A Communication

### Manual Testing

```bash
# Test basic communication
curl -X POST http://localhost:8083/api/a2a/kagent/k8s-agent \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tasks/send",
    "params": {
      "id": "test-123",
      "message": {
        "role": "user",
        "parts": [{
          "type": "text",
          "text": "Show cluster nodes"
        }]
      }
    },
    "id": 1
  }'
```

### Automated Testing

```bash
# Run kagent integration tests
./scripts/test-a2a-client.sh

# Run full integration test suite
make test-integration
```

## Expected Responses

### Successful Response
```json
{
  "jsonrpc": "2.0",
  "result": {
    "id": "test-123",
    "status": "completed",
    "message": {
      "role": "assistant",
      "parts": [{
        "type": "text",
        "text": "Cluster has 3 nodes: node1, node2, node3"
      }]
    }
  },
  "id": 1
}
```

### Error Response (Expected)
```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32603,
    "message": "Internal error",
    "data": "OpenAI quota exceeded"
  },
  "id": 1
}
```

## Troubleshooting

### Common Issues

1. **Connection Refused**
   - Ensure kagent sandbox is running
   - Check port 8083 is accessible

2. **OpenAI Quota Errors**
   - Add valid OpenAI API key to kagent secret
   - Check kagent agent reconciliation status

3. **Method Not Found**
   - Verify A2A method name (use `tasks/send`)
   - Check agent supports the method

### Debug Commands

```bash
# Check kagent sandbox status
kubectl get agents -n kagent

# View agent logs
kubectl logs -n kagent deployment/kagent-controller

# Test connectivity
curl -v http://localhost:8083/api/a2a/kagent/
```

## Next Steps

1. Implement full A2A client (Issue #10)
2. Add AgentCard discovery (Issue #3)
3. Enable streaming support (Issue #10)
4. Build avatar integration (Issues #12-15)

See the main [A2A Protocol Guide](../../docs/a2a-protocol.md) for implementation details.
