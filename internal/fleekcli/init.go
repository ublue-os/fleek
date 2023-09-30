/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/flake"
	"github.com/ublue-os/fleek/internal/fleekcli/usererr"
)

type initCmdFlags struct {
	apply bool
	force bool
	level string
}

func InitCommand() *cobra.Command {
	flags := initCmdFlags{}
	command := &cobra.Command{
		Use:     app.Trans("init.use"),
		Short:   app.Trans("init.short"),
		Long:    app.Trans("init.long"),
		Example: app.Trans("init.example"),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return initialize(cmd, args)
		},
	}
	command.Flags().BoolVarP(
		&flags.apply, app.Trans("init.applyFlag"), "a", false, app.Trans("init.applyFlagDescription"))
	command.Flags().BoolVarP(
		&flags.force, app.Trans("init.forceFlag"), "f", false, app.Trans("init.forceFlagDescription"))
	command.Flags().StringVar(
		&flags.level, app.Trans("init.levelFlag"), "default", app.Trans("init.levelFlagDescription"))
	return command
}

// initCmd represents the init command
func initialize(cmd *cobra.Command, _ []string) error {

	var verbose bool
	if cmd.Flag(app.Trans("fleek.verboseFlag")).Changed {
		verbose = true
	}
	var force bool
	if cmd.Flag(app.Trans("init.forceFlag")).Changed {
		force = true
	}
	cfg.Verbose = verbose

	fin.Description.Println(cmd.Short)

	loc := cmd.Flag(app.Trans("init.locationFlag")).Value.String()
	fl, err := flake.Load(cfg, app)
	cfg.FlakeDir = loc
	if err != nil {
		return usererr.WithUserMessage(err, app.Trans("flake.initializingTemplates"))
	}

	fl.Config.Bling = cmd.Flag(app.Trans("init.levelFlag")).Value.String()
	err = fl.Create(force, true)
	if err != nil {
		return usererr.WithUserMessage(err, app.Trans("flake.creating"))
	}

	if cmd.Flag(app.Trans("init.applyFlag")).Changed {
		err := fl.Apply()
		if err != nil {
			return usererr.WithUserMessage(err, app.Trans("init.applyFlag"))
		}
		fin.Logger.Info(app.Trans("global.completed"))

		return nil
	}
	fin.Logger.Info(app.Trans("init.complete"))

	return nil
}
