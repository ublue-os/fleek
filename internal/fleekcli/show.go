/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/core"
	"github.com/ublue-os/fleek/internal/ux"
)

func ShowCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   app.Trans("show.use"),
		Short: app.Trans("show.short"),
		Long:  app.Trans("show.long"),

		RunE: func(cmd *cobra.Command, args []string) error {
			return showFleek(cmd)
		},
	}
	return command
}

// initCmd represents the init command
func showFleek(cmd *cobra.Command) error {

	ux.Description.Println(cmd.Short)
	var b *core.Bling
	var err error
	switch f.config.Bling {
	case "high":
		b, err = core.HighBling()
		cobra.CheckErr(err)
	case "default":
		b, err = core.DefaultBling()
		cobra.CheckErr(err)
	case "low":
		b, err = core.LowBling()
		cobra.CheckErr(err)
	case "none":
		b, err = core.NoBling()
		cobra.CheckErr(err)
	}
	ux.Info.Println("["+b.Name+" Bling]", b.Description)

	var packages []string
	for n := range b.PackageMap {
		//fmt.Println(style.Render(n))
		packages = append(packages, n)
		//fmt.Println(style.Render(p.Description))
	}
	var programs []string
	for n := range b.ProgramMap {
		//fmt.Println(style.Render(n))
		programs = append(programs, n)
		//fmt.Println(style.Render(p.Description))
	}

	ux.ThreeColumnList(
		"["+b.Name+"] "+app.Trans("show.packages"), packages,
		"["+b.Name+"] "+app.Trans("show.managedPackages"), programs,
		app.Trans("show.userPackages"), f.config.Packages,
	)
	return nil
}
