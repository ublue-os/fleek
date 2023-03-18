package core

import (
	"os"
	"path/filepath"
)

// ConfigLocation returns the path for the
// fleek configuration file.
func ConfigLocation() (string, error) {
	hm, err := FlakeLocation()
	if err != nil {
		return "", err
	}
	return filepath.Join(hm, ".fleek.yml"), nil
}

// FlakeLocation returns the path where the
// interpolated flake will be created.
func FlakeLocation() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "home-manager"), nil
}

// MakeFlakeDir creates the directory that holds
// the interpolated flake.
func MakeFlakeDir() error {
	f, err := FlakeLocation()
	if err != nil {
		return err
	}
	return os.MkdirAll(f, 0755)
}
