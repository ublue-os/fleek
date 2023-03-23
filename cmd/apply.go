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
		app.Trans("apply.use"),
		app.Trans("apply.long"),
		app.Trans("apply.short"),
		apply,
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"dry-run",
			"d",
			app.Trans("apply.dryRun"),
			false,
		)).WithBoolFlag(
		cmdr.NewBoolFlag(
			"push",
			"p",
			app.Trans("apply.push"),
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
	if f.flakeStatus == FlakeBehind {
		cmdr.Error.Println(app.Trans("apply.behind"))
		return

	}
	if verbose {
		cmdr.Info.Println(app.Trans("apply.writingConfig"))
	}
	// only re-apply the templates if not `ejected`
	if !f.config.Ejected {
		if verbose {
			cmdr.Info.Println(app.Trans("apply.writingFlake"))
		}
		err := f.flake.Write()
		cobra.CheckErr(err)
		err = f.repo.Commit()
		if err != nil {
			cmdr.Error.Println(app.Trans("apply.commitError"), err)
		}
		cobra.CheckErr(err)

	}

	var dry bool
	if cmd.Flag("dry-run").Changed {
		dry = true
	}
	if !dry {
		cmdr.Info.Println(app.Trans("apply.applyingConfig"))
		err := f.flake.Apply()
		cobra.CheckErr(err)
		err = f.repo.Commit()
		if err != nil {
			cmdr.Error.Println(app.Trans("apply.commitError"), err)
		}
	} else {
		cmdr.Info.Println(app.Trans("apply.dryApplyingConfig"))
		err := f.flake.Check()
		cobra.CheckErr(err)
	}
	if push {
		cmdr.Info.Println(app.Trans("apply.pushing"))
		err := f.repo.Push()
		cobra.CheckErr(err)
	}

	cmdr.Success.Println(app.Trans("apply.done"))

}
