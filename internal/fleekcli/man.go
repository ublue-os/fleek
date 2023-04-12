package fleekcli

import (
	"fmt"
	"os"

	mcoral "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func ManCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "man",
		Short:                 "Generates Fleek's command line manpages",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			pterm.DisableColor()
			manPage, err := mcoral.NewManPage(1, cmd.Root())
			if err != nil {
				return err
			}

			_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
			return err
		},
	}

	return cmd
}
