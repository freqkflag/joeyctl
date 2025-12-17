package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// at top-level in cmd/root.go
var (
	DryRun bool
	Yes    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:          "joeyctl",
	Short:        "Homelab automation helpers",
	SilenceUsage: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&DryRun, "dry-run", false, "Print actions only—do not change files or restart services")
	rootCmd.PersistentFlags().BoolVar(&Yes, "yes", false, "Automatically confirm prompts (reserved for future use)")
}

func requireRoot() error {
	if os.Geteuid() != 0 {
		return fmt.Errorf("this command must be run as root")
	}
	return nil
}
