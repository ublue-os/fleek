package flake

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/cmdutil"
	fgit "github.com/ublue-os/fleek/internal/git"
)

const gitbin = "git"

func (f *Flake) gitOpen() (*git.Repository, error) {

	return git.PlainOpen(f.Config.UserFlakeDir())

}

func (f *Flake) Clone(repo string) error {
	if f.Config.Verbose {
		fin.Verbose.Printfln("Cloning %s to %s", repo, f.Config.UserFlakeDir())
	}
	cloneCmdline := []string{"clone", repo, f.Config.UserFlakeDir()}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	cmd := cmdutil.CommandTTY(gitbin, cloneCmdline...)
	cmd.Dir = home
	cmd.Env = os.Environ()
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("git clone: %w", err)
	}
	return nil
}

func (f *Flake) runGit(cmd string, cmdLine []string) error {
	command := cmdutil.CommandTTY(cmd, cmdLine...)
	command.Dir = f.Config.UserFlakeDir()
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
func (f *Flake) mayCommit(message string) error {
	git, err := f.IsGitRepo()
	if err != nil {
		fin.Logger.Error("git repo", fin.Logger.Args("error", err))
		return err
	}
	if git {
		fin.Logger.Debug("is git repo")
		// add
		fin.Logger.Debug("git will add")
		fin.Logger.Info(f.app.Trans("git.add"))

		err = f.add()
		if err != nil {
			fin.Logger.Error("git add", fin.Logger.Args("error", err))
			return err
		}
		if f.Config.Git.AutoCommit {
			fin.Logger.Debug("git will commit")
			fin.Logger.Info(f.app.Trans("git.commit"))

			err = f.commit(message)
			if err != nil {
				fin.Logger.Error("git commit", fin.Logger.Args("error", err))
				return err
			}

		}

		if f.Config.Git.AutoPush {
			fin.Logger.Debug("git will push")
			fin.Logger.Info(f.app.Trans("git.push"))
			err = f.push()
			if err != nil {
				fin.Logger.Error("git push", fin.Logger.Args("error", err))
				return err
			}
		}
	} else {
		fin.Logger.Info("skipping git")
		return nil
	}

	return nil
}
func (f *Flake) MayPull() error {
	git, err := f.IsGitRepo()
	if err != nil {
		fin.Logger.Error("check repo", fin.Logger.Args("error", err))
		return err
	}
	if git {
		fin.Logger.Info(f.app.Trans("git.commit"))
		fin.Logger.Debug("is git repo")

		if f.Config.Git.AutoPull {
			fin.Logger.Debug("git will pull")
			fin.Logger.Info(f.app.Trans("git.pull"))

			err = f.pull()
			if err != nil {
				fin.Logger.Error("git pull", fin.Logger.Args("error", err))
				return err
			}
		}

	} else {
		fin.Logger.Info("skipping git")
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

func (f *Flake) commit(message string) error {
	status, err := f.gitStatus()
	if err != nil {
		return errors.New("error parsing git status")
	}
	if status.Empty() {
		fin.Logger.Debug("git status is empty, skipping commit")
		return nil
	}
	if message == "" {
		message = "fleek: commit"
	}

	commitCmdLine := []string{"commit", "-m", message}
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
	cmd, buff := cmdutil.CommandTTYWithBuffer(gitbin, "status", "--ignored", "--porcelain=v2")
	cmd.Dir = f.Config.UserFlakeDir()
	cmd.Env = os.Environ()
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return fgit.ParseStatusPorcelainV2(buff.Bytes())

}

func (f *Flake) remote() (string, error) {
	var err error
	repo, err := f.gitOpen()
	if err != nil {
		fin.Logger.Error("flake git repo open", fin.Logger.Args("error", err))
		return "", fmt.Errorf("opening repository: %w", err)
	}

	list, err := repo.Remotes()
	if err != nil {
		fin.Logger.Error("flake git repo remotes", fin.Logger.Args("error", err))
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
func CloneRepository(repo string) (string, error) {

	dirname, err := os.MkdirTemp("", "fleek*")
	if err != nil {
		return "", err
	}
	cloneCmdline := []string{"clone", repo, dirname}
	command := cmdutil.CommandTTY(gitbin, cloneCmdline...)

	command.Env = os.Environ()
	err = command.Run()
	if err != nil {
		return "", fmt.Errorf("git clone: %w", err)
	}
	return dirname, nil
}
