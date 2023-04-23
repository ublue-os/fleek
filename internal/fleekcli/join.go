/*
Copyright © 2023 Brian Ketelsen <bketelsen@gmail.com>
*/
package fleekcli

import (
	"errors"
	"fmt"
	"os"

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

	dirName, err := flake.CloneRepository(args[0], &outBuffer)
	if err != nil {
		return err
	}
	if verbose {
		fmt.Println(outBuffer.String())
	}
	// read config
	config, err := fleek.ReadConfig(dirName)
	if err != nil {
		return err
	}

	_, err = os.Stat(config.FlakeDir)
	if err == nil {
		// exists
		return errors.New("target configuration directory already exists")
	}

	err = os.Rename(dirName, config.FlakeDir)
	if err != nil {
		return err
	}
	// move cloned repo
	fin.Info.Println(app.Trans("init.joining"))
	fl, err := flake.Load(config, app)
	if err != nil {
		return err
	}
	err = fl.Join(&outBuffer)
	if err != nil {
		return err
	}
	err = fl.Write("join new system", &outBuffer)
	if err != nil {
		fin.Debug.Printfln("flake write error: %s", err)
		return err
	}

	fin.Info.Println(app.Trans("join.complete"))

	return nil
}
