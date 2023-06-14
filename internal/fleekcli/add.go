package fleekcli

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/cache"
	"github.com/ublue-os/fleek/internal/flake"
)

type addCmdFlags struct {
	apply bool
}

func AddCommand() *cobra.Command {
	flags := addCmdFlags{}
	command := &cobra.Command{
		Use:     app.Trans("add.use"),
		Short:   app.Trans("add.short"),
		Long:    app.Trans("add.long"),
		Args:    cobra.MinimumNArgs(1),
		Example: app.Trans("add.example"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return add(cmd, args)
		},
	}
	command.Flags().BoolVarP(
		&flags.apply, app.Trans("add.applyFlag"), "a", false, app.Trans("add.applyFlagDescription"))

	return command
}

// initCmd represents the init command
func add(cmd *cobra.Command, args []string) error {

	fin.Description.Println(cmd.Short)
	err := mustConfig()
	if err != nil {
		return err
	}
	var apply bool
	if cmd.Flag(app.Trans("add.applyFlag")).Changed {
		apply = true
	}

	fl, err := flake.Load(cfg, app)
	if err != nil {
		return err
	}
	err = fl.MayPull()
	if err != nil {
		return err
	}
	pc, err := cache.New()
	if err != nil {
		fin.Error.Println(app.Trans("search.cacheError"))
		return err
	}
	var hits []cache.SearchResult
	var exactHits []cache.SearchResult

	var sb strings.Builder
	sb.WriteString("add packages: ")
	for _, p := range args {
		for i, pack := range pc.Packages {
			var hit bool
			if strings.Contains(i, p) {
				hit = true
			}
			if strings.Contains(pack.Name, p) {
				hit = true
			}
			if strings.Contains(pack.Description, p) {
				hit = true
			}
			firstPeriod := strings.Index(i, ".")
			sanitizedPackageName := i[firstPeriod+1:]
			secondPeriod := strings.Index(sanitizedPackageName, ".")
			sanitizedPackageName = sanitizedPackageName[secondPeriod+1:]
			if p == sanitizedPackageName {
				exactHits = append(exactHits, cache.SearchResult{Name: sanitizedPackageName, Package: pack})
			}
			if hit {
				hits = append(hits, cache.SearchResult{Name: sanitizedPackageName, Package: pack})
			}
		}
		if len(exactHits) == 1 {
			fin.Info.Println("Found exact match for " + p)
		}
		if len(exactHits) < 1 {
			if len(hits) > 0 {
				fin.Info.Println("Found " + fmt.Sprint(len(hits)) + " matche(s) for " + p)
				for _, hit := range hits {
					fin.Info.Println("\tName: ", hit.Name)
					fin.Info.Println("\tDescription: ", hit.Package.Description)
					fin.Warning.Printfln("\tRun `fleek add %s` to add it.", hit.Name)

				}
				return nil

			}
			fin.Info.Println("Found no matches for " + p + "!")
			return nil

		}
		fin.Info.Println("exact hits", len(exactHits))
		fin.Info.Println("possible matches", len(hits))

		fin.Info.Println(app.Trans("add.adding") + p)
		err = fl.Config.AddPackage(p)
		if err != nil {
			fin.Debug.Printfln("add package error: %s", err)
			return err
		}
		sb.WriteString(p + " ")

	}
	err = fl.Write(sb.String(), false, false)
	if err != nil {
		fin.Debug.Printfln("flake write error: %s", err)
		return err
	}

	if apply {
		fin.Info.Println(app.Trans("add.applying"))

		err = fl.Apply()
		if err != nil {
			if errors.Is(err, flake.ErrPackageConflict) {
				fin.Fatal.Println(app.Trans("global.errConflict"))
			}
			return err
		}
	} else {
		fin.Warning.Println(app.Trans("add.unapplied"))
	}
	fin.Success.Println(app.Trans("global.completed"))
	return nil
}
