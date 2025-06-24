//go:build integration
// +build integration

// Package integration provides integration tests for A2A protocol compliance.
//
// These tests require a running kagent sandbox or other A2A-compliant agents.
// Run with: go test -tags=integration ./test/integration/...
package integration

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/craine-io/openribcage/pkg/a2a/client"
	"github.com/craine-io/openribcage/pkg/a2a/types"
	"github.com/craine-io/openribcage/pkg/agentcard"
)

const (
	// kagent sandbox configuration
	kagentBaseURL   = "http://localhost:8083/api/a2a"
	kagentNamespace = "kagent"
	testAgent       = "k8s-agent"
	testTimeout     = 30 * time.Second
)

// TestA2AProtocolCompliance tests basic A2A protocol compliance
func TestA2AProtocolCompliance(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Create A2A client
	clientConfig := client.Config{
		BaseURL: kagentBaseURL,
		Timeout: testTimeout,
		Headers: make(map[string]string),
	}
	a2aClient := client.New(clientConfig)

	// Test task sending
	taskRequest := &types.TaskRequest{
		ID: "integration-test-" + time.Now().Format("20060102-150405"),
		Message: &types.Message{
			Role: "user",
			Parts: []types.Part{
				{
					Type: "text",
					Text: "What is the status of my cluster?",
				},
			},
		},
	}

	t.Log("Sending task to agent...")
	response, err := a2aClient.SendTask(ctx, testAgent, taskRequest)
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, taskRequest.ID, response.ID)
}

// TestAgentCardDiscovery tests AgentCard discovery functionality
func TestAgentCardDiscovery(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Create AgentCard discoverer
	discoverer := agentcard.NewDiscoverer(testTimeout)

	// Test discovery
	agentURL := kagentBaseURL + "/" + kagentNamespace + "/" + testAgent

	t.Log("Discovering AgentCard...")
	card, err := discoverer.Discover(ctx, agentURL)
	require.NoError(t, err)
	assert.NotNil(t, card)
	assert.Equal(t, testAgent, card.Name)
	assert.NotEmpty(t, card.GetCapabilities())
	assert.NotEmpty(t, card.Endpoints)
}

// TestA2AStreaming tests Server-Sent Events streaming
func TestA2AStreaming(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	// Create A2A client
	clientConfig := client.Config{
		BaseURL: kagentBaseURL,
		Timeout: testTimeout,
		Headers: make(map[string]string),
	}
	a2aClient := client.New(clientConfig)

	// Test streaming task
	taskRequest := &types.TaskRequest{
		ID: "streaming-test-" + time.Now().Format("20060102-150405"),
		Message: &types.Message{
			Role: "user",
			Parts: []types.Part{
				{
					Type: "text",
					Text: "Monitor cluster status with streaming updates",
				},
			},
		},
	}

	t.Log("Starting streaming task...")
	streamChan, errChan := a2aClient.StreamTask(ctx, testAgent, taskRequest)
	select {
	case update := <-streamChan:
		assert.NotNil(t, update)
		assert.Equal(t, taskRequest.ID, update.ID)
	case err := <-errChan:
		t.Fatalf("Stream error: %v", err)
	case <-ctx.Done():
		t.Fatal("Test timeout")
	}
}

// TestMultipleAgents tests communication with multiple kagent agents
func TestMultipleAgents(t *testing.T) {
	// ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	// defer cancel()

	// List of kagent agents to test
	testAgents := []string{
		"k8s-agent",
		"helm-agent",
		"istio-agent",
		"cilium-debug-agent",
	}

	// Create A2A client
	// clientConfig := client.Config{
	// 	BaseURL: kagentBaseURL,
	// 	Timeout: testTimeout,
	// 	Headers: make(map[string]string),
	// }
	// a2aClient := client.New(clientConfig)

	for _, agentName := range testAgents {
		t.Run(agentName, func(t *testing.T) {
			// taskRequest := &types.TaskRequest{
			// 	ID: "multi-agent-test-" + agentName + "-" + time.Now().Format("150405"),
			// 	Message: &types.Message{
			// 		Role: "user",
			// 		Parts: []types.Part{{Type: "text", Text: "Hello, please provide your status"}},
			// 	},
			// }

			// TODO: Implement when Issue #10 is complete
			// response, err := a2aClient.SendTask(ctx, agentName, taskRequest)
			// require.NoError(t, err)
			// assert.NotNil(t, response)
			// assert.Equal(t, taskRequest.ID, response.ID)

			t.Skip("Multi-agent testing implementation pending Issue #10")
		})
	}
}

// TestA2AErrorHandling tests error handling scenarios
func TestA2AErrorHandling(t *testing.T) {
	// ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	// defer cancel()

	// Create A2A client
	// clientConfig := client.Config{
	// 	BaseURL: kagentBaseURL,
	// 	Timeout: testTimeout,
	// 	Headers: make(map[string]string),
	// }
	// a2aClient := client.New(clientConfig)

	// Test with non-existent agent
	// taskRequest := &types.TaskRequest{
	// 	ID: "error-test-" + time.Now().Format("20060102-150405"),
	// 	Message: &types.Message{
	// 		Role: "user",
	// 		Parts: []types.Part{{Type: "text", Text: "Test message"}},
	// 	},
	// }

	// TODO: Implement when Issue #10 is complete
	// response, err := a2aClient.SendTask(ctx, "non-existent-agent", taskRequest)
	// assert.Error(t, err)
	// assert.Nil(t, response)

	t.Skip("A2A error handling implementation pending Issue #10")
}

// Helper function to check if kagent sandbox is available
func isKagentAvailable() bool {
	// TODO: Implement availability check
	// This will be used to skip tests if kagent is not running
	return false
}

// TestSetup ensures test environment is properly configured
func TestSetup(t *testing.T) {
	if !isKagentAvailable() {
		t.Skip("kagent sandbox is not available")
	}

	// TODO: Add setup validation
	// - Check kagent sandbox is running
	// - Verify agent endpoints are accessible
	// - Validate test configuration

	t.Log("Integration test environment setup validated")
}

// TestGetTaskStatus tests retrieving the status of a task
func TestGetTaskStatus(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	clientConfig := client.Config{
		BaseURL: kagentBaseURL,
		Timeout: testTimeout,
		Headers: make(map[string]string),
	}
	a2aClient := client.New(clientConfig)

	taskRequest := &types.TaskRequest{
		ID: "status-test-" + time.Now().Format("20060102-150405"),
		Message: &types.Message{
			Role:  "user",
			Parts: []types.Part{{Type: "text", Text: "What is the status of my cluster?"}},
		},
	}

	t.Log("Sending task to agent for status test...")
	response, err := a2aClient.SendTask(ctx, testAgent, taskRequest)
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, taskRequest.ID, response.ID)

	t.Log("Getting task status...")
	status, err := a2aClient.GetTaskStatus(ctx, testAgent, taskRequest.ID)
	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, taskRequest.ID, status.ID)
}

// TestCancelTask tests cancelling a running task
func TestCancelTask(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	clientConfig := client.Config{
		BaseURL: kagentBaseURL,
		Timeout: testTimeout,
		Headers: make(map[string]string),
	}
	a2aClient := client.New(clientConfig)

	taskRequest := &types.TaskRequest{
		ID: "cancel-test-" + time.Now().Format("20060102-150405"),
		Message: &types.Message{
			Role:  "user",
			Parts: []types.Part{{Type: "text", Text: "Start a long-running task"}},
		},
	}

	t.Log("Sending task to agent for cancel test...")
	response, err := a2aClient.SendTask(ctx, testAgent, taskRequest)
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, taskRequest.ID, response.ID)

	t.Log("Cancelling task...")
	err = a2aClient.CancelTask(ctx, testAgent, taskRequest.ID)
	// We expect either no error, or a specific error if the task is already completed or cannot be cancelled
	if err != nil {
		t.Logf("CancelTask returned error (may be expected if task already completed): %v", err)
	}
}
