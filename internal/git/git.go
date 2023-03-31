package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/ublue-os/fleek/internal/debug"
)

const gitbin = "git"

type FlakeRepo struct {
	RootDir string
	repo    *git.Repository
}

func (fr *FlakeRepo) open() error {
	var err error
	fr.repo, err = git.PlainOpen(fr.RootDir)
	return err

}
func NewFlakeRepo(root string) *FlakeRepo {
	frepo := &FlakeRepo{}
	frepo.RootDir = root

	return frepo
}

func (fr *FlakeRepo) IsValid() bool {
	var err error
	_, err = git.PlainOpen(fr.RootDir)
	if err != nil {
		if errors.Is(err, git.ErrRepositoryNotExists) {
			return false
		}
	}
	return true
}

func (fr *FlakeRepo) runGit(cmd string, cmdLine []string) ([]byte, error) {
	command := exec.Command(cmd, cmdLine...)
	command.Stdin = os.Stdin
	command.Dir = fr.RootDir
	command.Env = os.Environ()

	return command.Output()

}

func (fr *FlakeRepo) Commit() ([]byte, error) {

	addCmdline := []string{"add", "--all"}
	out, err := fr.runGit(gitbin, addCmdline)
	if err != nil {
		return out, fmt.Errorf("git add: %w", err)
	}
	commitCmdLine := []string{"commit", "-m", "fleek: commit"}
	commitOut, err := fr.runGit(gitbin, commitCmdLine)
	outStr := string(commitOut)
	if strings.Contains(outStr, "working tree clean") {
		return append(out, commitOut...), nil
	}
	return append(out, commitOut...), err

}

func (fr *FlakeRepo) Pull() ([]byte, error) {
	remote, err := fr.Remote()
	if err != nil {
		return []byte{}, err
	}
	// if no remote, no need to push
	if remote == "" {
		return []byte{}, err
	}
	pullCmdline := []string{"pull", "origin", "main"}
	out, err := fr.runGit(gitbin, pullCmdline)
	if err != nil {
		return out, fmt.Errorf("git pull: %w", err)
	}
	return out, nil
}

func (fr *FlakeRepo) LocalConfig(user, email string) ([]byte, error) {
	userCmdline := []string{"config", "user.name", user}
	out, err := fr.runGit(gitbin, userCmdline)
	if err != nil {
		return out, fmt.Errorf("git config: %w", err)
	}
	emailCmdline := []string{"config", "user.email", email}
	out2, err := fr.runGit(gitbin, emailCmdline)
	if err != nil {
		return append(out, out2...), fmt.Errorf("git config: %w", err)
	}
	return append(out, out2...), nil
}

func (fr *FlakeRepo) CreateRepo() ([]byte, error) {
	var output []byte
	initCmdLine := []string{"init"}
	out, err := fr.runGit(gitbin, initCmdLine)
	output = append(output, out...)
	if err != nil {
		return output, fmt.Errorf("git init: %w", err)
	}
	gitIgnore, err := os.Create(filepath.Join(fr.RootDir, ".gitignore"))
	if err != nil {
		return output, err
	}
	defer gitIgnore.Close()
	_, err = gitIgnore.WriteString("result")
	if err != nil {
		return output, err
	}
	configCmdLine := []string{"config", "pull.rebase", "true"}
	out, err = fr.runGit(gitbin, configCmdLine)
	output = append(output, out...)
	if err != nil {
		return output, fmt.Errorf("git config: %w", err)
	}
	return output, err
}
func (fr *FlakeRepo) SetRebase() ([]byte, error) {
	var output []byte

	configCmdLine := []string{"config", "pull.rebase", "true"}
	out, err := fr.runGit(gitbin, configCmdLine)
	output = append(output, out...)
	if err != nil {
		return output, fmt.Errorf("git config: %w", err)
	}
	return output, err
}
func (fr *FlakeRepo) Push() ([]byte, error) {
	remote, err := fr.Remote()
	if err != nil {
		return []byte{}, err
	}
	// if no remote, no need to push
	if remote == "" {
		return []byte{}, nil
	}
	pushCmdline := []string{"push", "origin", "main"}
	out, err := fr.runGit(gitbin, pushCmdline)
	if err != nil {
		return out, fmt.Errorf("git push: %w", err)
	}
	return out, nil

}

func (fr *FlakeRepo) Dirty() (bool, []byte, error) {

	var dirty bool

	cmd := exec.Command(gitbin, "status", "--porcelain")
	cmd.Dir = fr.RootDir
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
func (fr *FlakeRepo) AheadBehind(verbose bool) (bool, bool, []byte, error) {
	var output []byte
	remote, err := fr.Remote()
	if err != nil {
		debug.Log("getting remote: %s", err)
		return false, false, []byte{}, err
	}
	// if no remote, not ahead or behind
	if remote == "" {
		if verbose {
			fmt.Println("flake repo: no remote")
		}
		debug.Log("no remote")
		return false, false, []byte{}, nil
	}
	var ahead bool
	var behind bool
	fetch := exec.Command(gitbin, "fetch", "origin", "main")
	fetch.Dir = fr.RootDir

	fetch.Env = os.Environ()

	out, err := fetch.Output()
	output = append(output, out...)
	if err != nil {
		debug.Log("git fetch: %s", err)
		return false, false, output, fmt.Errorf("git fetch: %w", err)
	}
	cmd := exec.Command(gitbin, "status", "--ahead-behind")
	cmd.Env = os.Environ()
	cmd.Dir = fr.RootDir
	out, err = cmd.Output()
	output = append(output, out...)

	if err != nil {
		debug.Log("git status: %s", err)
		return false, false, output, fmt.Errorf("git status: %w", err)
	}

	outString := string(out)
	cleanOut := strings.TrimSpace(outString)
	if verbose {
		fmt.Println(cleanOut)
	}
	if len(cleanOut) > 0 {
		if strings.Contains(cleanOut, "ahead") {
			ahead = true
		}
		if strings.Contains(cleanOut, "behind") {
			behind = true
		}
		if strings.Contains(cleanOut, "diverged") {
			behind = true
			ahead = true
		}
	}

	return ahead, behind, output, nil
}

func (fr *FlakeRepo) RemoteAdd(remote string, name string) error {
	var err error
	err = fr.open()
	if err != nil {
		debug.Log("open flake repo: %s", err)
		return err
	}
	_, err = fr.repo.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{remote},
	})
	return err
}

func (fr *FlakeRepo) Remote() (string, error) {
	var err error
	err = fr.open()
	if err != nil {
		debug.Log("fr open: %s", err)

		return "", fmt.Errorf("opening repository: %w", err)
	}

	list, err := fr.repo.Remotes()
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
