package core

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	operatingSystems = []string{"linux", "darwin"}
	architectures    = []string{"aarch64", "x86_64"}
	shells           = []string{"bash", "zsh"}
	blingLevels      = []string{"low", "default", "high"}
	lowPackages      = []string{"htop", "git", "github-cli", "glab"}
	defaultPackages  = []string{"fzf", "ripgrep", "vscode"}
	highPackages     = []string{"lazygit", "jq", "yq", "neovim", "neofetch", "btop", "cheat"}
	lowPrograms      = []string{"starship"}
	defaultPrograms  = []string{"direnv"}
	highPrograms     = []string{"exa", "bat", "atuin", "zoxide"}
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
	Ejected    bool              `yaml:"ejected"`
	Systems    []System          `yaml:",flow"`
}
type GitConfig struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

type System struct {
	Hostname  string    `yaml:"hostname"`
	Username  string    `yaml:"username"`
	Arch      string    `yaml:"arch"`
	OS        string    `yaml:"os"`
	GitConfig GitConfig `yaml:"git"`
}

func (s System) HomeDir() string {
	base := "/home"
	if s.OS == "darwin" {
		base = "/Users"
	}
	return base + "/" + s.Username
}

func NewSystem(name, email string) (*System, error) {
	user, err := Username()
	if err != nil {
		return nil, err
	}
	host, err := Hostname()
	if err != nil {
		return nil, err
	}
	return &System{
		Hostname: host,
		Arch:     Arch(),
		OS:       runtime.GOOS,
		Username: user,
		GitConfig: GitConfig{
			Name:  name,
			Email: email,
		},
	}, nil
}

func (c Config) Validate() error {
	if !isValueInList(c.Shell, shells) {
		return errors.New("fleek.yml: invalid shell, valid shells are: " + strings.Join(shells, ", "))
	}
	if !isValueInList(c.Bling, blingLevels) {
		return errors.New("fleek.yml: invalid bling level, valid levels are: " + strings.Join(blingLevels, ", "))
	}
	for _, sys := range c.Systems {
		if !isValueInList(sys.Arch, architectures) {
			return errors.New("fleek.yml: invalid architecture, valid architectures are: " + strings.Join(architectures, ", "))
		}

		if !isValueInList(sys.OS, operatingSystems) {
			return errors.New("fleek.yml: invalid OS, valid operating systems are: " + strings.Join(operatingSystems, ", "))
		}
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

func (c *Config) Save() error {
	cfile, err := ConfigLocation()
	if err != nil {
		return err
	}
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
	return nil
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

func Ejected() (bool, error) {
	conf, err := ReadConfig()
	if err != nil {
		return false, err
	}
	return conf.Ejected, nil
}

func Clone(repo string) error {
	location, err := FlakeLocation()
	if err != nil {
		return err
	}

	clone := exec.Command("git", "clone", repo, location)
	clone.Stderr = os.Stderr
	clone.Stdin = os.Stdin
	clone.Stdout = os.Stdout
	clone.Env = os.Environ()

	return clone.Run()

}

// WriteSampleConfig creates the first fleek
// configuration file
func WriteSampleConfig(email, name string, force bool) error {

	aliases := make(map[string]string)
	aliases["cdfleek"] = "cd ~/.config/home-manager"
	sys, err := NewSystem(name, email)
	if err != nil {
		return err
	}
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
		Systems: []System{*sys},
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

// WriteEjectConfig updates the .fleek.yml file
// to indicated ejected status
func WriteEjectConfig() error {

	c := Config{
		Ejected: true,
	}
	cfile, err := ConfigLocation()
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

	err = os.WriteFile(cfile, n, 0755)
	if err != nil {
		return err
	}

	return nil
}
