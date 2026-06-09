package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	agentServerEndpoint string
	agentLabels         map[string]string
)

var agentStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the axis agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("starting axis agent")
		fmt.Printf("  server endpoint: %s\n", agentServerEndpoint)
		fmt.Printf("  labels: %v\n", agentLabels)
		return nil
	},
}

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage the axis agent",
}

func init() {
	agentCmd.AddCommand(agentStartCmd)

	agentStartCmd.Flags().StringVar(&agentServerEndpoint, "server-endpoint", "localhost:9090", "axis server endpoint")
	agentStartCmd.Flags().StringToStringVar(&agentLabels, "labels", map[string]string{}, "agent labels in key=value format")
}

func init() {
	rootCmd.AddCommand(agentCmd)
}
