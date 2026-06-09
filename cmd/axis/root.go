package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	cfgFile string
	cfg     *Config
)

// Config holds the full configuration.
type Config struct {
	Server ServerConfig `yaml:"server"`
	Agent  AgentConfig  `yaml:"agent"`
	Worker WorkerConfig `yaml:"worker"`
}

// ServerConfig holds server-specific configuration.
type ServerConfig struct {
	EtcdEndpoints []string `yaml:"etcdEndpoints"`
	ListenGrpc    string   `yaml:"listenGrpc"`
	LogLevel      string   `yaml:"logLevel"`
}

// AgentConfig holds agent-specific configuration.
type AgentConfig struct {
	ServerEndpoint string            `yaml:"serverEndpoint"`
	Labels         map[string]string `yaml:"labels"`
}

// WorkerConfig holds worker-specific configuration.
type WorkerConfig struct {
	StoreName string `yaml:"storeName"`
}

// rootCmd is the base command for the CLI.
var rootCmd = &cobra.Command{
	Use:   "axis",
	Short: "Axis — declarative IoT data platform",
	Long:  "Axis is a declarative IoT data platform driven by the control-plane model.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return loadConfig()
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/etc/axis/config.yaml", "path to config file")
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

// loadConfig reads the YAML config file (if it exists) and merges with defaults.
func loadConfig() error {
	cfg = &Config{
		Server: ServerConfig{
			EtcdEndpoints: []string{"http://127.0.0.1:2379"},
			ListenGrpc:    ":9090",
			LogLevel:      "info",
		},
		Agent: AgentConfig{
			ServerEndpoint: "localhost:9090",
			Labels:         map[string]string{},
		},
		Worker: WorkerConfig{
			StoreName: "default",
		},
	}

	data, err := os.ReadFile(cfgFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // no config file is fine
		}
		return fmt.Errorf("reading config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return fmt.Errorf("parsing config file: %w", err)
	}

	return nil
}
