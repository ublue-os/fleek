package fleekcli

import (
	"context"
	"os"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/debug"
	"github.com/ublue-os/fleek/internal/midcobra"
)

var (
	debugMiddleware *midcobra.DebugMiddleware = &midcobra.DebugMiddleware{}
	traceMiddleware *midcobra.TraceMiddleware = &midcobra.TraceMiddleware{}
	app             *fleek.App
	root            *cobra.Command
)

func Main() {

	code := Execute(context.Background(), os.Args[1:])
	os.Exit(code)
}
func Execute(ctx context.Context, args []string) int {
	defer debug.Recover()

	exe := midcobra.New(root)
	exe.AddMiddleware(traceMiddleware)
	exe.AddMiddleware(debugMiddleware)
	return exe.Execute(ctx, args)
}
func init() {
	app = fleek.NewApp()
	root = RootCmd()
	// Use https://github.com/pterm/pcli to style the output of cobra.
	fin.SetRepo("ublue-os/fleek")
	fin.SetRootCmd(root)
	fin.Setup()

	// Change global PTerm theme
	pterm.ThemeDefault.SectionStyle = *pterm.NewStyle(pterm.FgCyan)
}
