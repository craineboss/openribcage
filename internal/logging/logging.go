// Package logging provides structured logging for openribcage.
//
// This package sets up structured logging with proper formatting,
// log levels, and output destinations for the A2A protocol client.
package logging

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus with additional A2A-specific functionality
type Logger struct {
	*logrus.Logger
}

// NewLogger creates a new structured logger
func NewLogger(level, format, output string) (*Logger, error) {
	logger := logrus.New()

	// Set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.SetLevel(logLevel)

	// Set formatter
	switch strings.ToLower(format) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z",
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
		})
	case "text", "":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     false,
		})
	default:
		return nil, fmt.Errorf("unsupported log format: %s", format)
	}

	// Set output destination
	var writer io.Writer
	switch strings.ToLower(output) {
	case "stdout", "":
		writer = os.Stdout
	case "stderr":
		writer = os.Stderr
	default:
		// Try to open as file
		file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		writer = file
	}
	logger.SetOutput(writer)

	return &Logger{Logger: logger}, nil
}

// WithA2AContext adds A2A-specific context to log entries
func (l *Logger) WithA2AContext(agentID, taskID, method string) *logrus.Entry {
	return l.WithFields(logrus.Fields{
		"agent_id": agentID,
		"task_id":  taskID,
		"a2a_method": method,
		"component": "a2a-client",
	})
}

// WithAgentContext adds agent-specific context to log entries
func (l *Logger) WithAgentContext(agentID, agentName, agentURL string) *logrus.Entry {
	return l.WithFields(logrus.Fields{
		"agent_id":   agentID,
		"agent_name": agentName,
		"agent_url":  agentURL,
		"component":  "agent-registry",
	})
}

// WithStreamContext adds streaming context to log entries
func (l *Logger) WithStreamContext(streamID, agentID string) *logrus.Entry {
	return l.WithFields(logrus.Fields{
		"stream_id": streamID,
		"agent_id":  agentID,
		"component": "streaming",
	})
}

// WithAvatarContext adds avatar interface context to log entries
func (l *Logger) WithAvatarContext(avatarID, personaID string) *logrus.Entry {
	return l.WithFields(logrus.Fields{
		"avatar_id":  avatarID,
		"persona_id": personaID,
		"component":  "avatar-integration",
	})
}

// LogA2ARequest logs an A2A protocol request
func (l *Logger) LogA2ARequest(method, agentURL, taskID string) {
	l.WithA2AContext("", taskID, method).Infof("A2A request: %s -> %s", method, agentURL)
}

// LogA2AResponse logs an A2A protocol response
func (l *Logger) LogA2AResponse(method, agentURL, taskID string, success bool, duration string) {
	entry := l.WithA2AContext("", taskID, method).WithField("duration", duration)
	if success {
		entry.Infof("A2A response success: %s <- %s", method, agentURL)
	} else {
		entry.Errorf("A2A response error: %s <- %s", method, agentURL)
	}
}

// LogAgentDiscovery logs agent discovery events
func (l *Logger) LogAgentDiscovery(agentURL string, success bool, capabilities []string) {
	entry := l.WithField("agent_url", agentURL).WithField("capabilities", capabilities)
	if success {
		entry.Infof("Agent discovered successfully: %s", agentURL)
	} else {
		entry.Errorf("Agent discovery failed: %s", agentURL)
	}
}
