package fin

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/build"
)

// generateTitleString generates a pretty looking title string.
func generateTitleString(rootCmd *cobra.Command) string {
	return pterm.Sprintf("\n# %s | %s\n", pterm.Cyan(rootCmd.Name()), pterm.Green("v"+build.Version))
}
