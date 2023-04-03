package fleekcli

import (
	"io"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/debug"
	"github.com/ublue-os/fleek/internal/ux"
)

var f *Fleek

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
			var err error
			debug.Log("initializing fleek controller")
			f, err = initFleek(false)
			cobra.CheckErr(err)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if flags.quiet {
				cmd.SetErr(io.Discard)
			}

			debug.Log("repo status")
			if f.repo != nil {
				dirty, _, err := f.repo.Dirty()
				cobra.CheckErr(err)
				if dirty {
					ux.Warning.Println(app.Trans("fleek.dirty"))
				}
				debug.Log("getting remote status")
				ahead, behind, _, err := f.repo.AheadBehind(false)
				cobra.CheckErr(err)
				debug.Log("ahead: %v", ahead)
				debug.Log("behind: %v", behind)

				if ahead {
					debug.Log("remote status: %s", FlakeAhead.String())

					ux.Warning.Println("Remote Status: " + app.Trans("fleek.aheadStatus"))
				}
				if behind {
					debug.Log("remote status: %s", FlakeBehind.String())
					ux.Warning.Println("Remote Status: " + app.Trans("fleek.behindStatus"))

				}
				if ahead && behind {
					debug.Log("remote status: %s", FlakeDiverged.String())
					ux.Warning.Println("Remote Status: " + app.Trans("fleek.divergedStatus"))
				}
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
	gitGroup := &cobra.Group{
		ID:    "gitgroup",
		Title: app.Trans("global.gitGroup"),
	}
	command.AddGroup(initGroup, packageGroup, gitGroup, fleekGroup)
	addCmd := AddCommand()
	addCmd.GroupID = packageGroup.ID

	removeCmd := RemoveCommand()
	removeCmd.GroupID = packageGroup.ID

	syncCmd := SyncCmd()
	syncCmd.GroupID = gitGroup.ID

	showCmd := ShowCmd()
	showCmd.GroupID = fleekGroup.ID
	repoAddCmd := RepoAddCmd()
	repoAddCmd.GroupID = gitGroup.ID
	repoShowCmd := RepoShowCmd()
	repoShowCmd.GroupID = gitGroup.ID
	repoCmd := RepoCmd()
	repoCmd.GroupID = gitGroup.ID

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
	command.AddCommand(syncCmd)
	command.AddCommand(addCmd)
	command.AddCommand(removeCmd)
	command.AddCommand(applyCmd)
	command.AddCommand(updateCmd)
	command.AddCommand(repoCmd)
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
