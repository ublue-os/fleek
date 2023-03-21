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

func NewListPackagesCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		fleek.Trans("listPackages.use"),
		fleek.Trans("listPackages.long"),
		fleek.Trans("listPackages.short"),
		list,
	)
	return cmd
}

// initCmd represents the init command
func list(cmd *cobra.Command, args []string) {

	conf, err := core.ReadConfig()
	cobra.CheckErr(err)

	cmdr.Info.Println(fleek.Trans("listPackages.userBling"), strings.ToUpper(conf.Bling))
	switch conf.Bling {
	case "high":
		cmdr.Info.Println(fleek.Trans("listPackages.highBling"))
	case "default":
		cmdr.Info.Println(fleek.Trans("listPackages.defaultBling"))
	case "low":
		cmdr.Info.Println(fleek.Trans("listPackages.lowBling"))

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

	cmdr.Info.Println(fleek.Trans("listPackages.userInstalled"))

	for _, pkg := range conf.Packages {
		fmt.Printf("\t%s\n", pkg)
	}

}
