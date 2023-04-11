package fleek

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ublue-os/fleek/internal/ux"
	"gopkg.in/yaml.v3"
)

var (
	operatingSystems = []string{"linux", "darwin"}
	architectures    = []string{"aarch64", "x86_64"}
	shells           = []string{"bash", "zsh"}
	blingLevels      = []string{"none", "low", "default", "high"}
	LowPackages      = []string{"htop", "git", "github-cli", "glab"}
	DefaultPackages  = []string{"fzf", "ripgrep", "vscode"}
	HighPackages     = []string{"lazygit", "jq", "yq", "neovim", "neofetch", "btop", "cheat"}
	LowPrograms      = []string{"starship"}
	DefaultPrograms  = []string{"direnv"}
	HighPrograms     = []string{"exa", "bat", "atuin", "zoxide"}
)

// Config holds the options that will be
// merged into the home-manager flake.
type Config struct {
	Debug    bool   `yaml:"-"`
	Verbose  bool   `yaml:"-"`
	Force    bool   `yaml:"-"`
	Quiet    bool   `yaml:"-"`
	FlakeDir string `yaml:"flakedir"`
	Unfree   bool   `yaml:"unfree"`
	// bash or zsh
	Shell string `yaml:"shell"`
	// low, default, high
	Bling    string            `yaml:"bling"`
	Name     string            `yaml:"name"`
	Packages []string          `yaml:",flow"`
	Programs []string          `yaml:",flow"`
	Aliases  map[string]string `yaml:",flow"`
	Paths    []string          `yaml:"paths"`
	Ejected  bool              `yaml:"ejected"`
	Systems  []*System         `yaml:",flow"`
	Git      Git               `yaml:"git"`
}

func Levels() []string {
	return blingLevels
}

type Git struct {
	Enabled    bool `yaml:"enabled"`
	AutoCommit bool `yaml:"autocommit"`
	AutoPush   bool `yaml:"autopush"`
	AutoPull   bool `yaml:"autopull"`
}

type System struct {
	Hostname string `yaml:"hostname"`
	Username string `yaml:"username"`
	Arch     string `yaml:"arch"`
	OS       string `yaml:"os"`
}

func (s System) HomeDir() string {
	base := "/home"
	if s.OS == "darwin" {
		base = "/Users"
	}
	return base + "/" + s.Username
}

func NewSystem() (*System, error) {
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
	}, nil
}

var (
	ErrMissingFlakeDir        = errors.New("fleek.yml: missing `flakedir`")
	ErrInvalidShell           = errors.New("fleek.yml: invalid shell, valid shells are: " + strings.Join(shells, ", "))
	ErrInvalidBling           = errors.New("fleek.yml: invalid bling level, valid levels are: " + strings.Join(blingLevels, ", "))
	ErrorInvalidArch          = errors.New("fleek.yml: invalid architecture, valid architectures are: " + strings.Join(architectures, ", "))
	ErrInvalidOperatingSystem = errors.New("fleek.yml: invalid OS, valid operating systems are: " + strings.Join(operatingSystems, ", "))
	ErrPackageNotFound        = errors.New("package not found in configuration file")
	ErrProgramNotFound        = errors.New("program not found in configuration file")
)

func (c *Config) Validate() error {
	if c.FlakeDir == "" {
		return ErrMissingFlakeDir
	}
	if !isValueInList(c.Shell, shells) {
		return ErrInvalidShell
	}
	if !isValueInList(c.Bling, blingLevels) {
		return ErrInvalidBling
	}
	for _, sys := range c.Systems {
		if !isValueInList(sys.Arch, architectures) {
			return ErrorInvalidArch
		}

		if !isValueInList(sys.OS, operatingSystems) {
			return ErrInvalidOperatingSystem
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

func (c *Config) UserFlakeDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, c.FlakeDir)
}

func (c *Config) AddPackage(pack string) error {
	var found bool
	for _, p := range c.Packages {
		if p == pack {
			found = true
			break
		}
	}
	if found {
		return nil
	}
	c.Packages = append(c.Packages, pack)
	err := c.Validate()
	if err != nil {
		return err
	}
	return c.Save()
}
func (c *Config) RemovePackage(pack string) error {
	var index int
	var found bool
	for x, p := range c.Packages {
		if p == pack {
			index = x
			found = true
			break
		}
	}
	if found {
		c.Packages = append(c.Packages[:index], c.Packages[index+1:]...)
	} else {
		return ErrPackageNotFound
	}
	err := c.Validate()
	if err != nil {
		return err
	}
	return c.Save()
}
func (c *Config) RemoveProgram(prog string) error {
	var index int
	var found bool
	for x, p := range c.Programs {
		if p == prog {
			index = x
			found = true
			break
		}
	}
	if found {
		c.Programs = append(c.Programs[:index], c.Programs[index+1:]...)
	} else {
		return ErrProgramNotFound
	}
	err := c.Validate()
	if err != nil {
		return err
	}
	return c.Save()
}
func (c *Config) AddProgram(prog string) error {
	c.Programs = append(c.Programs, prog)
	err := c.Validate()
	if err != nil {
		return err
	}
	return c.Save()
}

func (c *Config) Save() error {
	cfile, err := c.Location()
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
	home, err := os.UserHomeDir()
	if err != nil {
		return c, err
	}
	csym := filepath.Join(home, ".fleek.yml")
	bb, err := os.ReadFile(csym)
	if err != nil {
		return c, err
	}
	err = yaml.Unmarshal(bb, c)
	if err != nil {
		return c, err
	}
	return c, nil
}

func (c *Config) WriteInitialConfig(force bool, symlink bool) error {
	aliases := make(map[string]string)
	aliases["fleeks"] = "cd " + c.UserFlakeDir()
	sys, err := NewSystem()
	if err != nil {
		ux.Debug.Printfln("new system err: %s ", err)
		return err
	}
	c.Unfree = true
	c.Name = "Fleek Configuration"
	c.Packages = []string{
		"helix",
	}
	c.Programs = []string{
		"dircolors",
	}
	c.Aliases = aliases
	c.Paths = []string{
		"$HOME/bin",
		"$HOME/.local/bin",
	}
	c.Systems = []*System{sys}

	cfile, err := c.Location()
	if err != nil {
		ux.Debug.Printfln("location err: %s ", err)
		return err
	}
	ux.Debug.Printfln("cfile: %s", cfile)

	_, err = os.Stat(cfile)

	ux.Debug.Printfln("stat err: %v ", err)
	ux.Debug.Printfln("force: %v ", force)

	if force || errors.Is(err, fs.ErrNotExist) {

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
		if symlink {
			csym := filepath.Join(home, ".fleek.yml")
			err = os.Symlink(cfile, csym)
			if err != nil {
				return err
			}
		}

	} else {
		return errors.New("cowardly refusing to overwrite config file without --force flag")
	}
	return nil
}

// WriteEjectConfig updates the .fleek.yml file
// to indicated ejected status
func (c *Config) Eject() error {

	c.Ejected = true

	cfile, err := c.Location()
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
