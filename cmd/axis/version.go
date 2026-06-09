package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Version = "dev"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print axis version information",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("axis %s\n", Version)
		fmt.Printf("go %s\n", runtime.Version())
		fmt.Printf("%s/%s\n", runtime.GOOS, runtime.GOARCH)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
