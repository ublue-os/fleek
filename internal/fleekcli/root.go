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
				dirty, _, err := f.repo.Dirty(false)
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
		Title: "Getting Started",
	}
	fleekGroup := &cobra.Group{
		ID:    "fleek",
		Title: "Configuration Commands",
	}

	packageGroup := &cobra.Group{
		ID:    "package",
		Title: "Package Management Commands",
	}
	gitGroup := &cobra.Group{
		ID:    "gitgroup",
		Title: "Git Commands",
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

	command.AddCommand(showCmd)
	command.AddCommand(syncCmd)
	command.AddCommand(addCmd)
	command.AddCommand(removeCmd)
	command.AddCommand(applyCmd)
	command.AddCommand(updateCmd)
	command.AddCommand(repoCmd)
	command.AddCommand(initCmd)
	command.AddCommand(ejectCmd)

	command.PersistentFlags().BoolVarP(
		&flags.quiet, app.Trans("fleek.quietFlag"), "q", false, app.Trans("fleek.quietFlagDescription"))
	command.PersistentFlags().BoolVarP(
		&flags.verbose, app.Trans("fleek.verboseFlag"), "v", false, app.Trans("fleek.verboseFlagDescription"))

	debugMiddleware.AttachToFlag(command.PersistentFlags(), app.Trans("fleek.debugFlag"))
	traceMiddleware.AttachToFlag(command.PersistentFlags(), app.Trans("fleek.traceFlag"))

	return command
}
