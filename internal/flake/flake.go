package flake

import (
	"embed"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/riywo/loginshell"
	app "github.com/ublue-os/fleek"
	"github.com/ublue-os/fleek/fin"
	"github.com/ublue-os/fleek/internal/fleek"
)

const nixbin = "nix"

var ErrPackageConflict = errors.New("package exists in fleek and nix profile")

type Flake struct {
	Templates map[string]*template.Template
	Config    *fleek.Config
	app       *app.App
}
type Data struct {
	Config   *fleek.Config
	UserName string
	Home     string
	Bling    *fleek.Bling
}
type SystemData struct {
	System fleek.System
	User   fleek.User
	BYOGit bool
}

func Load(cfg *fleek.Config, app *app.App) (*Flake, error) {
	if cfg.Verbose {
		fin.Verbose.Println(app.Trans("flake.initializingTemplates"))
	}
	tt := make(map[string]*template.Template)
	err := fs.WalkDir(templates, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".tmpl" {
			bb, err := templates.ReadFile(path)
			if err != nil {
				return err
			}
			tt[path] = template.Must(template.New(path).Parse(string(bb)))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Flake{
		Templates: tt,
		Config:    cfg,
		app:       app,
	}, nil
}

func (f *Flake) Update() error {
	fin.Info.Println(f.app.Trans("flake.update"))

	updateCmdLine := []string{"flake", "update"}
	err := f.runNix(nixbin, updateCmdLine)

	if err != nil {
		return err
	}
	err = f.mayCommit("fleek: update flake.lock")

	if err != nil {
		return err
	}
	return nil
}
func (f *Flake) Create(force bool, symlink bool) error {
	fin.Info.Println(f.app.Trans("init.writingConfigs"))
	err := f.ensureFlakeDir()
	if err != nil {
		return err
	}
	fullShell, err := loginshell.Shell()
	if err != nil {
		return err
	}
	shell := filepath.Base(fullShell)
	f.Config.Shell = shell
	err = f.Config.Validate()
	if err != nil {
		return err
	}
	var bling *fleek.Bling

	switch f.Config.Bling {
	case "high":
		bling, err = fleek.HighBling()
	case "default":
		bling, err = fleek.DefaultBling()
	case "low":
		bling, err = fleek.LowBling()
	case "none":
		bling, err = fleek.NoBling()
	default:
		bling, err = fleek.DefaultBling()
	}
	if err != nil {
		return err
	}

	fin.Info.Println(f.app.Trans("init.blingLevel", f.Config.Bling))
	err = f.Config.WriteInitialConfig(force, symlink)
	if err != nil {
		return err
	}
	// read the config again, because it may have changed
	loc := f.Config.UserFlakeDir()
	config, err := fleek.ReadConfig(loc)
	if err != nil {
		return err
	}
	f.Config = config
	sys, err := f.Config.CurrentSystem()
	if err != nil {
		return err
	}
	user := f.Config.UserForSystem(sys.Hostname)
	data := Data{
		Config: f.Config,
		Bling:  bling,
	}

	err = f.writeFile("templates/flake.nix.tmpl", "flake.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/README.md.tmpl", "README.md", data, force)
	if err != nil {
		return err
	}

	err = f.writeFile("templates/home.nix.tmpl", "home.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/aliases.nix.tmpl", "aliases.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/path.nix.tmpl", "path.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/programs.nix.tmpl", "programs.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/shell.nix.tmpl", "shell.nix", data, force)
	if err != nil {
		return err
	}

	err = f.writeSystem(*sys, "templates/host.nix.tmpl", force)
	if err != nil {
		return err
	}

	err = f.writeFile("templates/user.nix.tmpl", "user.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeUser(*sys, *user, "templates/user.nix.tmpl", force)
	if err != nil {
		return err
	}
	return nil

}
func (f *Flake) IsJoin() (bool, error) {
	// if the user has a flake.nix, but no fleek.yml, then we assume they want to join
	_, err := os.Stat(filepath.Join(f.Config.UserFlakeDir(), "flake.nix"))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	home, _ := os.UserHomeDir()
	_, err = os.Stat(filepath.Join(home, ".fleek.yml"))
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
func (f *Flake) Join() error {
	fin.Info.Println(f.app.Trans("init.writingConfigs"))

	// Symlink the yaml file to home
	cfile, err := f.Config.Location()
	if err != nil {
		fin.Debug.Printfln("location err: %s ", err)
		return err
	}
	fin.Debug.Printfln("init cfile: %s ", cfile)

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	csym := filepath.Join(home, ".fleek.yml")
	err = os.Symlink(cfile, csym)
	if err != nil {
		fin.Debug.Println("first symlink attempt failed")
		return err
	}

	err = f.Config.Validate()
	if err != nil {
		return err
	}
	fin.Debug.Println("new system")

	sys, err := fleek.NewSystem()
	if err != nil {
		return err
	}
	fin.Debug.Println("write system")

	//
	var found bool
	for _, s := range f.Config.Systems {
		if s.Hostname == sys.Hostname && s.Username == sys.Username && s.Arch == sys.Arch {
			fin.Debug.Println("system already exists")
			found = true
		}
	}
	if !found {
		f.Config.Systems = append(f.Config.Systems, sys)
		err = f.writeSystem(*sys, "templates/host.nix.tmpl", true)
		if err != nil {
			return err
		}

	}
	username, err := fleek.Username()
	if err != nil {
		return err
	}
	userFound := false
	for _, u := range f.Config.Users {
		if u.Username == username {
			userFound = true
		}
	}
	if !userFound {
		user, err := fleek.NewUser()
		if err != nil {
			return err
		}
		f.Config.Users = append(f.Config.Users, user)
		err = f.writeUser(*sys, *user, "templates/user.nix.tmpl", true)
		if err != nil {
			return err
		}
	}

	fin.Debug.Println("write config")

	err = f.Config.Save()
	if err != nil {
		fin.Debug.Println("config save failed")
		return err
	}
	git, err := f.IsGitRepo()
	if err != nil {
		return err
	}
	if git {
		fin.Warning.Println(f.app.Trans("git.warn"))
		err = f.setRebase()
		if err != nil {
			return err
		}
		err = f.mayCommit("fleek: new system")
		if err != nil {
			return err
		}
	}
	return nil

}

func (f *Flake) Check() error {
	checkCmdLine := []string{"run", "--impure", "home-manager/master", "build", "--impure", "--", "--flake", "."}
	err := f.runNix(nixbin, checkCmdLine)

	if err != nil {
		return err
	}
	return nil
}

// Write writes the applied flake configuration
func (f *Flake) Write(message string, writeHost, writeUser bool) error {
	force := true
	spinner, err := fin.Spinner().Start(f.app.Trans("flake.writing"))
	if err != nil {
		return err
	}

	var bling *fleek.Bling
	switch f.Config.Bling {
	case "high":
		bling, err = fleek.HighBling()
	case "default":
		bling, err = fleek.DefaultBling()
	case "low":
		bling, err = fleek.LowBling()
	case "none":
		bling, err = fleek.NoBling()
	default:
		bling, err = fleek.DefaultBling()
	}
	if err != nil {
		return err
	}

	data := Data{
		Config: f.Config,
		Bling:  bling,
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	err = f.ReadConfig(filepath.Join(home, f.Config.FlakeDir))
	if err != nil {
		return err
	}
	err = f.writeFile("templates/flake.nix.tmpl", "flake.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/README.md.tmpl", "README.md", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/.gitignore.tmpl", ".gitignore", data, force)
	if err != nil {
		return err
	}

	err = f.writeFile("templates/home.nix.tmpl", "home.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/aliases.nix.tmpl", "aliases.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/path.nix.tmpl", "path.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/programs.nix.tmpl", "programs.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/shell.nix.tmpl", "shell.nix", data, force)
	if err != nil {
		return err
	}
	sys, err := f.Config.CurrentSystem()
	if err != nil {
		return err
	}
	if writeHost {

		err = f.writeSystem(*sys, "templates/host.nix.tmpl", force)
		if err != nil {
			return err
		}
	}
	if writeHost {

		user := f.Config.UserForSystem(sys.Hostname)
		err = f.writeUser(*sys, *user, "templates/user.nix.tmpl", true)
		if err != nil {
			return err
		}
	}

	spinner.Success()
	err = f.mayCommit(message)

	if err != nil {
		return err
	}
	return nil

}

func (f *Flake) ensureFlakeDir() error {
	if f.Config.Verbose {
		fin.Verbose.Println(f.app.Trans("flake.ensureDir"))
	}
	err := f.Config.MakeFlakeDir()
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return nil
		}
	}
	return nil
}

func (f *Flake) ReadConfig(loc string) error {
	// load the new config
	config, err := fleek.ReadConfig(loc)
	if err != nil {
		return err
	}
	f.Config = config
	return nil
}
func (f *Flake) writeFile(template string, path string, d Data, force bool) error {
	fpath := filepath.Join(f.Config.UserFlakeDir(), path)
	err := os.MkdirAll(filepath.Dir(fpath), 0755)
	if err != nil {
		fin.Debug.Println("mkdir error", err)
	}
	_, err = os.Stat(fpath)

	if force || errors.Is(err, fs.ErrNotExist) {
		_, ok := f.Templates[template]

		if ok {
			file, err := os.Create(fpath)
			if err != nil {
				fin.Debug.Println("create error", err)
				return err
			}
			defer file.Close()

			if err = f.Templates[template].Execute(file, d); err != nil {
				fin.Debug.Println("template error", err)
				return err
			}
		} else {
			fin.Debug.Println("template not found", template)
			return errors.New("template not found")
		}
	} else {
		return errors.New("cowardly refusing to overwrite existing file without --force flag")
	}
	return nil
}

func (f *Flake) writeSystem(sys fleek.System, template string, force bool) error {
	var user *fleek.User
	var err error
	user = f.Config.UserForSystem(sys.Hostname)
	if user == nil {
		user, err = fleek.NewUser()
		if err != nil {
			return err
		}
		f.Config.Users = append(f.Config.Users, user)
		err = f.Config.Save()
		if err != nil {
			return err
		}
	}
	sysData := SystemData{
		System: sys,
		User:   *user,
	}
	if f.Config.BYOGit {
		sysData.BYOGit = true
	}

	hostPath := filepath.Join(f.Config.UserFlakeDir(), sys.Hostname)
	err = os.MkdirAll(hostPath, 0755)
	if err != nil {
		return err
	}
	fpath := filepath.Join(hostPath, sys.Hostname+".nix")
	_, err = os.Stat(fpath)
	if force || os.IsNotExist(err) {

		file, err := os.Create(fpath)
		if err != nil {
			return err
		}
		defer file.Close()

		if err = f.Templates[template].Execute(file, sysData); err != nil {
			return err
		}

	} else {
		return errors.New("cowardly refusing to overwrite existing file without --force flag")
	}

	return nil
}
func (f *Flake) writeUser(sys fleek.System, user fleek.User, template string, force bool) error {

	hostPath := filepath.Join(f.Config.UserFlakeDir(), sys.Hostname)
	err := os.MkdirAll(hostPath, 0755)
	if err != nil {
		return err
	}
	fpath := filepath.Join(hostPath, "user.nix")
	_, err = os.Stat(fpath)
	if force || os.IsNotExist(err) {

		file, err := os.Create(fpath)
		if err != nil {
			return err
		}
		defer file.Close()

		if err = f.Templates[template].Execute(file, user); err != nil {
			return err
		}

	} else {
		return errors.New("cowardly refusing to overwrite existing file without --force flag")
	}

	return nil
}
func (f *Flake) Apply() error {
	fin.Info.Println(f.app.Trans("flake.apply"))

	user, err := fleek.Username()

	if err != nil {
		return err
	}
	host, err := fleek.Hostname()
	if err != nil {
		return err
	}
	applyCmdLine := []string{"run", "--no-write-lock-file", "--impure", "home-manager/master", "--", "-b", "bak", "switch", "--flake", ".#" + user + "@" + host}
	err = f.runNix(nixbin, applyCmdLine)
	if err != nil {
		return err
	}
	return nil
}
func (f *Flake) runNix(cmd string, cmdLine []string) error {
	command := exec.Command(cmd, cmdLine...)
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Dir = f.Config.UserFlakeDir()
	command.Env = os.Environ()
	if f.Config.Unfree {
		command.Env = append(command.Env, "NIXPKGS_ALLOW_UNFREE=1")
	}

	return command.Run()

}

//go:embed all:templates
var templates embed.FS
