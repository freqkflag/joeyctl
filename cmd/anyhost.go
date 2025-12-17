package cmd

import (
	"github.com/freqkflag/joeyctl/internal/run"
	"github.com/freqkflag/joeyctl/internal/traefik"
	"github.com/spf13/cobra"
)

// anyhostCmd represents the anyhost command
var anyhostCmd = &cobra.Command{
	Use:   "anyhost",
	Short: "Write a catch-all router for any Host header",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}

		runner := run.Runner{DryRun: DryRun}
		msg, err := traefik.WriteCatchall(runner, traefik.CatchallOptions{
			DynamicDir: catchallDynamicDir,
			EntryPoint: catchallEntrypoint,
			BackendURL: catchallBackend,
			Mode:       traefik.ModeAnyhost,
		})
		if err != nil {
			return err
		}

		cmd.Println(msg)
		return nil
	},
}

func init() {
	catchallCmd.AddCommand(anyhostCmd)
	_ = anyhostCmd.MarkPersistentFlagRequired("backend")
}
