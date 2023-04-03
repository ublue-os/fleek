/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/debug"
	"github.com/ublue-os/fleek/internal/nix"
	"github.com/ublue-os/fleek/internal/ux"
)

type applyCmdFlags struct {
	push   bool
	dryRun bool
}

func ApplyCommand() *cobra.Command {
	flags := applyCmdFlags{}
	command := &cobra.Command{
		Use:     app.Trans("apply.use"),
		Short:   app.Trans("apply.short"),
		Long:    app.Trans("apply.long"),
		Example: app.Trans("apply.example"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return apply(cmd)
		},
	}
	command.Flags().BoolVarP(
		&flags.push, app.Trans("apply.pushFlag"), "a", false, app.Trans("apply.pushFlagDescription"))
	command.Flags().BoolVarP(
		&flags.dryRun, app.Trans("apply.dryRunFlag"), "d", false, app.Trans("apply.dryRunFlagDescription"))

	return command
}

func apply(cmd *cobra.Command) error {
	var verbose bool
	var push bool
	if cmd.Flag(app.Trans("fleek.verboseFlag")).Changed {
		verbose = true
	}
	ux.Description.Println(cmd.Short)

	if cmd.Flag(app.Trans("apply.pushFlag")).Changed {
		push = true
	}
	if f.flakeStatus == FlakeBehind {
		ux.Warning.Println(app.Trans("apply.behind"))
		return nil
	}
	if verbose {
		ux.Info.Println(app.Trans("apply.writingConfig"))
	}
	// only re-apply the templates if not `ejected`
	if !f.config.Ejected {
		if verbose {
			ux.Info.Println(cmd.Use, app.Trans("apply.writingFlake"))
		}
		flake, err := f.Flake()
		if err != nil {
			debug.Log("open flake error: %s", err)
			return err
		}
		err = flake.Write(false)
		if err != nil {
			debug.Log("write flake error: %s", err)
			return err
		}
		repo, err := f.Repo()
		if err != nil {
			debug.Log("get repo error: %s", err)
			return err
		}
		out, err := repo.Commit()
		if err != nil {
			debug.Log("git commit error: %s", err)
			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
	}
	var dry bool
	if cmd.Flag(app.Trans("apply.dryRunFlag")).Changed {
		dry = true
	}
	if !dry {
		ux.Info.Println(app.Trans("apply.applyingConfig"))
		flake, err := f.Flake()
		if err != nil {
			debug.Log("get flake error: %s", err)
			return err
		}
		out, err := flake.Apply()
		if err != nil {
			ux.Error.Println(string(out))

			if errors.Is(err, nix.ErrPackageConflict) {
				ux.Fatal.Println(app.Trans("global.errConflict"))
			}

			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
		r, err := f.Repo()
		if err != nil {
			debug.Log("get repo error: %s", err)
			return err
		}
		out, err = r.Commit()
		if err != nil {
			debug.Log("git commit error: %s", err)
			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
	} else {
		ux.Info.Println(app.Trans("apply.dryApplyingConfig"))
		flake, err := f.Flake()
		if err != nil {
			debug.Log("get flake error: %s", err)
			return err
		}
		out, err := flake.Check()
		if err != nil {
			debug.Log("flake check error: %s", err)
			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
	}
	if push {
		ux.Info.Println(app.Trans("apply.pushing"))
		repo, err := f.Repo()
		if err != nil {
			debug.Log("get repo error: %s", err)
			return err
		}
		out, err := repo.Push()
		if err != nil {
			debug.Log("git push error: %s", err)
			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
	}

	ux.Success.Println(app.Trans("apply.done"))
	return nil
}
