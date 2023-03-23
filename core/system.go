package core

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"os/user"
)

var ErrSysNotFound = errors.New("system not found")

func Runtime() string {
	arch := runtime.GOARCH
	os := runtime.GOOS
	var nixarch string
	switch arch {
	case "amd64":
		nixarch = "x86_64"
	case "arm64":
		nixarch = "aarch64"
	}
	return nixarch + "-" + os
}
func Arch() string {
	arch := runtime.GOARCH
	var nixarch string
	switch arch {
	case "amd64":
		nixarch = "x86_64"
	case "arm64":
		nixarch = "aarch64"
	}
	return nixarch
}

func Username() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return u.Username, nil
}

func Hostname() (string, error) {
	h, e := os.Hostname()
	if e != nil {
		return "", e
	}
	return h, nil
}

func CurrentSystem() (*System, error) {
	conf, err := ReadConfig()
	if err != nil {
		return nil, fmt.Errorf("reading config: %s", err)
	}
	host, err := Hostname()
	if err != nil {
		return nil, fmt.Errorf("getting hostname: %s", err)
	}
	user, err := Username()
	if err != nil {
		return nil, fmt.Errorf("getting username: %s", err)
	}
	for _, sys := range conf.Systems {
		if sys.Hostname == host {
			if sys.Username == user {
				return &sys, nil
			}
		}
	}
	return nil, ErrSysNotFound
}
