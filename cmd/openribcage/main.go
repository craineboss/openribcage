// Package main provides the openribcage A2A protocol client CLI application.
//
// openribcage is an A2A (Agent2Agent) protocol client designed to enable
// avatar-based interfaces for AI agent coordination. This CLI provides
// commands for agent discovery, communication, and coordination via the
// standardized A2A protocol.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"

	"github.com/craine-io/openribcage/internal/config"
	"github.com/craine-io/openribcage/pkg/a2a/client"
)

var (
	// Version information (set by build)
	version = "dev"
	commit  = "unknown"
	date    = "unknown"

	// Global flags
	configFile string
	verbose    bool
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
and parsing their capabilities. Can discover single agents or scan
multiple endpoints for agent registry building.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Starting A2A agent discovery...")
		
		// TODO: Implement agent discovery
		// This will be implemented in Issue #3
		if len(args) > 0 {
			logrus.Infof("Discovering agent at: %s", args[0])
		} else {
			logrus.Info("Scanning for agents...")
		}
		
		fmt.Println("Agent discovery functionality coming in Issue #3!")
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

	// Execute CLI
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
}
