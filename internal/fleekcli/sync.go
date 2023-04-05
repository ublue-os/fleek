package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/flake"
	"github.com/ublue-os/fleek/internal/ux"
)

func SyncCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   app.Trans("sync.use"),
		Short: app.Trans("sync.short"),
		Long:  app.Trans("sync.long"),

		RunE: func(cmd *cobra.Command, args []string) error {
			return sync(cmd)
		},
	}
	return command
}

// initCmd represents the init command
func sync(cmd *cobra.Command) error {
	ux.Description.Println(cmd.Short)
	err := mustConfig()
	if err != nil {
		return err
	}
	fl, err := flake.Load(cfg, app)
	if err != nil {
		return err
	}
	err = fl.Sync(app.Trans("flake.syncMessage"))
	if err != nil {
		return err
	}
	ux.Success.Println(app.Trans("global.completed"))
	return nil
}
