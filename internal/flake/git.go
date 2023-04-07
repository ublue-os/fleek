package flake

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/ublue-os/fleek/internal/debug"
	"github.com/ublue-os/fleek/internal/ux"
)

const gitbin = "git"

func (f *Flake) gitOpen() (*git.Repository, error) {

	return git.PlainOpen(f.Config.UserFlakeDir())

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
		ux.Info.Println(f.app.Trans("git.commit"))
		ux.Debug.Println("is git repo")
		// commit
		err = f.commit()
		if err != nil {
			ux.Debug.Printfln("git commit error: %s", err)
			return err
		}
		// add
		if f.Config.Git.AutoAdd {
			ux.Debug.Println("git will add")
			ux.Info.Println(f.app.Trans("git.add"))
			err = f.add()
			if err != nil {
				ux.Debug.Printfln("git add error: %s", err)
				return err
			}
		}
		if f.Config.Git.AutoPush {
			ux.Debug.Println("git will push")
			ux.Info.Println(f.app.Trans("git.push"))
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
		ux.Info.Println(f.app.Trans("git.commit"))
		ux.Debug.Println("is git repo")

		if f.Config.Git.AutoPull {
			ux.Debug.Println("git will pull")
			ux.Info.Println(f.app.Trans("git.pull"))
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

	commitCmdLine := []string{"commit", "-m", "fleek: commit"}
	err := f.runGit(gitbin, commitCmdLine)
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
	pullCmdline := []string{"pull", "origin", "main"}
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

func (f *Flake) Dirty() (bool, []byte, error) {

	var dirty bool

	cmd := exec.Command(gitbin, "status", "--porcelain")
	cmd.Dir = f.Config.UserFlakeDir()
	cmd.Env = os.Environ()
	out, err := cmd.Output()
	if err != nil {
		return false, out, fmt.Errorf("git status: %w", err)
	}

	outString := string(out)

	if len(outString) > 0 {
		lines := strings.Split(outString, "\n")
		for _, line := range lines {
			cleanLine := strings.TrimSpace(line)
			if cleanLine != "" {
				dirty = true
			}
		}

	}

	return dirty, out, nil
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

/*
from git-scm.com git-status docs

' ' = unmodified

M = modified

T = file type changed (regular file, symbolic link or submodule)

A = added

D = deleted

R = renamed

C = copied (if config option status.renames is set to "copies")

U = updated but unmerged

X          Y     Meaning
-------------------------------------------------
	 [AMD]   not updated
M        [ MTD]  updated in index
T        [ MTD]  type changed in index
A        [ MTD]  added to index
D                deleted from index
R        [ MTD]  renamed in index
C        [ MTD]  copied in index
[MTARC]          index and work tree matches
[ MTARC]    M    work tree changed since index
[ MTARC]    T    type changed in work tree since index
[ MTARC]    D    deleted in work tree
	    R    renamed in work tree
	    C    copied in work tree
-------------------------------------------------
D           D    unmerged, both deleted
A           U    unmerged, added by us
U           D    unmerged, deleted by them
U           A    unmerged, added by them
D           U    unmerged, deleted by us
A           A    unmerged, both added
U           U    unmerged, both modified
-------------------------------------------------
?           ?    untracked
!           !    ignored
-------------------------------------------------
*/
