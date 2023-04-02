/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/core"
	"github.com/ublue-os/fleek/internal/nix"
	"github.com/ublue-os/fleek/internal/ux"
	"github.com/vanilla-os/orchid/cmdr"
)

type initCmdFlags struct {
	apply    bool
	force    bool
	clone    string
	location string
}

func InitCommand() *cobra.Command {
	flags := initCmdFlags{}
	command := &cobra.Command{
		Use:   app.Trans("init.use"),
		Short: app.Trans("init.short"),
		Long:  app.Trans("init.long"),

		RunE: func(cmd *cobra.Command, args []string) error {
			return initialize(cmd)
		},
	}
	command.Flags().BoolVarP(
		&flags.apply, app.Trans("init.applyFlag"), "a", false, app.Trans("init.applyFlagDescription"))
	command.Flags().BoolVarP(
		&flags.force, app.Trans("init.forceFlag"), "f", false, app.Trans("init.forceFlagDescription"))
	command.Flags().StringVarP(
		&flags.clone, app.Trans("init.cloneFlag"), "c", "", app.Trans("init.cloneFlagDescription"))
	command.Flags().StringVarP(
		&flags.location, app.Trans("init.locationFlag"), "l", ".config/home-manager", app.Trans("init.locationFlagDescription"))

	return command
}

// initCmd represents the init command
func initialize(cmd *cobra.Command) error {
	var verbose bool
	if cmd.Flag(app.Trans("fleek.verboseFlag")).Changed {
		verbose = true
	}
	ux.Description.Println(cmd.Short)
	var upstream string
	loc := cmd.Flag(app.Trans("init.locationFlag")).Value.String()

	f.config = &core.Config{
		FlakeDir: loc,
	}
	f.flakeLocation = f.config.UserFlakeDir()
	if verbose {
		ux.Info.Println(app.Trans("init.flakeLocation"), f.flakeLocation)
	}
	if cmd.Flag(app.Trans("init.cloneFlag")).Changed {
		upstream = cmd.Flag(app.Trans("init.cloneFlag")).Value.String()

		// clone it
		spinner, err := ux.Spinner().Start(app.Trans("init.cloning"))
		if err != nil {
			return err
		}
		err = f.config.Clone(upstream)
		if err != nil {
			spinner.Fail()
			return err
		}
		r, err := f.Repo()
		if err != nil {
			return err
		}
		out, err := r.SetRebase()
		if err != nil {
			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
		spinner.Success()

		// prompt for git configuration
		email, err := cmdr.Prompt.Show(app.Trans("init.gitEmail"))
		if err != nil {
			return err
		}

		name, err := cmdr.Prompt.Show(app.Trans("init.gitName"))
		if err != nil {
			return err
		}
		out, err = r.LocalConfig(name, email)
		if err != nil {
			spinner.Fail()
			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
		spinner, err = ux.Spinner().Start(app.Trans("init.cloning"))
		if err != nil {
			spinner.Fail()
			return err
		}
		if cmd.Flag(app.Trans("init.applyFlag")).Changed {
			if err != nil {
				return err
			}
			// load the new config
			f.config, err = core.ReadConfig()
			if err != nil {
				return err
			}
			_, err = f.Flake()
			if err != nil {
				return err
			}

			_, err = f.Repo()
			if err != nil {
				return err
			}
			// only re-apply the templates if not `ejected`
			if !f.config.Ejected {
				if verbose {
					ux.Info.Println(app.Trans("apply.checkingSystem"))
				}
				var includeSystems bool
				// check to see if the current machine (system) is in the existing
				// configs. If not, create a new one and add it.
				_, err := core.CurrentSystem()
				if err != nil {
					if strings.Contains(err.Error(), "not") {
						ux.Info.Println(app.Trans("apply.newSystem"))

						// make a new system

						// create new system struct
						sys, err := core.NewSystem(email, name)
						if err != nil {
							return err
						}
						ux.Info.Printfln(app.Trans("init.newSystem", sys.Username, sys.Hostname))
						// get current config
						includeSystems = true
						// append new(current) system
						f.config.Systems = append(f.config.Systems, *sys)
						// save it
						err = f.config.Save()
						if err != nil {
							return err
						}
						repo, err := f.Repo()
						if err != nil {
							return err
						}
						out, err := repo.Commit()
						if verbose {
							ux.Info.Println(string(out))
						}
						if err != nil {
							return err
						}
					}
				}
				if verbose {
					ux.Info.Println(app.Trans("apply.writingFlake"))
				}
				err = f.flake.Write(includeSystems)
				if err != nil {
					return err
				}
				repo, err := f.Repo()
				if err != nil {
					return err
				}
				out, err := repo.Commit()
				if verbose {
					ux.Info.Println(string(out))
				}
				if err != nil {
					return err
				}

			}
			spinner.Success()

			ux.Info.Println(app.Trans("apply.applyingConfig"))
			out, err := f.flake.Apply()
			if err != nil {
				ux.Error.Println(string(out))

				if errors.Is(err, nix.ErrPackageConflict) {
					ux.Fatal.Println(app.Trans("global.errConflict"))
				}
				return err
			}
			if verbose {
				ux.Info.Println(string(out))
			}
			ux.Info.Println(app.Trans("apply.done"))
			return nil
		}
		ux.Info.Println(app.Trans("init.cloned"))

		return nil
	}
	ux.Info.Println(app.Trans("init.start"))
	var force bool
	if cmd.Flag("force").Changed {
		force = true
	}
	if verbose {
		ux.Info.Println(app.Trans("init.checkNix"))
	}

	ok := nix.CheckNix()
	if ok {
		email, err := cmdr.Prompt.Show(app.Trans("init.gitEmail"))
		if err != nil {
			return err
		}

		name, err := cmdr.Prompt.Show(app.Trans("init.gitName"))
		if err != nil {
			return err
		}
		if verbose {
			ux.Info.Println(app.Trans("init.writingConfigs"))
		}
		err = f.config.MakeFlakeDir()
		if err != nil {
			return err
		}
		ux.Info.Println("writing flake")
		err = core.WriteSampleConfig(loc, email, name, force)
		if err != nil {
			return err
		}
		ux.Info.Println("reading config")

		f.config, err = core.ReadConfig()
		if err != nil {
			return err
		}
		flake, err := f.Flake()
		if err != nil {
			return err
		}
		ux.Info.Println("init flake")

		err = flake.Init(force)
		if err != nil {
			return err
		}
		repo, err := f.Repo()
		if err != nil {
			return err
		}
		ux.Info.Println("create repo")

		out, err := repo.CreateRepo()
		if err != nil {
			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
		out, err = repo.LocalConfig(name, email)
		if err != nil {
			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
		out, err = repo.Commit()
		if err != nil {
			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
	} else {
		ux.Error.Println(app.Trans("init.nixNotFound"))
	}
	ux.Info.Println(app.Trans("init.complete"))
	return nil
}
