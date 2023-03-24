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
		Use:   app.Trans("remote.use"),
		Long:  app.Trans("remote.long"),
		Short: app.Trans("remote.short"),
	}
	cmdrcmd := &cmdr.Command{Command: cmd}
	return cmdrcmd
}
