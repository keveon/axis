package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	serverListenGrpc    string
	serverEtcdEndpoints []string
)

var serverStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the axis server",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("starting axis server")
		fmt.Printf("  grpc listen: %s\n", serverListenGrpc)
		fmt.Printf("  etcd endpoints: %v\n", serverEtcdEndpoints)
		// Simulate running
		time.Sleep(time.Second)
		return nil
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage the axis server",
}

func init() {
	serverCmd.AddCommand(serverStartCmd)

	serverStartCmd.Flags().StringVar(&serverListenGrpc, "listen-grpc", ":9090", "gRPC listen address")
	serverStartCmd.Flags().StringSliceVar(&serverEtcdEndpoints, "etcd-endpoints", []string{"http://127.0.0.1:2379"}, "etcd endpoints")
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
