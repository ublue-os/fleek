/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/flake"
)

func EjectCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   app.Trans("eject.use"),
		Short: app.Trans("eject.short"),
		Long:  app.Trans("eject.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return eject(cmd)
		},
	}

	return command
}

// initCmd represents the init command
func eject(cmd *cobra.Command) error {
	err := mustConfig()
	if err != nil {
		return err
	}
	fin.Description.Println(cmd.Short)

	fl, err := flake.Load(cfg, app)
	if err != nil {
		return err
	}
	fl.Config.Ejected = true
	fl.Config.Git.AutoCommit = false
	fl.Config.Git.AutoPull = false
	fl.Config.Git.AutoPush = false
	fl.Config.Git.Enabled = false
	for _, system := range fl.Config.Systems {
		// nix run --impure home-manager/master -- -b bak switch --flake .#bjk@ghanima
		fl.Config.Aliases["apply-"+system.Hostname] = fmt.Sprintf("nix run --impure home-manager/master -- -b bak switch --flake .#%s@%s", system.Username, system.Hostname)
		//fin.Info.Printfln("nix run --impure home-manager/master -- -b bak switch --flake .#%s@%s", system.Username, system.Hostname)
	}
	err = fl.Config.Save()
	if err != nil {
		return err
	}
	// reload config so it won't git push
	err = fl.ReadConfig()
	if err != nil {
		return err
	}
	err = fl.Write(true, "fleek: eject")
	if err != nil {
		return err
	}
	// TODO app trans
	fin.Info.Println(app.Trans("generate.runFlake"))

	for _, system := range fl.Config.Systems {
		// nix run --impure home-manager/master -- -b bak switch --flake .#bjk@ghanima
		fmt.Printf("nix run --impure home-manager/master -- -b bak switch --flake .#%s@%s\n", system.Username, system.Hostname)
		//fin.Info.Printfln("nix run --impure home-manager/master -- -b bak switch --flake .#%s@%s", system.Username, system.Hostname)
	}

	fin.Warning.Println(app.Trans("eject.complete"))
	return nil
}
