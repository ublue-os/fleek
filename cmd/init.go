/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/core"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewInitCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		fleek.Trans("init.use"),
		fleek.Trans("init.long"),
		fleek.Trans("init.short"),
		initialize,
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"force",
			"f",
			fleek.Trans("init.force"),
			false,
		)).WithStringFlag(
		cmdr.NewStringFlag(
			"clone",
			"c",
			fleek.Trans("init.clone"),
			"",
		))
	return cmd
}

// initCmd represents the init command
func initialize(cmd *cobra.Command, args []string) {
	var verbose bool
	if cmd.Flag("verbose").Changed {
		verbose = true
	}
	var repo string
	if cmd.Flag("clone").Changed {
		repo = cmd.Flag("clone").Value.String()

		// clone it
		err := core.Clone(repo)
		cobra.CheckErr(err)
		// return
		return

	}
	cmdr.Info.Println(fleek.Trans("init.start"))
	var force bool
	if cmd.Flag("force").Changed {
		force = true
	}
	if verbose {
		cmdr.Info.Println(fleek.Trans("init.checkNix"))
	}

	ok := core.CheckNix()
	if ok {
		email, err := cmdr.Prompt.Show("Git Config - enter your email address")
		cobra.CheckErr(err)

		name, err := cmdr.Prompt.Show("Git Config - enter your full name")
		cobra.CheckErr(err)
		if verbose {
			cmdr.Info.Println(fleek.Trans("init.writingConfigs"))
		}
		err = core.WriteSampleConfig(email, name, force)
		cobra.CheckErr(err)
		err = core.MakeFlakeDir()
		cobra.CheckErr(err)
		err = core.InitFlake(force)
		cobra.CheckErr(err)
	} else {
		cmdr.Error.Println(fleek.Trans("init.nixNotFound"))
	}
	cmdr.Info.Println(fleek.Trans("init.complete"))
}
