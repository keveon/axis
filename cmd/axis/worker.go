package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	workerStoreName string
)

var workerStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the axis worker",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("starting axis worker")
		fmt.Printf("  store name: %s\n", workerStoreName)
		return nil
	},
}

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Manage the axis worker",
}

func init() {
	workerCmd.AddCommand(workerStartCmd)

	workerStartCmd.Flags().StringVar(&workerStoreName, "store-name", "default", "store name")
}

func init() {
	rootCmd.AddCommand(workerCmd)
}
