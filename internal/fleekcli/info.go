/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/flake"
	"github.com/ublue-os/fleek/internal/fleek"
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

	fin.Description.Println(cmd.Short)
	err := mustConfig()
	if err != nil {
		return err
	}
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
	fin.Info.Println("["+b.Name+" Bling]", b.Description)

	needle := args[0]
	var found bool
	pkg, ok := b.PackageMap[needle]
	if ok {
		found = true

		fmt.Println(fin.TitleSectionPrinter(pkg.Name))
		fmt.Println(fin.DescriptionSectionPrinter(app.Trans("info.description")))
		fmt.Println(fin.ParagraphPrinter(pkg.Description))

	}
	prog, ok := b.ProgramMap[needle]
	if ok {
		found = true
		fmt.Println(fin.TitleSectionPrinter(prog.Name))
		fmt.Println(fin.DescriptionSectionPrinter(app.Trans("info.description")))
		fmt.Println(fin.ParagraphPrinter(prog.Description))

		if len(prog.Aliases) > 0 {
			fmt.Println(fin.DescriptionSectionPrinter(app.Trans("info.aliases")))
			var td pterm.TableData
			td = append(td, []string{"Alias", "Value", "Description"})

			for _, a := range prog.Aliases {
				td = append(td, []string{a.Key, a.Value, a.Description})
			}
			_ = fin.Table().WithHasHeader(true).WithHeaderRowSeparator("-").WithData(td).Render()

		}
	}
	if !found {
		fin.Warning.Println(needle, "-", app.Trans("info.notFound"))
	}
	return nil
}
