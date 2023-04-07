/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/flake"
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
			return update(cmd)
		},
	}
	command.Flags().BoolVarP(
		&flags.apply, app.Trans("update.applyFlag"), "a", false, app.Trans("update.applyFlagDescription"))

	return command
}

// initCmd represents the init command
func update(cmd *cobra.Command) error {
	ux.Description.Println(cmd.Short)
	err := mustConfig()
	if err != nil {
		return err
	}
	fl, err := flake.Load(cfg, app)
	if err != nil {
		return err
	}
	err = fl.MayPull()
	if err != nil {
		return err
	}

	if err := fl.Update(); err != nil {
		return err
	}
	if cmd.Flag(app.Trans("update.applyFlag")).Changed {
		if err := fl.Apply(); err != nil {
			return err
		}
	} else {
		ux.Warning.Println(app.Trans("update.needApply"))
	}

	ux.Success.Println(app.Trans("update.done"))
	return nil
}
