/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/flake"
	"github.com/ublue-os/fleek/internal/fleek"
	"github.com/ublue-os/fleek/internal/ux"
)

func InfoCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     app.Trans("info.use"),
		Short:   app.Trans("info.short"),
		Long:    app.Trans("info.long"),
		Example: app.Trans("info.example"),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return infoFleek(cmd, args)
		},
	}
	return command
}

func infoFleek(cmd *cobra.Command, args []string) error {

	ux.Description.Println(cmd.Short)

	fl, err := flake.Load(cfg, app)
	if err != nil {
		return err
	}
	var b *fleek.Bling

	switch fl.Config.Bling {
	case "high":
		b, err = fleek.HighBling()
		cobra.CheckErr(err)
	case "default":
		b, err = fleek.DefaultBling()
		cobra.CheckErr(err)
	case "low":
		b, err = fleek.LowBling()
		cobra.CheckErr(err)
	case "none":
		b, err = fleek.NoBling()
		cobra.CheckErr(err)
	}
	ux.Info.Println("["+b.Name+" Bling]", b.Description)

	needle := args[0]
	var found bool
	pkg, ok := b.PackageMap[needle]
	if ok {
		found = true
		ux.Info.Println(" -- " + pkg.Name + " --")
		ux.Description.Println(pkg.Description)
	}
	prog, ok := b.ProgramMap[needle]
	if ok {
		found = true
		ux.Info.Println(" -- " + prog.Name + " --")
		ux.Description.Println(prog.Description)
		if len(prog.Aliases) > 0 {

			ux.Info.Println(app.Trans("info.aliases"))
			for _, a := range prog.Aliases {
				ux.Description.Println("\t" + a.Description)
				ux.Info.Println("\t\t"+a.Key+": ", a.Value)
			}

		}
	}
	if !found {
		ux.Warning.Println(needle, "-", app.Trans("info.notFound"))
	}
	return nil
}
