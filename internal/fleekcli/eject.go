/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/ux"
	"github.com/vanilla-os/orchid/cmdr"
)

func EjectCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   app.Trans("eject.use"),
		Short: app.Trans("eject.short"),
		Long:  app.Trans("eject.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return eject(cmd)
		},
	}

	return command
}

// initCmd represents the init command
func eject(cmd *cobra.Command) error {
	ux.Description.Println(cmd.Short)

	ok, err := cmdr.Confirm.Show(app.Trans("eject.confirm"))
	if err != nil {
		return err
	}

	if ok {
		ux.Info.Println(app.Trans("eject.start"))
		flake, err := f.Flake()
		if err != nil {
			return err
		}
		err = flake.Write(true)
		if err != nil {
			return err
		}
		err = f.config.Eject()
		if err != nil {
			return err
		}
		ux.Success.Println(app.Trans("eject.complete"))
	}
	return nil
}
