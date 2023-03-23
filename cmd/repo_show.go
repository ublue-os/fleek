/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewRepoShowCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		app.Trans("reposhow.use"),
		app.Trans("reposhow.long"),
		app.Trans("reposhow.short"),
		show,
	)
	return cmd
}

// initCmd represents the init command
func show(cmd *cobra.Command, args []string) {

	repo, err := f.Repo()
	cobra.CheckErr(err)
	urls, err := repo.Remote()
	cobra.CheckErr(err)
	cmdr.Info.Println("configured:", f.config.Repository)
	cmdr.Info.Println("actual:", urls)

}
