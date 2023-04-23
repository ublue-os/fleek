/*
Copyright © 2023 Brian Ketelsen <bketelsen@gmail.com>
*/
package fleekcli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/flake"
	"github.com/ublue-os/fleek/internal/fleek"
	"github.com/ublue-os/fleek/internal/fleekcli/usererr"
	"github.com/ublue-os/fleek/internal/ux"
	"github.com/ublue-os/fleek/internal/xdg"
)

type initCmdFlags struct {
	apply       bool
	force       bool
	location    string
	level       string
	interactive bool
}

func InitCommand() *cobra.Command {
	flags := initCmdFlags{}
	command := &cobra.Command{
		Use:     app.Trans("init.use"),
		Short:   app.Trans("init.short"),
		Long:    app.Trans("init.long"),
		Example: app.Trans("init.example"),
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return initialize(cmd, args)
		},
	}
	command.Flags().BoolVarP(
		&flags.apply, app.Trans("init.applyFlag"), "a", false, app.Trans("init.applyFlagDescription"))
	command.Flags().BoolVarP(
		&flags.interactive, app.Trans("init.interactiveFlag"), "i", false, app.Trans("init.interactiveFlagDescription"))
	command.Flags().BoolVarP(
		&flags.force, app.Trans("init.forceFlag"), "f", false, app.Trans("init.forceFlagDescription"))
	command.Flags().StringVarP(
		&flags.location, app.Trans("init.locationFlag"), "l", xdg.DataSubpathRel("fleek"), app.Trans("init.locationFlagDescription"))
	command.Flags().StringVar(
		&flags.level, app.Trans("init.levelFlag"), "default", app.Trans("init.levelFlagDescription"))
	return command
}

// initCmd represents the init command
func initialize(cmd *cobra.Command, args []string) error {

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
	var interactive bool
	if cmd.Flag(app.Trans("init.interactiveFlag")).Changed {
		interactive = true
	}
	if interactive {
		config, err := initInteractive()
		if err != nil {
			return err
		}
		err = config.Validate()
		if err != nil {
			return err
		}
		err = config.Save()
		if err != nil {
			return err
		}
		fl, err := flake.Load(&config, app)
		if err != nil {
			return usererr.WithUserMessage(err, app.Trans("flake.initializingTemplates"))
		}
		err = fl.Create(true, true, true)
		if err != nil {
			return usererr.WithUserMessage(err, app.Trans("flake.creating"))
		}

		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		cfile := filepath.Join(fl.Config.FlakeDir, ".fleek.yml")
		csym := filepath.Join(home, ".fleek.yml")
		err = os.Symlink(cfile, csym)
		if err != nil {
			return err
		}
		fin.Warning.Println(app.Trans("init.finalize", fl.Config.FlakeDir))

		return nil

	}

	loc := cmd.Flag(app.Trans("init.locationFlag")).Value.String()
	fl, err := flake.Load(cfg, app)
	cfg.FlakeDir = loc
	if err != nil {
		return usererr.WithUserMessage(err, app.Trans("flake.initializingTemplates"))
	}
	if len(args) > 0 {
		err = fl.Clone(args[0], &outBuffer)
		if err != nil {
			return err
		}
	}

	join, err := fl.IsJoin()
	if err != nil {
		return err
	}
	if join {
		fin.Info.Println(app.Trans("init.joining"))
		err := fl.Join(&outBuffer)
		if err != nil {
			return err
		}
		err = fl.Write("join new system", &outBuffer)
		if err != nil {
			fin.Debug.Printfln("flake write error: %s", err)
			return err
		}

	} else {
		fl.Config.Bling = cmd.Flag(app.Trans("init.levelFlag")).Value.String()
		err = fl.Create(force, false, true)
		if err != nil {
			return usererr.WithUserMessage(err, app.Trans("flake.creating"))
		}
	}
	if cmd.Flag(app.Trans("init.applyFlag")).Changed {
		err := fl.Apply(&outBuffer)
		if err != nil {
			fmt.Println(outBuffer.String())
			return usererr.WithUserMessage(err, app.Trans("init.applyFlag"))
		}
		fin.Info.Println(app.Trans("global.completed"))

		return nil

	}

	fin.Info.Println(app.Trans("init.complete"))

	return nil
}

func initInteractive() (fleek.Config, error) {
	config := fleek.Config{}
	// Prompt for bling level
	choices := fleek.Levels()
	prompt := "Choose your Bling Level"
	level, err := ux.PromptSingle(prompt, choices)
	if err != nil {
		return config, err
	}
	config.Bling = level

	// Prompt for location
	prompt = "Configuration Directory Location: $HOME/"
	choices = []string{
		xdg.ConfigSubpathRel("fleek"),
		xdg.ConfigSubpathRel("home-manager"),
		xdg.DataSubpathRel("fleek"),
		"fleek",
	}
	loc, err := ux.PromptSingle(prompt, choices)
	if err != nil {
		return config, err
	}
	config.FlakeDir = loc
	config.Aliases = make(map[string]string)
	config.Aliases["fleeks"] = "cd ~/" + config.FlakeDir

	// Prompt for shell
	shell, err := fleek.UserShell()
	if err != nil {
		return config, err
	}
	use, err := ux.Confirm("Use detected shell: " + shell)
	if err != nil {
		return config, err
	}
	if use {
		config.Shell = shell
	} else {
		prompt = "Shell"
		choices = []string{
			"bash",
			"zsh",
			"fish",
		}
		shell, err := ux.PromptSingle(prompt, choices)
		if err != nil {
			return config, err
		}
		config.Shell = shell
	}

	user, err := fleek.NewUser()
	if err != nil {
		return config, err
	}
	config.Users = append(config.Users, user)

	// system
	sys, err := fleek.NewSystem()
	if err != nil {
		return config, err
	}
	config.Systems = append(config.Systems, sys)

	return config, nil
}
