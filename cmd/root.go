package cmd

import (
	"embed"
	"errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/core"
	"github.com/ublue-os/fleek/git"
	"github.com/ublue-os/fleek/nix"
	"github.com/vanilla-os/orchid/cmdr"
)

var fleek *cmdr.App
var flake *nix.Flake
var config *core.Config
var repo *git.FlakeRepo
var firstrun bool
var flakeLocation string
var ahead bool
var behind bool

const (
	verboseFlag string = "verbose"
	syncFlag    string = "sync"
)

func New(version string, fs embed.FS) *cmdr.App {
	fleek = cmdr.NewApp("fleek", version, fs)
	return fleek
}
func NewRootCommand(version string) *cmdr.Command {
	root := cmdr.NewCommand(
		fleek.Trans("fleek.use"),
		fleek.Trans("fleek.long"),
		fleek.Trans("fleek.short"),
		nil).
		WithPersistentBoolFlag(
			cmdr.NewBoolFlag(
				verboseFlag,
				"v",
				fleek.Trans("fleek.verboseFlag"),
				false)).
		WithPersistentBoolFlag(
			cmdr.NewBoolFlag(
				syncFlag,
				"s",
				fleek.Trans("fleek.syncFlag"),
				false))

	root.Version = version
	root.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		ok := nix.CheckNix()
		if !ok {
			cmdr.Error.Println(fleek.Trans("fleek.installNix"))
			os.Exit(1)
		}

		// set up config and flake before each command
		var err error
		config, err = core.ReadConfig()
		if errors.Is(err, os.ErrNotExist) {
			firstrun = true
			return
		}

		flakeLocation, err = core.FlakeLocation()
		cobra.CheckErr(err)

		flake, err = nix.NewFlake(flakeLocation, config)
		cobra.CheckErr(err)

		repo = git.NewFlakeRepo(flakeLocation)

		cmdr.Info.Println(fleek.Trans("fleek.gitStatus"))

		dirty, err := repo.Dirty()
		cobra.CheckErr(err)
		if dirty {
			cmdr.Warning.Println(fleek.Trans("fleek.dirty"))
		}
		ahead, behind, err = repo.AheadBehind()
		cobra.CheckErr(err)
		if ahead {
			cmdr.Warning.Println(fleek.Trans("fleek.ahead"))
		}
		if behind {
			cmdr.Warning.Println(fleek.Trans("fleek.behind"))
		}
		if cmd.Flag("sync").Changed && behind {
			cmdr.Info.Println(fleek.Trans("fleek.pull"))
			err = repo.Pull()
			cobra.CheckErr(err)
			behind = false
		}

	}
	root.PersistentPostRun = func(cmd *cobra.Command, args []string) {

		repo = git.NewFlakeRepo(flakeLocation)

		cmdr.Info.Println(fleek.Trans("fleek.gitStatus"))

		dirty, err := repo.Dirty()
		cobra.CheckErr(err)
		if dirty {
			cmdr.Warning.Println(fleek.Trans("fleek.dirty"))
		}
		ahead, behind, err = repo.AheadBehind()
		cobra.CheckErr(err)
		if ahead {
			cmdr.Warning.Println(fleek.Trans("fleek.ahead"))
		}
		if behind {
			cmdr.Warning.Println(fleek.Trans("fleek.behind"))
		}
		if cmd.Flag("sync").Changed && ahead {
			cmdr.Info.Println(fleek.Trans("fleek.push"))
			err = repo.Push()
			cobra.CheckErr(err)

		}

	}
	return root
}
