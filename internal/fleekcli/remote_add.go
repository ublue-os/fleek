/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/debug"
	"github.com/ublue-os/fleek/internal/ux"
)

type remoteAddCmdFlags struct {
	name string
}

func RepoAddCmd() *cobra.Command {
	flags := remoteAddCmdFlags{}
	command := &cobra.Command{
		Use:   app.Trans("remoteadd.use"),
		Short: app.Trans("remoteadd.short"),
		Long:  app.Trans("remoteadd.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return remoteadd(cmd, args)
		},
	}
	command.Flags().StringVarP(
		&flags.name, app.Trans("remoteadd.applyFlag"), "n", "origin", app.Trans("reemoteadd.name"))

	return command
}

// initCmd represents the init command
func remoteadd(cmd *cobra.Command, args []string) error {
	var verbose bool
	if cmd.Flag(app.Trans("fleek.verboseFlag")).Changed {
		verbose = true
	}
	ux.Description.Println(cmd.Short)
	f.config.Repository = args[0]
	if verbose {
		ux.Info.Println(app.Trans("remoteAdd.saving"))
	}
	err := f.config.Save()
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
	repo, err := f.Repo()
	if err != nil {
		debug.Log("getting repo error: %s", err)
		ux.Error.Println(err)
		return err
	}
	err = repo.RemoteAdd(args[0], name)
	if err != nil {
		debug.Log("adding remote repo error: %s", err)
		ux.Error.Println(err)
		return err
	}
	out, err := repo.Commit()
	if err != nil {
		debug.Log("repo commit error: %s", err)
		ux.Error.Println(err)
		return err
	}
	if verbose {
		ux.Info.Println(string(out))
	}
	ux.Info.Println(f.config.Repository)
	ux.Success.Println(app.Trans("remoteAdd.done"))
	return nil
}
