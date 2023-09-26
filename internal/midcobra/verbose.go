// Copyright 2022 Jetpack Technologies Inc and contributors. All rights reserved.
// Use of this source code is governed by the license in the LICENSE file.

package midcobra

import (
	"errors"
	"os"
	"os/exec"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/debug"
	"github.com/ublue-os/fleek/internal/fleekcli/usererr"
	"github.com/ublue-os/fleek/internal/ux"
	"github.com/ublue-os/fleek/internal/verbose"
)

type VerboseMiddleware struct {
	executionID string // uuid
	flag        *pflag.Flag
}

var _ Middleware = (*VerboseMiddleware)(nil)

func (d *VerboseMiddleware) AttachToFlag(flags *pflag.FlagSet, flagName string) {
	flags.Bool(
		flagName,
		false,
		"Show full stack traces on errors",
	)
	d.flag = flags.Lookup(flagName)
	d.flag.Hidden = true
}

func (d *VerboseMiddleware) preRun(_ *cobra.Command, _ []string) {
	if d == nil {
		return
	}

	strVal := ""
	if d.flag.Changed {
		strVal = d.flag.Value.String()
	} else {
		strVal = os.Getenv("FLEEK_VERBOSE")
	}
	if enabled, _ := strconv.ParseBool(strVal); enabled {
		verbose.Enable()
	}
}

func (d *VerboseMiddleware) postRun(cmd *cobra.Command, _ []string, runErr error) {
	if runErr == nil {
		return
	}
	if usererr.HasUserMessage(runErr) {
		if usererr.IsWarning(runErr) {
			ux.Fwarning(cmd.ErrOrStderr(), runErr.Error())
			return
		}
		color.New(color.FgRed).Fprintf(cmd.ErrOrStderr(), "\nError: %s\n\n", runErr.Error())
	} else {
		color.New(color.FgRed).Fprintf(cmd.ErrOrStderr(), "Error: %v\n\n", runErr)
	}

	st := debug.EarliestStackTrace(runErr)
	var exitErr *exec.ExitError
	if errors.As(runErr, &exitErr) {
		fin.Debug.Printfln("Command stderr: %s\n", exitErr.Stderr)
	}
	fin.Debug.Printfln("\nExecutionID:%s\n%+v\n", d.executionID, st)
}

func (d *VerboseMiddleware) withExecutionID(execID string) Middleware {
	d.executionID = execID
	return d
}
