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
		app.Trans("listPackages.use"),
		app.Trans("listPackages.long"),
		app.Trans("listPackages.short"),
		list,
	)
	return cmd
}

// initCmd represents the init command
func list(cmd *cobra.Command, args []string) {

	cmdr.Info.Println(app.Trans("listPackages.userBling"), strings.ToUpper(f.config.Bling))
	switch f.config.Bling {
	case "high":
		cmdr.Info.Println(app.Trans("listPackages.highBling"))
	case "default":
		cmdr.Info.Println(app.Trans("listPackages.defaultBling"))
	case "low":
		cmdr.Info.Println(app.Trans("listPackages.lowBling"))

	}
	if f.config.Bling == "high" {
		for _, pkg := range core.HighPackages {
			fmt.Printf("\t%s\n", pkg)
		}
	}
	if f.config.Bling == "default" || f.config.Bling == "high" {

		for _, pkg := range core.DefaultPackages {
			f("\t%s\n", pkg)
		}
	}
	for _, pkg := range core.LowPackages {
		fmt.Printf("\t%s\n", pkg)
	}

	cmdr.Info.Println(app.Trans("listPackages.userInstalled"))

	for _, pkg := range f.config.Packages {
		fmt.Printf("\t%s\n", pkg)
	}

}
