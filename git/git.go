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

func (fr *FlakeRepo) Commit() error {

	addCmdline := []string{"add", "--all"}
	_, err := fr.runGit(gitbin, addCmdline)
	if err != nil {
		return fmt.Errorf("git add: %s", err)
	}
	commitCmdLine := []string{"commit", "-m", "fleek: commit"}
	out, err := fr.runGit(gitbin, commitCmdLine)
	outStr := string(out)
	if strings.Contains(outStr, "working tree clean") {
		return nil
	}
	return err

}

func (fr *FlakeRepo) Pull() error {
	pullCmdline := []string{"pull", "origin", "main"}
	_, err := fr.runGit(gitbin, pullCmdline)
	if err != nil {
		return fmt.Errorf("git add: %s", err)
	}
	return nil
}

func (fr *FlakeRepo) CreateRepo() error {
	initCmdLine := []string{"init"}
	_, err := fr.runGit(gitbin, initCmdLine)
	if err != nil {
		return fmt.Errorf("git init: %s", err)
	}
	gitIgnore, err := os.Create(filepath.Join(fr.RootDir, ".gitignore"))
	if err != nil {
		return err
	}
	defer gitIgnore.Close()
	_, err = gitIgnore.WriteString("result")

	return err
}
func (fr *FlakeRepo) Push() error {
	pushCmdline := []string{"push", "origin", "main"}
	_, err := fr.runGit(gitbin, pushCmdline)
	if err != nil {
		return fmt.Errorf("git push: %s", err)
	}
	return nil

}

func (fr *FlakeRepo) Dirty() (bool, error) {

	var dirty bool

	cmd := exec.Command(gitbin, "status", "--porcelain")
	cmd.Dir = fr.RootDir
	cmd.Env = os.Environ()
	out, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("git status: %s", err)
	}

	outString := string(out)

	if len(outString) > 0 {
		lines := strings.Split(outString, "\n")
		for _, line := range lines {
			cleanLine := strings.TrimSpace(line)
			if cleanLine != "" {
				parts := strings.Split(cleanLine, " ")
				var local string
				var remote string
				if len(parts[0]) == 1 {
					local = parts[0]
				}
				if len(parts[0]) == 2 {
					remote = parts[0][:1]
				}
				fmt.Printf("git status: %s\n", parts[1])

				fmt.Printf("\tlocal: %s\n", local)
				if len(remote) > 0 {
					fmt.Printf("\tremote: %s\n", remote)
				}
				dirty = true
			}
		}

	}

	return dirty, nil
}
func (fr *FlakeRepo) AheadBehind() (bool, bool, error) {

	var ahead bool
	var behind bool

	cmd := exec.Command(gitbin, "status", "--ahead-behind")
	cmd.Env = os.Environ()
	cmd.Dir = fr.RootDir
	out, err := cmd.Output()
	if err != nil {
		return false, false, fmt.Errorf("git status: %s", err)
	}

	outString := string(out)
	cleanOut := strings.TrimSpace(outString)

	if len(cleanOut) > 0 {
		if strings.Contains(cleanOut, "ahead") {
			ahead = true
		}
	}
	fetch := exec.Command(gitbin, "fetch", "origin", "main")
	fetch.Dir = fr.RootDir

	fetch.Env = os.Environ()

	err = fetch.Run()
	if err != nil {
		return false, false, fmt.Errorf("git fetch: %s", err)
	}
	pfcmd := exec.Command(gitbin, "status", "--ahead-behind")
	pfcmd.Dir = fr.RootDir
	pfcmd.Env = os.Environ()

	postFetchout, err := pfcmd.Output()
	if err != nil {
		return false, false, fmt.Errorf("git status: %s", err)
	}

	foutString := string(postFetchout)
	cleanFetchOut := strings.TrimSpace(foutString)

	if len(cleanFetchOut) > 0 {
		if strings.Contains(cleanOut, "behind") {
			behind = true
		}
	}
	return ahead, behind, nil
}

func (fr *FlakeRepo) RemoteAdd(remote string, name string) error {
	var err error
	err = fr.open()
	if err != nil {
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
		return "", fmt.Errorf("opening repository: %s", err)
	}

	list, err := fr.repo.Remotes()
	if err != nil {
		return "", fmt.Errorf("getting remotes	: %s", err)
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
