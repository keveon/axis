package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	configRole           string
	configServerEndpoint string
	configLabels         map[string]string
	configOutput         string
)

// configCmd is the parent for config subcommands.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate and inspect axis configuration",
}

// configGenerateCmd generates a default config file for the given role.
var configGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a default configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := Config{
			Server: ServerConfig{
				EtcdEndpoints: []string{"http://127.0.0.1:2379"},
				ListenGrpc:    ":9090",
				LogLevel:      "info",
			},
			Agent: AgentConfig{
				ServerEndpoint: configServerEndpoint,
				Labels:         configLabels,
			},
			Worker: WorkerConfig{
				StoreName: "default",
			},
		}

		data, err := yaml.Marshal(&c)
		if err != nil {
			return fmt.Errorf("marshalling config: %w", err)
		}

		if configOutput == "" || configOutput == "-" {
			fmt.Println(string(data))
			return nil
		}

		if err := os.WriteFile(configOutput, data, 0644); err != nil {
			return fmt.Errorf("writing config file: %w", err)
		}

		fmt.Printf("config written to %s\n", configOutput)
		return nil
	},
}

// configShowCmd reads the loaded config and prints the effective configuration.
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the effective merged configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("marshalling config: %w", err)
		}
		fmt.Println(string(data))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configGenerateCmd)
	configCmd.AddCommand(configShowCmd)

	configGenerateCmd.Flags().StringVar(&configRole, "role", "server", "role to generate config for (server, agent, worker)")
	configGenerateCmd.Flags().StringVar(&configServerEndpoint, "server-endpoint", "localhost:9090", "server endpoint (for agent config)")
	configGenerateCmd.Flags().StringToStringVar(&configLabels, "labels", map[string]string{}, "agent labels in key=value format")
	configGenerateCmd.Flags().StringVar(&configOutput, "output", "-", "output file path (default stdout)")
}
