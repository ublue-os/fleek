package fleek

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/cmdutil"
	"github.com/ublue-os/fleek/internal/ux"
	"github.com/ublue-os/fleek/internal/xdg"
	"gopkg.in/yaml.v3"
)

var (
	operatingSystems = []string{"linux", "darwin"}
	architectures    = []string{"aarch64", "x86_64"}
	shells           = []string{"bash", "zsh"}
	blingLevels      = []string{"none", "low", "default", "high"}
	LowPackages      = []string{"htop", "git", "github-cli", "glab"}
	DefaultPackages  = []string{"fzf", "ripgrep", "vscode", "just"}
	HighPackages     = []string{"lazygit", "jq", "yq", "neovim", "neofetch", "btop", "cheat"}
	LowPrograms      = []string{"starship"}
	DefaultPrograms  = []string{"direnv"}
	HighPrograms     = []string{"eza", "bat", "atuin", "zoxide"}
)

var systemAliases = map[string]string{
	"update-fleek":         "nix run https://getfleek.dev/latest.tar.gz -- update",
	"latest-fleek-version": "nix run https://getfleek.dev/latest.tar.gz -- version",
}

// Config holds the options that will be
// merged into the home-manager flake.
type Config struct {
	MinVersion string `yaml:"min_version"`
	Debug      bool   `yaml:"-"`
	Verbose    bool   `yaml:"-"`
	Force      bool   `yaml:"-"`
	Quiet      bool   `yaml:"-"`
	FlakeDir   string `yaml:"flakedir"`
	Unfree     bool   `yaml:"unfree"`
	// bash or zsh
	Shell string `yaml:"shell"`
	// low, default, high
	Bling    string              `yaml:"bling"`
	Name     string              `yaml:"name"`
	Overlays map[string]*Overlay `yaml:",flow"`
	Packages []string            `yaml:",flow"`
	Programs []string            `yaml:",flow"`
	// issue 211, remove or block bling packages
	Blocklist []string          `yaml:"blocklist,flow"`
	Aliases   map[string]string `yaml:",flow"`
	Paths     []string          `yaml:"paths"`
	Ejected   bool              `yaml:"ejected"`
	// issue 200 - disable any git integration
	BYOGit      bool      `yaml:"byo_git"`
	Systems     []*System `yaml:",flow"`
	Git         Git       `yaml:"git"`
	Users       []*User   `yaml:",flow"`
	Track       string    `yaml:"track"`
	AllowBroken bool      `yaml:"allow_broken"`
	AutoGC      bool      `yaml:"auto_gc"`
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
	Home     string `yaml:"home"`
	User     *User  `yaml:"user"`
}

type User struct {
	Username          string `yaml:"username"`
	Name              string `yaml:"name"`
	Email             string `yaml:"email"`
	SSHPublicKeyFile  string `yaml:"ssh_public_key_file"`
	SSHPrivateKeyFile string `yaml:"ssh_private_key_file"`
}

type Overlay struct {
	URL    string `yaml:"url"`
	Follow bool   `yaml:"follow"`
}

func (u User) HomeDir(s System) string {
	if s.Home != "" {
		return s.Home
	}
	base := "/home"
	if s.OS == "darwin" {
		base = "/Users"
	}
	return base + "/" + u.Username
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

// CollectGarbage runs nix-collect-garbage
func CollectGarbage() error {
	command := cmdutil.CommandTTY("nix-collect-garbage", "-d")
	command.Stderr = io.Discard
	command.Stdout = io.Discard
	command.Env = os.Environ()

	return command.Run()

}
func NewUser() (*User, error) {
	fin.Logger.Info("Enter User Details for Git Configuration:")
	user := &User{}
	var use bool

	envname := os.Getenv("FLEEK_USER_NAME")
	if envname == "" {
		name, err := Name()
		if err != nil {
			return user, err
		}
		// Prompt for name
		name = strings.TrimSpace(name)
		if name != "" {

			fin.Logger.Info("Detected your name: " + name)
			use, err = ux.Confirm("Use detected name: " + name)
			if err != nil {
				return user, err
			}
		}
		if use {
			user.Name = name
		} else {
			prompt := "Name"
			iname, err := ux.Input(prompt, name, "Your Name")
			if err != nil {
				return user, err
			}
			user.Name = iname
		}
	} else {
		user.Name = envname
	}

	// It doesn't make sense to change the username,
	// so just use the detected one
	uname, err := Username()
	if err != nil {
		return user, err
	}
	user.Username = uname

	envmail := os.Getenv("FLEEK_USER_EMAIL")
	if envmail == "" {
		// email

		cmd := "git"
		cmdLine := []string{"config", "--global", "user.email"}
		command := exec.Command(cmd, cmdLine...)
		command.Stdin = os.Stdin

		command.Env = os.Environ()
		var email string
		bb, err := command.Output()
		if err != nil {
			// get the email manually
			prompt := "Email"
			email, err = ux.Input(prompt, "", "Your Email Address")
			if err != nil {
				return user, err
			}
			user.Email = email
		} else {
			email = strings.TrimSpace(string(bb))
			use, err = ux.Confirm("Use detected email: " + email)
			if err != nil {
				return user, err
			}
			if use {
				user.Email = email
			} else {
				prompt := "Email"
				uemail, err := ux.Input(prompt, "", "Your Email Address")
				if err != nil {
					return user, err
				}
				user.Email = uemail
			}
		}
	} else {
		user.Email = envmail
	}
	envpubkey := os.Getenv("FLEEK_USER_PUBKEY")
	envprivkey := os.Getenv("FLEEK_USER_PRIVKEY")
	if (envpubkey == "") && (envprivkey == "") {

		// ssh keys
		privateKey := ""
		publicKey := ""

		// find and add ssh keys
		sshDir := filepath.Join(os.Getenv("HOME"), ".ssh")
		sshFiles, err := os.ReadDir(sshDir)
		hasSSH := true
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				hasSSH = false
			} else {
				return user, err
			}
		}
		if hasSSH {
			candidates := []string{}
			for _, f := range sshFiles {
				if strings.HasSuffix(f.Name(), ".pub") {
					candidates = append(candidates, f.Name())
				}
			}
			if len(candidates) > 0 {
				key, err := ux.PromptSingle("Choose Git SSH Key", candidates)
				if err != nil {
					return user, err
				}
				privateKey = strings.Replace(key, ".pub", "", 1)
				privateKey = filepath.Join("~", ".ssh", privateKey)
				publicKey = filepath.Join("~", ".ssh", key)
				user.SSHPrivateKeyFile = privateKey
				user.SSHPublicKeyFile = publicKey
			}
		}
	} else {
		user.SSHPrivateKeyFile = envprivkey
		user.SSHPublicKeyFile = envpubkey
	}

	return user, nil
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

func (c *Config) Tracks() string {
	if c.Track != "" {
		return c.Track
	}
	return "nixos-unstable"
}
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

func (c *Config) UniqueSystems() []string {
	var m = make(map[string]bool)
	var systems = []string{}

	for _, sys := range c.Systems {
		syskey := sys.Arch + "-" + sys.OS
		if m[syskey] {
			continue
		}
		systems = append(systems, syskey)
		m[syskey] = true
	}
	sort.Strings(systems)
	return systems

}

func (c *Config) UserFlakeDir() string {
	home, _ := os.UserHomeDir()
	// if for some reason the flakedir key is
	// missing, try loading the default location
	if c.FlakeDir == "" {
		return filepath.Join(home, xdg.DataSubpathRel("fleek"))
	}
	return filepath.Join(home, c.FlakeDir)
}

func (c *Config) UserForSystem(system string) *User {
	var userSystem *System
	for _, sys := range c.Systems {
		if sys.Hostname == system {
			userSystem = sys
		}
	}
	if userSystem.User != nil {
		return userSystem.User
	}
	// legacy unmigrated users
	if c.Users != nil {
		for _, u := range c.Users {
			if u.Username == userSystem.Username {
				return u
			}
		}
	}
	return nil
}

func (c *Config) AllAliases() map[string]string {
	for k, v := range systemAliases {
		c.Aliases[k] = v
	}
	return c.Aliases
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
// pointed to in the $HOME/.fleek.yml symlink
func ReadConfig(loc string) (*Config, error) {
	c := &Config{}
	home, err := os.UserHomeDir()
	if err != nil {
		return c, err
	}
	if loc == "" {

		csym := filepath.Join(home, ".fleek.yml")
		loc = csym
	} else {
		if strings.HasPrefix(loc, home) {
			loc = filepath.Join(loc, ".fleek.yml")
		} else {
			loc = filepath.Join(home, loc, ".fleek.yml")
		}
	}
	bb, err := os.ReadFile(loc)
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
	systemAliases["fleeks"] = "cd ~/" + c.FlakeDir
	sys, err := NewSystem()
	if err != nil {
		fin.Logger.Debug("new system", fin.Logger.Args("error", err))
		return err
	}
	user, err := NewUser()
	if err != nil {
		fin.Logger.Debug("new user", fin.Logger.Args("error", err))

		return err
	}
	sys.User = user
	c.Unfree = true
	c.AutoGC = true
	c.Name = "Fleek Configuration"
	c.Packages = []string{
		"helix",
	}
	c.Programs = []string{
		"dircolors",
	}
	c.Aliases = systemAliases
	c.Paths = []string{
		"$HOME/bin",
		"$HOME/.local/bin",
	}
	c.Systems = []*System{sys}
	c.MinVersion = "0.8.4"
	c.Track = "nixos-unstable"
	c.BYOGit = false
	c.Git.Enabled = true
	c.Git.AutoCommit = true
	c.Git.AutoPull = true
	c.Git.AutoPush = true

	cfile, err := c.Location()
	if err != nil {
		fin.Logger.Debug("location err", fin.Logger.Args("error", err))

		return err
	}
	fin.Logger.Debug("config", fin.Logger.Args("file", cfile))

	_, err = os.Stat(cfile)
	fin.Logger.Debug("stat", fin.Logger.Args("error", err))
	fin.Logger.Debug("force", fin.Logger.Args("value", force))

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
			// ignore the error. Delete if it exists
			_ = os.Remove(filepath.Join(home, ".fleek.yml"))
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

func (c *Config) AsVersion() (*version.Version, error) {
	return version.NewVersion(c.MinVersion)
}

// Needs migration checks to see if the host directory
// has a file with the same name as the host.
// e.g. ./beast/beast.nix
func (c *Config) NeedsMigration() bool {
	for _, s := range c.Systems {
		systemDir := filepath.Join(c.UserFlakeDir(), s.Hostname)
		systemFile := filepath.Join(systemDir, s.Hostname+".nix")
		// beast/beast.nix
		if Exists(systemFile) {
			fin.Logger.Warn("Found unmigrated system file:", fin.Logger.Args("file", systemFile))

			return true
		}

		hostFile := filepath.Join(systemDir, "host.nix")
		// beast/host.nix
		if Exists(hostFile) {
			fin.Logger.Warn("Found unmigrated system file:", fin.Logger.Args("file", hostFile))
			return true
		}
		hostFile = filepath.Join(systemDir, "user.nix")
		// beast/user.nix
		if Exists(hostFile) {
			fin.Logger.Warn("Found unmigrated system file:", fin.Logger.Args("file", hostFile))
			return true
		}
		if s.User == nil {
			fin.Logger.Warn("Found unmigrated system users")
			return true
		}

	}
	return false
}

func (c *Config) Migrate() error {
	for _, s := range c.Systems {
		systemDir := filepath.Join(c.UserFlakeDir(), s.Hostname)
		systemFile := filepath.Join(systemDir, s.Hostname+".nix")
		// beast/beast.nix
		if Exists(systemFile) {
			fin.Logger.Warn("Migrating system file", fin.Logger.Args("file", systemFile))
			userFile := filepath.Join(systemDir, s.Username+".nix")
			err := Move(systemFile, userFile)
			if err != nil {
				return err
			}
		}
		hostFile := filepath.Join(systemDir, "user.nix")
		// beast/user.nix -> beast/custom.nix
		if Exists(hostFile) {
			fin.Logger.Warn("Migrating system file", fin.Logger.Args("file", hostFile))
			newHostFile := filepath.Join(systemDir, "custom.nix")

			err := Move(hostFile, newHostFile)
			if err != nil {
				return err
			}
		}

		hostFile = filepath.Join(systemDir, "host.nix")
		// beast/host.nix -> beast/custom.nix
		if Exists(hostFile) {
			fin.Logger.Warn("Migrating system file", fin.Logger.Args("file", hostFile))
			newHostFile := filepath.Join(systemDir, "custom.nix")

			err := Move(hostFile, newHostFile)
			if err != nil {
				return err
			}
		}
		if s.User == nil {
			fin.Logger.Warn("Migrating users", fin.Logger.Args("hostname", s.Hostname))
			sysuser := c.UserForSystem(s.Hostname)

			s.User = sysuser
			err := c.Save()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
