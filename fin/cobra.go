package fin

import (
	"io"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ublue-os/fleek/internal/build"
)

// HelpFunc is a drop in replacement for spf13/cobra `HelpFunc`
func HelpFunc() func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, strings []string) {
		var ret string

		ret += generateTitleString(rootCmd)
		ret += generateUsageTemplate(cmd)
		ret += generateDescriptionTemplate(cmd.Long)
		ret += generateExamplesTemplate(cmd)
		ret += generateCommandsTemplate(cmd.Commands())
		ret += generateFlagsTemplate(cmd.Flags())
		ret += "\n"

		pterm.Print(ret)
	}
}

// FlagErrorFunc is a drop in replacement for spf13/cobra `FlagErrorFunc`
func FlagErrorFunc() func(*cobra.Command, error) error {
	return func(cmd *cobra.Command, err error) error {
		var ret string

		ret += generateTitleString(rootCmd)
		ret += generateUsageTemplate(cmd)
		ret += generateDescriptionTemplate(cmd.Long)
		ret += generateCommandsTemplate(cmd.Commands())
		ret += generateFlagsTemplate(cmd.Flags())
		ret += "\n\n"
		ret += pterm.Error.WithShowLineNumber(false).Sprintln(err)

		pterm.Print(ret)
		return nil
	}
}

// GlobalNormalizationFunc is a drop in replacement for spf13/cobra `GlobalNormalizationFunc`
func GlobalNormalizationFunc() func(f *pflag.FlagSet, name string) pflag.NormalizedName {
	return rootCmd.GlobalNormalizationFunc()
}

// HelpTemplate is a drop in replacement for spf13/cobra `HelpTemplate`
func HelpTemplate() string {
	return ""
}

// UsageFunc is a drop in replacement for spf13/cobra `UsageFunc`
func UsageFunc() func(*cobra.Command) error {
	return rootCmd.UsageFunc()
}

// UsageTemplate is a drop in replacement for spf13/cobra `UsageTemplate`
func UsageTemplate() string {
	return ""
}

// VersionTemplate is a drop in replacement for spf13/cobra `VersionTemplate`
func VersionTemplate() string {
	return pterm.Info.Sprintfln("%s is on version: %s", rootCmd.Name(), pterm.Magenta(build.Version))
}

// Err is a drop in replacement for spf13/cobra `Err`
func Err() io.Writer {
	return errorWriter{}
}

type errorWriter struct{}

func (errorWriter) Write(p []byte) (n int, err error) {
	pterm.Error.WithMessageStyle(pterm.NewStyle()).WithShowLineNumber(false).Println(string(p))
	return len(p), nil
}

// PcliOut is a drop in replacement for spf13/cobra `SetOut`
func PcliOut() io.Writer {
	return pcliOut{}
}

type pcliOut struct{}

func (pcliOut) Write(p []byte) (n int, err error) {
	str := string(p)
	if strings.Contains(str, " is on version: ") {
		pterm.Print(str)
	}
	return len(p), nil
}
