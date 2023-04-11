/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/flake"
	"github.com/ublue-os/fleek/internal/fleekcli/usererr"
	"github.com/ublue-os/fleek/internal/xdg"
)

type generateCmdFlags struct {
	apply    bool
	force    bool
	location string
	level    string
}

func GenerateCommand() *cobra.Command {
	flags := generateCmdFlags{}
	command := &cobra.Command{
		Use:   app.Trans("generate.use"),
		Short: app.Trans("generate.short"),
		Long:  app.Trans("generate.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generate(cmd)
		},
	}
	command.Flags().BoolVarP(
		&flags.apply, app.Trans("generate.applyFlag"), "a", false, app.Trans("generate.applyFlagDescription"))
	command.Flags().BoolVarP(
		&flags.force, app.Trans("generate.forceFlag"), "f", false, app.Trans("generate.forceFlagDescription"))
	command.Flags().StringVarP(
		&flags.location, app.Trans("generate.locationFlag"), "l", xdg.ConfigSubpathRel("fleek"), app.Trans("generate.locationFlagDescription"))
	command.Flags().StringVar(
		&flags.level, app.Trans("generate.levelFlag"), "default", app.Trans("generate.levelFlagDescription"))

	return command
}

// initCmd represents the init command
func generate(cmd *cobra.Command) error {
	var verbose bool
	if cmd.Flag(app.Trans("fleek.verboseFlag")).Changed {
		verbose = true
	}
	var force bool
	if cmd.Flag(app.Trans("generate.forceFlag")).Changed {
		force = true
	}
	cfg.Verbose = verbose

	fin.Description.Println(cmd.Short)

	loc := cmd.Flag(app.Trans("generate.locationFlag")).Value.String()
	fl, err := flake.Load(cfg, app)
	cfg.FlakeDir = loc
	if err != nil {
		return usererr.WithUserMessage(err, app.Trans("flake.initializingTemplates"))
	}

	fl.Config.Bling = cmd.Flag(app.Trans("generate.levelFlag")).Value.String()
	fin.Info.Println("Bling level:", fl.Config.Bling)
	err = fl.Create(force, false)
	if err != nil {
		return usererr.WithUserMessage(err, app.Trans("flake.creating"))
	}

	fin.Info.Printfln(app.Trans("generate.complete"), loc)
	fl.Config.Ejected = true
	fl.Config.Git.AutoCommit = false
	fl.Config.Git.AutoPull = false
	fl.Config.Git.AutoPush = false
	fl.Config.Git.Enabled = false
	for _, system := range fl.Config.Systems {
		// nix run --impure home-manager/master -- -b bak switch --flake .#bjk@ghanima
		fl.Config.Aliases["apply-"+system.Hostname] = fmt.Sprintf("nix run --impure home-manager/master -- -b bak switch --flake .#%s@%s", system.Username, system.Hostname)
		//fin.Info.Printfln("nix run --impure home-manager/master -- -b bak switch --flake .#%s@%s", system.Username, system.Hostname)
	}
	err = fl.Config.Save()
	if err != nil {
		return err
	}
	fin.Info.Println("writing,", fl.Config.Bling)
	err = fl.Write(true, "fleek: generate")
	if err != nil {
		return err
	}

	if cmd.Flag(app.Trans("generate.applyFlag")).Changed {
		err := fl.Apply()
		if err != nil {
			return usererr.WithUserMessage(err, app.Trans("generate.applyFlag"))
		}
		fin.Info.Println(app.Trans("global.completed"))

		return nil
	}
	// TODO app trans
	fin.Info.Println("Run the following commands from the flake directory to apply your changes:")

	for _, system := range fl.Config.Systems {
		// nix run --impure home-manager/master -- -b bak switch --flake .#bjk@ghanima
		fmt.Printf("nix run --impure home-manager/master -- -b bak switch --flake .#%s@%s\n", system.Username, system.Hostname)
		//fin.Info.Printfln("nix run --impure home-manager/master -- -b bak switch --flake .#%s@%s", system.Username, system.Hostname)
	}

	return nil
}
