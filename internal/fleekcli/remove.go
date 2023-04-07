/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/flake"
	"github.com/ublue-os/fleek/internal/ux"
)

type removeCmdFlags struct {
	apply bool
}

func RemoveCommand() *cobra.Command {
	flags := removeCmdFlags{}
	command := &cobra.Command{
		Use:     app.Trans("remove.use"),
		Short:   app.Trans("remove.short"),
		Long:    app.Trans("remove.long"),
		Example: app.Trans("remove.example"),
		Args:    cobra.MinimumNArgs(1),

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
	err := mustConfig()
	if err != nil {
		return err
	}
	var apply bool
	if cmd.Flag(app.Trans("remove.applyFlag")).Changed {
		apply = true
	}

	fl, err := flake.Load(cfg, app)
	if err != nil {
		return err
	}
	err = fl.MayPull()
	if err != nil {
		return err
	}

	var sb strings.Builder

	sb.WriteString("remove packages: ")

	for _, p := range args {

		if verbose {
			ux.Verbose.Printfln(app.Trans("remove.config"), p)
		}
		err = fl.Config.RemovePackage(p)
		if err != nil {
			ux.Debug.Printfln("remove package error: %s", err)
			return err
		}
		sb.WriteString(p + " ")

	}
	err = fl.Write(false)
	if err != nil {
		ux.Debug.Printfln("flake write error: %s", err)
		return err
	}

	if apply {
		if verbose {
			ux.Info.Println(app.Trans("remove.applying"))
		}
		err = fl.Apply()
		if err != nil {
			if errors.Is(err, flake.ErrPackageConflict) {
				ux.Fatal.Println(app.Trans("global.errConflict"))
			}
			return err
		}
	} else {
		ux.Warning.Println(app.Trans("remove.notApplied"))
	}

	ux.Success.Println(app.Trans("remove.done"))
	return nil
}
