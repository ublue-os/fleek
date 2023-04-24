package flake

import (
	"embed"
	"errors"
	"io"
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

func (f *Flake) Update(outWriter io.Writer) error {
	spinner, err := fin.Spinner().Start(f.app.Trans("flake.update"))
	if err != nil {
		return err
	}
	updateCmdLine := []string{"run", ".#update"}
	err = f.runNix(nixbin, updateCmdLine, outWriter)

	if err != nil {
		return err
	}

	spinner.Success()
	err = f.mayCommit("fleek: update flake.lock", outWriter)

	if err != nil {
		return err
	}
	return nil
}
func (f *Flake) Create(force bool, skipConfigWrite bool, symlink bool) error {
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

	data := Data{
		Config: f.Config,
		Bling:  bling,
	}

	if !skipConfigWrite {
		err = f.Config.WriteInitialConfig(f.Config.Force, symlink)
		if err != nil {
			return err
		}
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

	// users directory
	err = f.copyFile("default.nix", "users", force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/users/config.nix.tmpl", "users/config.nix", data, force)
	if err != nil {
		return err
	}
	// home/bling/default directory
	err = f.copyFile("default.nix", "home/bling/default", force)
	if err != nil {
		return err
	}
	err = f.copyFile("direnv.nix", "home/bling/default", force)
	if err != nil {
		return err
	}
	err = f.copyFile("gh.nix", "home/bling/default", force)
	if err != nil {
		return err
	}
	err = f.copyFile("packages.nix", "home/bling/default", force)
	if err != nil {
		return err
	}
	// home/bling/high directory
	err = f.copyFile("default.nix", "home/bling/high", force)
	if err != nil {
		return err
	}
	err = f.copyFile("bat.nix", "home/bling/high", force)
	if err != nil {
		return err
	}
	err = f.copyFile("exa.nix", "home/bling/high", force)
	if err != nil {
		return err
	}
	err = f.copyFile("packages.nix", "home/bling/high", force)
	if err != nil {
		return err
	}
	err = f.copyFile("programs.nix", "home/bling/high", force)
	if err != nil {
		return err
	}
	// home/bling/low directory
	err = f.copyFile("default.nix", "home/bling/low", force)
	if err != nil {
		return err
	}
	err = f.copyFile("packages.nix", "home/bling/low", force)
	if err != nil {
		return err
	}
	err = f.copyFile("programs.nix", "home/bling/low", force)
	if err != nil {
		return err
	}

	// home/bling/none directory
	err = f.copyFile("default.nix", "home/bling/none", force)
	if err != nil {
		return err
	}
	err = f.copyFile("bash.nix", "home/bling/none", force)
	if err != nil {
		return err
	}
	err = f.copyFile("path.nix", "home/bling/none", force)
	if err != nil {
		return err
	}
	err = f.copyFile("zsh.nix", "home/bling/none", force)
	if err != nil {
		return err
	}
	// home/global directory
	err = f.copyFile("default.nix", "home/global", force)
	if err != nil {
		return err
	}
	err = f.copyFile("git.nix", "home/global", force)
	if err != nil {
		return err
	}
	// home/hosts directory
	sys, err := f.Config.CurrentSystem()
	if err != nil {
		return err
	}

	err = f.writeSystem(*sys, "templates/home/hosts/host.nix.tmpl", force)
	if err != nil {
		return err
	}

	// home/users/{user}/custom.nix
	for _, user := range f.Config.Users {
		err = f.writeUser(*user, "templates/home/users/user.nix.tmpl", force)
		if err != nil {
			return err
		}
	}

	// home/users/fleek
	err = f.writeFile("templates/home/users/fleek/aliases.nix.tmpl", "home/users/fleek/aliases.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/home/users/fleek/default.nix.tmpl", "home/users/fleek/default.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/home/users/fleek/packages.nix.tmpl", "home/users/fleek/packages.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/home/users/fleek/path.nix.tmpl", "home/users/fleek/path.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/home/users/fleek/programs.nix.tmpl", "home/users/fleek/programs.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/home/users/fleek/shell.nix.tmpl", "home/users/fleek/shell.nix", data, force)
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
func (f *Flake) Join(out io.Writer) error {
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
		err = f.writeSystem(*sys, "templates/home/hosts/host.nix.tmpl", true)
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
		err = f.writeUser(*user, "templates/home/users/user.nix.tmpl", true)
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
		err = f.setRebase(out)
		if err != nil {
			return err
		}
		err = f.mayCommit("fleek: new system", out)
		if err != nil {
			return err
		}
	}
	return nil

}

func (f *Flake) Check(out io.Writer) error {
	user, err := fleek.Username()

	if err != nil {
		return err
	}
	host, err := fleek.Hostname()
	if err != nil {
		return err
	}
	checkCmdLine := []string{"build", ".#homeConfigurations." + "\"" + user + "@" + host + "\"" + ".activationPackage"}

	err = f.runNix(nixbin, checkCmdLine, out)

	if err != nil {
		return err
	}
	return nil
}

// Write writes the applied flake configuration
func (f *Flake) Write(message string, out io.Writer) error {
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

	// users directory
	err = f.copyFile("default.nix", "users", force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/users/config.nix.tmpl", "users/config.nix", data, force)
	if err != nil {
		return err
	}
	// home/bling/default directory
	err = f.copyFile("default.nix", "home/bling/default", force)
	if err != nil {
		return err
	}
	err = f.copyFile("direnv.nix", "home/bling/default", force)
	if err != nil {
		return err
	}
	err = f.copyFile("gh.nix", "home/bling/default", force)
	if err != nil {
		return err
	}
	err = f.copyFile("packages.nix", "home/bling/default", force)
	if err != nil {
		return err
	}
	// home/bling/high directory
	err = f.copyFile("default.nix", "home/bling/high", force)
	if err != nil {
		return err
	}
	err = f.copyFile("bat.nix", "home/bling/high", force)
	if err != nil {
		return err
	}
	err = f.copyFile("exa.nix", "home/bling/high", force)
	if err != nil {
		return err
	}
	err = f.copyFile("packages.nix", "home/bling/high", force)
	if err != nil {
		return err
	}
	err = f.copyFile("programs.nix", "home/bling/high", force)
	if err != nil {
		return err
	}
	// home/bling/low directory
	err = f.copyFile("default.nix", "home/bling/low", force)
	if err != nil {
		return err
	}
	err = f.copyFile("packages.nix", "home/bling/low", force)
	if err != nil {
		return err
	}
	err = f.copyFile("programs.nix", "home/bling/low", force)
	if err != nil {
		return err
	}

	// home/bling/none directory
	err = f.copyFile("default.nix", "home/bling/none", force)
	if err != nil {
		return err
	}
	err = f.copyFile("bash.nix", "home/bling/none", force)
	if err != nil {
		return err
	}
	err = f.copyFile("path.nix", "home/bling/none", force)
	if err != nil {
		return err
	}
	err = f.copyFile("zsh.nix", "home/bling/none", force)
	if err != nil {
		return err
	}
	// home/global directory
	err = f.copyFile("default.nix", "home/global", force)
	if err != nil {
		return err
	}
	err = f.copyFile("git.nix", "home/global", force)
	if err != nil {
		return err
	}

	// home/users/fleek
	err = f.writeFile("templates/home/users/fleek/aliases.nix.tmpl", "home/users/fleek/aliases.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/home/users/fleek/default.nix.tmpl", "home/users/fleek/default.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/home/users/fleek/packages.nix.tmpl", "home/users/fleek/packages.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/home/users/fleek/path.nix.tmpl", "home/users/fleek/path.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/home/users/fleek/programs.nix.tmpl", "home/users/fleek/programs.nix", data, force)
	if err != nil {
		return err
	}
	err = f.writeFile("templates/home/users/fleek/shell.nix.tmpl", "home/users/fleek/shell.nix", data, force)
	if err != nil {
		return err
	}

	err = f.writeFile("templates/README.md.tmpl", "README.md", data, force)
	if err != nil {
		return err
	}

	spinner.Success()
	err = f.mayCommit(message, out)

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
func (f *Flake) copyFile(fname string, path string, force bool) error {

	filePath := filepath.Join(f.Config.UserFlakeDir(), path)
	err := os.MkdirAll(filePath, 0755)
	if err != nil {
		return err
	}
	templatePath := filepath.Join("templates", path, fname)

	fpath := filepath.Join(filePath, fname)
	_, err = os.Stat(fpath)
	if force || os.IsNotExist(err) {
		fin, err := templates.Open(templatePath)
		if err != nil {
			return err
		}
		defer fin.Close()

		fout, err := os.Create(fpath)
		if err != nil {
			return err
		}
		defer fout.Close()

		_, err = io.Copy(fout, fin)

		if err != nil {
			return err
		}

	} else {
		return errors.New("cowardly refusing to overwrite existing file without --force flag")
	}
	return nil
}
func (f *Flake) writeSystem(sys fleek.System, template string, force bool) error {

	hostPath := filepath.Join(f.Config.UserFlakeDir(), "home", "hosts")
	err := os.MkdirAll(hostPath, 0755)
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

		if err = f.Templates[template].Execute(file, sys); err != nil {
			return err
		}

	} else {
		return errors.New("cowardly refusing to overwrite existing file without --force flag")
	}

	return nil
}
func (f *Flake) writeUser(user fleek.User, template string, force bool) error {

	userPath := filepath.Join(f.Config.UserFlakeDir(), "home", "users", user.Username)
	err := os.MkdirAll(userPath, 0755)
	if err != nil {
		return err
	}
	fpath := filepath.Join(userPath, "custom.nix")
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
func (f *Flake) Apply(outWriter io.Writer) error {
	fin.Info.Println(f.app.Trans("flake.apply"))

	user, err := fleek.Username()

	if err != nil {
		return err
	}
	host, err := fleek.Hostname()
	if err != nil {
		return err
	}

	applyCmdLine := []string{"run", "--no-write-lock-file", "--impure", "home-manager", "--", "-b", "bak", "switch", "--flake", ".#" + user + "@" + host}
	err = f.runNix(nixbin, applyCmdLine, outWriter)
	if err != nil {
		return err
	}
	return nil
}
func (f *Flake) runNix(cmd string, cmdLine []string, out io.Writer) error {
	command := exec.Command(cmd, cmdLine...)
	command.Stdin = os.Stdin
	command.Stderr = out
	command.Stdout = out
	command.Dir = f.Config.UserFlakeDir()
	command.Env = os.Environ()
	if f.Config.Unfree {
		command.Env = append(command.Env, "NIXPKGS_ALLOW_UNFREE=1")
	}

	return command.Run()

}

//go:embed all:templates
var templates embed.FS
