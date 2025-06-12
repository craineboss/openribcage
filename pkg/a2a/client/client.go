// Package client provides A2A protocol client functionality.
//
// This package implements the core A2A (Agent2Agent) protocol client that
// enables openribcage to communicate with A2A-compliant agent frameworks.
//
// TODO: Full implementation will be completed in Issue #10
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/craine-io/openribcage/pkg/a2a/types"
)

// Config holds A2A client configuration
type Config struct {
	BaseURL string            `json:"base_url"`
	Timeout time.Duration     `json:"timeout"`
	Headers map[string]string `json:"headers"`
}

// Client represents an A2A protocol client
type Client struct {
	config     Config
	logger     *logrus.Logger
	httpClient *http.Client
}

// New creates a new A2A protocol client
func New(config Config) *Client {
	return &Client{
		config: config,
		logger: logrus.New(),
		httpClient: &http.Client{
			Timeout: config.Timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
	}
}

// Init initializes the A2A client package
// This function is called from cmd/openribcage/main.go
func Init() error {
	logrus.Debug("A2A client package initialized")
	return nil
}

// SendMessage sends a message to an A2A agent
func (c *Client) SendMessage(ctx context.Context, agentID string, message *types.Message) (*types.TaskResponse, error) {
	// Construct the request URL
	url := fmt.Sprintf("%s/%s", c.config.BaseURL, agentID)

	// Create the JSON-RPC request
	req := &types.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  types.A2AMethods.MessageSend,
		Params: map[string]interface{}{
			"message": message,
		},
		ID: uuid.New().String(), // Generate a unique request ID
	}

	// Marshal the request to JSON
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range c.config.Headers {
		httpReq.Header.Set(k, v)
	}

	// Send the request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read and parse response
	var jsonResp types.JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for JSON-RPC error
	if jsonResp.Error != nil {
		return nil, fmt.Errorf("JSON-RPC error: %s (code: %d)", jsonResp.Error.Message, jsonResp.Error.Code)
	}

	// Parse the result into a TaskResponse
	var taskResp types.TaskResponse
	if err := json.Unmarshal(jsonResp.Result, &taskResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task response: %w", err)
	}

	return &taskResp, nil
}

// SendTask sends a task to an A2A agent
func (c *Client) SendTask(ctx context.Context, agentID string, req *types.TaskRequest) (*types.TaskResponse, error) {
	// Construct the request URL
	url := fmt.Sprintf("%s/%s", c.config.BaseURL, agentID)

	// Create the JSON-RPC request
	jsonReq := &types.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  types.A2AMethods.TasksSend,
		Params: map[string]interface{}{
			"id":      req.ID,
			"message": req.Message,
		},
		ID: uuid.New().String(),
	}

	// Marshal the request to JSON
	reqBody, err := json.Marshal(jsonReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range c.config.Headers {
		httpReq.Header.Set(k, v)
	}

	// Send the request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read and parse response
	var jsonResp types.JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for JSON-RPC error
	if jsonResp.Error != nil {
		return nil, fmt.Errorf("JSON-RPC error: %s (code: %d)", jsonResp.Error.Message, jsonResp.Error.Code)
	}

	// Parse the result into a TaskResponse
	var taskResp types.TaskResponse
	if err := json.Unmarshal(jsonResp.Result, &taskResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task response: %w", err)
	}

	return &taskResp, nil
}

// StreamTask sends a task with streaming response
// TODO: Implement in Issue #10
func (c *Client) StreamTask(ctx context.Context, agentID string, req *types.TaskRequest) (<-chan *types.StreamResponse, <-chan error) {
	responseChan := make(chan *types.StreamResponse)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errorChan)
		errorChan <- fmt.Errorf("StreamTask not yet implemented - see Issue #10")
	}()

	return responseChan, errorChan
}

// GetTaskStatus retrieves the status of a task
func (c *Client) GetTaskStatus(ctx context.Context, agentID, taskID string) (*types.TaskStatus, error) {
	// Construct the request URL
	url := fmt.Sprintf("%s/%s", c.config.BaseURL, agentID)

	// Create the JSON-RPC request
	jsonReq := &types.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  types.A2AMethods.TasksStatus,
		Params: map[string]interface{}{
			"id": taskID,
		},
		ID: uuid.New().String(),
	}

	// Marshal the request to JSON
	reqBody, err := json.Marshal(jsonReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range c.config.Headers {
		httpReq.Header.Set(k, v)
	}

	// Send the request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read and parse response
	var jsonResp types.JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for JSON-RPC error
	if jsonResp.Error != nil {
		return nil, fmt.Errorf("JSON-RPC error: %s (code: %d)", jsonResp.Error.Message, jsonResp.Error.Code)
	}

	// Parse the result into a TaskStatus
	var taskStatus types.TaskStatus
	if err := json.Unmarshal(jsonResp.Result, &taskStatus); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task status: %w", err)
	}

	return &taskStatus, nil
}

// CancelTask cancels a running task
func (c *Client) CancelTask(ctx context.Context, agentID, taskID string) error {
	// Construct the request URL
	url := fmt.Sprintf("%s/%s", c.config.BaseURL, agentID)

	// Create the JSON-RPC request
	jsonReq := &types.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  types.A2AMethods.TasksCancel,
		Params: map[string]interface{}{
			"id": taskID,
		},
		ID: uuid.New().String(),
	}

	// Marshal the request to JSON
	reqBody, err := json.Marshal(jsonReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range c.config.Headers {
		httpReq.Header.Set(k, v)
	}

	// Send the request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read and parse response
	var jsonResp types.JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Check for JSON-RPC error
	if jsonResp.Error != nil {
		return fmt.Errorf("JSON-RPC error: %s (code: %d)", jsonResp.Error.Message, jsonResp.Error.Code)
	}

	return nil
}

// Ping tests connectivity to an A2A agent
// TODO: Implement in Issue #10
func (c *Client) Ping(ctx context.Context, agentURL string) error {
	return fmt.Errorf("Ping not yet implemented - see Issue #10")
}
