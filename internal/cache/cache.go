package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ublue-os/fleek/internal/core"
	"github.com/ublue-os/fleek/internal/debug"
	"github.com/ublue-os/fleek/internal/xdg"
)

const nixbin = "nix"

type PackageCache struct {
	location string
	Packages PackageList
}

type PackageList map[string]Package

type Package struct {
	Description string `json:"description"`
	Name        string `json:"pname"`
	Version     string `json:"version"`
}

var cacheName = "packages.json"

func New() (*PackageCache, error) {
	cacheDir := xdg.CacheSubpath("fleek")
	debug.Log("package cache: %s", cacheDir)

	pc := &PackageCache{
		location: cacheDir,
	}
	if _, err := os.Stat(cacheDir); errors.Is(err, fs.ErrNotExist) {
		err := core.MkdirAll(cacheDir)
		if err != nil {
			return pc, err
		}
	}
	if pc.valid() {
		debug.Log("package list exists")
		// read it
		bb, err := os.ReadFile(pc.cacheFile())
		if err != nil {
			return pc, err
		}
		var plist PackageList
		debug.Log("unmarshal package list")
		err = json.Unmarshal(bb, &plist)
		if err != nil {
			return pc, err
		}
		pc.Packages = plist

	} else {
		err := pc.Update()
		if err != nil {
			return pc, err
		}
	}
	return pc, nil
}

func (pc *PackageCache) valid() bool {
	_, err := os.Stat(pc.cacheFile())
	return !errors.Is(err, fs.ErrNotExist)

}
func (pc *PackageCache) cacheFile() string {
	return filepath.Join(pc.location, cacheName)
}
func (pc *PackageCache) Update() error {
	debug.Log("updating package list")
	// get it
	bb, err := pc.packageIndex()
	if err != nil {
		return err
	}
	debug.Log("writing cache file: %s", pc.cacheFile())

	err = os.WriteFile(pc.cacheFile(), bb, 0755)
	if err != nil {
		return err
	}
	var plist PackageList
	debug.Log("unmarshal package list")
	err = json.Unmarshal(bb, &plist)
	if err != nil {
		return err
	}
	pc.Packages = plist
	return nil
}
func (pc *PackageCache) runNix(cmd string, cmdLine []string) ([]byte, error) {
	command := exec.Command(cmd, cmdLine...)
	command.Stdin = os.Stdin
	command.Env = os.Environ()

	return command.Output()

}

func (pc *PackageCache) packageIndex() ([]byte, error) {
	// nix search nixpkgs --json
	indexCmdLine := []string{"search", "nixpkgs", "--json"}
	out, err := pc.runNix(nixbin, indexCmdLine)
	if err != nil {
		return out, fmt.Errorf("nix search: %w", err)
	}

	return out, nil
}
