package core

import (
	"fmt"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
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
