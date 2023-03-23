/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewRemoveCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		app.Trans("remove.use"),
		app.Trans("remove.long"),
		app.Trans("remove.short"),
		remove,
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"program",
			"p",
			app.Trans("remove.program"),
			false,
		)).WithBoolFlag(
		cmdr.NewBoolFlag(
			"apply",
			"a",
			app.Trans("remove.apply"),
			false,
		))
	cmd.Args = cobra.MinimumNArgs(1)
	return cmd
}

// initCmd represents the init command
func remove(cmd *cobra.Command, args []string) {
	var verbose bool
	if cmd.Flag("verbose").Changed {
		verbose = true
	}

	var apply bool
	if cmd.Flag("apply").Changed {
		apply = true
	}
	if verbose {
		cmdr.Info.Println(app.Trans("remove.applying"))
	}

	var err error

	for _, p := range args {
		if cmd.Flag("program").Changed {
			err = f.config.RemoveProgram(p)
			cobra.CheckErr(err)
		} else {
			err = f.config.RemovePackage(p)
			cobra.CheckErr(err)
		}

	}
	if apply {
		err = f.flake.Apply()
		cobra.CheckErr(err)
	}

	cmdr.Info.Println(app.Trans("remove.done"))
}
