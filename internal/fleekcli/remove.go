/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/debug"
	"github.com/ublue-os/fleek/internal/nix"
	"github.com/ublue-os/fleek/internal/ux"
)

type removeCmdFlags struct {
	apply bool
}

func RemoveCommand() *cobra.Command {
	flags := removeCmdFlags{}
	command := &cobra.Command{
		Use:   app.Trans("remove.use"),
		Short: app.Trans("remove.short"),
		Long:  app.Trans("remove.long"),
		Args:  cobra.MinimumNArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			return remove(cmd, args)
		},
	}
	command.Flags().BoolVarP(
		&flags.apply, app.Trans("remove.applyFlag"), "a", false, app.Trans("remove.applyFlagDescription"))

	return command
}

// initCmd represents the init command
func remove(cmd *cobra.Command, args []string) error {
	var verbose bool
	if cmd.Flag(app.Trans("fleek.verboseFlag")).Changed {
		verbose = true
	}
	ux.Description.Println(cmd.Short)
	var apply bool
	if cmd.Flag(app.Trans("remove.applyFlag")).Changed {
		apply = true
	}

	var err error

	for _, p := range args {
		/*		if cmd.Flag("program").Changed {
					err = f.config.RemoveProgram(p)
					cobra.CheckErr(err)
				} else {
		*/
		ux.Info.Println(app.Trans("remove.config"))
		err = f.config.RemovePackage(p)
		if err != nil {
			debug.Log("remove package error: %s", err)
			ux.Error.Println(err)
			return err
		}
		//	}

	}
	if apply {
		if verbose {
			ux.Info.Println(app.Trans("remove.applying"))
		}
		flake, err := f.Flake()
		if err != nil {
			debug.Log("get flake error: %s", err)
			return err
		}
		err = flake.Write(false)
		if err != nil {
			debug.Log("flake write error: %s", err)
			return err
		}
		repo, err := f.Repo()
		if err != nil {
			debug.Log("get repo error: %s", err)
			return err
		}
		spinner, err := ux.Spinner().Start(app.Trans("remove.applying"))
		if err != nil {
			return err
		}
		out, err := repo.Commit()
		if err != nil {
			debug.Log("commit error: %s", err)
			spinner.Fail()
			return err
		}
		spinner.Success()
		if verbose {
			ux.Info.Println(string(out))
		}
		spinner, err = ux.Spinner().Start(app.Trans("global.commit"))
		if err != nil {
			return err
		}
		out, err = flake.Apply()
		if err != nil {
			spinner.Fail()
			ux.Error.Println(string(out))

			if errors.Is(err, nix.ErrPackageConflict) {
				ux.Fatal.Println(app.Trans("global.errConflict"))
			}

			return err
		}
		spinner.Success()
		if verbose {
			ux.Info.Println(string(out))
		}
	}

	ux.Success.Println(app.Trans("remove.done"))
	return nil
}
