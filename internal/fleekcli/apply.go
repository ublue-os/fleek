/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/flake"
	"github.com/ublue-os/fleek/internal/ux"
)

type applyCmdFlags struct {
	dryRun bool
}

func ApplyCommand() *cobra.Command {
	flags := applyCmdFlags{}
	command := &cobra.Command{
		Use:     app.Trans("apply.use"),
		Short:   app.Trans("apply.short"),
		Long:    app.Trans("apply.long"),
		Example: app.Trans("apply.example"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return apply(cmd)
		},
	}
	command.Flags().BoolVarP(
		&flags.dryRun, app.Trans("apply.dryRunFlag"), "d", false, app.Trans("apply.dryRunFlagDescription"))

	return command
}

func apply(cmd *cobra.Command) error {
	ux.Description.Println(cmd.Short)
	err := mustConfig()
	if err != nil {
		return err
	}

	var dry bool
	if cmd.Flag(app.Trans("apply.dryRunFlag")).Changed {
		dry = true
	}
	fl, err := flake.Load(cfg, app)
	if err != nil {
		return err
	}
	if err := fl.Write(true); err != nil {
		return err
	}
	if !dry {
		if err := fl.Apply(); err != nil {
			return err
		}
	} else {
		ux.Info.Println(app.Trans("apply.dryApplyingConfig"))
		if bb, err := fl.Check(); err != nil {
			if err != nil {
				ux.Warning.Println(string(bb))
				return err
			}
		}
	}
	ux.Success.Println(app.Trans("global.completed"))
	return nil
}
