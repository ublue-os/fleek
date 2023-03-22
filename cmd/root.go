package cmd

import (
	"embed"
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
var flakeLocation string

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
		cobra.CheckErr(err)
		flakeLocation, err = core.FlakeLocation()
		flake, err = nix.NewFlake(flakeLocation, config)
		cobra.CheckErr(err)
		repo = git.NewFlakeRepo(flakeLocation)

	}
	return root
}
