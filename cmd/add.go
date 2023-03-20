/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/core"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewAddCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		fleek.Trans("add.use"),
		fleek.Trans("add.long"),
		fleek.Trans("add.short"),
		add,
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"program",
			"p",
			fleek.Trans("add.program"),
			false,
		)).WithBoolFlag(
		cmdr.NewBoolFlag(
			"apply",
			"a",
			fleek.Trans("add.apply"),
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
		cmdr.Info.Println(fleek.Trans("add.applying"))
	}

	conf, err := core.ReadConfig()
	cobra.CheckErr(err)
	for _, p := range args {
		conf.AddPackage(p)

	}
	if apply {
		err = core.ApplyFlake()
		cobra.CheckErr(err)
	}

	cmdr.Info.Println(fleek.Trans("add.done"))
}
