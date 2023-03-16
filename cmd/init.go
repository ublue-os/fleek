/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ublue/fleek/core"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize fleek configuration.",
	Long: `Initialize fleek configuration by creating $HOME/.fleek.yml
and a persistence directory for your configs.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Setting up Fleek.")
		var force bool
		if cmd.Flag("force").Changed {
			force = true
		}
		fmt.Println("Checking nix configuration.")

		ok := core.CheckNix()
		if ok {
			fmt.Println("Writing fleek configuration file.")
			err := core.WriteSampleConfig(force)
			cobra.CheckErr(err)
			err = core.MakeFlakeDir()
			cobra.CheckErr(err)
			err = core.InitFlake(force)
			cobra.CheckErr(err)
		} else {
			fmt.Println("Is nix installed?")
		}
		fmt.Println("Done. \n\nEdit ~/.fleek.yml to your taste and run `fleek apply`")

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("force", "f", false, "Overwrite configs if they exist")
}
