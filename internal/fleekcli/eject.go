/*
Copyright © 2023 Brian Ketelsen <bketelsen@gmail.com>
*/
package fleekcli

import (
	"fmt"
	"os"
	"path/filepath"

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

	fl.Config.Aliases["fleek-apply"] = "nix run"

	err = fl.Config.Save()
	if err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	// reload config so it won't git push
	err = fl.ReadConfig(filepath.Join(home, fl.Config.FlakeDir))
	if err != nil {
		return err
	}
	err = fl.Write("fleek: eject")
	if err != nil {
		return err
	}
	// TODO app trans
	fin.Info.Println(app.Trans("generate.runFlake"))

	fmt.Printf("nix run")

	fin.Warning.Println(app.Trans("eject.complete"))
	return nil
}
