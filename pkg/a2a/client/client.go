// Package client provides A2A protocol client functionality.
//
// This package implements the core A2A (Agent2Agent) protocol client that
// enables openribcage to communicate with A2A-compliant agent frameworks.
//
// TODO: Full implementation will be completed in Issue #10
package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	// TODO: Add HTTP client, connection pool, etc. in Issue #10
}

// New creates a new A2A protocol client
func New(config Config) *Client {
	return &Client{
		config: config,
		logger: logrus.New(),
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// Init initializes the A2A client package
// This function is called from cmd/openribcage/main.go
func Init() error {
	// TODO: Implement package initialization in Issue #10
	// This might include:
	// - Setting up default configurations
	// - Initializing connection pools
	// - Registering A2A method handlers
	logrus.Debug("A2A client package initialized (scaffolding)")
	return nil
}

// SendTask sends a task to an A2A agent
// TODO: Implement in Issue #10
func (c *Client) SendTask(ctx context.Context, agentID string, req *types.TaskRequest) (*types.TaskResponse, error) {
	return nil, fmt.Errorf("SendTask not yet implemented - see Issue #10")
}

// StreamTask sends a task with streaming response
func (c *Client) StreamTask(ctx context.Context, agentID string, req *types.TaskRequest) (<-chan *types.StreamResponse, <-chan error) {
	out := make(chan *types.StreamResponse)
	errs := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errs)

		// Construct the request URL
		url := fmt.Sprintf("%s/%s", c.config.BaseURL, agentID)

		// Create the JSON-RPC request
		jsonReq := &types.JSONRPCRequest{
			JSONRPC: "2.0",
			Method:  types.A2AMethods.TasksStream,
			Params: map[string]interface{}{
				"id":      req.ID,
				"message": req.Message,
			},
			ID: uuid.New().String(),
		}

		reqBody, err := json.Marshal(jsonReq)
		if err != nil {
			errs <- fmt.Errorf("failed to marshal request: %w", err)
			return
		}

		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
		if err != nil {
			errs <- fmt.Errorf("failed to create request: %w", err)
			return
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Accept", "text/event-stream")
		for k, v := range c.config.Headers {
			httpReq.Header.Set(k, v)
		}

		resp, err := c.httpClient.Do(httpReq)
		if err != nil {
			errs <- fmt.Errorf("request failed: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errs <- fmt.Errorf("unexpected status: %s", resp.Status)
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data:") {
				data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
				if data == "" {
					continue
				}
				var streamResp types.StreamResponse
				if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
					errs <- fmt.Errorf("failed to unmarshal stream response: %w", err)
					return
				}
				select {
				case out <- &streamResp:
				case <-ctx.Done():
					errs <- ctx.Err()
					return
				}
			}
		}
		if err := scanner.Err(); err != nil {
			errs <- fmt.Errorf("scanner error: %w", err)
		}
	}()

	return out, errs
}

// GetTaskStatus retrieves the status of a task
// TODO: Implement in Issue #10
func (c *Client) GetTaskStatus(ctx context.Context, agentID, taskID string) (*types.TaskStatus, error) {
	return nil, fmt.Errorf("GetTaskStatus not yet implemented - see Issue #10")
}

// CancelTask cancels a running task
// TODO: Implement in Issue #10
func (c *Client) CancelTask(ctx context.Context, agentID, taskID string) error {
	return fmt.Errorf("CancelTask not yet implemented - see Issue #10")
}

// Ping tests connectivity to an A2A agent
// TODO: Implement in Issue #10
func (c *Client) Ping(ctx context.Context, agentURL string) error {
	return fmt.Errorf("Ping not yet implemented - see Issue #10")
}
