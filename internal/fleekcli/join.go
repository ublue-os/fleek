/*
Copyright © 2023 Brian Ketelsen <bketelsen@gmail.com>
*/
package fleekcli

import (
	"errors"
	"os"

	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/flake"
	"github.com/ublue-os/fleek/internal/fleek"
)

type joinCmdFlags struct {
	apply bool
}

func JoinCommand() *cobra.Command {
	flags := joinCmdFlags{}
	command := &cobra.Command{
		Use:     app.Trans("join.use"),
		Short:   app.Trans("join.short"),
		Long:    app.Trans("join.long"),
		Example: app.Trans("join.example"),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return join(cmd, args)
		},
	}
	command.Flags().BoolVarP(
		&flags.apply, app.Trans("join.applyFlag"), "a", false, app.Trans("join.applyFlagDescription"))
	return command
}

// joinCmd represents the join command
func join(cmd *cobra.Command, args []string) error {

	var verbose bool
	if cmd.Flag(app.Trans("fleek.verboseFlag")).Changed {
		verbose = true
	}

	cfg.Verbose = verbose

	fin.Description.Println(cmd.Short)

	dirName, err := flake.CloneRepository(args[0])
	if err != nil {
		return err
	}

	// read config
	config, err := fleek.ReadConfig(dirName)
	if err != nil {
		return err
	}

	_, err = os.Stat(config.UserFlakeDir())
	if err == nil {
		// exists
		return errors.New("target configuration directory already exists")
	}
	// move cloned repo
	err = cp.Copy(dirName, config.UserFlakeDir())
	if err != nil {
		return err
	}
	fin.Info.Println(app.Trans("init.joining"))
	fl, err := flake.Load(config, app)
	if err != nil {
		return err
	}
	err = fl.Join()
	if err != nil {
		return err
	}
	// reload config and flake
	config, err = fleek.ReadConfig(config.UserFlakeDir())
	if err != nil {
		return err
	}
	fl, err = flake.Load(config, app)
	if err != nil {
		return err
	}
	err = fl.Write("join new system", false)
	if err != nil {
		fin.Debug.Printfln("flake write error: %s", err)
		return err
	}

	if cmd.Flag(app.Trans("join.applyFlag")).Changed {
		err = fl.Apply()
		if err != nil {
			return err
		}
		return nil
	}
	fin.Info.Println(app.Trans("join.complete"))

	return nil
}
