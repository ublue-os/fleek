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
		fleek.Trans("remove.use"),
		fleek.Trans("remove.long"),
		fleek.Trans("remove.short"),
		remove,
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"program",
			"p",
			fleek.Trans("remove.program"),
			false,
		)).WithBoolFlag(
		cmdr.NewBoolFlag(
			"apply",
			"a",
			fleek.Trans("remove.apply"),
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
		cmdr.Info.Println(fleek.Trans("remove.applying"))
	}

	var err error

	for _, p := range args {
		if cmd.Flag("program").Changed {
			err = config.RemoveProgram(p)
			cobra.CheckErr(err)
		} else {
			err = config.RemovePackage(p)
			cobra.CheckErr(err)
		}

	}
	if apply {
		err = flake.Apply()
		cobra.CheckErr(err)
	}

	cmdr.Info.Println(fleek.Trans("remove.done"))
}
