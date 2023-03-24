/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewRepoAddCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		app.Trans("remoteadd.use"),
		app.Trans("remoteadd.long"),
		app.Trans("remoteadd.short"),
		remoteadd,
	).WithStringFlag(cmdr.NewStringFlag(
		"name",
		"n",
		app.Trans("reemoteadd.name"),
		"origin",
	))
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

// initCmd represents the init command
func remoteadd(cmd *cobra.Command, args []string) {

	f.config.Repository = args[0]
	err := f.config.Save()
	cobra.CheckErr(err)
	// now actually add the remote
	name := cmd.Flag("name").Value.String()
	repo, err := f.Repo()
	cobra.CheckErr(err)
	repo.RemoteAdd(args[0], name)
	err = repo.Commit()
	cobra.CheckErr(err)

	cmdr.Info.Println(f.config.Repository)
}
