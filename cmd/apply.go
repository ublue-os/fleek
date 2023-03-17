/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/core"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewApplyCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		fleek.Trans("apply.use"),
		fleek.Trans("apply.long"),
		fleek.Trans("apply.short"),
		apply,
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"dry-run",
			"d",
			fleek.Trans("apply.dryRun"),
			false,
		))
	return cmd
}

func apply(cmd *cobra.Command, args []string) {

	var verbose bool
	if cmd.Flag("verbose").Changed {
		verbose = true
	}
	if verbose {
		cmdr.Info.Println(fleek.Trans("apply.writingConfig"))
	}
	err := core.WriteFlake()
	cobra.CheckErr(err)
	if verbose {
		cmdr.Info.Println(fleek.Trans("apply.writingFlake"))
	}
	err = core.CheckFlake()
	cobra.CheckErr(err)
	var dry bool
	if cmd.Flag("dry-run").Changed {
		dry = true
	}
	if !dry {
		cmdr.Info.Println(fleek.Trans("apply.applyingConfig"))
		err = core.ApplyFlake()
		cobra.CheckErr(err)
	} else {
		cmdr.Info.Println(fleek.Trans("apply.dryApplyingConfig"))
	}
	cmdr.Success.Println(fleek.Trans("apply.done"))

}
