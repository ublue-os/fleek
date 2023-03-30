package fleekcli

import (
	"errors"
	"fmt"
	"path/filepath"

	"io/fs"
	"os"

	"github.com/spf13/cobra"
	"github.com/ublue-os/fleek/internal/core"
	"github.com/ublue-os/fleek/internal/debug"
	"github.com/ublue-os/fleek/internal/git"
	"github.com/ublue-os/fleek/internal/nix"
	"github.com/ublue-os/fleek/internal/ux"
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
	return [...]string{"None", "Exists", "Dirty", "Behind", "Ahead", "Diverged"}[f]
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
type Fleek struct {
	flake         *nix.Flake
	config        *core.Config
	repo          *git.FlakeRepo
	configStatus  ConfigStatus
	flakeStatus   FlakeStatus
	gitStatus     GitStatus
	flakeLocation string
}

func initFleek(verbose bool) (*Fleek, error) {
	f := &Fleek{}
	var err error

	// set up config

	config, err := core.ReadConfig()
	if err != nil {
		debug.Log("read config : %s", err)
		if errors.Is(err, os.ErrNotExist) {
			debug.Log("no config found")
			if verbose {
				ux.Info.Println(app.Trans("fleek.noConfigFound"))
			}
			f.configStatus = ConfigNone
		} else {
			return f, err
		}
	} else {
		debug.Log("config found")
		f.configStatus = ConfigExists
	}
	f.flakeStatus = FlakeNone
	f.gitStatus = GitNone
	debug.Log("config status: %s", f.configStatus.String())
	if f.configStatus != ConfigNone {
		debug.Log("config isn't none")

		if config != nil {
			debug.Log("validating config")
			err = config.Validate()
			if err != nil {
				debug.Log("config validation error: %s", err)
				if errors.Is(err, core.ErrMissingFlakeDir) {
					ux.Info.Println(app.Trans("fleek.migrating"))
					// get previous default flake location

					defaultFlakeDir := filepath.Join(".config", "home-manager")
					config.FlakeDir = defaultFlakeDir
					// now save the config
					debug.Log("saving fixed config")
					err2 := config.Save()
					if err2 != nil {
						debug.Log("error saving fixed config: %s", err)
						return f, err2
					}
					ux.Success.Println(app.Trans("fleek.migrated"))

				} else {
					debug.Log("error validating config: %s", err)
					return f, err
				}
			}
			f.config = config
			f.flakeLocation = f.config.UserFlakeDir()
			debug.Log("flake location: %s", f.flakeLocation)

			// setup flake
			flake, err := nix.NewFlake(f.flakeLocation, f.config)
			if err != nil {
				debug.Log("new flake error: %s", err)
				return f, err
			}
			debug.Log("flake struct created")
			exists, err := flake.Exists()
			if err != nil {
				debug.Log("flake exists error: %s", err)
				if errors.Is(err, fs.ErrNotExist) {
					debug.Log("no flake directory found at %s", f.flakeLocation)
					if verbose {
						ux.Info.Println(app.Trans("fleek.noFlakeFound"))
					}
					debug.Log("setting flakeStatus: %s", FlakeNone.String())
					f.flakeStatus = FlakeNone
				}
			}
			debug.Log("flake exists: %v", exists)
			if exists {
				debug.Log("setting flakeStatus: %s", FlakeExists.String())
				f.flakeStatus = FlakeExists
			}
			f.flake = flake

			// setup repo
			f.repo = git.NewFlakeRepo(f.flakeLocation)
			exists = f.repo.IsValid()
			debug.Log("repo exists: %v", exists)
			if exists {
				if verbose {
					ux.Info.Println(app.Trans("fleek.validRepository"))
				}
				f.gitStatus = GitExists
			} else {
				if verbose {
					ux.Info.Println(app.Trans("fleek.notValidRepository"))
				}
			}

			dirty, out, err := f.repo.Dirty(false)
			cobra.CheckErr(err)
			if dirty {
				f.flakeStatus = FlakeDirty
			}
			if verbose {
				ux.Info.Println(string(out))
			}
			remote, err := f.repo.Remote()
			debug.Log("remote found: %s", remote)
			cobra.CheckErr(err)
			if remote != "" {
				f.gitStatus = GitHasRemote
				if verbose {
					ux.Info.Println(app.Trans("fleek.remoteFound"), remote)
				}
			} else {
				if verbose {
					ux.Info.Println(app.Trans("fleek.noRemote"))
				}
			}

			// TODO: refactor this out to reuse
			if f.gitStatus == GitHasRemote {
				if verbose {
					ux.Info.Println(app.Trans("fleek.gettingStatus"), remote)
				}
				debug.Log("getting remote status")
				ahead, behind, _, err := f.repo.AheadBehind(verbose)
				if err != nil {
					debug.Log("remote status error: %v", err)
					return f, err
				}
				debug.Log("ahead: %v", ahead)
				debug.Log("behind: %v", behind)

				if ahead {
					debug.Log("remote status: %s", FlakeAhead.String())

					if verbose {
						ux.Info.Println(app.Trans("fleek.aheadStatus"))
					}
					f.flakeStatus = FlakeAhead
				}
				if behind {
					debug.Log("remote status: %s", FlakeBehind.String())
					if verbose {
						ux.Info.Println(app.Trans("fleek.behindStatus"))
					}
					f.flakeStatus = FlakeBehind
				}
				if ahead && behind {
					debug.Log("remote status: %s", FlakeDiverged.String())
					if verbose {
						ux.Info.Println(app.Trans("fleek.divergedStatus"))
					}
					f.flakeStatus = FlakeDiverged
				}

			} else {
				debug.Log("no remote configured")

				if verbose {
					ux.Info.Println(app.Trans("fleek.skippingStatus"), remote)
				}
			}
		}
	}

	return f, nil
}

func (f *Fleek) Flake() (*nix.Flake, error) {
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
func (f *Fleek) Repo() (*git.FlakeRepo, error) {
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

func (f *Fleek) Sanity() ([]string, error) {
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
