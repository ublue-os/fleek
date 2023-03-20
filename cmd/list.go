/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/core"
	"github.com/vanilla-os/orchid/cmdr"
)

func NewListCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		fleek.Trans("list.use"),
		fleek.Trans("list.long"),
		fleek.Trans("list.short"),
		list,
	)
	return cmd
}

// initCmd represents the init command
func list(cmd *cobra.Command, args []string) {

	conf, err := core.ReadConfig()
	cobra.CheckErr(err)

	cmdr.Info.Println(fleek.Trans("list.userBling"), strings.ToUpper(conf.Bling))
	switch conf.Bling {
	case "high":
		cmdr.Info.Println(fleek.Trans("list.highBling"))
	case "default":
		cmdr.Info.Println(fleek.Trans("list.defaultBling"))
	case "low":
		cmdr.Info.Println(fleek.Trans("list.lowBling"))

	}
	if conf.Bling == "high" {
		for _, pkg := range core.HighPackages {
			fmt.Printf("\t%s\n", pkg)
		}
	}
	if conf.Bling == "default" || conf.Bling == "high" {

		for _, pkg := range core.DefaultPackages {
			fmt.Printf("\t%s\n", pkg)
		}
	}
	for _, pkg := range core.LowPackages {
		fmt.Printf("\t%s\n", pkg)
	}

	cmdr.Info.Println(fleek.Trans("list.userInstalled"))

	for _, pkg := range conf.Packages {
		fmt.Printf("\t%s\n", pkg)
	}

}
