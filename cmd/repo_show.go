/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/core"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewRepoShowCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		fleek.Trans("reposhow.use"),
		fleek.Trans("reposhow.long"),
		fleek.Trans("reposhow.short"),
		show,
	)
	return cmd
}

// initCmd represents the init command
func show(cmd *cobra.Command, args []string) {

	conf, err := core.ReadConfig()
	cobra.CheckErr(err)
	urls, err := core.Remote()
	cobra.CheckErr(err)
	cmdr.Info.Println("config file says:", conf.Repository)
	cmdr.Info.Println("git says:", urls)

}
