/*
Copyright © 2023 Brian Ketelsen <bketelsen@gmail.com>
*/
package fleekcli

import (
	"errors"
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/flake"
	"github.com/ublue-os/fleek/internal/fleek"
	"gopkg.in/yaml.v3"
)

func JoinCommand() *cobra.Command {
	command := &cobra.Command{
		Use:     app.Trans("join.use"),
		Short:   app.Trans("join.short"),
		Long:    app.Trans("join.long"),
		Example: app.Trans("join.example"),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return join(cmd, args)
		},
	}
	return command
}

// joinCmd represents the join command
func join(cmd *cobra.Command, args []string) error {

	var verbose bool
	if cmd.Flag(app.Trans("fleek.verboseFlag")).Changed {
		verbose = true
	}

	cfg.Verbose = verbose

	fin.Description.Println(cmd.Short)

	dirName, err := flake.CloneRepository(args[0])
	if err != nil {
		return err
	}

	// read config
	config, err := readConfigFromGitClone(dirName)
	if err != nil {
		return err
	}

	_, err = os.Stat(config.UserFlakeDir())
	if err == nil {
		// exists
		return errors.New("target configuration directory already exists")
	}
	// move cloned repo
	err = cp.Copy(dirName, config.UserFlakeDir())
	if err != nil {
		return err
	}
	// reload config and flake
	config, err = fleek.ReadConfig(config.UserFlakeDir())
	if err != nil {
		return err
	}
	migrate := config.NeedsMigration()
	if migrate {
		fin.Logger.Info("Migration required")
		err := config.Migrate()
		if err != nil {
			fin.Logger.Error("migrating host files", fin.Logger.Args("error", err))
			os.Exit(1)
		}
		fl, err := flake.Load(config, app)
		if err != nil {
			fin.Logger.Error("loading flake", fin.Logger.Args("error", err))
			os.Exit(1)
		}

		// Symlink the yaml file to home
		cfile, err := fl.Config.Location()
		if err != nil {
			fin.Logger.Error("config location", fin.Logger.Args("error", err))
			return err
		}
		fin.Logger.Debug("config", fin.Logger.Args("file", cfile))

		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		csym := filepath.Join(home, ".fleek.yml")
		err = os.Symlink(cfile, csym)
		if err != nil {
			fin.Logger.Debug("symlink  failed")
			return err
		}
		err = fl.Write("update host and user files", true, false)
		if err != nil {
			fin.Logger.Error("writing flake", fin.Logger.Args("error", err))
			os.Exit(1)
		}
	}

	fin.Logger.Info(app.Trans("init.joining"))
	fl, err := flake.Load(config, app)
	if err != nil {
		return err
	}
	err = fl.Join()
	if err != nil {
		return err
	}
	// reload config and flake
	config, err = fleek.ReadConfig(config.UserFlakeDir())
	if err != nil {
		return err
	}
	fl, err = flake.Load(config, app)
	if err != nil {
		return err
	}
	err = fl.Write("join new system", true, true)
	if err != nil {
		fin.Logger.Error("flake write", fin.Logger.Args("error", err))
		return err
	}

	err = fl.Apply()
	if err != nil {
		return err
	}
	fin.Logger.Info(app.Trans("join.complete"))

	return nil
}

func readConfigFromGitClone(loc string) (*fleek.Config, error) {
	c := &fleek.Config{}
	loc = filepath.Join(loc, ".fleek.yml")
	bb, err := os.ReadFile(loc)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(bb, c)
	if err != nil {
		return c, err
	}
	return c, nil
}
