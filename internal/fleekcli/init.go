/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/flake"
	"github.com/ublue-os/fleek/internal/fleekcli/usererr"
	"github.com/ublue-os/fleek/internal/ux"
)

type initCmdFlags struct {
	apply    bool
	force    bool
	location string
	level    string
}

func InitCommand() *cobra.Command {
	flags := initCmdFlags{}
	command := &cobra.Command{
		Use:     app.Trans("init.use"),
		Short:   app.Trans("init.short"),
		Long:    app.Trans("init.long"),
		Example: app.Trans("init.example"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return initialize(cmd)
		},
	}
	command.Flags().BoolVarP(
		&flags.apply, app.Trans("init.applyFlag"), "a", false, app.Trans("init.applyFlagDescription"))
	command.Flags().BoolVarP(
		&flags.force, app.Trans("init.forceFlag"), "f", false, app.Trans("init.forceFlagDescription"))
	command.Flags().StringVarP(
		&flags.location, app.Trans("init.locationFlag"), "l", "Sync/fleek", app.Trans("init.locationFlagDescription"))
	command.Flags().StringVar(
		&flags.level, app.Trans("init.levelFlag"), "default", app.Trans("init.levelFlagDescription"))
	return command
}

// initCmd represents the init command
func initialize(cmd *cobra.Command) error {

	var verbose bool
	if cmd.Flag(app.Trans("fleek.verboseFlag")).Changed {
		verbose = true
	}
	var force bool
	if cmd.Flag(app.Trans("init.forceFlag")).Changed {
		force = true
	}
	cfg.Verbose = verbose

	ux.Description.Println(cmd.Short)

	fl, err := flake.Load(cfg, app)
	if err != nil {
		return usererr.WithUserMessage(err, app.Trans("flake.initializingTemplates"))
	}
	loc := cmd.Flag(app.Trans("init.locationFlag")).Value.String()
	cfg.FlakeDir = loc

	join, err := fl.IsJoin()
	if err != nil {
		return err
	}
	if join {
		ux.Info.Println(app.Trans("init.joining"))
		err := fl.Join()
		if err != nil {
			return err
		}

	} else {
		fl.Config.Bling = cmd.Flag(app.Trans("init.levelFlag")).Value.String()
		err = fl.Create(force)
		if err != nil {
			return usererr.WithUserMessage(err, app.Trans("flake.creating"))
		}
	}

	if cmd.Flag(app.Trans("init.applyFlag")).Changed {
		err := fl.Apply()
		if err != nil {
			return usererr.WithUserMessage(err, app.Trans("init.applyFlag"))
		}
		ux.Info.Println(app.Trans("global.complete"))

		return nil
	}
	ux.Info.Println(app.Trans("init.complete"))

	return nil
}
