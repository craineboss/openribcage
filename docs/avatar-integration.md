# Avatar Interface Integration Guide

This document explains how openribcage integrates with avatar-based interfaces to provide natural conversation experiences with A2A agents.

## Overview

Avatar integration transforms A2A protocol data into natural, conversational interfaces where users interact with AI agents through personalized avatars that reflect agent capabilities and personalities.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                 Avatar Interface Layer                          │
│         (React/Vue/Angular + Avatar Components)                 │
├─────────────────────────────────────────────────────────────────┤
│              Avatar Integration API                              │
│           (Persona Mapping + Stream Processing)                 │
├─────────────────────────────────────────────────────────────────┤
│              openribcage A2A Client                             │
│           (AgentCard Discovery + JSON-RPC 2.0)                  │
├─────────────────────────────────────────────────────────────────┤
│              A2A Protocol (JSON-RPC 2.0 / SSE)                  │
│                  (Standard agent communication)                 │
└─────────────────────────────────────────────────────────────────┘
```

## AgentCard to Avatar Persona Mapping

### Persona Generation

```go
// Map AgentCard to avatar persona
mapper := avatar.NewPersonaMapper()
persona, err := mapper.MapToPersona(agentCard)
if err != nil {
    log.Fatalf("Persona mapping failed: %v", err)
}

fmt.Printf("Avatar: %s\n", persona.Name)
fmt.Printf("Style: %s\n", persona.VisualTraits.Style)
fmt.Printf("Personality: %+v\n", persona.Personality)
```

### Example Persona Mapping

```json
{
  "id": "k8s-agent-persona",
  "name": "Kubernetes Commander",
  "description": "Your expert Kubernetes cluster manager",
  "personality": {
    "formal": 0.7,
    "enthusiasm": 0.6,
    "technical": 0.9,
    "helpfulness": 0.8,
    "creativity": 0.4
  },
  "visual_traits": {
    "style": "technical",
    "color_theme": "blue",
    "icon": "kubernetes-logo",
    "avatar": "professional-tech"
  },
  "communication_style": {
    "tone": "professional",
    "vocabulary": "technical",
    "greeting": "Hello! I'm ready to help manage your Kubernetes cluster."
  }
}
```

## Real-time Avatar Updates

### A2A Streaming Integration

```go
// Process A2A stream updates for avatar
integrator := avatar.NewStreamIntegrator()

// Subscribe to agent stream
streamChan, errChan := a2aClient.StreamTask(ctx, "k8s-agent", taskRequest)

for {
    select {
    case update := <-streamChan:
        // Convert A2A update to avatar update
        avatarUpdate, err := integrator.ProcessStreamUpdate(ctx, update)
        if err != nil {
            log.Warnf("Stream processing failed: %v", err)
            continue
        }
        
        // Send to avatar interface
        sendToAvatarUI(avatarUpdate)
        
    case err := <-errChan:
        log.Errorf("Stream error: %v", err)
        return
    }
}
```

### Frontend Integration Example

```jsx
import React, { useState, useEffect } from 'react';

function AvatarInterface({ agentId }) {
  const [avatar, setAvatar] = useState(null);
  const [messages, setMessages] = useState([]);
  const [status, setStatus] = useState('idle');
  
  useEffect(() => {
    // Load avatar persona
    fetch(`/api/avatars/${agentId}`)
      .then(res => res.json())
      .then(setAvatar);
    
    // Subscribe to real-time updates
    const eventSource = new EventSource(`/api/avatars/${agentId}/stream`);
    eventSource.onmessage = (event) => {
      const update = JSON.parse(event.data);
      handleAvatarUpdate(update);
    };
    
    return () => eventSource.close();
  }, [agentId]);
  
  const handleAvatarUpdate = (update) => {
    switch (update.update_type) {
      case 'status_change':
        setStatus(update.data.status);
        break;
      case 'message':
        setMessages(prev => [...prev, update.data]);
        break;
      case 'animation':
        triggerAnimation(update.animations);
        break;
    }
  };
  
  const sendMessage = async (text) => {
    const response = await fetch(`/api/avatars/${agentId}/message`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ message: text })
    });
    
    const result = await response.json();
    setMessages(prev => [...prev, result]);
  };
  
  if (!avatar) return <div>Loading avatar...</div>;
  
  return (
    <div className="avatar-interface">
      <AvatarDisplay 
        persona={avatar}
        status={status}
        className={`avatar-${avatar.visual_traits.style}`}
      />
      <MessageHistory messages={messages} />
      <MessageInput onSend={sendMessage} />
    </div>
  );
}
```

## Implementation Status

### 🔄 Current Development (Issues #12-15)

- **Issue #12**: A2A-to-avatar integration requirements analysis
- **Issue #13**: AgentCard to avatar persona mapping system
- **Issue #14**: A2A streaming integration for real-time updates
- **Issue #15**: Avatar-A2A communication protocols

### 📋 Implementation Plan

1. **Persona Mapping Engine** - Convert AgentCard data to avatar characteristics
2. **Stream Integration** - Real-time A2A updates to avatar interface
3. **UI Adaptation** - Dynamic interface based on agent capabilities
4. **Multi-Agent Coordination** - Avatar team management

## Testing Avatar Integration

### Unit Tests

```go
func TestPersonaMapping(t *testing.T) {
    card := &types.AgentCard{
        Name: "k8s-agent",
        Capabilities: []string{"kubernetes-management"},
    }
    
    mapper := avatar.NewPersonaMapper()
    persona, err := mapper.MapToPersona(card)
    
    assert.NoError(t, err)
    assert.Equal(t, "technical", persona.VisualTraits.Style)
    assert.Greater(t, persona.Personality.Technical, 0.5)
}
```

### Integration Tests

```bash
# Test avatar API
curl -X POST http://localhost:8080/api/avatars/k8s-agent/message \
  -H "Content-Type: application/json" \
  -d '{"message": "Show cluster status"}'

# Test streaming
curl -N http://localhost:8080/api/avatars/k8s-agent/stream
```

## Next Steps

1. **Complete Issue #12** - Requirements analysis
2. **Implement Issue #13** - Persona mapping system
3. **Build Issue #14** - Streaming integration
4. **Create Issue #15** - Communication protocols
5. **Test with AAMI** - Avatar Agency Management Interface

## References

- [AAMI Repository](https://github.com/craine-io/aami) (Private)
- [A2A Protocol Specification](https://github.com/google-a2a/A2A)
- [Avatar Interface Design Patterns](./examples/avatar-patterns.md)
