package fleekcli

import (
	"io"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/flake"
	"github.com/ublue-os/fleek/internal/fleek"
	"github.com/ublue-os/fleek/internal/fleekcli/usererr"
	"github.com/ublue-os/fleek/internal/vercheck"
	"github.com/ublue-os/fleek/internal/xdg"
)

var cfg *fleek.Config
var cfgFound bool

type rootCmdFlags struct {
	quiet    bool
	verbose  bool
	location string
}

func RootCmd() *cobra.Command {
	flags := rootCmdFlags{}
	command := &cobra.Command{
		Use:   app.Trans("fleek.use"),
		Short: app.Trans("fleek.short"),
		Long:  app.Trans("fleek.long"),

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if flags.quiet {
				cmd.SetErr(io.Discard)
			}
			vercheck.CheckVersion(cmd.ErrOrStderr(), cmd.CommandPath())
			fin.Logger.Debug("debug enabled")
			info, ok := debug.ReadBuildInfo()
			if ok {

				fin.Logger.Trace(info.String())

			}

			err := flake.ForceProfile()
			if err != nil {
				fin.Logger.Error("Nix can't list profiles.")
				os.Exit(1)
			}
			// try to get the config, which may not exist yet
			c, err := fleek.ReadConfig(flags.location)
			if err == nil {

				fin.Logger.Debug(app.Trans("fleek.configLoaded"), fin.Logger.Args("location", flags.location))

				cfg = c
				cfgFound = true
			} else {
				cfg = &fleek.Config{}
				cfgFound = false
			}
			if cfg != nil {
				cfg.Quiet = flags.quiet
				cfg.Verbose = flags.verbose
				fin.Logger.Debug("git",
					fin.Logger.Args(
						"autopush", cfg.Git.AutoPush,
						"autocommit", cfg.Git.AutoCommit,
						"autopull", cfg.Git.AutoPull,
					))
				if cfg.Ejected {
					if cmd.Name() != app.Trans("apply.use") {
						fin.Logger.Error(app.Trans("eject.ejected"))
						os.Exit(1)
					}
				}

				migrate := cfg.NeedsMigration()
				if migrate {
					fin.Logger.Warn("Migration required")
					err := cfg.Migrate()
					if err != nil {
						fin.Logger.Error("migrating host files", fin.Logger.Args("error", err))
						os.Exit(1)
					}
					fl, err := flake.Load(cfg, app)
					if err != nil {
						fin.Logger.Error("loading flake", fin.Logger.Args("error", err))
						os.Exit(1)
					}
					err = fl.Write("update host and user files", true, false)
					if err != nil {
						fin.Logger.Error(" writing flake:", fin.Logger.Args("error", err))
						os.Exit(1)
					}
				}

			}

		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if flags.quiet {
				cmd.SetErr(io.Discard)
			}

			if cfg.AutoGC {
				fin.Logger.Info("Running nix-collect-garbage")
				// we don't care too much if there's an error here
				_ = fleek.CollectGarbage()
			}

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	initGroup := &cobra.Group{
		ID:    "init",
		Title: app.Trans("global.initGroup"),
	}
	fleekGroup := &cobra.Group{
		ID:    "fleek",
		Title: app.Trans("global.fleekGroup"),
	}

	packageGroup := &cobra.Group{
		ID:    "package",
		Title: app.Trans("global.packageGroup"),
	}

	command.AddGroup(initGroup, packageGroup, fleekGroup)
	addCmd := AddCommand()
	addCmd.GroupID = packageGroup.ID

	removeCmd := RemoveCommand()
	removeCmd.GroupID = packageGroup.ID

	showCmd := ShowCmd()
	showCmd.GroupID = fleekGroup.ID

	applyCmd := ApplyCommand()
	applyCmd.GroupID = fleekGroup.ID

	updateCmd := UpdateCommand()
	updateCmd.GroupID = packageGroup.ID

	initCmd := InitCommand()
	initCmd.GroupID = initGroup.ID
	joinCmd := JoinCommand()
	joinCmd.GroupID = initGroup.ID
	ejectCmd := EjectCommand()
	ejectCmd.GroupID = fleekGroup.ID
	generateCmd := GenerateCommand()
	generateCmd.GroupID = fleekGroup.ID
	searchCmd := SearchCommand()
	searchCmd.GroupID = packageGroup.ID

	infoCmd := InfoCommand()
	infoCmd.GroupID = packageGroup.ID
	writeCmd := WriteCommand()
	writeCmd.GroupID = fleekGroup.ID
	manCmd := ManCommand()

	docsCmd := genDocsCmd()
	command.AddCommand(docsCmd)
	command.AddCommand(manCmd)
	command.AddCommand(showCmd)

	command.AddCommand(addCmd)
	command.AddCommand(removeCmd)
	command.AddCommand(applyCmd)
	command.AddCommand(updateCmd)

	command.AddCommand(initCmd)
	command.AddCommand(joinCmd)

	command.AddCommand(ejectCmd)
	command.AddCommand(searchCmd)
	command.AddCommand(infoCmd)
	command.AddCommand(generateCmd)
	command.AddCommand(writeCmd)
	command.AddCommand(VersionCmd())

	command.PersistentFlags().BoolVarP(
		&flags.quiet, app.Trans("fleek.quietFlag"), "q", false, app.Trans("fleek.quietFlagDescription"))
	command.PersistentFlags().StringVarP(
		&flags.location, app.Trans("init.locationFlag"), "l", xdg.DataSubpathRel("fleek"), app.Trans("init.locationFlagDescription"))

	verboseMiddleware.AttachToFlag(command.PersistentFlags(), app.Trans("fleek.verboseFlag"))

	debugMiddleware.AttachToFlag(command.PersistentFlags(), app.Trans("fleek.debugFlag"))
	traceMiddleware.AttachToFlag(command.PersistentFlags(), app.Trans("fleek.traceFlag"))

	return command
}

func mustConfig() error {

	if !cfgFound {
		return usererr.New("configuration files not found, run `fleek init`")
	}
	return nil
}
