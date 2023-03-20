/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/core"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewRepoAddCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		fleek.Trans("repoadd.use"),
		fleek.Trans("repoadd.long"),
		fleek.Trans("repoadd.short"),
		remoteadd,
	).WithStringFlag(cmdr.NewStringFlag(
		"name",
		"n",
		fleek.Trans("repoadd.name"),
		"origin",
	))
	cmd.Args = cobra.ExactArgs(1)
	return cmd
}

// initCmd represents the init command
func remoteadd(cmd *cobra.Command, args []string) {

	conf, err := core.ReadConfig()
	cobra.CheckErr(err)
	conf.Repository = args[0]
	err = conf.Save()
	cobra.CheckErr(err)
	// now actually add the remote
	name := cmd.Flag("name").Value.String()
	core.RemoteAdd(args[0], name)
	err = core.Commit()
	cobra.CheckErr(err)

	cmdr.Info.Println(conf.Repository)
}
