package cmd

import (
	"github.com/freqkflag/joeyctl/internal/run"
	"github.com/freqkflag/joeyctl/internal/traefik"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove any existing catch-all configuration files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}

		runner := run.Runner{DryRun: DryRun}
		msg, err := traefik.RemoveCatchalls(runner, catchallDynamicDir, catchallDomain)
		if err != nil {
			return err
		}

		cmd.Println(msg)
		return nil
	},
}

func init() {
	catchallCmd.AddCommand(removeCmd)
}
