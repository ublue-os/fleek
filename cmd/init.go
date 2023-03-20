/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

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
		)).WithBoolFlag(
		cmdr.NewBoolFlag(
			"apply",
			"a",
			fleek.Trans("init.apply"),
			false,
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
		if cmd.Flag("apply").Changed {
			// only re-apply the templates if not `ejected`
			if ejected, _ := core.Ejected(); !ejected {
				if verbose {

					cmdr.Info.Println(fleek.Trans("apply.checkingSystem"))
				}
				// check to see if the current machine (system) is in the existing
				// configs. If not, create a new one and add it.
				_, err := core.CurrentSystem()
				if err != nil {
					if strings.Contains(err.Error(), "not") {
						cmdr.Info.Println(fleek.Trans("apply.newSystem"))

						//make a new system

						// prompt for git configuration
						email, err := cmdr.Prompt.Show("Git Config - enter your email address")
						cobra.CheckErr(err)

						name, err := cmdr.Prompt.Show("Git Config - enter your full name")
						cobra.CheckErr(err)

						// create new system struct
						sys, err := core.NewSystem(email, name)
						cobra.CheckErr(err)
						cmdr.Info.Println("New System: %s@%s", sys.Username, sys.Hostname)
						// get current config
						conf, err := core.ReadConfig()
						cobra.CheckErr(err)

						// append new(current) system
						conf.Systems = append(conf.Systems, *sys)
						// save it
						err = conf.Save()
						cobra.CheckErr(err)

					}
				}

				if verbose {
					cmdr.Info.Println(fleek.Trans("apply.writingFlake"))
				}
				err = core.WriteFlake()
				cobra.CheckErr(err)

			}
			cmdr.Info.Println(fleek.Trans("apply.applyingConfig"))
			err := core.ApplyFlake()
			cobra.CheckErr(err)
			cmdr.Success.Println(fleek.Trans("apply.done"))

			return
		}

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
		err = core.MakeFlakeDir()
		cobra.CheckErr(err)
		err = core.WriteSampleConfig(email, name, force)
		cobra.CheckErr(err)

		err = core.InitFlake(force)
		cobra.CheckErr(err)
	} else {
		cmdr.Error.Println(fleek.Trans("init.nixNotFound"))
	}
	cmdr.Info.Println(fleek.Trans("init.complete"))
}
