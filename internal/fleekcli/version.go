package fleekcli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/build"
	"github.com/ublue-os/fleek/internal/vercheck"
)

type versionFlags struct {
	verbose bool
}

func VersionCmd() *cobra.Command {
	flags := versionFlags{}
	command := &cobra.Command{
		Use:   app.Trans("version.use"),
		Short: app.Trans("version.short"),
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return versionCmdFunc(cmd, args, flags)
		},
	}

	command.Flags().BoolVarP(&flags.verbose, app.Trans("version.flagVerbose"), "v", false, // value
		app.Trans("version.flagVerboseDescription"),
	)

	return command
}

func versionCmdFunc(cmd *cobra.Command, _ []string, flags versionFlags) error {
	w := cmd.OutOrStdout()
	v := getVersionInfo()
	if flags.verbose {
		fmt.Fprintf(w, app.Trans("version.version"), v.Version)

		fmt.Fprintf(w, app.Trans("version.platform"), v.Platform)
		fmt.Fprintf(w, app.Trans("version.commit"), v.Commit)
		fmt.Fprintf(w, app.Trans("version.time"), v.CommitDate)
		fmt.Fprintf(w, app.Trans("version.go"), v.GoVersion)
	} else {
		fmt.Fprintf(w, "%v\n", v.Version)
	}
	return nil
}

type versionInfo struct {
	Version      string
	IsPrerelease bool
	Platform     string
	Commit       string
	CommitDate   string
	GoVersion    string
}

func getVersionInfo() *versionInfo {
	v := &versionInfo{
		Version:    build.Version,
		Platform:   fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH),
		Commit:     build.Commit,
		CommitDate: build.CommitDate,
		GoVersion:  runtime.Version(),
	}

	return v
}
