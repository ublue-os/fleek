/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/ux"
)

func RepoShowCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   app.Trans("remoteshow.use"),
		Short: app.Trans("remoteshow.short"),
		Long:  app.Trans("remoteshow.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return show(cmd, args)
		},
	}
	return command
}

// initCmd represents the init command
func show(cmd *cobra.Command, args []string) error {
	ux.Description.Println(cmd.Short)

	repo, err := f.Repo()
	if err != nil {
		return err
	}
	urls, err := repo.Remote()
	if err != nil {
		return err
	}
	ux.Info.Println(app.Trans("remoteshow.configured"), f.config.Repository)
	ux.Info.Println(app.Trans("remoteshow.actual"), urls)
	ux.Description.Println(app.Trans("remoteshow.done"))

	return nil
}
