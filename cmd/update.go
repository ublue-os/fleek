/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewUpdateCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		app.Trans("update.use"),
		app.Trans("update.long"),
		app.Trans("update.short"),
		update,
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"apply",
			"a",
			app.Trans("update.apply"),
			false,
		))
	return cmd
}

// initCmd represents the init command
func update(cmd *cobra.Command, args []string) {
	cmdr.Info.Println(app.Trans("update.start"))

	flake, err := f.Flake()
	cobra.CheckErr(err)
	flake.Update()
	if cmd.Flag("apply").Changed {
		cmdr.Info.Println(app.Trans("update.apply"))

		err = flake.Apply()
		cobra.CheckErr(err)
	} else {
		cmdr.Info.Println(app.Trans("update.needApply"))

	}
	cmdr.Success.Println(app.Trans("update.done"))
}
