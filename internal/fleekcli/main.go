package fleekcli

import (
	"context"
	"os"

	"github.com/ublue-os/fleek"
	"github.com/ublue-os/fleek/internal/debug"
	"github.com/ublue-os/fleek/internal/midcobra"
)

var (
	debugMiddleware *midcobra.DebugMiddleware = &midcobra.DebugMiddleware{}
	traceMiddleware *midcobra.TraceMiddleware = &midcobra.TraceMiddleware{}
	app             *fleek.App
)

func Main() {
	app = fleek.NewApp()

	code := Execute(context.Background(), os.Args[1:])
	os.Exit(code)
}
func Execute(ctx context.Context, args []string) int {
	defer debug.Recover()
	exe := midcobra.New(RootCmd())
	exe.AddMiddleware(traceMiddleware)
	exe.AddMiddleware(debugMiddleware)
	return exe.Execute(ctx, args)
}
