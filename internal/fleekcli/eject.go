/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
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
func eject(_ *cobra.Command) error {
	err := mustConfig()
	if err != nil {
		return err
	}
	return nil
}
