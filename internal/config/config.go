// Package config provides configuration management for openribcage.
//
// This package handles loading and managing configuration from various
// sources including files, environment variables, and command line flags.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// Config holds the application configuration
type Config struct {
	// Server configuration
	Server ServerConfig `yaml:"server" json:"server"`
	
	// A2A client configuration
	A2A A2AConfig `yaml:"a2a" json:"a2a"`
	
	// Logging configuration
	Logging LoggingConfig `yaml:"logging" json:"logging"`
	
	// Registry configuration
	Registry RegistryConfig `yaml:"registry" json:"registry"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Host         string        `yaml:"host" json:"host"`
	Port         int           `yaml:"port" json:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" json:"write_timeout"`
	TLS          TLSConfig     `yaml:"tls" json:"tls"`
}

// A2AConfig holds A2A client configuration
type A2AConfig struct {
	Timeout         time.Duration     `yaml:"timeout" json:"timeout"`
	RetryAttempts   int               `yaml:"retry_attempts" json:"retry_attempts"`
	RetryDelay      time.Duration     `yaml:"retry_delay" json:"retry_delay"`
	DefaultHeaders  map[string]string `yaml:"default_headers" json:"default_headers"`
	StreamTimeout   time.Duration     `yaml:"stream_timeout" json:"stream_timeout"`
	DiscoveryHosts  []string          `yaml:"discovery_hosts" json:"discovery_hosts"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level" json:"level"`
	Format string `yaml:"format" json:"format"`
	Output string `yaml:"output" json:"output"`
}

// RegistryConfig holds agent registry configuration
type RegistryConfig struct {
	CleanupInterval time.Duration `yaml:"cleanup_interval" json:"cleanup_interval"`
	StaleThreshold  time.Duration `yaml:"stale_threshold" json:"stale_threshold"`
	MaxAgents       int           `yaml:"max_agents" json:"max_agents"`
}

// TLSConfig holds TLS configuration
type TLSConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	CertFile string `yaml:"cert_file" json:"cert_file"`
	KeyFile  string `yaml:"key_file" json:"key_file"`
}

var (
	// Global configuration instance
	globalConfig *Config
	logger       = logrus.New()
)

// Init initializes the configuration system
func Init(configFile string) error {
	// Set default configuration
	globalConfig = &Config{
		Server: ServerConfig{
			Host:         "localhost",
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			TLS: TLSConfig{
				Enabled: false,
			},
		},
		A2A: A2AConfig{
			Timeout:        30 * time.Second,
			RetryAttempts:  3,
			RetryDelay:     1 * time.Second,
			StreamTimeout:  5 * time.Minute,
			DefaultHeaders: make(map[string]string),
			DiscoveryHosts: []string{},
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "text",
			Output: "stdout",
		},
		Registry: RegistryConfig{
			CleanupInterval: 5 * time.Minute,
			StaleThreshold:  10 * time.Minute,
			MaxAgents:       100,
		},
	}

	// Load configuration file if specified
	if configFile != "" {
		if err := loadConfigFile(configFile); err != nil {
			return fmt.Errorf("failed to load config file: %w", err)
		}
	} else {
		// Try to load from default locations
		defaultPaths := []string{
			"./openribcage.yaml",
			"./config/openribcage.yaml",
			filepath.Join(os.Getenv("HOME"), ".openribcage.yaml"),
			"/etc/openribcage/config.yaml",
		}

		for _, path := range defaultPaths {
			if _, err := os.Stat(path); err == nil {
				if err := loadConfigFile(path); err != nil {
					logger.Warnf("Failed to load config from %s: %v", path, err)
				} else {
					logger.Infof("Loaded configuration from: %s", path)
					break
				}
			}
		}
	}

	// Override with environment variables
	loadEnvironmentVariables()

	return nil
}

// Get returns the global configuration
func Get() *Config {
	return globalConfig
}

// loadConfigFile loads configuration from a YAML file
func loadConfigFile(filename string) error {
	// TODO: Implement YAML configuration loading
	// This will be implemented when needed
	// Use gopkg.in/yaml.v3 for YAML parsing
	return fmt.Errorf("YAML config loading not yet implemented")
}

// loadEnvironmentVariables loads configuration from environment variables
func loadEnvironmentVariables() {
	// Server configuration
	if host := os.Getenv("OPENRIBCAGE_HOST"); host != "" {
		globalConfig.Server.Host = host
	}

	// A2A configuration
	if timeout := os.Getenv("OPENRIBCAGE_A2A_TIMEOUT"); timeout != "" {
		if duration, err := time.ParseDuration(timeout); err == nil {
			globalConfig.A2A.Timeout = duration
		}
	}

	// Logging configuration
	if level := os.Getenv("OPENRIBCAGE_LOG_LEVEL"); level != "" {
		globalConfig.Logging.Level = level
	}

	// Add more environment variable mappings as needed
}
