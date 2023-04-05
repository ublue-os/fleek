/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/debug"
	"github.com/ublue-os/fleek/internal/flake"
	"github.com/ublue-os/fleek/internal/ux"
)

type remoteAddCmdFlags struct {
	name string
}

func RepoAddCmd() *cobra.Command {
	flags := remoteAddCmdFlags{}
	command := &cobra.Command{
		Use:     app.Trans("remoteadd.use"),
		Short:   app.Trans("remoteadd.short"),
		Long:    app.Trans("remoteadd.long"),
		Args:    cobra.ExactArgs(1),
		Example: app.Trans("remoteadd.example"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return remoteadd(cmd, args)
		},
	}
	command.Flags().StringVarP(
		&flags.name, app.Trans("remoteadd.nameFlag"), "n", "origin", app.Trans("remoteadd.nameFlagDescription"))

	return command
}

// initCmd represents the init command
func remoteadd(cmd *cobra.Command, args []string) error {
	var verbose bool
	if cmd.Flag(app.Trans("fleek.verboseFlag")).Changed {
		verbose = true
	}
	ux.Description.Println(cmd.Short)

	fl, err := flake.Load(cfg, app)
	if err != nil {
		return err
	}

	fl.Config.Repository = args[0]
	if verbose {
		ux.Info.Println(app.Trans("remoteAdd.saving"))
	}
	err = fl.Config.Save()
	if err != nil {
		debug.Log("save config error: %s", err)
		ux.Error.Println(err)
		return err
	}
	// now actually add the remote
	if verbose {
		ux.Info.Println(app.Trans("remoteAdd.addingRemote"))
	}
	name := cmd.Flag("name").Value.String()

	err = fl.RemoteAdd(args[0], name)
	if err != nil {
		debug.Log("adding remote repo error: %s", err)
		ux.Error.Println(err)
		return err
	}
	err = fl.Commit("fleek: add remote repository")
	if err != nil {
		debug.Log("commit error: %s", err)
		return err
	}

	ux.Info.Println(fl.Config.Repository)
	ux.Success.Println(app.Trans("remoteAdd.done"))
	return nil
}
