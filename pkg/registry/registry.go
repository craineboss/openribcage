// Package registry provides A2A agent registry and management.
//
// This package implements the agent registry that maintains information
// about discovered A2A agents, their capabilities, status, and provides
// agent lookup and management functionality.
package registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/craine-io/openribcage/pkg/a2a/types"
)

// Registry manages discovered A2A agents
type Registry struct {
	mu      sync.RWMutex
	agents  map[string]*types.Agent
	logger  *logrus.Logger
	cleanup time.Duration
}

// NewRegistry creates a new agent registry
func NewRegistry(cleanupInterval time.Duration) *Registry {
	return &Registry{
		agents:  make(map[string]*types.Agent),
		logger:  logrus.New(),
		cleanup: cleanupInterval,
	}
}

// Register registers a new agent in the registry
func (r *Registry) Register(agent *types.Agent) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Infof("Registering agent: %s (%s)", agent.Name, agent.URL)

	// TODO: Implement agent registration
	// This will be implemented in Issue #3
	// 1. Validate agent data
	// 2. Check for duplicates
	// 3. Add to registry
	// 4. Update timestamps

	r.agents[agent.ID] = agent
	return nil
}

// Unregister removes an agent from the registry
func (r *Registry) Unregister(agentID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Infof("Unregistering agent: %s", agentID)

	if _, exists := r.agents[agentID]; !exists {
		return fmt.Errorf("agent not found: %s", agentID)
	}

	delete(r.agents, agentID)
	return nil
}

// Get retrieves an agent by ID
func (r *Registry) Get(agentID string) (*types.Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agent, exists := r.agents[agentID]
	if !exists {
		return nil, fmt.Errorf("agent not found: %s", agentID)
	}

	return agent, nil
}

// List returns all registered agents
func (r *Registry) List() []*types.Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	agents := make([]*types.Agent, 0, len(r.agents))
	for _, agent := range r.agents {
		agents = append(agents, agent)
	}

	return agents
}

// FindByCapability finds agents with specific capabilities
func (r *Registry) FindByCapability(capability string) []*types.Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matches []*types.Agent
	for _, agent := range r.agents {
		if agent.Card != nil {
			caps := agent.Card.Capabilities
			if (capability == "streaming" && caps.Streaming) ||
				(capability == "pushNotifications" && caps.PushNotifications) ||
				(capability == "stateTransitionHistory" && caps.StateTransitionHistory) {
				matches = append(matches, agent)
			}
		}
	}

	return matches
}

// UpdateStatus updates an agent's status
func (r *Registry) UpdateStatus(agentID string, status types.AgentStatus) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	agent, exists := r.agents[agentID]
	if !exists {
		return fmt.Errorf("agent not found: %s", agentID)
	}

	agent.Status = status
	agent.LastSeen = time.Now()

	r.logger.Debugf("Updated agent %s status to %s", agentID, status)
	return nil
}

// StartCleanup starts the cleanup goroutine for stale agents
func (r *Registry) StartCleanup(ctx context.Context) {
	ticker := time.NewTicker(r.cleanup)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.cleanupStaleAgents()
		}
	}
}

// cleanupStaleAgents removes agents that haven't been seen recently
func (r *Registry) cleanupStaleAgents() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	staleThreshold := 5 * time.Minute

	for id, agent := range r.agents {
		if now.Sub(agent.LastSeen) > staleThreshold {
			r.logger.Warnf("Removing stale agent: %s", agent.Name)
			delete(r.agents, id)
		}
	}
}
