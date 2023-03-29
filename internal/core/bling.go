package core

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

type Bling struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Packages    []string `yaml:"packages"`
	Programs    []string `yaml:"programs"`
	PackageMap  map[string]*Package
	ProgramMap  map[string]*Program
}

var (
	//go:embed none.yml
	none []byte
	//go:embed low.yml
	low []byte
	//go:embed default.yml
	dflt []byte
	//go:embed high.yml
	high []byte
)

func loadBling(bytes []byte) (*Bling, error) {

	var b Bling

	err := yaml.Unmarshal(bytes, &b)
	if err != nil {
		return &b, nil
	}
	progs, err := LoadPrograms()
	if err != nil {
		return &b, err
	}
	pkgs, err := LoadPackages()
	if err != nil {
		return &b, err
	}
	b.PackageMap = make(map[string]*Package, len(pkgs))
	b.ProgramMap = make(map[string]*Program, len(progs))
	for _, pkg := range pkgs {
		b.PackageMap[pkg.Name] = pkg
	}
	for _, prog := range progs {
		b.ProgramMap[prog.Name] = prog
	}
	return &b, nil
}

func NoBling() (*Bling, error) {

	return loadBling(none)
}
func LowBling() (*Bling, error) {
	return loadBling(low)
}
func DefaultBling() (*Bling, error) {
	return loadBling(dflt)
}
func HighBling() (*Bling, error) {
	return loadBling(high)
}
