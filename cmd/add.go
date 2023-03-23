/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewAddCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		app.Trans("add.use"),
		app.Trans("add.long"),
		app.Trans("add.short"),
		add,
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"program",
			"p",
			app.Trans("add.program"),
			false,
		)).WithBoolFlag(
		cmdr.NewBoolFlag(
			"apply",
			"a",
			app.Trans("add.apply"),
			false,
		))
	cmd.Args = cobra.MinimumNArgs(1)
	return cmd
}

// initCmd represents the init command
func add(cmd *cobra.Command, args []string) {
	var verbose bool
	if cmd.Flag("verbose").Changed {
		verbose = true
	}

	var apply bool
	if cmd.Flag("apply").Changed {
		apply = true
	}
	if verbose {
		cmdr.Info.Println(app.Trans("add.applying"))
	}

	var err error

	for _, p := range args {
		if cmd.Flag("program").Changed {
			err = f.config.AddProgram(p)
			cobra.CheckErr(err)
		} else {
			err = f.config.AddPackage(p)
			cobra.CheckErr(err)
		}

	}
	if apply {
		flake, err := f.Flake()
		cobra.CheckErr(err)
		err = flake.Apply()
		cobra.CheckErr(err)
	}

	cmdr.Info.Println(app.Trans("add.done"))
}
