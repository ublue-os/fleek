/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/core"
	"github.com/ublue-os/fleek/nix"
	"github.com/vanilla-os/orchid/cmdr"
)

var (
	locationFlag string = "location"
)

func NewInitCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		app.Trans("init.use"),
		app.Trans("init.long"),
		app.Trans("init.short"),
		initialize,
	).WithBoolFlag(
		cmdr.NewBoolFlag(
			"force",
			"f",
			app.Trans("init.force"),
			false,
		)).WithStringFlag(
		cmdr.NewStringFlag(
			"clone",
			"c",
			app.Trans("init.clone"),
			"",
		)).WithBoolFlag(
		cmdr.NewBoolFlag(
			"apply",
			"a",
			app.Trans("init.apply"),
			false,
		)).
		WithStringFlag(
			cmdr.NewStringFlag(
				locationFlag,
				"l",
				app.Trans("init.locationFlag"),
				".config/home-manager"))
	return cmd
}

// initCmd represents the init command
func initialize(cmd *cobra.Command, args []string) {
	var verbose bool
	if cmd.Flag("verbose").Changed {
		verbose = true
	}
	var upstream string
	loc := cmd.Flag(locationFlag).Value.String()
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	// hack!
	floc := filepath.Join(home, loc)
	f.config = &core.Config{
		FlakeDir: floc,
	}
	f.flakeLocation = f.config.FlakeDir

	if cmd.Flag("clone").Changed {
		upstream = cmd.Flag("clone").Value.String()

		// clone it
		err := f.config.Clone(upstream)
		cobra.CheckErr(err)
		if cmd.Flag("apply").Changed {
			// load the new config
			f.config, err = core.ReadConfig()
			cobra.CheckErr(err)
			_, err := f.Flake()
			cobra.CheckErr(err)
			_, err = f.Repo()
			cobra.CheckErr(err)

			// only re-apply the templates if not `ejected`
			if !f.config.Ejected {
				if verbose {
					cmdr.Info.Println(app.Trans("apply.checkingSystem"))
				}
				var includeSystems bool
				// check to see if the current machine (system) is in the existing
				// configs. If not, create a new one and add it.
				_, err := core.CurrentSystem()
				if err != nil {
					if strings.Contains(err.Error(), "not") {
						cmdr.Info.Println(app.Trans("apply.newSystem"))

						// make a new system

						// prompt for git configuration
						email, err := cmdr.Prompt.Show("Git Config - enter your email address")
						cobra.CheckErr(err)

						name, err := cmdr.Prompt.Show("Git Config - enter your full name")
						cobra.CheckErr(err)

						// create new system struct
						sys, err := core.NewSystem(email, name)
						cobra.CheckErr(err)
						cmdr.Info.Printfln("New System: %s@%s", sys.Username, sys.Hostname)
						// get current config
						includeSystems = true
						// append new(current) system
						f.config.Systems = append(f.config.Systems, *sys)
						// save it
						err = f.config.Save()
						cobra.CheckErr(err)
					}
				}

				if verbose {
					cmdr.Info.Println(app.Trans("apply.writingFlake"))
				}
				err = f.flake.Write(includeSystems)
				cobra.CheckErr(err)

			}
			cmdr.Info.Println(app.Trans("apply.applyingConfig"))
			err = f.flake.Apply()
			cobra.CheckErr(err)
			cmdr.Success.Println(app.Trans("apply.done"))
			return
		}
		cmdr.Info.Println(app.Trans("init.cloned"))

		return
	}
	cmdr.Info.Println(app.Trans("init.start"))
	var force bool
	if cmd.Flag("force").Changed {
		force = true
	}
	if verbose {
		cmdr.Info.Println(app.Trans("init.checkNix"))
	}

	ok := nix.CheckNix()
	if ok {
		email, err := cmdr.Prompt.Show("Git Config - enter your email address")
		cobra.CheckErr(err)

		name, err := cmdr.Prompt.Show("Git Config - enter your full name")
		cobra.CheckErr(err)
		if verbose {
			cmdr.Info.Println(app.Trans("init.writingConfigs"))
		}
		err = f.config.MakeFlakeDir()
		cobra.CheckErr(err)

		err = core.WriteSampleConfig(floc, email, name, force)
		cobra.CheckErr(err)
		f.config, err = core.ReadConfig()
		cobra.CheckErr(err)
		flake, err := f.Flake()
		cobra.CheckErr(err)
		err = flake.Init(force)
		cobra.CheckErr(err)
		repo, err := f.Repo()
		cobra.CheckErr(err)
		err = repo.CreateRepo()
		cobra.CheckErr(err)
		err = repo.LocalConfig(name, email)
		cobra.CheckErr(err)
		err = repo.Commit()
		cobra.CheckErr(err)
	} else {
		cmdr.Error.Println(app.Trans("init.nixNotFound"))
	}
	cmdr.Info.Println(app.Trans("init.complete"))
}
