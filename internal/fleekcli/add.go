/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/debug"
	"github.com/ublue-os/fleek/internal/ux"
)

type addCmdFlags struct {
	apply bool
}

func AddCommand() *cobra.Command {
	flags := addCmdFlags{}
	command := &cobra.Command{
		Use:   app.Trans("add.use"),
		Short: app.Trans("add.short"),
		Long:  app.Trans("add.long"),
		Args:  cobra.MinimumNArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			return add(cmd, args)
		},
	}
	command.Flags().BoolVarP(
		&flags.apply, app.Trans("add.applyFlag"), "a", false, app.Trans("add.applyFlagDescription"))

	return command
}

// initCmd represents the init command
func add(cmd *cobra.Command, args []string) error {
	var verbose bool
	if cmd.Flag(app.Trans("fleek.verboseFlag")).Changed {
		verbose = true
	}
	ux.Description.Println(cmd.Short)

	var apply bool
	if cmd.Flag(app.Trans("add.applyFlag")).Changed {
		apply = true
	}

	var err error

	for _, p := range args {

		ux.Info.Println(app.Trans("add.adding") + p)

		err = f.config.AddPackage(p)
		if err != nil {
			debug.Log("add package error: %s", err)
			return err
		}

	}

	if apply {
		ux.Info.Println(app.Trans("add.applying"))
		flake, err := f.Flake()
		if err != nil {
			debug.Log("get flake error: %s", err)
			return err
		}
		err = flake.Write(false)
		if err != nil {
			debug.Log("flake write error: %s", err)
			return err
		}
		repo, err := f.Repo()
		if err != nil {
			debug.Log("get repo error: %s", err)
			return err
		}
		out, err := repo.Commit()
		if err != nil {
			debug.Log("commit error: %s", err)
			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
		out, err = flake.Apply()
		if err != nil {
			debug.Log("flake apply error: %s", err)
			return err
		}
		if verbose {
			ux.Info.Println(string(out))
		}
	}

	ux.Success.Println(app.Trans("add.done"))
	return nil
}
