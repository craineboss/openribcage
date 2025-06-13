// Package types defines the A2A protocol types and data structures.
//
// This package contains all the types, structs, and interfaces used
// in A2A (Agent2Agent) protocol communication, including JSON-RPC 2.0
// message structures, AgentCard formats, and streaming response types.
package types

import (
	"encoding/json"
	"time"
)

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      interface{} `json:"id"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
	ID      interface{}     `json:"id"`
}

// JSONRPCError represents a JSON-RPC 2.0 error
type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// TaskRequest represents an A2A task request
type TaskRequest struct {
	ID      string   `json:"id"`
	Message *Message `json:"message"`
}

// TaskResponse represents an A2A task response
type TaskResponse struct {
	ID      string   `json:"id"`
	Message *Message `json:"message,omitempty"`
	Status  string   `json:"status"`
	Error   string   `json:"error,omitempty"`
}

// TaskStatus represents the status of an A2A task
type TaskStatus struct {
	ID          string    `json:"id"`
	Status      string    `json:"status"`
	Progress    float64   `json:"progress,omitempty"`
	StartedAt   time.Time `json:"started_at,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	Error       string    `json:"error,omitempty"`
}

// Message represents an A2A message with role and parts
type Message struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

// Part represents a message part (text, file, or data)
type Part struct {
	Type string      `json:"type"`
	Text string      `json:"text,omitempty"`
	Data interface{} `json:"data,omitempty"`
	File *FilePart   `json:"file,omitempty"`
}

// FilePart represents a file attachment in a message
type FilePart struct {
	Name     string `json:"name"`
	MimeType string `json:"mime_type"`
	Size     int64  `json:"size"`
	URL      string `json:"url,omitempty"`
	Content  []byte `json:"content,omitempty"`
}

// StreamResponse represents a streaming response from an A2A agent
type StreamResponse struct {
	ID        string      `json:"id"`
	Timestamp time.Time   `json:"timestamp"`
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
	Done      bool        `json:"done,omitempty"`
}

// AgentCard represents an A2A agent card (.well-known/agent.json)
type AgentCard struct {
	Name               string               `json:"name"`
	Description        string               `json:"description"`
	URL                string               `json:"url"`
	Version            string               `json:"version"`
	Capabilities       []string             `json:"capabilities"`
	Authentication     *AgentAuthentication `json:"authentication,omitempty"`
	DefaultInputModes  []string             `json:"defaultInputModes,omitempty"`
	DefaultOutputModes []string             `json:"defaultOutputModes,omitempty"`
	Skills             []AgentSkill         `json:"skills,omitempty"`
	Endpoints          []Endpoint           `json:"endpoints,omitempty"`
	Metadata           interface{}          `json:"metadata,omitempty"`
}

// Endpoint represents an A2A agent endpoint
type Endpoint struct {
	Type        string            `json:"type"`
	URL         string            `json:"url"`
	Methods     []string          `json:"methods"`
	Description string            `json:"description,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
}

// Agent represents a discovered and registered A2A agent
type Agent struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	URL          string      `json:"url"`
	Card         *AgentCard  `json:"card"`
	Status       AgentStatus `json:"status"`
	LastSeen     time.Time   `json:"last_seen"`
	DiscoveredAt time.Time   `json:"discovered_at"`
}

// AgentStatus represents the status of an agent
type AgentStatus string

const (
	AgentStatusOnline      AgentStatus = "online"
	AgentStatusOffline     AgentStatus = "offline"
	AgentStatusError       AgentStatus = "error"
	AgentStatusDiscovering AgentStatus = "discovering"
)

// A2AMethods contains the standard A2A protocol methods
var A2AMethods = struct {
	TasksSend     string
	TasksStream   string
	TasksStatus   string
	TasksCancel   string
	MessageSend   string
	MessageStream string
}{
	TasksSend:     "tasks/send",
	TasksStream:   "tasks/sendSubscribe",
	TasksStatus:   "tasks/status",
	TasksCancel:   "tasks/cancel",
	MessageSend:   "message/send",
	MessageStream: "message/stream",
}

// AgentAuthentication represents authentication requirements for an A2A agent
// (copied from agentcard.go)
type AgentAuthentication struct {
	Type   string                 `json:"type"`
	Config map[string]interface{} `json:"config,omitempty"`
}

// AgentSkill represents a skill provided by an A2A agent
// (copied from agentcard.go)
type AgentSkill struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
