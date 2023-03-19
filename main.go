package main

import (
	"embed"

	"github.com/ublue-os/fleek/cmd"
	"github.com/vanilla-os/orchid/cmdr"
)

var (
	Version = "42"
)

//go:embed locales/*.yml
var fs embed.FS
var fleek *cmdr.App

func main() {
	fleek = cmd.New(Version, fs)

	root := cmd.NewRootCommand(Version)
	fleek.CreateRootCommand(root)

	apply := cmd.NewApplyCommand()
	root.AddCommand(apply)

	init := cmd.NewInitCommand()
	root.AddCommand(init)

	eject := cmd.NewEjectCommand()
	root.AddCommand(eject)

	// run the app
	err := fleek.Run()
	if err != nil {
		cmdr.Error.Println(err)
	}
}
