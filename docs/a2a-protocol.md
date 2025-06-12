# A2A Protocol Implementation Guide

This document provides comprehensive guidance for implementing the A2A (Agent2Agent) protocol in openribcage, covering all aspects from basic JSON-RPC 2.0 communication to advanced streaming patterns.

## Overview

The A2A protocol is a standards-based communication protocol for AI agent coordination, built on JSON-RPC 2.0 with standardized method calls, AgentCard discovery, and Server-Sent Events streaming.

## A2A Protocol Foundation

### JSON-RPC 2.0 Transport

All A2A communication uses JSON-RPC 2.0 over HTTP(S):

```json
{
  "jsonrpc": "2.0",
  "method": "tasks/send",
  "params": {
    "id": "task-123",
    "message": {
      "role": "user",
      "parts": [
        {
          "type": "text",
          "text": "What is the status of my cluster?"
        }
      ]
    }
  },
  "id": 1
}
```

### Standard A2A Methods

The A2A protocol defines these core method calls:

#### Task Management
- **`tasks/send`** - Send a task to an agent
- **`tasks/sendSubscribe`** - Send task with streaming response
- **`tasks/status`** - Check task execution status
- **`tasks/cancel`** - Cancel a running task

#### Message Handling (Legacy)
- **`message/send`** - Send message to agent
- **`message/stream`** - Stream messages to agent

## Message Structure

### Message Format

```go
type Message struct {
    Role  string `json:"role"`  // "user", "assistant", "system"
    Parts []Part `json:"parts"`
}

type Part struct {
    Type string      `json:"type"`           // "text", "file", "data"
    Text string      `json:"text,omitempty"`
    Data interface{} `json:"data,omitempty"`
    File *FilePart   `json:"file,omitempty"`
}
```

### Content Types

#### TextPart
Simple text content for natural language communication:

```json
{
  "type": "text",
  "text": "Show me the cluster status"
}
```

#### FilePart
File attachments with metadata:

```json
{
  "type": "file",
  "file": {
    "name": "config.yaml",
    "mime_type": "application/yaml",
    "size": 1024,
    "url": "https://example.com/file.yaml"
  }
}
```

#### DataPart
Structured data for complex requests:

```json
{
  "type": "data",
  "data": {
    "query": "SELECT * FROM pods",
    "format": "json"
  }
}
```

## Implementation in openribcage

### Client Architecture

```go
// pkg/a2a/client/client.go
type Client struct {
    baseURL    string
    httpClient *http.Client
    auth       *auth.Authenticator
    logger     *logrus.Logger
}

// Core A2A methods
func (c *Client) SendTask(ctx context.Context, agentID string, req *types.TaskRequest) (*types.TaskResponse, error)
func (c *Client) StreamTask(ctx context.Context, agentID string, req *types.TaskRequest) (<-chan *types.StreamResponse, <-chan error)
func (c *Client) GetTaskStatus(ctx context.Context, agentID, taskID string) (*types.TaskStatus, error)
func (c *Client) CancelTask(ctx context.Context, agentID, taskID string) error
```

### Error Handling

A2A protocol uses standard JSON-RPC 2.0 error codes:

```go
type JSONRPCError struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

// Standard error codes
const (
    ParseError     = -32700
    InvalidRequest = -32600
    MethodNotFound = -32601
    InvalidParams  = -32602
    InternalError  = -32603
)
```

## Streaming Communication

### Server-Sent Events (SSE)

A2A uses SSE for real-time streaming responses:

```go
// Subscribe to agent stream
streamChan, errChan := client.StreamTask(ctx, "k8s-agent", taskRequest)

for {
    select {
    case response := <-streamChan:
        // Process streaming response
        handleStreamResponse(response)
    case err := <-errChan:
        // Handle stream error
        log.Errorf("Stream error: %v", err)
        return
    case <-ctx.Done():
        return
    }
}
```

### Stream Response Format

```go
type StreamResponse struct {
    ID        string      `json:"id"`
    Timestamp time.Time   `json:"timestamp"`
    Type      string      `json:"type"`     // "progress", "data", "complete"
    Data      interface{} `json:"data"`
    Done      bool        `json:"done,omitempty"`
}
```

## Authentication Patterns

### Bearer Token Authentication

```go
creds := &auth.Credentials{
    Type:  auth.AuthTypeBearer,
    Token: "your-api-token",
}

client := NewClient(Config{
    BaseURL: "http://localhost:8083/api/a2a",
    Auth:    creds,
})
```

### API Key Authentication

```go
creds := &auth.Credentials{
    Type:   auth.AuthTypeAPIKey,
    APIKey: "your-api-key",
}
```

## Testing with kagent

### Setting up kagent Sandbox

```bash
# Clone kagent sandbox
git clone https://github.com/craine-io/istio-envoy-sandboxes
cd istio-envoy-sandboxes/k3d-sandboxes/kagent-sandbox

# Setup cluster with A2A endpoints
./scripts/cluster-setup-k3d-kagent-everything.sh

# Test A2A connectivity
curl http://localhost:8083/api/a2a/kagent/k8s-agent/.well-known/agent.json
```

### Example A2A Calls

#### Send Task to k8s-agent

```bash
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

#### Stream Task Response

```bash
curl -X POST http://localhost:8083/api/a2a/kagent/k8s-agent \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tasks/sendSubscribe",
    "params": {
      "id": "stream-123",
      "message": {
        "role": "user",
        "parts": [{
          "type": "text",
          "text": "Monitor cluster status"
        }]
      }
    },
    "id": 1
  }'
```

## Performance Considerations

### Connection Pooling

```go
client := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
    Timeout: 30 * time.Second,
}
```

### Request Timeouts

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

response, err := client.SendTask(ctx, agentID, request)
```

### Retry Mechanisms

```go
func (c *Client) sendWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
    var lastErr error
    for i := 0; i < c.maxRetries; i++ {
        resp, err := c.httpClient.Do(req)
        if err == nil && resp.StatusCode < 500 {
            return resp, nil
        }
        lastErr = err
        time.Sleep(time.Duration(i+1) * time.Second)
    }
    return nil, lastErr
}
```

## Implementation Roadmap

### Phase 1: Basic A2A Client (Issues #3, #4)
- [ ] AgentCard discovery implementation
- [ ] JSON-RPC 2.0 client foundation
- [ ] Basic task/message methods
- [ ] Error handling and logging

### Phase 2: Advanced Features (Issue #10)
- [ ] Server-Sent Events streaming
- [ ] Connection management and pooling
- [ ] Authentication integration
- [ ] Comprehensive method implementation

### Phase 3: Production Ready (Issue #11)
- [ ] Performance optimization
- [ ] Advanced error recovery
- [ ] Metrics and monitoring
- [ ] Load balancing and failover

## Compliance Validation

### A2A Protocol Checklist

- [ ] JSON-RPC 2.0 transport layer
- [ ] Standard method implementation
- [ ] AgentCard discovery support
- [ ] SSE streaming capability
- [ ] Proper error code handling
- [ ] Authentication integration
- [ ] Content type support (TextPart, FilePart, DataPart)

### Testing Checklist

- [ ] Unit tests for all A2A methods
- [ ] Integration tests with kagent
- [ ] Streaming response handling
- [ ] Error scenario coverage
- [ ] Performance benchmarking
- [ ] Multi-agent coordination

## References

- [A2A Protocol Specification](https://github.com/google-a2a/A2A)
- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)
- [Server-Sent Events Standard](https://html.spec.whatwg.org/multipage/server-sent-events.html)
- [kagent A2A Implementation](https://github.com/solo-io/kagent)

---

*This document will be updated as A2A client implementation progresses through Issues #3, #4, #10, and #11.*
