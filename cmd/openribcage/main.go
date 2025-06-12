// Package main provides the openribcage A2A protocol client CLI application.
//
// openribcage is an A2A (Agent2Agent) protocol client designed to enable
// avatar-based interfaces for AI agent coordination. This CLI provides
// commands for agent discovery, communication, and coordination via the
// standardized A2A protocol.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/craine-io/openribcage/internal/config"
	"github.com/craine-io/openribcage/pkg/a2a/client"
	"github.com/craine-io/openribcage/pkg/agentcard"
)

var (
	// Version information (set by build)
	version = "dev"
	commit  = "unknown"
	date    = "unknown"

	// Global flags
	configFile string
	verbose    bool

	// Discovery flags
	discoveryTimeout time.Duration
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "openribcage",
	Short: "A2A Protocol Client for Avatar Interfaces",
	Long: `openribcage is an A2A (Agent2Agent) protocol client designed to enable
avatar-based interfaces for AI agent coordination.

It provides standardized agent discovery, communication, and orchestration
capabilities through the proven A2A protocol standard, enabling natural
conversation with AI agents through avatar interfaces.`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Set up logging
		if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		}
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	},
}

// discoverCmd represents the discover command
var discoverCmd = &cobra.Command{
	Use:   "discover [agent-url]",
	Short: "Discover A2A agents and their capabilities",
	Long: `Discover A2A-compliant agents by scanning for AgentCard endpoints
and parsing their capabilities. Returns agent information in JSON format.

Examples:
  # Discover a single agent
  openribcage discover http://localhost:8083/api/a2a/kagent/k8s-agent
  
  # Discover with verbose logging
  openribcage discover -v http://localhost:8083/api/a2a/kagent/k8s-agent`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		agentURL := args[0]

		logrus.Infof("Discovering A2A agent at: %s", agentURL)

		// Create AgentCard discoverer
		discoverer := agentcard.NewDiscoverer(discoveryTimeout)

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), discoveryTimeout)
		defer cancel()

		// Discover the agent
		card, err := discoverer.Discover(ctx, agentURL)
		if err != nil {
			logrus.Errorf("Agent discovery failed: %v", err)
			os.Exit(1)
		}

		// Output AgentCard as JSON
		output, err := json.MarshalIndent(card, "", "  ")
		if err != nil {
			logrus.Errorf("Failed to marshal AgentCard to JSON: %v", err)
			os.Exit(1)
		}

		fmt.Println(string(output))
		logrus.Infof("Successfully discovered agent: %s (version: %s)", card.Name, card.Version)
	},
}

// communicateCmd represents the communicate command
var communicateCmd = &cobra.Command{
	Use:   "communicate [agent-url] [message]",
	Short: "Communicate with A2A agents",
	Long: `Send messages to A2A agents using the JSON-RPC 2.0 protocol.
Supports both single messages and streaming communication patterns.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		agentURL := args[0]
		message := args[1]

		logrus.Infof("Communicating with agent: %s", agentURL)
		logrus.Infof("Message: %s", message)

		// TODO: Implement A2A communication
		// This will be implemented in Issue #4 and #10
		fmt.Println("A2A communication functionality coming in Issue #10!")
	},
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start openribcage A2A client server",
	Long: `Start the openribcage server to provide A2A client services
for avatar interfaces. Exposes REST API and WebSocket endpoints
for real-time agent communication.`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Starting openribcage A2A client server...")

		// TODO: Implement server mode
		// This will provide API for avatar interfaces
		fmt.Println("Server mode functionality coming in Phase 2!")
	},
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.openribcage.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Discovery command flags
	discoverCmd.Flags().DurationVar(&discoveryTimeout, "timeout", 30*time.Second, "discovery timeout duration")

	// Add subcommands
	rootCmd.AddCommand(discoverCmd)
	rootCmd.AddCommand(communicateCmd)
	rootCmd.AddCommand(serveCmd)
}

func main() {
	// Initialize configuration
	if err := config.Init(configFile); err != nil {
		logrus.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Initialize A2A client
	if err := client.Init(); err != nil {
		logrus.Fatalf("Failed to initialize A2A client: %v", err)
	}

	// Initialize AgentCard package
	if err := agentcard.Init(); err != nil {
		logrus.Fatalf("Failed to initialize AgentCard package: %v", err)
	}

	// Execute CLI
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
}
