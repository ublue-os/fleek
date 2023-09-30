package fleekcli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/build"
	"github.com/ublue-os/fleek/internal/envir"
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

	command.AddCommand(selfUpdateCmd())
	return command
}

func selfUpdateCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "update",
		Short: "Update fleek launcher and binary",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return vercheck.SelfUpdate(cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}

	return command
}
func versionCmdFunc(cmd *cobra.Command, _ []string, flags versionFlags) error {
	w := cmd.OutOrStdout()
	v := getVersionInfo()
	lv, err := latestVersion()
	if err != nil {
		fin.Logger.Warn("Unable to check latest version", fin.Logger.Args("error", err))
	}
	if flags.verbose {
		fmt.Fprintf(w, app.Trans("version.version"), v.Version)

		fmt.Fprintf(w, app.Trans("version.platform"), v.Platform)
		fmt.Fprintf(w, app.Trans("version.commit"), v.Commit)
		fmt.Fprintf(w, app.Trans("version.time"), v.CommitDate)
		fmt.Fprintf(w, app.Trans("version.go"), v.GoVersion)
		fmt.Fprintf(w, "Launcher:    %v\n", v.LauncherVersion)
		fmt.Fprintf(w, "Latest Version:    %v\n", lv)
		fmt.Fprintf(w, "Upgrade available: %v\n", isNewFleekAvailable(v.Version, lv))

	} else {
		fmt.Fprintf(w, "%v\n", v.Version)
	}
	return nil
}

type versionInfo struct {
	Version         string
	IsPrerelease    bool
	Platform        string
	Commit          string
	CommitDate      string
	GoVersion       string
	LauncherVersion string
}

func getVersionInfo() *versionInfo {
	v := &versionInfo{
		Version:         build.Version,
		Platform:        fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH),
		Commit:          build.Commit,
		CommitDate:      build.CommitDate,
		GoVersion:       runtime.Version(),
		LauncherVersion: os.Getenv(envir.LauncherVersion),
	}

	return v
}

func latestVersion() (string, error) {
	res, err := http.Get("https://releases.getfleek.dev/fleek/stable/version")
	if err != nil {
		return "unknown", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "unknown", err
	}
	sb := string(body)
	return strings.TrimSpace(sb), nil
}

// isNewFleekAvailable returns true if a new fleek CLI binary version is available.
func isNewFleekAvailable(current, latest string) bool {
	if latest == "" {
		return false
	}
	if strings.Contains(current, "0.0.0-dev") {
		return false
	}
	return vercheck.SemverCompare(current, latest) < 0
}
