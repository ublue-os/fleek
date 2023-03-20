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
		Use:   fleek.Trans("repo.use"),
		Long:  fleek.Trans("repo.long"),
		Short: fleek.Trans("repo.short"),
	}
	cmdrcmd := &cmdr.Command{Command: cmd}
	return cmdrcmd
}
