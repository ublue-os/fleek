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

func NewListProgramsCommand() *cmdr.Command {
	cmd := cmdr.NewCommandRun(
		app.Trans("listPrograms.use"),
		app.Trans("listPrograms.long"),
		app.Trans("listPrograms.short"),
		listPrograms,
	)
	return cmd
}

// initCmd represents the init command
func listPrograms(cmd *cobra.Command, args []string) {

	conf, err := core.ReadConfig()
	cobra.CheckErr(err)

	cmdr.Info.Println(app.Trans("listPrograms.userBling"), strings.ToUpper(conf.Bling))
	switch conf.Bling {
	case "high":
		cmdr.Info.Println(app.Trans("listPrograms.highBling"))
	case "default":
		cmdr.Info.Println(app.Trans("listPrograms.defaultBling"))
	case "low":
		cmdr.Info.Println(app.Trans("listPrograms.lowBling"))

	}
	if conf.Bling == "high" {
		for _, pkg := range core.HighPrograms {
			fmt.Printf("\t%s\n", pkg)
		}
	}
	if conf.Bling == "default" || conf.Bling == "high" {

		for _, pkg := range core.DefaultPrograms {
			fmt.Printf("\t%s\n", pkg)
		}
	}
	for _, pkg := range core.LowPrograms {
		fmt.Printf("\t%s\n", pkg)
	}

	cmdr.Info.Println(app.Trans("listPrograms.userInstalled"))

	for _, pkg := range conf.Programs {
		fmt.Printf("\t%s\n", pkg)
	}

}
