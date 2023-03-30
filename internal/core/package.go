package core

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

// Package is an binary or library that
// can be installed
type Package struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

var (
	//go:embed packages.yml
	packages []byte
)

func LoadPackages() ([]*Package, error) {
	var pp []*Package
	err := yaml.Unmarshal(packages, &pp)
	if err != nil {
		return pp, nil
	}
	return pp, nil
}
