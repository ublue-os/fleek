/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
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
		)).WithBoolFlag(
		cmdr.NewBoolFlag(
			"push",
			"p",
			fleek.Trans("apply.push"),
			false,
		))
	return cmd
}

func apply(cmd *cobra.Command, args []string) {

	var verbose bool
	var push bool
	if cmd.Flag("verbose").Changed {
		verbose = true
	}

	if cmd.Flag("push").Changed {
		push = true
	}
	if behind {
		cmdr.Error.Println(fleek.Trans("apply.behind"))
		return

	}
	if verbose {
		cmdr.Info.Println(fleek.Trans("apply.writingConfig"))
	}
	// only re-apply the templates if not `ejected`
	if !config.Ejected {
		if verbose {
			cmdr.Info.Println(fleek.Trans("apply.writingFlake"))
		}
		err := flake.Write()
		cobra.CheckErr(err)
		err = repo.Commit()
		if err != nil {
			cmdr.Error.Println(fleek.Trans("apply.commitError"), err)
		}
		cobra.CheckErr(err)

	}

	var dry bool
	if cmd.Flag("dry-run").Changed {
		dry = true
	}
	if !dry {
		cmdr.Info.Println(fleek.Trans("apply.applyingConfig"))
		err := flake.Apply()
		cobra.CheckErr(err)
		err = repo.Commit()
		if err != nil {
			cmdr.Error.Println(fleek.Trans("apply.commitError"), err)
		}
	} else {
		cmdr.Info.Println(fleek.Trans("apply.dryApplyingConfig"))
		err := flake.Check()
		cobra.CheckErr(err)
	}
	if push {
		cmdr.Info.Println(fleek.Trans("apply.pushing"))
		err := repo.Push()
		cobra.CheckErr(err)
	}

	cmdr.Success.Println(fleek.Trans("apply.done"))

}
