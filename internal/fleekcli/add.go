package fleekcli

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/debug"
	"github.com/ublue-os/fleek/internal/flake"
	"github.com/ublue-os/fleek/internal/ux"
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

	ux.Description.Println(cmd.Short)
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
	var sb strings.Builder
	sb.WriteString("add packages: ")
	for _, p := range args {
		ux.Info.Println(app.Trans("add.adding") + p)
		err = fl.Config.AddPackage(p)
		if err != nil {
			debug.Log("add package error: %s", err)
			return err
		}
		sb.WriteString(p + " ")

	}
	err = fl.Write(false)
	if err != nil {
		debug.Log("flake write error: %s", err)
		return err
	}

	if apply {
		ux.Info.Println(app.Trans("add.applying"))

		err = fl.Apply()
		if err != nil {
			if errors.Is(err, flake.ErrPackageConflict) {
				ux.Fatal.Println(app.Trans("global.errConflict"))
			}
			return err
		}
	} else {
		ux.Warning.Println(app.Trans("add.unapplied"))
	}
	ux.Success.Println(app.Trans("global.completed"))
	return nil
}
