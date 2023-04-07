package fleekcli

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/fleek"
	"github.com/ublue-os/fleek/internal/fleekcli/usererr"
	"github.com/ublue-os/fleek/internal/ux"
)

var cfg *fleek.Config
var cfgFound bool

type rootCmdFlags struct {
	quiet   bool
	verbose bool
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
			ux.Debug.Println("debug enabled")
			// try to get the config, which may not exist yet
			c, err := fleek.ReadConfig()
			if err == nil {
				if flags.verbose {
					ux.Info.Println(app.Trans("fleek.configLoaded"))
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
				ux.Debug.Printfln("git autopush: %v", cfg.Git.AutoPush)
				ux.Debug.Printfln("git autocommit: %v", cfg.Git.AutoCommit)
				ux.Debug.Printfln("git autopull: %v", cfg.Git.AutoPull)

			}

		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if flags.quiet {
				cmd.SetErr(io.Discard)
			}
			ux.Debug.Printfln("git autopush: %v", cfg.Git.AutoPush)
			ux.Debug.Printfln("git autocommit: %v", cfg.Git.AutoCommit)
			ux.Debug.Printfln("git autopull: %v", cfg.Git.AutoPull)

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

	ejectCmd := EjectCommand()
	ejectCmd.GroupID = fleekGroup.ID

	searchCmd := SearchCommand()
	searchCmd.GroupID = packageGroup.ID

	infoCmd := InfoCommand()
	infoCmd.GroupID = packageGroup.ID
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
	command.AddCommand(ejectCmd)
	command.AddCommand(searchCmd)
	command.AddCommand(infoCmd)

	command.AddCommand(VersionCmd())

	command.PersistentFlags().BoolVarP(
		&flags.quiet, app.Trans("fleek.quietFlag"), "q", false, app.Trans("fleek.quietFlagDescription"))
	command.PersistentFlags().BoolVarP(
		&flags.verbose, app.Trans("fleek.verboseFlag"), "v", false, app.Trans("fleek.verboseFlagDescription"))

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
