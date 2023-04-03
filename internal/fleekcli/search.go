/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/cache"
	"github.com/ublue-os/fleek/internal/ux"
)

type searchCmdFlags struct {
	update bool
	fuzzy  bool
}

func SearchCommand() *cobra.Command {
	flags := searchCmdFlags{}
	command := &cobra.Command{
		Use:     app.Trans("search.use"),
		Short:   app.Trans("search.short"),
		Long:    app.Trans("search.long"),
		Args:    cobra.ExactArgs(1),
		Example: app.Trans("search.example"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return search(cmd, args)
		},
	}
	command.Flags().BoolVarP(
		&flags.update, app.Trans("search.updateFlag"), "u", false, app.Trans("search.updateFlagDescription"))
	command.Flags().BoolVarP(
		&flags.fuzzy, app.Trans("search.fuzzyFlag"), "f", false, app.Trans("search.fuzzyFlagDescription"))

	return command
}

// initCmd represents the init command
func search(cmd *cobra.Command, args []string) error {
	ux.Description.Println(cmd.Short)

	var update bool
	if cmd.Flag(app.Trans("search.updateFlag")).Changed {
		update = true
	}
	var fuzzy bool
	if cmd.Flag(app.Trans("search.fuzzyFlag")).Changed {
		fuzzy = true
	}
	if fuzzy {
		ux.Info.Println(app.Trans("search.fuzzyEnabled"))
	}

	needle := args[0]
	spinner, err := ux.Spinner().Start(app.Trans("search.openingCache"))
	if err != nil {
		return err
	}
	pc, err := cache.New()
	if err != nil {
		_ = spinner.Stop()
		ux.Error.Println(app.Trans("search.cacheError"))
		return err
	}
	spinner.Success()
	if update {
		spinner, err := ux.Spinner().Start(app.Trans("search.updatingCache"))
		if err != nil {
			return err
		}
		err = pc.Update()
		if err != nil {
			_ = spinner.Stop()
			ux.Error.Println(app.Trans("search.cacheError"))
			return err
		}
		spinner.Success()
	}
	var hits []cache.Package
	var exactHits []cache.Package
	for i, p := range pc.Packages {
		var hit bool
		if fuzzy {
			if strings.Contains(i, needle) {
				hit = true
			}
			if strings.Contains(p.Name, needle) {
				hit = true
			}
			if strings.Contains(p.Description, needle) {
				hit = true
			}
		}

		if p.Name == needle {
			exactHits = append(exactHits, p)
		}
		if hit {
			hits = append(hits, p)
		}
	}

	// print inexact matches first so exact matches
	// are at the bottom of the output
	if fuzzy {
		if len(hits) == 0 {
			ux.Warning.Println(app.Trans("search.noResults"))
		} else {
			ux.Info.Println(app.Trans("search.fuzzyMatches"))
			_ = ux.Table().WithHasHeader(true).WithData(toTableDataWithHeader(hits)).Render()
		}
	}
	if len(exactHits) == 0 {
		// don't display if we already displayed fuzzy results
		if !fuzzy {
			ux.Warning.Println(app.Trans("search.noResultsExact"))
		}
	} else {
		ux.Info.Println(app.Trans("search.exactMatches"))
		_ = ux.Table().WithHasHeader(true).WithData(toTableDataWithHeader(exactHits)).Render()

	}

	if len(exactHits) > 0 {
		for _, h := range exactHits {
			// TODO: i18n
			ux.Info.Printfln(app.Trans("search.try", h.Name, h.Name))
		}
	}

	ux.Success.Println(app.Trans("global.completed"))
	return nil
}

func toTableDataWithHeader(pp []cache.Package) pterm.TableData {
	var table pterm.TableData

	header := []string{app.Trans("search.package"), app.Trans("search.version"), app.Trans("search.description")}
	table = append(table, header)

	for _, p := range pp {
		row := []string{p.Name, p.Version, p.Description}
		table = append(table, row)
	}
	return table
}
