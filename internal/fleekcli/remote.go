/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
)

func RepoCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   app.Trans("remote.use"),
		Short: app.Trans("remote.short"),
		Long:  app.Trans("remote.long"),
	}
	command.AddCommand(RepoShowCmd())
	command.AddCommand(RepoAddCmd())

	return command
}
