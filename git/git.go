package git

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/ublue-os/fleek/core"

	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

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

func (fr *FlakeRepo) Worktree() (*git.Worktree, error) {
	var err error
	err = fr.open()
	if err != nil {
		return nil, fmt.Errorf("opening worktree: %s", err)
	}
	w, err := fr.repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("opening worktree: %s", err)
	}
	return w, err
}

func (fr *FlakeRepo) Commit() error {
	w, err := fr.Worktree()
	if err != nil {
		return fmt.Errorf("unable to open git worktree")
	}
	err = w.AddGlob("*")
	if err != nil {
		return fmt.Errorf("add glob: %s", err)

	}

	sys, err := core.CurrentSystem()
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

func (fr *FlakeRepo) CreateRepo() error {
	var err error
	dotGit := filepath.Join(fr.RootDir, ".git")
	store := osfs.New(dotGit)
	_, err = git.Init(filesystem.NewStorage(store, cache.NewObjectLRUDefault()), store)
	if err != nil {
		return err
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
	var err error
	err = fr.open()
	if err != nil {
		return fmt.Errorf("opening repository: %s", err)
	}
	return fr.repo.Push(&git.PushOptions{})

}

func (fr *FlakeRepo) Dirty() (bool, error) {

	var err error
	err = fr.open()
	if err != nil {
		return false, fmt.Errorf("opening repository: %s", err)
	}
	w, err := fr.Worktree()
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
