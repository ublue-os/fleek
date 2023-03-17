package main

import (
	"embed"

	"github.com/ublue-os/fleek/cmd"
	"github.com/vanilla-os/orchid/cmdr"
)

var (
	Version = "1.7.0-1"
)

//go:embed locales/*.yml
var fs embed.FS
var fleek *cmdr.App

func main() {
	fleek = cmd.New(Version, fs)

	// root command
	root := cmd.NewRootCommand(Version)
	// root command
	fleek.CreateRootCommand(root)

	apply := cmd.NewApplyCommand()
	root.AddCommand(apply)

	init := cmd.NewInitCommand()
	root.AddCommand(init)
	/*
		enter := cmd.NewEnterCommand()
		root.AddCommand(cmd.AddContainerFlags(enter))

		export := cmd.NewExportCommand()
		root.AddCommand(cmd.AddContainerFlags(export))

		initialize := cmd.NewInitializeCommand()
		root.AddCommand(cmd.AddContainerFlags(initialize))

		install := cmd.NewInstallCommand()
		root.AddCommand(cmd.AddContainerFlags(install))

		list := cmd.NewListCommand()
		root.AddCommand(cmd.AddContainerFlags(list))

		purge := cmd.NewPurgeCommand()
		root.AddCommand(cmd.AddContainerFlags(purge))

		remove := cmd.NewRemoveCommand()
		root.AddCommand(cmd.AddContainerFlags(remove))

		run := cmd.NewRunCommand()
		root.AddCommand(cmd.AddContainerFlags(run))

		search := cmd.NewSearchCommand()
		root.AddCommand(cmd.AddContainerFlags(search))

		show := cmd.NewShowCommand()
		root.AddCommand(cmd.AddContainerFlags(show))

		unexport := cmd.NewUnexportCommand()
		root.AddCommand(cmd.AddContainerFlags(unexport))

		upgrade := cmd.NewUpgradeCommand()
		root.AddCommand(cmd.AddContainerFlags(upgrade))

		update := cmd.NewUpdateCommand()
		root.AddCommand(cmd.AddContainerFlags(update))
		// run the app
	*/
	err := fleek.Run()
	if err != nil {
		cmdr.Error.Println(err)
	}
}
