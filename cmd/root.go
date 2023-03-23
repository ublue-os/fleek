package cmd

import (
	"embed"
	"os"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/nix"
	"github.com/vanilla-os/orchid/cmdr"
)

var app *cmdr.App
var f *fleek

const (
	verboseFlag string = "verbose"
	syncFlag    string = "sync"
	nixGC       string = "garbage-collect"
)

func New(version string, fs embed.FS) *cmdr.App {
	app = cmdr.NewApp("fleek", version, fs)
	return app
}
func NewRootCommand(version string) *cmdr.Command {
	root := cmdr.NewCommand(
		app.Trans("fleek.use"),
		app.Trans("fleek.long"),
		app.Trans("fleek.short"),
		nil).
		WithPersistentBoolFlag(
			cmdr.NewBoolFlag(
				verboseFlag,
				"v",
				app.Trans("fleek.verboseFlag"),
				false)).
		WithPersistentBoolFlag(
			cmdr.NewBoolFlag(
				syncFlag,
				"s",
				app.Trans("fleek.syncFlag"),
				false)).
		WithPersistentBoolFlag(
			cmdr.NewBoolFlag(
				nixGC,
				"g",
				app.Trans("fleek.nixGarbage"),
				false))

	root.Version = version
	root.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		ok := nix.CheckNix()
		if !ok {
			cmdr.Error.Println(app.Trans("fleek.installNix"))
			os.Exit(1)
		}

		var err error
		f, err = initFleek()
		cobra.CheckErr(err)

		/*
			dirty, err := repo.Dirty()
			cobra.CheckErr(err)
			if dirty {
				cmdr.Warning.Println(app.Trans("fleek.dirty"))
			}
			ahead, behind, err = repo.AheadBehind()
			cobra.CheckErr(err)
			if ahead {
				cmdr.Warning.Println(app.Trans("fleek.ahead"))
			}
			if behind {
				cmdr.Warning.Println(app.Trans("fleek.behind"))
			}

		*/
		if cmd.Flag(syncFlag).Changed && f.flakeStatus == FlakeBehind {
			cmdr.Info.Println(app.Trans("fleek.pull"))
			r, err := f.Repo()
			cobra.CheckErr(err)
			r.Pull()
			cobra.CheckErr(err)
			f.flakeStatus = FlakeExists

		}

	}
	root.PersistentPostRun = func(cmd *cobra.Command, args []string) {
		repo, err := f.Repo()
		cobra.CheckErr(err)
		dirty, err := repo.Dirty()
		cobra.CheckErr(err)
		if dirty {
			f.flakeStatus = FlakeDirty
			cmdr.Warning.Println(app.Trans("fleek.dirty"))
		}
		ahead, behind, err := f.repo.AheadBehind()
		cobra.CheckErr(err)
		if ahead {
			f.flakeStatus = FlakeAhead
			// only show the warning if we're not already
			// planning to sync
			if !cmd.Flag("sync").Changed {
				cmdr.Warning.Println(app.Trans("fleek.ahead"))
			}
		}
		if behind {
			f.flakeStatus = FlakeBehind
			cmdr.Warning.Println(app.Trans("fleek.behind"))
		}
		if cmd.Flag(nixGC).Changed {
			f, err := f.Flake()
			cobra.CheckErr(err)
			err = f.GC()
			cobra.CheckErr(err)

		}
		if cmd.Flag("sync").Changed && f.flakeStatus == FlakeAhead {
			cmdr.Info.Println(app.Trans("fleek.push"))
			repo, err := f.Repo()
			cobra.CheckErr(err)
			err = repo.Push()
			cobra.CheckErr(err)

		}
		/*
			repo = git.NewFlakeRepo(flakeLocation)

			cmdr.Info.Println(app.Trans("fleek.gitStatus"))

			dirty, err := repo.Dirty()
			cobra.CheckErr(err)
			if dirty {
				cmdr.Warning.Println(app.Trans("fleek.dirty"))
			}
			ahead, behind, err = repo.AheadBehind()
			cobra.CheckErr(err)
			if ahead {
				cmdr.Warning.Println(app.Trans("fleek.ahead"))
			}
			if behind {
				cmdr.Warning.Println(app.Trans("fleek.behind"))
			}
			if cmd.Flag("sync").Changed && ahead {
				cmdr.Info.Println(app.Trans("fleek.push"))
				err = repo.Push()
				cobra.CheckErr(err)

			}
		*/

	}
	return root
}
