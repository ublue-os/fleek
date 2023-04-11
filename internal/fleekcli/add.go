package fleekcli

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/flake"
)

type addCmdFlags struct {
	apply bool
}

func AddCommand() *cobra.Command {
	flags := addCmdFlags{}
	command := &cobra.Command{
		Use:     app.Trans("add.use"),
		Short:   app.Trans("add.short"),
		Long:    app.Trans("add.long"),
		Args:    cobra.MinimumNArgs(1),
		Example: app.Trans("add.example"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return add(cmd, args)
		},
	}
	command.Flags().BoolVarP(
		&flags.apply, app.Trans("add.applyFlag"), "a", false, app.Trans("add.applyFlagDescription"))

	return command
}

// initCmd represents the init command
func add(cmd *cobra.Command, args []string) error {

	fin.Description.Println(cmd.Short)
	err := mustConfig()
	if err != nil {
		return err
	}
	var apply bool
	if cmd.Flag(app.Trans("add.applyFlag")).Changed {
		apply = true
	}

	fl, err := flake.Load(cfg, app)
	if err != nil {
		return err
	}
	err = fl.MayPull()
	if err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString("add packages: ")
	for _, p := range args {
		fin.Info.Println(app.Trans("add.adding") + p)
		err = fl.Config.AddPackage(p)
		if err != nil {
			fin.Debug.Printfln("add package error: %s", err)
			return err
		}
		sb.WriteString(p + " ")

	}
	err = fl.Write(false, sb.String())
	if err != nil {
		fin.Debug.Printfln("flake write error: %s", err)
		return err
	}

	if apply {
		fin.Info.Println(app.Trans("add.applying"))

		err = fl.Apply()
		if err != nil {
			if errors.Is(err, flake.ErrPackageConflict) {
				fin.Fatal.Println(app.Trans("global.errConflict"))
			}
			return err
		}
	} else {
		fin.Warning.Println(app.Trans("add.unapplied"))
	}
	fin.Success.Println(app.Trans("global.completed"))
	return nil
}
