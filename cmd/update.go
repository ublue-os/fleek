/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/core"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewUpdateCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		fleek.Trans("update.use"),
		fleek.Trans("update.long"),
		fleek.Trans("update.short"),
		update,
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"apply",
			"a",
			fleek.Trans("update.apply"),
			false,
		))
	return cmd
}

// initCmd represents the init command
func update(cmd *cobra.Command, args []string) {
	cmdr.Info.Println(fleek.Trans("update.start"))

	err := core.UpdateFlake()
	cobra.CheckErr(err)
	if cmd.Flag("apply").Changed {
		cmdr.Info.Println(fleek.Trans("update.apply"))

		err = core.ApplyFlake()
		cobra.CheckErr(err)
	} else {
		cmdr.Info.Println(fleek.Trans("update.needApply"))

	}
	cmdr.Success.Println(fleek.Trans("update.done"))
}
