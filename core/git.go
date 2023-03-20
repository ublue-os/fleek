package core

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/cache"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

func Worktree() (*git.Worktree, error) {
	flake, err := FlakeLocation()
	if err != nil {
		return nil, err
	}
	r, err := git.PlainOpen(flake)
	if err != nil {
		return nil, fmt.Errorf("opening repository: %s", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, fmt.Errorf("opening worktree: %s", err)
	}
	return w, err
}

func Commit() error {
	w, err := Worktree()
	if err != nil {
		return fmt.Errorf("unable to open git worktree")
	}
	err = w.AddGlob("*")
	if err != nil {
		return fmt.Errorf("add glob: %s", err)

	}

	sys, err := CurrentSystem()
	if err != nil {
		return fmt.Errorf("can't commit without system config: %s", err)
	}

	_, err = w.Commit("fleek: update configs", &git.CommitOptions{
		Author: &object.Signature{
			Name:  sys.GitConfig.Name,
			Email: sys.GitConfig.Email,
			When:  time.Now(),
		},
	})
	if err != nil {
		return fmt.Errorf("commit repository: %s", err)

	}
	return nil
}

func CreateRepo() error {
	floc, err := FlakeLocation()
	if err != nil {
		return err
	}
	dotGit := filepath.Join(floc, ".git")
	store := osfs.New(dotGit)
	_, err = git.Init(filesystem.NewStorage(store, cache.NewObjectLRUDefault()), store)
	if err != nil {
		return err
	}
	gitIgnore, err := os.Create(filepath.Join(floc, ".gitignore"))
	if err != nil {
		return err
	}
	defer gitIgnore.Close()
	_, err = gitIgnore.WriteString("result")

	return err
}
func Push() error {
	flake, err := FlakeLocation()
	if err != nil {
		return err
	}
	r, err := git.PlainOpen(flake)
	if err != nil {
		return fmt.Errorf("opening repository: %s", err)
	}
	return r.Push(&git.PushOptions{})

}

func Dirty() (bool, error) {
	flake, err := FlakeLocation()
	if err != nil {
		return false, err
	}
	r, err := git.PlainOpen(flake)
	if err != nil {
		return false, fmt.Errorf("opening repository: %s", err)
	}

	w, err := r.Worktree()
	if err != nil {
		return false, fmt.Errorf("unable to open git worktree")
	}

	status, err := w.Status()
	if err != nil {
		return false, fmt.Errorf("status: %s", err)

	}
	for f, s := range status {
		fmt.Println("file: ", f)
		fmt.Println("staging: ", s.Staging)
		fmt.Println("worktree: ", s.Worktree)
		fmt.Println("extra: ", s.Extra)
	}
	if len(status) > 0 {
		return true, nil
	}

	return false, nil
}

func RemoteAdd(remote string, name string) error {
	flake, err := FlakeLocation()
	if err != nil {
		return err
	}
	r, err := git.PlainOpen(flake)
	if err != nil {
		return fmt.Errorf("opening repository: %s", err)
	}
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: name,
		URLs: []string{remote},
	})
	return err
}

func Remote() (string, error) {
	flake, err := FlakeLocation()
	if err != nil {
		return "", err
	}
	r, err := git.PlainOpen(flake)
	if err != nil {
		return "", fmt.Errorf("opening repository: %s", err)
	}

	list, err := r.Remotes()
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
