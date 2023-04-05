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
	apply          bool
	force          bool
	clone          string
	promptPassword bool
	location       string
	branch         string
	privateKey     string
	level          string
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
		&flags.clone, app.Trans("init.cloneFlag"), "c", "", app.Trans("init.cloneFlagDescription"))
	command.Flags().StringVarP(
		&flags.branch, app.Trans("init.branchFlag"), "b", "main", app.Trans("init.branchFlagDescription"))
	command.Flags().StringVarP(
		&flags.location, app.Trans("init.locationFlag"), "l", ".config/home-manager", app.Trans("init.locationFlagDescription"))

	command.Flags().StringVar(
		&flags.level, app.Trans("init.levelFlag"), "default", app.Trans("init.levelFlagDescription"))
	command.Flags().StringVarP(
		&flags.privateKey, app.Trans("init.privateKeyFlag"), "k", ".ssh/id_rsa", app.Trans("init.privateKeyFlagDescription"))
	command.Flags().BoolVarP(
		&flags.promptPassword, app.Trans("init.promptPasswordFlag"), "p", false, app.Trans("init.promptPasswordFlagDescription"))
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
	var upstream string
	var branch string
	prompt := cmd.Flag(app.Trans("init.promptPasswordFlag")).Changed

	loc := cmd.Flag(app.Trans("init.locationFlag")).Value.String()
	keyfile := cmd.Flag(app.Trans("init.privateKeyFlag")).Value.String()

	cfg.FlakeDir = loc

	if cmd.Flag(app.Trans("init.cloneFlag")).Changed {
		upstream = cmd.Flag(app.Trans("init.cloneFlag")).Value.String()
		branch = cmd.Flag(app.Trans("init.branchFlag")).Value.String()
		err := fl.Clone(upstream, branch, keyfile, prompt)
		if err != nil {
			return usererr.WithUserMessage(err, app.Trans("flake.cloning", upstream))
		}
	} else {
		fl.Config.Bling = cmd.Flag(app.Trans("init.levelFlag")).Value.String()
		err := fl.Create(force)
		if err != nil {
			return usererr.WithUserMessage(err, app.Trans("flake.creating", upstream))
		}
	}
	if cmd.Flag(app.Trans("init.applyFlag")).Changed {
		err := fl.Apply("")
		if err != nil {
			return usererr.WithUserMessage(err, app.Trans("init.applyFlag"))
		}
		ux.Info.Println(app.Trans("global.complete"))

		return nil
	}
	ux.Info.Println(app.Trans("init.complete"))

	return nil
}
