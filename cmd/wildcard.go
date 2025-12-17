package cmd

import (
	"github.com/freqkflag/joeyctl/internal/run"
	"github.com/freqkflag/joeyctl/internal/traefik"
	"github.com/spf13/cobra"
)

// wildcardCmd represents the wildcard command
var wildcardCmd = &cobra.Command{
	Use:   "wildcard",
	Short: "Write a wildcard catch-all router for a specific domain",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := requireRoot(); err != nil {
			return err
		}

		runner := run.Runner{DryRun: DryRun}
		msg, err := traefik.WriteCatchall(runner, traefik.CatchallOptions{
			Domain:     catchallDomain,
			DynamicDir: catchallDynamicDir,
			EntryPoint: catchallEntrypoint,
			BackendURL: catchallBackend,
			Priority:   catchallPriority,
			Mode:       traefik.ModeWildcard,
		})
		if err != nil {
			return err
		}

		cmd.Println(msg)
		return nil
	},
}

func init() {
	catchallCmd.AddCommand(wildcardCmd)
	_ = wildcardCmd.MarkPersistentFlagRequired("backend")
}
