package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"io/fs"
	"os"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/core"
	"github.com/ublue-os/fleek/git"
	"github.com/ublue-os/fleek/nix"
	"github.com/vanilla-os/orchid/cmdr"
)

type FlakeStatus int

const (
	FlakeNone FlakeStatus = iota
	FlakeExists
	FlakeDirty
	FlakeBehind
	FlakeAhead
	FlakeDiverged
)

func (f FlakeStatus) String() string {
	return [...]string{"None", "Dirty", "Behind", "Ahead", "Diverged"}[f]
}

type ConfigStatus int

const (
	ConfigNone ConfigStatus = iota
	ConfigExists
)

func (c ConfigStatus) String() string {
	return [...]string{"None", "Exists"}[c]
}

type GitStatus int

const (
	GitNone GitStatus = iota
	GitExists
	GitHasRemote
)

func (g GitStatus) String() string {
	return [...]string{"Repository Not Initialized", "Repository Exists", "Repository Has Remote"}[g]
}

// Fleek is the controller for the command
// line experience and holds state for all
// the commands.
type fleek struct {
	flake         *nix.Flake
	config        *core.Config
	repo          *git.FlakeRepo
	configStatus  ConfigStatus
	flakeStatus   FlakeStatus
	gitStatus     GitStatus
	flakeLocation string
}

func initFleek(verbose bool) (*fleek, error) {
	f := &fleek{}
	// set up config
	var err error
	config, err := core.ReadConfig()

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if verbose {
				cmdr.Info.Println(app.Trans("fleek.noConfigFound"))
			}
			f.configStatus = ConfigNone
		} else {
			return f, err
		}
	} else {
		f.configStatus = ConfigExists
	}
	f.flakeStatus = FlakeNone
	f.gitStatus = GitNone
	if f.configStatus != ConfigNone {
		if config != nil {
			err = config.Validate()
			if err != nil {
				if errors.Is(err, core.ErrMissingFlakeDir) {
					cmdr.Info.Println("Migrating .fleek.yml to current version")
					// get previous default flake location

					defaultFlakeDir := filepath.Join(".config", "home-manager")
					config.FlakeDir = defaultFlakeDir
					// now save the config
					err2 := config.Save()
					if err2 != nil {
						return f, err2
					}
					cmdr.Success.Println("Migrated .fleek.yml ")

				} else {
					return f, err
				}
			}
			f.config = config
			f.flakeLocation = f.config.UserFlakeDir()

			// setup flake
			flake, err := nix.NewFlake(f.flakeLocation, f.config)
			if err != nil {
				return f, err
			}
			exists, err := flake.Exists()
			if err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					if verbose {
						cmdr.Info.Println(app.Trans("fleek.noFlakeFound"))
					}
					f.flakeStatus = FlakeNone
				}
			}
			if exists {
				f.flakeStatus = FlakeExists
			}
			f.flake = flake

			// setup repo
			f.repo = git.NewFlakeRepo(f.flakeLocation)
			exists = f.repo.IsValid()
			if exists {
				if verbose {
					cmdr.Info.Println(app.Trans("fleek.validRepository"))
				}
				f.gitStatus = GitExists
			} else {
				if verbose {
					cmdr.Info.Println(app.Trans("fleek.notValidRepository"))
				}
			}

			remote, err := f.repo.Remote()
			cobra.CheckErr(err)
			if remote != "" {
				f.gitStatus = GitHasRemote
				if verbose {
					cmdr.Info.Println(app.Trans("fleek.remoteFound"), remote)
				}
			} else {
				if verbose {
					cmdr.Info.Println(app.Trans("fleek.noRemote"))
				}
			}
			if f.gitStatus == GitHasRemote {
				if verbose {
					cmdr.Info.Println(app.Trans("fleek.gettingStatus"), remote)
				}
				ahead, behind, err := f.repo.AheadBehind(verbose)
				if err != nil {
					return f, err
				}
				if ahead {
					if verbose {
						cmdr.Info.Println(app.Trans("fleek.aheadStatus"))
					}
					f.flakeStatus = FlakeAhead
				}
				if behind {
					if verbose {
						cmdr.Info.Println(app.Trans("fleek.behindStatus"))
					}
					f.flakeStatus = FlakeBehind
				}
				if ahead && behind {
					if verbose {
						cmdr.Info.Println(app.Trans("fleek.divergedStatus"))
					}
					f.flakeStatus = FlakeDiverged
				}

			} else {
				if verbose {
					cmdr.Info.Println(app.Trans("fleek.skippingStatus"), remote)
				}
			}
		}
	}

	return f, nil
}

func (f *fleek) Flake() (*nix.Flake, error) {
	if f.flake != nil {
		return f.flake, nil
	}
	// setup flake
	flake, err := nix.NewFlake(f.flakeLocation, f.config)
	if err != nil {
		return nil, err
	}
	exists, err := flake.Exists()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			f.flakeStatus = FlakeNone
		}
	}
	if exists {
		f.flakeStatus = FlakeExists
	}

	f.flake = flake
	return f.flake, nil
}
func (f *fleek) Repo() (*git.FlakeRepo, error) {
	if f.repo != nil {
		return f.repo, nil
	}
	f.repo = git.NewFlakeRepo(f.flakeLocation)
	exists := f.repo.IsValid()
	if exists {
		f.gitStatus = GitExists
	}
	return f.repo, nil
}

func (f *fleek) Sanity() ([]string, error) {
	// check nix
	// check flakes
	msgs := []string{}
	shell, err := core.UserShell()
	if err != nil {
		return []string{}, err
	}
	configuredShell := f.config.Shell
	if shell != configuredShell {
		msgs = append(msgs, "----Configuration Mismatch----")
		msgs = append(msgs, fmt.Sprintf("~/.fleek.yml configured shell is %s, but user configured shell is %s", configuredShell, shell))
		msgs = append(msgs, "Consult your operating system documentation on how to change your login shell")

	}
	return msgs, nil
}
