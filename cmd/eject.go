/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/core"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewEjectCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		fleek.Trans("eject.use"),
		fleek.Trans("eject.long"),
		fleek.Trans("eject.short"),
		eject,
	)
	return cmd
}

// initCmd represents the init command
func eject(cmd *cobra.Command, args []string) {

	ok, err := cmdr.Confirm.Show(fleek.Trans("eject.confirm"))
	cobra.CheckErr(err)

	if ok {
		cmdr.Info.Println(fleek.Trans("eject.start"))
		err := flake.Write()
		cobra.CheckErr(err)
		err = core.WriteEjectConfig()
		cobra.CheckErr(err)
		cmdr.Info.Println(fleek.Trans("eject.complete"))
	}

}
