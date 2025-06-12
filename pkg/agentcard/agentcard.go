// Package agentcard provides A2A AgentCard discovery and parsing.
//
// This package implements the AgentCard discovery system that enables
// openribcage to automatically find and register A2A-compliant agents
// via .well-known/agent.json endpoints.
package agentcard

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/craine-io/openribcage/pkg/a2a/types"
)

// Discoverer handles AgentCard discovery and validation
type Discoverer struct {
	client     *http.Client
	logger     *logrus.Logger
	timeout    time.Duration
	maxRetries int
	retryDelay time.Duration
}

// NewDiscoverer creates a new AgentCard discoverer
func NewDiscoverer(timeout time.Duration) *Discoverer {
	return &Discoverer{
		client: &http.Client{
			Timeout: timeout,
		},
		logger:     logrus.New(),
		timeout:    timeout,
		maxRetries: 3,
		retryDelay: time.Second * 2,
	}
}

// Discover discovers an AgentCard from an agent URL
func (d *Discoverer) Discover(ctx context.Context, agentURL string) (*types.AgentCard, error) {
	d.logger.Debugf("Discovering AgentCard from: %s", agentURL)

	// 1. Construct .well-known/agent.json URL
	agentCardURL := BuildAgentCardURL(agentURL)
	d.logger.Debugf("AgentCard URL: %s", agentCardURL)

	// 2. Make HTTP GET request with retry logic
	data, err := d.fetchWithRetry(ctx, agentCardURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch AgentCard from %s: %w", agentCardURL, err)
	}

	// 3. Parse JSON response into AgentCard
	var card types.AgentCard
	if err := json.Unmarshal(data, &card); err != nil {
		return nil, fmt.Errorf("failed to parse AgentCard JSON: %w", err)
	}

	// 4. Validate AgentCard format
	if err := d.Validate(&card); err != nil {
		return nil, fmt.Errorf("AgentCard validation failed: %w", err)
	}

	d.logger.Infof("Successfully discovered AgentCard: %s (version: %s)", card.Name, card.Version)
	return &card, nil
}

// fetchWithRetry performs HTTP GET with retry logic
func (d *Discoverer) fetchWithRetry(ctx context.Context, url string) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= d.maxRetries; attempt++ {
		if attempt > 0 {
			d.logger.Debugf("Retry attempt %d/%d for %s", attempt, d.maxRetries, url)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(d.retryDelay):
				// Continue with retry
			}
		}

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %w", err)
			continue
		}

		// Set appropriate headers for AgentCard discovery
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", "openribcage/1.0 (A2A-Protocol-Client)")

		resp, err := d.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("HTTP request failed: %w", err)
			continue
		}
		defer resp.Body.Close()

		// Check for successful response
		if resp.StatusCode == http.StatusOK {
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				lastErr = fmt.Errorf("failed to read response body: %w", err)
				continue
			}
			return data, nil
		}

		// Handle specific HTTP status codes
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, fmt.Errorf("AgentCard not found (404) at %s", url)
		case http.StatusUnauthorized:
			return nil, fmt.Errorf("unauthorized access (401) to %s", url)
		case http.StatusForbidden:
			return nil, fmt.Errorf("forbidden access (403) to %s", url)
		default:
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
		}

		// Don't retry on client errors (4xx)
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return nil, lastErr
		}
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", d.maxRetries+1, lastErr)
}

// Validate validates an AgentCard format and content
func (d *Discoverer) Validate(card *types.AgentCard) error {
	d.logger.Debugf("Validating AgentCard: %s", card.Name)

	// 1. Check required fields
	if card.Name == "" {
		return fmt.Errorf("agent name is required")
	}
	if card.Version == "" {
		return fmt.Errorf("agent version is required")
	}

	// 2. Validate each endpoint (only if endpoints exist)
	for i, endpoint := range card.Endpoints {
		if err := d.validateEndpoint(&endpoint); err != nil {
			return fmt.Errorf("invalid endpoint %d: %w", i, err)
		}
	}

	// 3. Validate capabilities format (basic check)
	// Capabilities is now an object, not a slice. Check that at least one field is present (all are bools, so just check for presence).
	// Optionally, you could check that the struct is not the zero value (all false), but that's not required by spec.
	// So, no loop needed here.

	d.logger.Debugf("AgentCard validation successful: %s", card.Name)
	return nil
}

// validateEndpoint validates a single endpoint
func (d *Discoverer) validateEndpoint(endpoint *types.Endpoint) error {
	// Validate URL format
	if endpoint.URL == "" {
		return fmt.Errorf("endpoint URL is required")
	}

	parsedURL, err := url.Parse(endpoint.URL)
	if err != nil {
		return fmt.Errorf("invalid endpoint URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("endpoint URL must use http or https scheme")
	}

	// Validate endpoint type
	validTypes := []string{"a2a", "streaming", "webhook"}
	if !contains(validTypes, endpoint.Type) {
		return fmt.Errorf("unsupported endpoint type: %s (supported: %v)", endpoint.Type, validTypes)
	}

	// Validate A2A methods for a2a endpoints
	if endpoint.Type == "a2a" {
		for _, method := range endpoint.Methods {
			if !isValidA2AMethod(method) {
				return fmt.Errorf("invalid A2A method: %s", method)
			}
		}
	}

	return nil
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// isValidA2AMethod checks if a method is a valid A2A protocol method
func isValidA2AMethod(method string) bool {
	validMethods := []string{
		types.A2AMethods.TasksSend,
		types.A2AMethods.TasksStream,
		types.A2AMethods.TasksStatus,
		types.A2AMethods.TasksCancel,
		types.A2AMethods.MessageSend,
		types.A2AMethods.MessageStream,
	}
	return contains(validMethods, method)
}

// Parse parses AgentCard JSON data
func (d *Discoverer) Parse(data []byte) (*types.AgentCard, error) {
	var card types.AgentCard
	if err := json.Unmarshal(data, &card); err != nil {
		return nil, fmt.Errorf("failed to parse AgentCard JSON: %w", err)
	}

	if err := d.Validate(&card); err != nil {
		return nil, fmt.Errorf("AgentCard validation failed: %w", err)
	}

	return &card, nil
}

// BuildAgentCardURL constructs the AgentCard URL from a base agent URL
func BuildAgentCardURL(agentURL string) string {
	// Clean and parse the URL
	agentURL = strings.TrimSpace(agentURL)
	if agentURL == "" {
		return ""
	}

	// Ensure URL has a scheme
	if !strings.HasPrefix(agentURL, "http://") && !strings.HasPrefix(agentURL, "https://") {
		agentURL = "http://" + agentURL
	}

	// Parse URL to handle it properly
	parsedURL, err := url.Parse(agentURL)
	if err != nil {
		// Fallback to simple string concatenation
		return strings.TrimSuffix(agentURL, "/") + "/.well-known/agent.json"
	}

	// Construct the .well-known path
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/") + "/.well-known/agent.json"
	return parsedURL.String()
}

// Init initializes the agentcard package
func Init() error {
	// Package initialization - currently no special setup needed
	return nil
}
