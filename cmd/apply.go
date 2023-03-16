/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ublue/fleek/core"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply nix home-manager result based on your ~/.fleek.yml configuration",
	Long:  `Apply nix home-manager result based on your ~/.fleek.yml configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Writing Configuration")
		err := core.WriteFlake()
		cobra.CheckErr(err)
		fmt.Println("Compiling Configuration")
		err = core.CheckFlake()
		cobra.CheckErr(err)
		var dry bool
		if cmd.Flag("dry-run").Changed {
			dry = true
		}
		if !dry {
			fmt.Println("Applying Configuration")
			err = core.ApplyFlake()
			cobra.CheckErr(err)
		} else {
			fmt.Println("Dry Run, Not applying Configuration")
		}
		fmt.Println("Done.")
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	applyCmd.Flags().BoolP("dry-run", "d", false, "dry run - compile but don't apply")
}
