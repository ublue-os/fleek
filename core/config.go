package core

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the options that will be
// merged into the home-manager flake.
type Config struct {
	Unfree     bool              `yaml:"unfree"`
	Repository string            `yaml:"repo"`
	Name       string            `yaml:"name"`
	Packages   []string          `yaml:",flow"`
	Programs   []string          `yaml:",flow"`
	Aliases    map[string]string `yaml:",flow"`
	Paths      []string          `yaml:"paths"`
	Me         Me                `yaml:"me"`
}
type Me struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

// ReadConfig returns the configuration data
// stored in $HOME/.fleek.yml
func ReadConfig() (*Config, error) {
	c := &Config{}
	cfile, err := ConfigLocation()
	if err != nil {
		return c, err
	}
	bb, err := os.ReadFile(cfile)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(bb, c)
	if err != nil {
		return c, err
	}
	return c, nil
}

// WriteSampleConfig creates the first fleek
// configuration file
func WriteSampleConfig(email, name string, force bool) error {
	aliases := make(map[string]string)
	aliases["cdfleek"] = "cd ~/.config/home-manager"
	c := Config{
		Unfree:     true,
		Name:       "My Fleek Configuration",
		Repository: "git@github.com/my/homeconfig",
		Packages: []string{
			"neovim",
			"fzf",
			"nixfmt",
			"lazygit",
			"ripgrep",
			"jq",
			"dive",
			"htop",
			"yq",
			"vscode",
			"go_1_20",
			"gnumake",
			"gcc",
			"statix",
			"rustup",
			"goreleaser",
		},
		Programs: []string{
			"direnv",
			"starship",
			"atuin",
			"gh",
			"zellij",
		},
		Aliases: aliases,
		Paths: []string{
			"$HOME/bin",
			"$HOME/.local/bin",
		},
		Me: Me{
			Name:  name,
			Email: email,
		},
	}
	cfile, err := ConfigLocation()
	if err != nil {
		return err
	}
	_, err = os.Stat(cfile)

	if force || os.IsNotExist(err) {

		cfg, err := os.Create(cfile)
		if err != nil {
			return err
		}
		bb, err := yaml.Marshal(&c)
		if err != nil {
			return err
		}
		m := make(map[interface{}]interface{})
		err = yaml.Unmarshal(bb, &m)
		if err != nil {
			return err
		}
		n, err := yaml.Marshal(&m)
		if err != nil {
			return err
		}
		// convert to string to get `-` style lists
		sbb := string(n)
		_, err = cfg.WriteString(sbb)
		if err != nil {
			return err
		}
	} else {
		return errors.New("cowardly refusing to overwrite config file without --force flag")
	}
	return nil
}
