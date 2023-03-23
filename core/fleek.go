package core

import (
	"os"
	"path/filepath"
)

// ConfigLocation returns the path for the
// fleek configuration file.
func (c *Config) Location() (string, error) {
	return filepath.Join(c.FlakeDir, ".fleek.yml"), nil
}

// MakeFlakeDir creates the directory that holds
// the interpolated flake.
func (c *Config) MakeFlakeDir() error {

	return os.MkdirAll(c.FlakeDir, 0755)
}
