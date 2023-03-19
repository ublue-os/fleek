package cmd

import (
	"embed"
	"os"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/core"
	"github.com/vanilla-os/orchid/cmdr"
)

var fleek *cmdr.App

const (
	verboseFlag string = "verbose"
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
				false))

	root.Version = version
	root.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		ok := core.CheckNix()
		if !ok {
			cmdr.Error.Println(fleek.Trans("fleek.installNix"))
			os.Exit(1)
		}
	}
	return root
}
