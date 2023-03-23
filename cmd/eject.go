/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewEjectCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		app.Trans("eject.use"),
		app.Trans("eject.long"),
		app.Trans("eject.short"),
		eject,
	)
	return cmd
}

// initCmd represents the init command
func eject(cmd *cobra.Command, args []string) {

	ok, err := cmdr.Confirm.Show(app.Trans("eject.confirm"))
	cobra.CheckErr(err)

	if ok {
		cmdr.Info.Println(app.Trans("eject.start"))
		flake, err := f.Flake()
		cobra.CheckErr(err)
		err = flake.Write(true)
		cobra.CheckErr(err)
		err = f.config.Eject()
		cobra.CheckErr(err)
		cmdr.Info.Println(app.Trans("eject.complete"))
	}

}
