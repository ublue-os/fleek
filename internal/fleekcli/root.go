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
			fin.Debug.Println("debug enabled")
			info, ok := debug.ReadBuildInfo()
			if ok {

				fin.Debug.Println(info.String())

			}

			err := flake.ForceProfile()
			if err != nil {
				fin.Error.Println("Nix can't list profiles.")
				os.Exit(1)
			}
			// try to get the config, which may not exist yet
			c, err := fleek.ReadConfig(flags.location)
			if err == nil {
				if flags.verbose {
					fin.Info.Println(app.Trans("fleek.configLoaded"))
				}
				cfg = c
				cfgFound = true
			} else {
				cfg = &fleek.Config{}
				cfgFound = false
			}
			if cfg != nil {
				cfg.Quiet = flags.quiet
				cfg.Verbose = flags.verbose
				fin.Debug.Printfln("git autopush: %v", cfg.Git.AutoPush)
				fin.Debug.Printfln("git autocommit: %v", cfg.Git.AutoCommit)
				fin.Debug.Printfln("git autopull: %v", cfg.Git.AutoPull)
				if cfg.Ejected {
					if cmd.Name() != app.Trans("apply.use") {
						fin.Error.Println(app.Trans("eject.ejected"))
						os.Exit(1)
					}
				}

				migrate := cfg.NeedsMigration()
				if migrate {
					fin.Info.Println("Migration required")
					err := cfg.Migrate()
					if err != nil {
						fin.Error.Println("error migrating host files:", err)
						os.Exit(1)
					}
					fl, err := flake.Load(cfg, app)
					if err != nil {
						fin.Error.Println("error loading flake:", err)
						os.Exit(1)
					}
					err = fl.Write("update host and user files", true, false)
					if err != nil {
						fin.Error.Println("error writing flake:", err)
						os.Exit(1)
					}
				}

			}

		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if flags.quiet {
				cmd.SetErr(io.Discard)
			}
			fin.Debug.Printfln("git autopush: %v", cfg.Git.AutoPush)
			fin.Debug.Printfln("git autocommit: %v", cfg.Git.AutoCommit)
			fin.Debug.Printfln("git autopull: %v", cfg.Git.AutoPull)
			fin.Debug.Printfln("auto gc: %v", cfg.AutoGC)

			if cfg.AutoGC {
				fin.Info.Println("Running nix-collect-garbage")
				// we don't care too much if there's an error here
				fleek.CollectGarbage()
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
	command.PersistentFlags().BoolVarP(
		&flags.verbose, app.Trans("fleek.verboseFlag"), "v", false, app.Trans("fleek.verboseFlagDescription"))
	command.PersistentFlags().StringVarP(
		&flags.location, app.Trans("init.locationFlag"), "l", xdg.DataSubpathRel("fleek"), app.Trans("init.locationFlagDescription"))

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
