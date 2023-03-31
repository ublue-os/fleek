package core

import (
	_ "embed"

	"gopkg.in/yaml.v3"
)

var (
	//go:embed programs.yml
	programs []byte
)

func LoadPrograms() ([]*Program, error) {
	var pp []*Program
	err := yaml.Unmarshal(programs, &pp)
	if err != nil {
		return pp, err
	}
	return pp, nil
}

// Program is an application that is installed
// and also has configuration attached.
type Program struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	ConfigLines []ConfigLine `yaml:"config_lines,flow"`
	Aliases     []Alias      `yaml:",flow"`
}

// ConfigLine represents a line of a program's
// configuration. It must be in the form of
// `key = value“ such that multiple config lines for
// the same program can be appended together to create
// a valid configuration.
type ConfigLine struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Alias struct {
	Key         string `yaml:"key"`
	Value       string `yaml:"value"`
	Description string `yaml:"description"`
}
