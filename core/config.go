package core

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	shells          = []string{"bash", "zsh"}
	blingLevels     = []string{"low", "default", "high"}
	lowPackages     = []string{"htop"}
	defaultPackages = []string{"fzf", "ripgrep", "vscode"}
	highPackages    = []string{"lazygit", "jq", "yq", "neovim", "neofetch", "btop", "cheat"}
	lowPrograms     = []string{"starship"}
	defaultPrograms = []string{"gh", "direnv"}
	highPrograms    = []string{"exa", "bat", "atuin", "zoxide"}
)

// Config holds the options that will be
// merged into the home-manager flake.
type Config struct {
	Unfree bool `yaml:"unfree"`
	// bash or zsh
	Shell string `yaml:"shell"`
	// low, default, high
	Bling      string            `yaml:"bling"`
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

func (c Config) Validate() error {
	if !isValueInList(c.Shell, shells) {
		return errors.New("fleek.yml: invalid shell, valid shells are: " + strings.Join(shells, ", "))
	}
	if !isValueInList(c.Bling, blingLevels) {
		return errors.New("fleek.yml: invalid bling level, valid levels are: " + strings.Join(blingLevels, ", "))
	}
	return nil
}

func isValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
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
		Shell:      "bash",
		Bling:      "default",
		Name:       "My Fleek Configuration",
		Repository: "git@github.com/my/homeconfig",
		Packages: []string{
			"helix",
		},
		Programs: []string{
			"dircolors",
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
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		csym := filepath.Join(home, ".fleek.yml")
		err = os.Symlink(cfile, csym)
		if err != nil {
			return err
		}
	} else {
		return errors.New("cowardly refusing to overwrite config file without --force flag")
	}
	return nil
}
