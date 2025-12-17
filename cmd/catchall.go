package cmd

import (
	"github.com/spf13/cobra"
)

var (
	catchallDomain     string
	catchallDynamicDir string
	catchallEntrypoint string
	catchallBackend    string
	catchallPriority   int
)

// catchallCmd represents the catchall command
var catchallCmd = &cobra.Command{
	Use:   "catchall",
	Short: "Manage Traefik catch-all routers",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	traefikCmd.AddCommand(catchallCmd)
	catchallCmd.PersistentFlags().StringVar(&catchallDomain, "domain", "cultofjoey.com", "Domain to protect with the wildcard catch-all")
	catchallCmd.PersistentFlags().StringVar(&catchallDynamicDir, "dynamic-dir", "/etc/traefik/dynamic", "Traefik dynamic configuration directory")
	catchallCmd.PersistentFlags().StringVar(&catchallEntrypoint, "entrypoint", "websecure", "Traefik entrypoint for created routers")
	catchallCmd.PersistentFlags().StringVar(&catchallBackend, "backend", "", "Backend URL (for example http://192.168.12.24:80)")
	catchallCmd.PersistentFlags().IntVar(&catchallPriority, "priority", 1, "Router priority (wildcard only)")
}
