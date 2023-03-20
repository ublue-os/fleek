package core

import (
	"os"
	"runtime"

	"os/user"
)

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
