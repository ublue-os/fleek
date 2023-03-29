package fleekcli

import (
	"context"
	"os"

	"github.com/ublue-os/fleek"
	"github.com/ublue-os/fleek/internal/debug"
	"github.com/ublue-os/fleek/internal/midcobra"
)

/*
func Main() {

	fleek = cmd.New(Version, fs)

	root := cmd.NewRootCommand(Version)
	fleek.CreateRootCommand(root)

	apply := cmd.NewApplyCommand()
	root.AddCommand(apply)

	init := cmd.NewInitCommand()
	root.AddCommand(init)

	eject := cmd.NewEjectCommand()
	root.AddCommand(eject)

	add := cmd.NewAddCommand()
	root.AddCommand(add)

	remove := cmd.NewRemoveCommand()
	root.AddCommand(remove)
	repo := cmd.NewRepoCommand()
	root.AddCommand(repo)
	reposhow := cmd.NewRepoShowCommand()
	repo.AddCommand(reposhow)
	repoadd := cmd.NewRepoAddCommand()
	repo.AddCommand(repoadd)

	list := cmd.NewListCommand()
	listPkgs := cmd.NewListPackagesCommand()
	listProgs := cmd.NewListProgramsCommand()

	list.AddCommand(listPkgs)
	list.AddCommand(listProgs)
	root.AddCommand(list)

	update := cmd.NewUpdateCommand()
	root.AddCommand(update)

	show := cmd.NewShowCommand()
	root.AddCommand(show)
	// run the app
	err := fleek.Run()

		if err != nil {
			cmdr.Error.Println(err)
		}

}
*/

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
