package fleek

import (
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pterm/pterm"
	"github.com/ublue-os/fleek/fin"
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
	Users    []*User           `yaml:",flow"`
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
type User struct {
	Username          string `yaml:"username"`
	Name              string `yaml:"name"`
	Email             string `yaml:"email"`
	SSHPublicKeyFile  string `yaml:"ssh_public_key_file"`
	SSHPrivateKeyFile string `yaml:"ssh_private_key_file"`
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
func NewUser() (*User, error) {
	fin.Info.Println("Enter User Details for Git Configuration:")
	user := &User{}
	name, err := Name()
	if err != nil {
		return user, err
	}
	// Prompt for name
	var use bool
	use, err = ux.Confirm("Use detected name: " + name)
	if err != nil {
		return user, err
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
	// It doesn't make sense to change the username,
	// so just use the detected one
	uname, err := Username()
	if err != nil {
		return user, err
	}
	user.Username = uname

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
	err = MkdirAll(filepath.Dir(cfile))
	if err != nil {
		if !errors.Is(err, fs.ErrExist) {
			return err
		}
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
func ReadConfig(loc string) (*Config, error) {
	c := &Config{}

	if loc == "" {

		home, err := os.UserHomeDir()
		if err != nil {
			return c, err
		}
		csym := filepath.Join(home, ".fleek.yml")
		loc = csym
	} else {
		loc = filepath.Join(loc, ".fleek.yml")
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
	aliases := make(map[string]string)
	aliases["fleeks"] = "cd ~/" + c.FlakeDir
	sys, err := NewSystem()
	if err != nil {
		fin.Debug.Printfln("new system err: %s ", err)
		return err
	}
	user, err := NewUser()
	if err != nil {
		fin.Debug.Printfln("new user err: %s ", err)
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
	c.Users = []*User{user}
	c.Git.Enabled = true
	c.Git.AutoCommit = true
	c.Git.AutoPull = true
	c.Git.AutoPush = true

	cfile, err := c.Location()
	if err != nil {
		fin.Debug.Printfln("location err: %s ", err)
		return err
	}
	fin.Debug.Printfln("cfile: %s", cfile)

	_, err = os.Stat(cfile)

	fin.Debug.Printfln("stat err: %v ", err)
	fin.Debug.Printfln("force: %v ", force)

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

func (c *Config) NeedsMigration() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	count := 0
	aliases := filepath.Join(home, c.FlakeDir, "aliases.nix")
	_, err = os.Stat(aliases)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			count = count + 1
		}
	}
	pathnix := filepath.Join(home, c.FlakeDir, "path.nix")
	_, err = os.Stat(pathnix)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			count = count + 1
		}
	}
	return count < 1

}
func (c *Config) MigrateV2() error {
	if len(c.Users) < 1 {
		user, err := NewUser()
		if err != nil {
			fin.Debug.Printfln("new user err: %s ", err)
			return err
		}
		c.Users = append(c.Users, user)
		err = c.Save()
		if err != nil {
			return err
		}
		// NEED TO LOOK THROUGH SYSTEMS
		// and create users referenced there
		for _, sys := range c.Systems {
			u := sys.Username
			found := false
			for _, uu := range c.Users {
				if uu.Email == u {
					found = true
				}
			}
			if !found {
				// create new user
				email, err := pterm.DefaultInteractiveTextInput.Show(sys.Hostname + " - Enter your email for git")
				if err != nil {
					return err
				}
				name, err := pterm.DefaultInteractiveTextInput.Show(sys.Hostname + " - Enter your name for git")
				if err != nil {
					return err
				}
				// TODO gracefully ask for keys
				user := &User{
					Email:    email,
					Name:     name,
					Username: sys.Username,
				}
				// save
				c.Users = append(c.Users, user)
				err = c.Save()
				if err != nil {
					return err
				}
			}
		}

	}

	// move user config
	// from ./user.nix
	// to home/users/{user}/custom.nix
	uname, err := Username()
	if err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	newDir := filepath.Join(home, c.FlakeDir, "home", "users", uname)
	fin.Debug.Println("newDir: ", newDir)
	err = os.MkdirAll(newDir, 0755)
	if err != nil {
		return err
	}
	oldLocation := filepath.Join(home, c.FlakeDir, "user.nix")
	newLocation := filepath.Join(newDir, "custom.nix")
	err = os.Rename(oldLocation, newLocation)
	if err != nil {
		return err
	}
	// move system config
	hostDir := filepath.Join(home, c.FlakeDir, "home", "hosts")
	err = os.MkdirAll(hostDir, 0755)
	if err != nil {
		return err
	}
	for _, sys := range c.Systems {
		oldLocation := filepath.Join(home, c.FlakeDir, sys.Hostname, "user.nix")
		newLocation := filepath.Join(hostDir, sys.Hostname+".nix")
		err = os.Rename(oldLocation, newLocation)
		if err != nil {
			return err
		}
		oldDir := filepath.Join(home, c.FlakeDir, sys.Hostname)
		err = os.RemoveAll(oldDir)
		if err != nil {
			return err
		}

	}
	// config files
	aliases := filepath.Join(home, c.FlakeDir, "aliases.nix")
	err = os.Remove(aliases)
	if err != nil {
		return err
	}
	homenix := filepath.Join(home, c.FlakeDir, "home.nix")
	err = os.Remove(homenix)
	if err != nil {
		return err
	}
	pathnix := filepath.Join(home, c.FlakeDir, "path.nix")
	err = os.Remove(pathnix)
	if err != nil {
		return err
	}
	prognix := filepath.Join(home, c.FlakeDir, "programs.nix")
	err = os.Remove(prognix)
	if err != nil {
		return err
	}
	shellnix := filepath.Join(home, c.FlakeDir, "shell.nix")
	err = os.Remove(shellnix)
	if err != nil {
		return err
	}
	flakelock := filepath.Join(home, c.FlakeDir, "flake.lock")
	err = os.Remove(flakelock)
	if err != nil {
		return err
	}
	return nil
}
func (c *Config) ActiveUser() (*User, error) {
	u, err := Username()
	if err != nil {
		return nil, err
	}
	for _, user := range c.Users {
		if user.Username == u {
			return user, nil
		}
	}
	return nil, nil
}
