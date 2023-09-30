package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/cmdutil"
	"github.com/ublue-os/fleek/internal/fleek"
	"github.com/ublue-os/fleek/internal/xdg"
)

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
type SearchResult struct {
	Name    string  `json:"name"`
	Package Package `json:"package"`
}

var cacheName = "packages.json"

func New() (*PackageCache, error) {
	cacheDir := xdg.CacheSubpath("fleek")
	fin.Logger.Debug("package cache", fin.Logger.Args("dir", cacheDir))

	pc := &PackageCache{
		location: cacheDir,
	}
	if _, err := os.Stat(cacheDir); errors.Is(err, fs.ErrNotExist) {
		err := fleek.MkdirAll(cacheDir)
		if err != nil {
			return pc, err
		}
	}
	if pc.valid() {
		fin.Logger.Debug("package list exists")
		// read it
		bb, err := os.ReadFile(pc.cacheFile())
		if err != nil {
			return pc, err
		}
		var plist PackageList
		fin.Logger.Debug("unmarshal package list")
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
	fin.Logger.Debug("updating package list")
	// get it
	bb, err := pc.packageIndex()
	if err != nil {
		return err
	}
	fin.Logger.Debug("writing cache file", fin.Logger.Args("file", pc.cacheFile()))

	err = os.WriteFile(pc.cacheFile(), bb, 0755)
	if err != nil {
		return err
	}
	var plist PackageList
	fin.Logger.Debug("unmarshal package list")
	err = json.Unmarshal(bb, &plist)
	if err != nil {
		return err
	}
	pc.Packages = plist
	return nil
}

func (pc *PackageCache) packageIndex() ([]byte, error) {
	args := []string{"search", "nixpkgs", "--json"}
	cmd, buf := cmdutil.CommandTTYWithBufferNoOut("nix", args...)
	cmd.Env = os.Environ()
	// nix search nixpkgs --json
	err := cmd.Run()
	if err != nil {
		return buf.Bytes(), fmt.Errorf("nix search: %w", err)
	}

	return buf.Bytes(), nil
}
