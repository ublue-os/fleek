package flake

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"

	"github.com/go-git/go-git/v5"
	"github.com/ublue-os/fleek/internal/debug"
	fgit "github.com/ublue-os/fleek/internal/git"
	"github.com/ublue-os/fleek/internal/ux"
)

const gitbin = "git"

func (f *Flake) gitOpen() (*git.Repository, error) {

	return git.PlainOpen(f.Config.UserFlakeDir())

}

func (f *Flake) Clone(repo string) error {
	cloneCmdline := []string{"clone", repo, f.Config.UserFlakeDir()}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	command := exec.Command(gitbin, cloneCmdline...)
	command.Stdin = os.Stdin
	command.Dir = home
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = os.Environ()
	err = command.Run()
	if err != nil {
		return fmt.Errorf("git clone: %w", err)
	}
	return nil
}

func (f *Flake) runGit(cmd string, cmdLine []string) error {
	command := exec.Command(cmd, cmdLine...)
	command.Stdin = os.Stdin
	command.Dir = f.Config.UserFlakeDir()
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = os.Environ()
	return command.Run()
}

func (f *Flake) IsGitRepo() (bool, error) {

	loc, err := f.Config.GitLocation()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(loc)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func (f *Flake) mayCommit() error {
	git, err := f.IsGitRepo()
	if err != nil {
		ux.Debug.Printfln("git repo error: %s", err)
		return err
	}
	if git {
		ux.Debug.Println("is git repo")
		// add
		ux.Debug.Println("git will add")
		if f.Config.Verbose {
			ux.Verbose.Println(f.app.Trans("git.add"))
		}
		err = f.add()
		if err != nil {
			ux.Debug.Printfln("git add error: %s", err)
			return err
		}
		if f.Config.Git.AutoCommit {
			ux.Debug.Println("git will commit")
			if f.Config.Verbose {
				ux.Verbose.Println(f.app.Trans("git.commit"))
			}
			err = f.commit()
			if err != nil {
				ux.Debug.Printfln("git commit error: %s", err)
				return err
			}
			// commit
			err = f.commit()
			if err != nil {
				ux.Debug.Printfln("git commit error: %s", err)
				return err
			}
		}

		if f.Config.Git.AutoPush {
			ux.Debug.Println("git will push")
			if f.Config.Verbose {
				ux.Verbose.Println(f.app.Trans("git.push"))
			}
			err = f.push()
			if err != nil {
				ux.Debug.Printfln("git push error: %s", err)
				return err
			}
		}
	} else {
		ux.Debug.Println("skipping git")
		return nil
	}

	return nil
}
func (f *Flake) MayPull() error {
	git, err := f.IsGitRepo()
	if err != nil {
		ux.Debug.Printfln("git repo error: %s", err)
		return err
	}
	if git {
		if f.Config.Verbose {
			ux.Verbose.Println(f.app.Trans("git.commit"))
		}
		ux.Debug.Println("is git repo")

		if f.Config.Git.AutoPull {
			ux.Debug.Println("git will pull")
			if f.Config.Verbose {
				ux.Verbose.Println(f.app.Trans("git.pull"))
			}
			err = f.pull()
			if err != nil {
				ux.Debug.Printfln("git pull error: %s", err)
				return err
			}
		}

	} else {
		ux.Debug.Println("skipping git")
		return nil
	}

	return nil
}

func (f *Flake) add() error {
	addCmdline := []string{"add", "--all"}
	err := f.runGit(gitbin, addCmdline)
	if err != nil {
		return fmt.Errorf("git add: %w", err)
	}
	return nil
}

func (f *Flake) commit() error {
	status, err := f.gitStatus()
	if err != nil {
		return err
	}
	if status.Empty() {
		ux.Debug.Println("git status is empty, skipping commit")
		return nil
	}

	commitCmdLine := []string{"commit", "-m", "fleek: commit"}
	err = f.runGit(gitbin, commitCmdLine)
	if err != nil {
		return fmt.Errorf("git commit: %w", err)
	}
	return err
}

func (f *Flake) pull() error {
	remote, err := f.remote()
	if err != nil {
		return err
	}
	// if no remote, no need to push
	if remote == "" {
		return err
	}
	// totally stole --autostash --rebase from chezmoi, thanks twpayne
	pullCmdline := []string{"pull", "--autostash", "--rebase", "origin", "main"}
	err = f.runGit(gitbin, pullCmdline)
	if err != nil {
		return fmt.Errorf("git pull: %w", err)
	}
	return nil
}

func (f *Flake) setRebase() error {

	configCmdLine := []string{"config", "pull.rebase", "true"}
	err := f.runGit(gitbin, configCmdLine)
	if err != nil {
		return fmt.Errorf("git config: %w", err)
	}
	return err
}
func (f *Flake) push() error {
	remote, err := f.remote()
	if err != nil {
		return err
	}
	// if no remote, no need to push
	if remote == "" {
		return nil
	}
	pushCmdline := []string{"push", "origin", "main"}
	err = f.runGit(gitbin, pushCmdline)
	if err != nil {
		return fmt.Errorf("git push: %w", err)
	}
	return nil

}

func (f *Flake) gitStatus() (*fgit.Status, error) {
	// git status --ignored --porcelain=v2
	cmd := exec.Command(gitbin, "status", "--ignored", "--porcelain=v2")
	cmd.Dir = f.Config.UserFlakeDir()
	cmd.Env = os.Environ()
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return fgit.ParseStatusPorcelainV2(out)

}

func (f *Flake) remote() (string, error) {
	var err error
	repo, err := f.gitOpen()
	if err != nil {
		debug.Log("fr open: %s", err)

		return "", fmt.Errorf("opening repository: %w", err)
	}

	list, err := repo.Remotes()
	if err != nil {
		return "", fmt.Errorf("getting remotes	: %w", err)
	}
	var urls string
	for _, r := range list {
		for _, upstream := range r.Config().URLs {
			urls = urls + upstream + "\n"

		}

	}
	return urls, nil
}
