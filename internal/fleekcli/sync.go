/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/debug"
	"github.com/ublue-os/fleek/internal/ux"
)

func SyncCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   app.Trans("sync.use"),
		Short: app.Trans("sync.short"),
		Long:  app.Trans("sync.long"),

		RunE: func(cmd *cobra.Command, args []string) error {
			return sync(cmd)
		},
	}
	return command
}

// initCmd represents the init command
func sync(cmd *cobra.Command) error {
	var verbose bool
	if cmd.Flag(app.Trans("fleek.verboseFlag")).Changed {
		verbose = true
	}
	ux.Description.Println(cmd.Short)

	if verbose {
		ux.Info.Println(app.Trans("sync.flakeStatus") + f.flakeStatus.String())
		ux.Info.Println(app.Trans("sync.remoteStatus") + f.gitStatus.String())
	}
	spinner, err := ux.Spinner().Start(app.Trans("sync.gettingLocalStatus"))
	if err != nil {
		return err
	}
	dirty, out, err := f.repo.Dirty()
	if err != nil {
		spinner.Fail()
		return err
	}
	if dirty {
		f.flakeStatus = FlakeDirty
	}
	if verbose {
		ux.Info.Println(string(out))
	}
	ux.Info.Println(f.flakeStatus.String())

	spinner.Success()
	// winging it on the logic here.
	if f.flakeStatus == FlakeDirty {
		spinner, err := ux.Spinner().Start(app.Trans("sync.gitCommit"))
		if err != nil {
			return err
		}
		out, err := f.repo.Commit()
		if err != nil {
			spinner.Fail()
			return err
		}
		spinner.Success()
		if verbose {
			ux.Info.Println(string(out))
		}
	}
	spinner, err = ux.Spinner().Start(app.Trans("sync.gettingRemoteStatus"))
	if err != nil {
		return err
	}
	ahead, behind, out, err := f.repo.AheadBehind(false)
	if err != nil {
		debug.Log("git status error: %s", err)
		spinner.Fail()
		return err
	}
	if verbose {
		ux.Info.Println(string(out))
	}
	spinner.Success()
	debug.Log("ahead: %v behind: %v", ahead, behind)
	if ahead {
		debug.Log("remote status: %s", FlakeAhead.String())

		if verbose {
			ux.Info.Println(app.Trans("sync.remoteStatus") + app.Trans("fleek.aheadStatus"))
		}
		f.flakeStatus = FlakeAhead
	}
	if behind {
		debug.Log("remote status: %s", FlakeBehind.String())
		if verbose {
			ux.Info.Println(app.Trans("sync.remoteStatus") + app.Trans("fleek.behindStatus"))
		}
		f.flakeStatus = FlakeBehind
	}
	if ahead && behind {
		debug.Log("remote status: %s", FlakeDiverged.String())
		if verbose {
			ux.Info.Println(app.Trans("sync.remoteStatus") + app.Trans("fleek.divergedStatus"))
		}
		f.flakeStatus = FlakeDiverged
	}
	switch f.flakeStatus {
	case FlakeAhead:
		if verbose {
			ux.Info.Println(app.Trans("sync.gitPush"))
		}
		out, err := f.repo.Push()
		if err != nil {
			debug.Log("git push: %s", err)
			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
	case FlakeBehind:
		if verbose {
			ux.Info.Println(app.Trans("sync.gitPull"))
		}
		out, err := f.repo.Pull()
		if err != nil {
			debug.Log("git pull: %s", err)
			if verbose {
				ux.Error.Println(string(out))
			}
			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
	case FlakeDiverged:
		if verbose {
			ux.Info.Println(app.Trans("sync.gitPull"))
		}
		out, err := f.repo.Pull()
		if err != nil {
			debug.Log("git pull: %s", err)
			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
		if verbose {
			ux.Info.Println(app.Trans("sync.gitPush"))
		}
		out, err = f.repo.Push()
		if err != nil {
			debug.Log("git push: %s", err)
			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
	default:
		ux.Info.Println(app.Trans("sync.noChanges"))
	}
	ux.Success.Println(app.Trans("global.completed"))
	return nil
}
