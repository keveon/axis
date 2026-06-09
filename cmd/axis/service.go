package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

var serviceRole string

const systemdUnitTemplate = `[Unit]
Description=Axis {{.Role}}
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/axis {{.Role}} start --config /etc/axis/config.yaml
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
`

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage systemd service for axis",
}

var serviceInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install systemd unit file for the given role",
	RunE: func(cmd *cobra.Command, args []string) error {
		unitPath := fmt.Sprintf("/etc/systemd/system/axis-%s.service", serviceRole)

		tmpl, err := template.New("unit").Parse(systemdUnitTemplate)
		if err != nil {
			return fmt.Errorf("parsing template: %w", err)
		}

		f, err := os.Create(unitPath)
		if err != nil {
			return fmt.Errorf("creating unit file: %w", err)
		}
		defer f.Close()

		if err := tmpl.Execute(f, map[string]string{"Role": serviceRole}); err != nil {
			return fmt.Errorf("writing unit file: %w", err)
		}

		fmt.Printf("installed %s\n", unitPath)
		return nil
	},
}

var serviceUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove systemd unit file for the given role",
	RunE: func(cmd *cobra.Command, args []string) error {
		unitPath := fmt.Sprintf("/etc/systemd/system/axis-%s.service", serviceRole)

		if err := os.Remove(unitPath); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("unit file not found: %s", unitPath)
			}
			return fmt.Errorf("removing unit file: %w", err)
		}

		fmt.Printf("removed %s\n", unitPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(serviceInstallCmd)
	serviceCmd.AddCommand(serviceUninstallCmd)

	serviceCmd.PersistentFlags().StringVar(&serviceRole, "role", "server", "role (server, agent, worker)")
}
