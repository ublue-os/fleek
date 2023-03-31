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

type updateCmdFlags struct {
	apply bool
}

func UpdateCommand() *cobra.Command {
	flags := updateCmdFlags{}
	command := &cobra.Command{
		Use:   app.Trans("update.use"),
		Short: app.Trans("update.short"),
		Long:  app.Trans("update.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return update(cmd, args)
		},
	}
	command.Flags().BoolVarP(
		&flags.apply, app.Trans("update.applyFlag"), "a", false, app.Trans("update.applyFlagDescription"))

	return command
}

// initCmd represents the init command
func update(cmd *cobra.Command, args []string) error {
	var verbose bool
	if cmd.Flag(app.Trans("fleek.verboseFlag")).Changed {
		verbose = true
	}
	ux.Description.Println(cmd.Short)

	flake, err := f.Flake()
	if err != nil {
		debug.Log("flake open error: %s", err)
		return err
	}
	spinner, err := ux.Spinner().Start(app.Trans("update.start"))
	if err != nil {
		return err
	}
	out, err := flake.Update()
	if err != nil {
		debug.Log("flake update error: %s", err)
		spinner.Fail()
		return err
	}
	if verbose {
		ux.Info.Println(string(out))
	}
	spinner.Success()
	if cmd.Flag("apply").Changed {
		spinner, err := ux.Spinner().Start(app.Trans("global.applying"))
		if err != nil {
			return err
		}
		out, err := flake.Apply()
		if err != nil {
			spinner.Fail()
			ux.Error.Println(string(out))

			if errors.Is(err, nix.ErrPackageConflict) {
				ux.Fatal.Println(app.Trans("global.errConflict"))
			}

			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
		spinner.Success()
	} else {
		ux.Warning.Println(app.Trans("update.needApply"))

	}
	ux.Success.Println(app.Trans("update.done"))
	return nil
}
