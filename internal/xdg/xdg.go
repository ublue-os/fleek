package xdg

import (
	"os"
	"path/filepath"
)

func DataSubpath(subpath string) string {
	return filepath.Join(dataDir(), subpath)
}
func DataSubpathRel(subpath string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}
	dp := filepath.Join(dataDir(), subpath)
	rp, _ := filepath.Rel(home, dp)
	return rp
}

func ConfigSubpath(subpath string) string {
	return filepath.Join(configDir(), subpath)
}
func ConfigSubpathRel(subpath string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}
	cp := filepath.Join(configDir(), subpath)
	rp, _ := filepath.Rel(home, cp)
	return rp
}

func CacheSubpath(subpath string) string {
	return filepath.Join(cacheDir(), subpath)
}

func CacheSubpathRel(subpath string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}
	cp := filepath.Join(cacheDir(), subpath)
	rp, _ := filepath.Rel(home, cp)
	return rp
}
func StateSubpath(subpath string) string {
	return filepath.Join(stateDir(), subpath)
}
func StateSubpathRel(subpath string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}
	sp := filepath.Join(stateDir(), subpath)
	rp, _ := filepath.Rel(home, sp)
	return rp
}

func dataDir() string   { return resolveDir("XDG_DATA_HOME", ".local/share") }
func configDir() string { return resolveDir("XDG_CONFIG_HOME", ".config") }
func cacheDir() string  { return resolveDir("XDG_CACHE_HOME", ".cache") }
func stateDir() string  { return resolveDir("XDG_STATE_HOME", ".local/state") }

func resolveDir(envvar string, defaultPath string) string {
	dir := os.Getenv(envvar)
	if dir != "" {
		return dir
	}

	home, err := os.UserHomeDir()
	if err != nil {
		home = "~"
	}

	return filepath.Join(home, defaultPath)
}
