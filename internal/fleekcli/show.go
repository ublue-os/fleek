/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/core"
	"github.com/ublue-os/fleek/internal/ux"
)

// add json flag
// add level override flag
type showCmdFlags struct {
	json  bool
	level string
}

func ShowCmd() *cobra.Command {
	flags := showCmdFlags{}

	command := &cobra.Command{
		Use:     app.Trans("show.use"),
		Short:   app.Trans("show.short"),
		Long:    app.Trans("show.long"),
		Example: app.Trans("show.example"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return showFleek(cmd)
		},
	}
	command.Flags().BoolVarP(
		&flags.json, app.Trans("show.jsonFlag"), "j", false, app.Trans("show.jsonFlagDescription"))
	command.Flags().StringVarP(
		&flags.level, app.Trans("show.levelFlag"), "l", "", app.Trans("show.levelFlagDescription"))
	return command
}

// initCmd represents the init command
func showFleek(cmd *cobra.Command) error {
	var showJSON bool
	var level string

	if cmd.Flag(app.Trans("show.jsonFlag")).Changed {
		showJSON = true
	}

	if cmd.Flag(app.Trans("show.levelFlag")).Changed {
		level = cmd.Flag(app.Trans("show.levelFlag")).Value.String()
	} else {
		level = f.config.Bling
	}
	if !showJSON {
		ux.Description.Println(cmd.Short)
	}
	var b *core.Bling
	var err error

	switch level {
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
	default:
		ux.Error.Println(app.Trans("show.invalidLevel", level))
		return nil
	}

	if !showJSON {
		ux.Info.Println("["+b.Name+" Bling]", b.Description)
	}

	if showJSON {
		bb, err := json.Marshal(b)
		if err != nil {
			return err
		}
		fmt.Println(string(bb))
		return nil
	}

	ux.ThreeColumnList(
		"["+b.Name+"] "+app.Trans("show.packages"), b.Packages,
		"["+b.Name+"] "+app.Trans("show.managedPackages"), b.Programs,
		app.Trans("show.userPackages"), f.config.Packages,
	)

	return nil
}
