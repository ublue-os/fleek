/*
Copyright © 2023 Brian Ketelsen <bketelsen@gmail.com>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/flake"
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
	fin.Description.Println(cmd.Short)
	err := mustConfig()
	if err != nil {
		return err
	}
	fl, err := flake.Load(cfg, app)
	if err != nil {
		return err
	}
	if cfg.Ejected {
		if err := fl.Apply(); err != nil {
			return err
		}
		fin.Success.Println(app.Trans("global.completed"))
	}
	var dry bool
	if cmd.Flag(app.Trans("apply.dryRunFlag")).Changed {
		dry = true
	}

	err = fl.MayPull()
	if err != nil {
		return err
	}

	if err := fl.Write("fleek: apply"); err != nil {
		return err
	}
	if !dry {
		if err := fl.Apply(); err != nil {
			return err
		}
	} else {
		fin.Info.Println(app.Trans("apply.dryApplyingConfig"))
		if bb, err := fl.Check(); err != nil {
			if err != nil {
				fin.Warning.Println(string(bb))
				return err
			}
		}
	}
	fin.Success.Println(app.Trans("global.completed"))
	return nil
}
