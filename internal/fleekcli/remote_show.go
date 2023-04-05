/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/flake"
	"github.com/ublue-os/fleek/internal/ux"
)

func RepoShowCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   app.Trans("remoteshow.use"),
		Short: app.Trans("remoteshow.short"),
		Long:  app.Trans("remoteshow.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return show(cmd)
		},
	}
	return command
}

// initCmd represents the init command
func show(cmd *cobra.Command) error {
	ux.Description.Println(cmd.Short)
	err := mustConfig()
	if err != nil {
		return err
	}
	fl, err := flake.Load(cfg, app)
	if err != nil {
		return err
	}

	urls, err := fl.Remote()
	if err != nil {
		return err
	}
	ux.Info.Println(app.Trans("remoteshow.configured"), fl.Config.Repository)
	ux.Info.Println(app.Trans("remoteshow.actual"), urls)
	ux.Description.Println(app.Trans("remoteshow.done"))

	return nil
}
