/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package fleekcli

import (
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/flake"
)

func UpdateCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   app.Trans("update.use"),
		Short: app.Trans("update.short"),
		Long:  app.Trans("update.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return update(cmd)
		},
	}

	return command
}

// initCmd represents the init command
func update(cmd *cobra.Command) error {
	fin.Description.Println(cmd.Short)
	err := mustConfig()
	if err != nil {
		return err
	}
	fl, err := flake.Load(cfg, app)
	if err != nil {
		return err
	}
	err = fl.MayPull()
	if err != nil {
		return err
	}
	// Write the templates from the latest version of fleek
	// to get any possible changes to the templates
	if err := fl.WriteTemplates(); err != nil {
		return err
	}

	// update the nix flake lock
	if err := fl.Update(); err != nil {
		return err
	}
	// We just updated the flake lock, which might pull a new
	// version of fleek or other deps in. Update the system templates to
	// get new fixes without having to update/apply twice
	if err := fl.WriteTemplates(); err != nil {
		return err
	}
	if err := fl.Apply(); err != nil {
		return err
	}

	fin.Success.Println(app.Trans("update.done"))
	return nil
}
