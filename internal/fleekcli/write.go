/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/flake"
)

// WriteCommand is an internal hidden command that
// gets run after an update.
func WriteCommand() *cobra.Command {
	command := &cobra.Command{
		Hidden: true,
		Use:    app.Trans("write.use"),
		Short:  app.Trans("write.short"),
		Long:   app.Trans("write.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return write(cmd)
		},
	}

	return command
}

// writeCmd represents the write command
func write(cmd *cobra.Command) error {
	fin.Description.Println(cmd.Short)
	err := mustConfig()
	if err != nil {
		return err
	}
	fl, err := flake.Load(cfg, app)
	if err != nil {
		return err
	}

	err = fl.Write("flake update", true, false)
	if err != nil {
		fin.Logger.Error("flake write", fin.Logger.Args("error", err))

		return err
	}

	fin.Success.Println(app.Trans("write.done"))
	return nil
}
