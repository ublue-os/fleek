package fleek

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"os/user"

	"github.com/ublue-os/fleek/fin"
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

func Name() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return u.Name, nil
}

func Hostname() (string, error) {
	override := os.Getenv("FLEEK_HOST_OVERRIDE")
	if override != "" {
		return override, nil
	}
	h, e := os.Hostname()
	if e != nil {
		return "", e
	}
	return h, nil
}

func (c *Config) CurrentSystem() (*System, error) {

	host, err := Hostname()
	if err != nil {
		return nil, fmt.Errorf("getting hostname: %w", err)
	}
	user, err := Username()
	if err != nil {
		return nil, fmt.Errorf("getting username: %w", err)
	}
	for _, sys := range c.Systems {
		if sys.Hostname == host {
			if sys.Username == user {
				return sys, nil
			}
		}
	}
	return nil, ErrSysNotFound
}

func UserShell() (string, error) {
	// modified from https://github.com/captainsafia/go-user-shell/blob/master/user_shell.go
	// MIT License
	// Copyright (c) 2017 Safia Abdalla
	var shell string
	switch runtime.GOOS {
	case "windows":
		if os.Getenv("COMSPEC") != "" {
			shell = os.Getenv("COMSPEC")
		} else {
			shell = "/cmd.exe"
		}
	case "darwin":
		if os.Getenv("SHELL") != "" {
			shell = os.Getenv("SHELL")
		} else {
			shell = "/bin/zsh"
		}
	default:
		if os.Getenv("SHELL") != "" {
			shell = os.Getenv("SHELL")
		} else {
			shell = "/bin/sh"
		}
	}
	if strings.Contains(shell, "zsh") {
		shell = "zsh"
	}
	if strings.Contains(shell, "bash") {
		shell = "bash"
	}
	fin.Logger.Debug("shell", fin.Logger.Args("detected", shell))
	return shell, nil

}

func MkdirAll(path string) error {
	return os.Mkdir(path, 0755)
}
