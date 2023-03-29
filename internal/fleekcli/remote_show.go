/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/vanilla-os/orchid/cmdr"
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
	repo, err := f.Repo()
	cobra.CheckErr(err)
	urls, err := repo.Remote()
	cobra.CheckErr(err)
	cmdr.Info.Println("configured:", f.config.Repository)
	cmdr.Info.Println("actual:", urls)
	return nil
}
