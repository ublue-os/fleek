/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewRepoCommand() *cmdr.Command {
	cmd := &cobra.Command{
		Use:   app.Trans("repo.use"),
		Long:  app.Trans("repo.long"),
		Short: app.Trans("repo.short"),
	}
	cmdrcmd := &cmdr.Command{Command: cmd}
	return cmdrcmd
}
