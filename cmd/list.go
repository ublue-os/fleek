/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewListCommand() *cmdr.Command {
	cmd := &cobra.Command{
		Use:   fleek.Trans("list.use"),
		Long:  fleek.Trans("list.long"),
		Short: fleek.Trans("list.short"),
	}
	cmdrcmd := &cmdr.Command{Command: cmd}
	return cmdrcmd
}
